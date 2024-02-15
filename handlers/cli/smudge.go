package cli

import (
	"bufio"
	"io"
	"log/slog"
	"os"

	"github.com/prskr/git-age/core/ports"
	"github.com/prskr/git-age/core/services"
	"github.com/prskr/git-age/infrastructure"
)

type SmudgeCliHandler struct {
	KeysFlag        `embed:""`
	Opener          ports.FileOpener `kong:"-"`
	FileToCleanPath string           `arg:"" name:"file" help:"Path to the file to clean"`
}

func (h *SmudgeCliHandler) Run() error {
	if err := requireStdin(); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	if isEncrypted, err := h.Opener.IsEncrypted(reader); err != nil {
		return err
	} else if !isEncrypted {
		slog.Warn("expected age-encrypted file, but got plaintext. Copying to stdout.")
		_, err = io.Copy(os.Stdout, reader)
		return err
	}

	decryptedReader, err := h.Opener.OpenFile(reader)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, decryptedReader)

	return err
}

func (h *SmudgeCliHandler) AfterApply() (err error) {
	h.Opener, err = services.NewAgeSealer(services.WithIdentities(infrastructure.NewIdentities(h.Keys)))
	return err
}
