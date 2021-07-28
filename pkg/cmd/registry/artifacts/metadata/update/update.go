package update

import (
	"github.com/redhat-developer/app-services-cli/pkg/localize"

	"github.com/redhat-developer/app-services-cli/pkg/cmd/flag"
	flagutil "github.com/redhat-developer/app-services-cli/pkg/cmdutil/flags"

	"github.com/redhat-developer/app-services-cli/pkg/iostreams"

	"github.com/redhat-developer/app-services-cli/pkg/logging"

	"github.com/spf13/cobra"

	"github.com/redhat-developer/app-services-cli/internal/config"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
)

type Options struct {
	// TODO find better name for this :D
	name string

	outputFormat string

	interactive bool

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewMetadataUpdateCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "update-metadata",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.name = args[0]
				// TODO check if file exist and can be read
			}

			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			return runCreate(opts)
		},
	}

	// TODO internationalize
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.MustLocalize("registry.cmd.flag.output.description"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

// nolint:funlen
func runCreate(opts *Options) error {
	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	logger.Info("")

	// connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	// if err != nil {
	// 	return err
	// }
	// TODO make network call to get the registry info
	// response, _, err := connection.API().ServiceRegistryMgmt().CreateRegistry(context.Background()).

	// TODO
	logger.Info("created")

	// switch opts.outputFormat {
	// case "json":
	// 	data, _ := json.MarshalIndent(response, "", cmdutil.DefaultJSONIndent)
	// 	_ = dump.JSON(opts.IO.Out, data)
	// case "yaml", "yml":
	// 	data, _ := yaml.Marshal(response)
	// 	_ = dump.YAML(opts.IO.Out, data)
	// }

	return nil
}
