package mdcli

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"

	mdlib "github.com/spinnaker/md-lib-go"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

// Print is a command line interface for printing a non-anchored view of the MD config file
func Print(opts *CommandOptions) (int, error) {
	configPath := filepath.Join(opts.ConfigDir, opts.ConfigFile)
	if _, err := os.Stat(configPath); err != nil {
		return 1, err
	}

	cli := mdlib.NewClient(
		mdlib.WithBaseURL(opts.BaseURL),
		mdlib.WithHTTPClient(opts.HTTPClient),
	)

	mdProcessor := mdlib.NewDeliveryConfigProcessor(
		mdlib.WithDirectory(opts.ConfigDir),
		mdlib.WithFile(opts.ConfigFile),
		mdlib.WithLogger(opts.Logger),
	)

	valErr, err := mdProcessor.Validate(cli)
	if err != nil {
		if valErr != nil {
			opts.Logger.Errorf("%s\nReason: %s", valErr.Error, valErr.Message)
			return 1, err
		}
		return 1, err
	}
	opts.Logger.Noticef("PASSED VALIDATION")

	bytes, err := mdProcessor.GetDeliveryConfigYamlDeAnchored()
	if err != nil {
		opts.Logger.Errorf("%s\n", err)
		return 1, err
	}

	opts.Logger.Noticef("%s", "\n"+string(bytes))

	return 0, nil
}
