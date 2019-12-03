package sbis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const (
	defaultBaseURL = "https://api.sbis.ru/"
	userAgent      = "go-sbis"
	defaultHost    = "api.sbis.ru"
)

type Client struct {
	clientMu  sync.Mutex
	client    *http.Client
	baseURL   *url.URL
	userAgent string
	ctx       context.Context
	verbose   bool
	logger    zerolog.Logger
	config    *AuthConfig
	sid       string
	inn       string

	Authorization         *Authorization
	ListKKTbyOrganization *ListKKTbyOrganization
	ListOfFiscalDriver    *ListOfFiscalDriver
	ListOfFiscalDoc       *ListOfFiscalDoc
}

func NewClient(options ...Option) (*Client, error) {
	c := &Client{
		logger: zerolog.New(os.Stdout).With().Timestamp().Str("service", "sbis").Logger(),
	}

	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	if c.config == nil {
		c.config = c.setDefaultConfig()
	}
	if len(c.inn) == 0 {
		c.inn = os.Getenv("SBIS_INN")
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c.baseURL = baseURL

	httpClient := &http.Client{}
	c.client = httpClient

	c.userAgent = userAgent

	if len(c.config.AppClientID) == 0 {
		c.logger.Fatal().Msgf(ErrAuthConfigNotFound.Error())
		return nil, ErrAuthConfigNotFound
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	// services
	c.Authorization = NewAuthorization(c)
	c.ListKKTbyOrganization = NewListKKTbyOrganization(c)

	return c, nil
}

func (c *Client) NewRequest(withAuth bool, method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if withAuth {
		resp, err := c.Authorization.GetSID()
		if err != nil {
			return nil, err
		}
		c.sid = resp.Sid
	}

	req.Header.Set("Host", defaultHost)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(c.sid) > 0 {
		req.Header.Set("Cookie", "sid="+c.sid)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

func (s *Client) setDefaultConfig() *AuthConfig {
	return &AuthConfig{
		AppClientID: os.Getenv("SBIS_APP_CLIENT_ID"),
		Login:       os.Getenv("SBIS_LOGIN"),
		Password:    os.Getenv("SBIS_PASSWORD"),
	}
}
