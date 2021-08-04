package update

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redhat-developer/app-services-cli/pkg/cmdutil"
	"github.com/redhat-developer/app-services-cli/pkg/connection"
	"github.com/redhat-developer/app-services-cli/pkg/dump"
	"github.com/redhat-developer/app-services-cli/pkg/localize"
	registryinstanceclient "github.com/redhat-developer/app-services-sdk-go/registryinstance/apiv1internal/client"
	"gopkg.in/yaml.v2"

	"github.com/redhat-developer/app-services-cli/pkg/cmd/flag"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/util"
	flagutil "github.com/redhat-developer/app-services-cli/pkg/cmdutil/flags"

	"github.com/redhat-developer/app-services-cli/pkg/iostreams"

	"github.com/redhat-developer/app-services-cli/pkg/logging"

	"github.com/spf13/cobra"

	"github.com/redhat-developer/app-services-cli/internal/config"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
)

type Options struct {
	artifact string
	group    string

	file         string
	artifactType string

	registryID   string
	outputFormat string

	IO         *iostreams.IOStreams
	Config     config.IConfig
	Connection factory.ConnectionFunc
	Logger     func() (logging.Logger, error)
	localizer  localize.Localizer
}

func NewUpdateCommand(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		Connection: f.Connection,
		Logger:     f.Logger,
		localizer:  f.Localizer,
	}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates new artifact",
		Long:    "Creates new artifact from file or directly standard input",
		Example: "",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			validOutputFormats := flagutil.ValidOutputFormats
			if opts.outputFormat != "" && !flagutil.IsValidInput(opts.outputFormat, validOutputFormats...) {
				return flag.InvalidValueError("output", opts.outputFormat, validOutputFormats...)
			}

			if len(args) > 0 {
				opts.file = args[0]
			}

			if opts.registryID != "" {
				return runUpdate(opts)
			}

			cfg, err := opts.Config.Load()
			if err != nil {
				return err
			}

			if opts.artifactType != "" {
				if _, err = registryinstanceclient.NewArtifactTypeFromValue(opts.artifactType); err != nil {
					return fmt.Errorf("Invalid artifact type. Please use one of following values: " + util.GetAllowedArtifactTypeEnumValuesAsString())
				}
			}

			if !cfg.HasServiceRegistry() {
				return fmt.Errorf("No service Registry selected. Use 'rhoas service-registry use' use to select your registry")
			}

			opts.registryID = fmt.Sprint(cfg.Services.ServiceRegistry.InstanceID)
			return runUpdate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "json", opts.localizer.MustLocalize("registry.cmd.flag.output.description"))
	cmd.Flags().StringVarP(&opts.file, "file", "f", "", "File location of the artifact")

	cmd.Flags().StringVarP(&opts.artifact, "artifact", "a", "", "Id of the artifact")
	cmd.Flags().StringVarP(&opts.group, "group", "g", "", "Id of the artifact")
	cmd.Flags().StringVarP(&opts.artifactType, "type", "t", "", "Type of artifact. Choose from:  "+util.GetAllowedArtifactTypeEnumValuesAsString())
	cmd.Flags().StringVarP(&opts.registryID, "registryId", "", "", "Id of the registry to be used. By default uses currently selected registry.")

	flagutil.EnableOutputFlagCompletion(cmd)

	return cmd
}

func runUpdate(opts *Options) error {

	logger, err := opts.Logger()
	if err != nil {
		return err
	}

	conn, err := opts.Connection(connection.DefaultConfigRequireMasAuth)
	if err != nil {
		return err
	}

	dataAPI, _, err := conn.API().ServiceRegistryInstance(opts.registryID)
	if err != nil {
		return err
	}

	if opts.group == "" {
		logger.Info("Group was not specified. Using 'default' artifacts group.")
		opts.group = util.DefaultArtifactGroup
	}

	var specifiedFile *os.File
	if opts.file != "" {
		logger.Info("Opening file: " + opts.file)
		specifiedFile, err = os.Open(opts.file)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Reading file content from stdin")
		specifiedFile, err = util.CreateFileFromStdin()
		if err != nil {
			return err
		}
	}

	ctx := context.Background()
	request := dataAPI.ArtifactsApi.UpdateArtifact(ctx, opts.group, opts.artifact)
	request = request.Body(specifiedFile)
	metadata, _, err := request.Execute()
	if err != nil {
		return err
	}
	logger.Info("Artifact updated")

	switch opts.outputFormat {
	case "json":
		data, _ := json.MarshalIndent(metadata, "", cmdutil.DefaultJSONIndent)
		_ = dump.JSON(opts.IO.Out, data)
	case "yaml", "yml":
		data, _ := yaml.Marshal(metadata)
		_ = dump.YAML(opts.IO.Out, data)
	}

	return nil
}
