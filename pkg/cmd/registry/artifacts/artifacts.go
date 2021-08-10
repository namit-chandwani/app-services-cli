package artifacts

import (
	"github.com/redhat-developer/app-services-cli/pkg/cmd/factory"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/create"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/delete"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/get"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/list"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/crud/update"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/download"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/metadata"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/registry/artifacts/versions"
	"github.com/spf13/cobra"
)

func NewArtifactsCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Manage Service Registry Artifacts commands",
		Long: `Apicurio Registry Artifacts enables developers to manage and share the structure of their data. 
				For example, client applications can dynamically push or pull the latest updates to or from the registry without needing to redeploy.
				Apicurio Registry also enables developers to create rules that govern how registry content can evolve over time. 
				For example, this includes rules for content validation and version compatibility.
				
				Registry commands enable client applications to manage the artifacts in the registry. 
				This set of commands provide create, read, update, and delete operations for schema and API artifacts, rules, versions, and metadata.`,
		Example: `
		## Create artifact in my-group from schema.json file
		rhoas service-registry artifacts create my-group schema.json

		## List Artifacts
		rhoas service-registry artifacts list my-group
		`,
		Args: cobra.MinimumNArgs(1),
	}

	// add sub-commands
	cmd.AddCommand(
		// CRUD
		create.NewCreateCommand(f),
		get.NewGetCommand(f),
		delete.NewDeleteCommand(f),
		list.NewListCommand(f),
		update.NewUpdateCommand(f),

		// Misc
		metadata.NewMetadataCommand(f),
		versions.NewVersionsCommand(f),
		download.NewDownloadCommand(f),
	)

	return cmd
}
