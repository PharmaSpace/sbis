package sbis

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Sbis struct {
	ctx     context.Context
	verbose bool
	logger  zerolog.Logger
	config  *AuthConfig
}

func NewSbis(options ...Option) (*Sbis, error) {
	s := &Sbis{
		logger: zerolog.New(os.Stdout).With().Timestamp().Str("service", "sbis").Logger(),
	}

	for _, o := range options {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	if s.config == nil {
		s.config = s.setDefaultConfig()
	}

	if len(s.config.AppClientID) == 0 {
		s.logger.Fatal().Msgf(ErrAuthConfigNotFound.Error())
		return nil, ErrAuthConfigNotFound
	}

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	return s, nil
}

func (s Sbis) setDefaultConfig() *AuthConfig {
	return &AuthConfig{
		AppClientID: os.Getenv("SBIS_APP_CLIENT_ID"),
		Login:       os.Getenv("SBIS_LOGIN"),
		Password:    os.Getenv("SBIS_PASSWORD"),
	}
}
