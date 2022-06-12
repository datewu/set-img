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
	log.Debug().Int("len", len(c.challenges)).Msg("challenges")
	for _, ch := range c.challenges {
		err = c.checkBearerToken(ctx, ch)
		if err == nil {
			return true, nil
		}
		log.Err(err).Msg("checkBearerToken failed")
		log.Debug().Str("scheme", ch.Scheme).Interface("params", ch.Parameters).
			Msg("challenge detail")
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

	log.Debug().Str("method", authReq.Method).Str("url", authReq.URL.Redacted()).
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
	log.Debug().Msgf("%s", string(body))
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
	log.Debug().Interface("headers", resp.Header).Msg("headers")
	c.challenges = parseAuthHeader(resp.Header)
	return nil
}

func parseAuthHeader(header http.Header) []challenge {
	challenges := []challenge{}
	for _, h := range header[http.CanonicalHeaderKey("WWW-Authenticate")] {
		log.Debug().Str("header", h).Msg("values in WWW-Authenticate")
		p, v := consumeParams(h)
		if v != "" {
			challenges = append(challenges, challenge{Scheme: v, Parameters: p})
		}
	}
	return challenges
}

func consumeParams(v string) (params map[string]string, newv string) {
	for {
		var name, value string
		name, value, v = consumeParam(v)
		if name == "" {
			break
		}
		// Use only the first occurrence of each param name.
		// This is required in some other places that don't use consumeParams,
		// but it seems like reasonable behavior in general.
		if _, seen := params[name]; seen {
			continue
		}
		if params == nil {
			params = make(map[string]string)
		}
		params[name] = value
	}
	return params, v
}

func consumeParam(v string) (name, value, newv string) {
	v = skipWSAnd(v, ';')
	name, v = consumeItem(v)
	if name == "" {
		return "", "", v
	}
	name = strings.ToLower(name)
	v = skipWS(v)
	if peek(v) == '=' {
		v = skipWS(v[1:])
		value, v = consumeItemOrQuoted(v)
	}
	return name, value, v
}

func skipWSAnd(v string, and byte) string {
	for v != "" && (v[0] == ' ' || v[0] == '\t' || v[0] == and) {
		v = v[1:]
	}
	return v
}

func skipWS(v string) string {
	for v != "" && (v[0] == ' ' || v[0] == '\t') {
		v = v[1:]
	}
	return v
}

// consumeItem returns the item from the beginning of v, and the rest of v.
// An item is a run of text up to whitespace, comma, semicolon, or equal sign.
// Callers should check that the item is non-empty if they need to make progress.
func consumeItem(v string) (item, newv string) {
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case ' ', '\t', ',', ';', '=':
			return v[:i], v[i:]
		}
	}
	return v, ""
}

func peek(v string) byte {
	if v == "" {
		return 0
	}
	return v[0]
}

func consumeItemOrQuoted(v string) (text, newv string) {
	if peek(v) == '"' {
		text, newv = consumeQuoted(v)
		return
	}
	return consumeItem(v)
}

func consumeQuoted(v string) (text, newv string) {
	return consumeDelimited(v, '"', '"')
}

func consumeDelimited(v string, opener, closer byte) (text, newv string) {
	if peek(v) != opener {
		return "", v
	}
	v = v[1:]

	// In the common case, when there are no quoted pairs,
	// we can simply slice the string between the outermost delimiters.
	nesting := 1
	i := 0
	for ; i < len(v); i++ {
		switch v[i] {
		case closer:
			nesting--
			if nesting == 0 {
				return v[:i], v[i+1:]
			}
		case opener:
			nesting++
		case '\\': // start of a quoted pair
			goto buffered
		}
	}
	// We've reached the end of v, but nesting is still > 0.
	// This is an unterminated string.
	return v, ""

buffered:
	// Once we have encountered a quoted pair, we have to unquote into a buffer.
	b := &strings.Builder{}
	b.WriteString(v[:i])
	quoted := false
	for ; i < len(v); i++ {
		switch {
		case quoted:
			b.WriteByte(v[i])
			quoted = false
		case v[i] == closer:
			nesting--
			if nesting == 0 {
				return b.String(), v[i+1:]
			}
			b.WriteByte(v[i])
		case v[i] == opener:
			nesting++
			b.WriteByte(v[i])
		case v[i] == '\\':
			quoted = true
		default:
			b.WriteByte(v[i])
		}
	}
	// We've reached the end of v, but nesting is still > 0.
	// This is an unterminated string.
	return b.String(), ""
}
