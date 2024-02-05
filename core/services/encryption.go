package services

import (
	"io"

	"filippo.io/age"
	"github.com/prskr/git-age/core/ports"
)

var (
	_ ports.FileOpener = (*AgeSealer)(nil)
	_ ports.FileSealer = (*AgeSealer)(nil)
)

type AgeSealerOption func(*AgeSealer) error

func WithRecipients(r ports.Recipients) AgeSealerOption {
	return func(sealer *AgeSealer) error {
		recipients, err := r.All()
		if err != nil {
			return err
		}

		sealer.Recipients = recipients
		return nil
	}
}

func WithIdentities(id ports.Identities) AgeSealerOption {
	return func(sealer *AgeSealer) error {
		identities, err := id.All()
		if err != nil {
			return err
		}

		sealer.Identities = identities
		return nil
	}
}

func NewAgeSealer(opts ...AgeSealerOption) (*AgeSealer, error) {
	sealer := new(AgeSealer)

	for _, opt := range opts {
		if err := opt(sealer); err != nil {
			return nil, err
		}
	}

	return sealer, nil
}

type AgeSealer struct {
	Recipients []age.Recipient
	Identities []age.Identity
}

func (h *AgeSealer) CanOpen() bool {
	return len(h.Identities) > 0
}

func (h *AgeSealer) CanSeal() bool {
	return len(h.Recipients) > 0
}

func (h *AgeSealer) AddRecipients(r ...age.Recipient) {
	h.Recipients = append(h.Recipients, r...)
}

func (h *AgeSealer) AddIdentities(identities ...age.Identity) {
	h.Identities = append(h.Identities, identities...)
}

func (h *AgeSealer) OpenFile(reader io.Reader) (io.Reader, error) {
	return age.Decrypt(reader, h.Identities...)
}

func (h *AgeSealer) SealFile(dst io.Writer) (io.WriteCloser, error) {
	return age.Encrypt(dst, h.Recipients...)
}
