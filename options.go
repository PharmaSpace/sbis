package sbis

import (
	"context"
)

type Option func(sbis *Client) error

func Config(config *AuthConfig) Option {
	return func(r *Client) error {
		r.config = config
		return nil
	}
}

func SetInn(inn string) Option {
	return func(r *Client) error {
		r.inn = inn
		return nil
	}
}

func SetAuthConfig(config *AuthConfig) Option {
	return func(r *Client) error {
		r.config = config
		return nil
	}
}

func Context(ctx context.Context) Option {
	return func(r *Client) error {
		r.ctx = ctx
		return nil
	}
}

func Verbose() Option {
	return func(r *Client) error {
		r.verbose = true
		return nil
	}
}
