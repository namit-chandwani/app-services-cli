package describe

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	flagutil "github.com/redhat-developer/app-services-cli/pkg/cmdutil/flags"
	"github.com/redhat-developer/app-services-cli/pkg/connection"
	"github.com/redhat-developer/app-services-cli/pkg/iostreams"
	"github.com/redhat-developer/app-services-cli/pkg/localize"
	"github.com/redhat-developer/app-services-cli/pkg/serviceregistry"

	"github.com/redhat-developer/app-services-cli/pkg/cmd/flag"

	"github.com/redhat-developer/app-services-cli/internal/config"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
	"github.com/redhat-developer/app-services-cli/pkg/dump"
	srsmgmtv1 "github.com/redhat-developer/app-services-sdk-go/registrymgmt/apiv1/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Options struct {
	id           string
	name         string
	outputFormat string

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	localizer  localize.Localizer
}

// NewDescribeCommand describes a service instance, either by passing an `--id flag`
// or by using the service instance set in the config, if any
func NewDescribeCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		Config:     f.Config,
		Connection: f.Connection,
		IO:         f.IOStreams,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   f.Localizer.MustLocalize("registry.cmd.describe.shortDescription"),
		Long:    f.Localizer.MustLocalize("registry.cmd.describe.longDescription"),
		Example: f.Localizer.MustLocalize("registry.cmd.describe.example"),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			if len(args) > 0 {
				opts.name = args[0]
			}

			if opts.name != "" && opts.id != "" {
				return errors.New(opts.localizer.MustLocalize("service.error.idAndNameCannotBeUsed"))
			}

			if opts.id != "" || opts.name != "" {
				return runDescribe(opts)
			}

			cfg, err := opts.Config.Load()
			if err != nil {
				return err
			}

			var registryConfig *config.ServiceRegistryConfig
			if cfg.Services.ServiceRegistry == registryConfig || cfg.Services.ServiceRegistry.InstanceID == "" {
				return errors.New(opts.localizer.MustLocalize("registry.common.error.noServiceSelected"))
			}

			opts.id = cfg.Services.ServiceRegistry.InstanceID

			return runDescribe(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.MustLocalize("registry.cmd.flag.output.description"))
	cmd.Flags().StringVar(&opts.id, "id", "", opts.localizer.MustLocalize("registry.common.flag.id"))

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runDescribe(opts *Options) error {
	connection, err := opts.Connection(connection.DefaultConfigSkipMasAuth)
	if err != nil {
		return err
	}

	api := connection.API()

	var registry *srsmgmtv1.RegistryRest
	ctx := context.Background()
	if opts.name != "" {
		registry, _, err = serviceregistry.GetServiceRegistryByName(ctx, api.ServiceRegistryMgmt(), opts.name)
		if err != nil {
			return err
		}
	} else {
		registry, _, err = serviceregistry.GetServiceRegistryByID(ctx, api.ServiceRegistryMgmt(), opts.id)
		if err != nil {
			return err
		}
	}

	return printService(opts.IO.Out, registry, opts.outputFormat)
}

func printService(w io.Writer, registry interface{}, outputFormat string) error {
	switch outputFormat {
	case dump.YAMLFormat, dump.YMLFormat:
		data, err := yaml.Marshal(registry)
		if err != nil {
			return err
		}
		return dump.YAML(w, data)
	default:
		data, err := json.Marshal(registry)
		if err != nil {
			return err
		}
		return dump.JSON(w, data)
	}
}
