package github

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/set-img/internal/auth/utils"
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
	jsonlog.Debug("challenge", map[string]string{"len": strconv.Itoa(len(c.challenges))})
	for _, ch := range c.challenges {
		err = c.checkBearerToken(ctx, ch)
		if err == nil {
			return true, nil
		}
		jsonlog.Err(err, ch.Parameters)
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

	jsonlog.Debug("checkBearerToken going to request", map[string]string{
		"method": authReq.Method, "url": authReq.URL.Redacted()})
	res, err := c.client.Do(authReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status code")
	}
	body, err := io.ReadAll(res.Body)
	jsonlog.Debug("response", map[string]string{"body": string(body)})
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
		jsonlog.Debug("values in WWW-Authenticate", map[string]string{"header": h})
		p, v := utils.ConsumeParams(h)
		if v != "" {
			challenges = append(challenges, challenge{Scheme: v, Parameters: p})
		}
	}
	return challenges
}
