package sbis

import (
	"context"
)

type Option func(sbis *Sbis) error

func Config(config *AuthConfig) Option {
	return func(r *Sbis) error {
		r.config = config
		return nil
	}
}

func Context(ctx context.Context) Option {
	return func(r *Sbis) error {
		r.ctx = ctx
		return nil
	}
}

func Verbose() Option {
	return func(r *Sbis) error {
		r.verbose = true
		return nil
	}
}
