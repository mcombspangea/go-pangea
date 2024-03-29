package redact

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	RedactText(ctx context.Context, input *TextInput) (*pangea.PangeaResponse[TextOutput], error)
	RedactStructured(ctx context.Context, input *StructuredInput) (*pangea.PangeaResponse[StructuredOutput], error)
}

type Redact struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*Redact, error) {
	cli := &Redact{
		Client: pangea.NewClient("redact", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*Redact) error
