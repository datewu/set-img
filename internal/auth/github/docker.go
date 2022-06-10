package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Octet types from RFC 7230.
type octetType byte

var octetTypes [256]octetType

const (
	isToken octetType = 1 << iota
	isSpace
)

type auth struct {
	Username string
	Password string
}

type challenge struct {
	// Scheme is the auth-scheme according to RFC 7235
	Scheme string

	// Parameters are the auth-params according to RFC 7235
	Parameters map[string]string
}

// microDockerClient only check Bearer tokens.
type microDockerClient struct {
	userAgent string
	auth
	scope         authScope
	challenges    []challenge
	client        http.Client
	registryToken string
}

type authScope struct {
	remoteName string
	actions    string
}

type bearerToken struct {
	Token          string    `json:"token"`
	AccessToken    string    `json:"access_token"`
	ExpiresIn      int       `json:"expires_in"`
	IssuedAt       time.Time `json:"issued_at"`
	expirationTime time.Time
}

func newMicroDockerClient(name, pwd string) *microDockerClient {
	m := new(microDockerClient)
	m.auth.Username = name
	m.auth.Password = pwd
	m.client = http.Client{
		Timeout: time.Second * 3,
	}
	return m
}

func (c *microDockerClient) CheckToken(ctx context.Context, name, pwd string) (bool, error) {
	//	microCli.userAgent = "micro-docker-client"
	err := c.detectPropertiesHelper(ctx)
	if err != nil {
		return false, err
	}
	for _, ch := range c.challenges {
		err = c.checkBearerToken(ctx, ch)
		if err == nil {
			return true, nil
		}
	}
	return false, errors.New("not implemented")
}

func (c *microDockerClient) checkBearerToken(ctx context.Context, challenge challenge) error {
	realm, ok := challenge.Parameters["realm"]
	if !ok {
		return errors.New("missing realm in bearer auth challenge")
	}

	authReq, err := http.NewRequestWithContext(ctx, http.MethodGet, realm, nil)
	if err != nil {
		return err
	}

	params := authReq.URL.Query()
	if c.auth.Username != "" {
		params.Add("account", c.auth.Username)
	}

	if service, ok := challenge.Parameters["service"]; ok && service != "" {
		params.Add("service", service)
	}

	// for _, scope := range scopes {
	// 	if scope.remoteName != "" && scope.actions != "" {
	// 		params.Add("scope", fmt.Sprintf("repository:%s:%s", scope.remoteName, scope.actions))
	// 	}
	// }

	authReq.URL.RawQuery = params.Encode()

	if c.auth.Username != "" && c.auth.Password != "" {
		authReq.SetBasicAuth(c.auth.Username, c.auth.Password)
	}
	if c.userAgent != "" {
		authReq.Header.Add("User-Agent", c.userAgent)
	}

	log.Info().Str("method", authReq.Method).Str("url", authReq.URL.Redacted()).
		Msg("checkBearerToken going to request")
	res, err := c.client.Do(authReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status code")
	}
	body, err := io.ReadAll(res.Body)
	log.Info().Msgf("%s", string(body))
	return nil
}

// detectPropertiesHelper performs the work of detectProperties which executes
// it at most once.
func (c *microDockerClient) detectPropertiesHelper(ctx context.Context) error {
	resp, err := http.Get(ghcrRegistry)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return errors.New("not 200 and not 401")
	}
	c.challenges = parseAuthHeader(resp.Header)
	return nil
}

func parseAuthHeader(header http.Header) []challenge {
	challenges := []challenge{}
	for _, h := range header[http.CanonicalHeaderKey("WWW-Authenticate")] {
		v, p := parseValueAndParams(h)
		if v != "" {
			challenges = append(challenges, challenge{Scheme: v, Parameters: p})
		}
	}
	return challenges
}

// NOTE: This is not a fully compliant parser per RFC 7235:
// Most notably it does not support more than one challenge within a single header
// Some of the whitespace parsing also seems noncompliant.
// But it is clearly better than what we used to have…
func parseValueAndParams(header string) (value string, params map[string]string) {
	params = make(map[string]string)
	value, s := expectToken(header)
	if value == "" {
		return
	}
	value = strings.ToLower(value)
	s = "," + skipSpace(s)
	for strings.HasPrefix(s, ",") {
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" {
			return
		}
		if !strings.HasPrefix(s, "=") {
			return
		}
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" {
			return
		}
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	}
	return
}

func expectToken(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isToken == 0 {
			break
		}
	}
	return s[:i], s[i:]
}

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectTokenOrQuoted(s string) (value string, rest string) {
	if !strings.HasPrefix(s, "\"") {
		return expectToken(s)
	}
	s = s[1:]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ {
				b := s[i]
				switch {
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				}
			}
			return "", ""
		}
	}
	return "", ""
}
