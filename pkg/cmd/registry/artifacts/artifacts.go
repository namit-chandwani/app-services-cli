package artifacts

import (
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/create"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/delete"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/get"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/list"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/update"
	metadataGet "github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/metadata/get"
	metadataUpdate "github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/metadata/update"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/state"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/versions"

	"github.com/spf13/cobra"
)

func NewArtifactsCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "artifacts",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
	}

	// add sub-commands
	cmd.AddCommand(
		// CRUD
		create.NewCreateCommand(f),
		get.NewGetCommand(f),
		delete.NewDeleteCommand(f),
		list.NewListCommand(f),
		update.NewUpdateCommand(f),

		// Metadata
		metadataGet.NewMetadataGetCommand(f),
		metadataUpdate.NewMetadataUpdateCommand(f),

		// Misc
		state.NewStateUpdateCommand(f),
		versions.NewVersionsGetCommand(f),

		// Rules
		// TODO
	)

	return cmd
}
