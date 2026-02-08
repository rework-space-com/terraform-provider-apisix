package model

import (
	"context"
	"strings"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SecretResourceModel struct {
	ID    types.String     `tfsdk:"id"`
	Vault *SecretVaultType `tfsdk:"vault"`
	AWS   *SecretAWSType   `tfsdk:"aws"`
	GCP   *SecretGCPType   `tfsdk:"gcp"`
}

var SecretSchema = schema.Schema{
	Description: "Manages APISIX Secret Management",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the Secret.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"vault": SecretVaultSchemaAttribute,
		"aws":   SecretAWSSchemaAttribute,
		"gcp":   SecretGCPSchemaAttribute,
	},
}

func SecretFromTerraformToApi(ctx context.Context, terraformDataModel *SecretResourceModel) (secretManager api_client.SecretManager, apiDataModel api_client.Secret) {
	if terraformDataModel.Vault != nil {
		apiDataModel = SecretVaultFromTerraformToApi(ctx, terraformDataModel.ID.ValueStringPointer(), terraformDataModel.Vault)
		secretManager = api_client.Vault
	} else if terraformDataModel.AWS != nil {
		apiDataModel = SecretAWSFromTerraformToApi(ctx, terraformDataModel.ID.ValueStringPointer(), terraformDataModel.AWS)
		secretManager = api_client.AWS
	} else if terraformDataModel.GCP != nil {
		apiDataModel = SecretGCPFromTerraformToApi(ctx, terraformDataModel.ID.ValueStringPointer(), terraformDataModel.GCP)
		secretManager = api_client.GCP
	}

	tflog.Debug(ctx, "Result of SecretFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return secretManager, apiDataModel
}

func SecretFromApiToTerraform(ctx context.Context, apiDataModel api_client.Secret) (terraformDataModel SecretResourceModel) {
	apiID := apiDataModel.GetID()
	parts := strings.Split(apiID, "/")

	if len(parts) != 2 {
		tflog.Error(ctx, "API returned an invalid ID format", map[string]any{"id": apiID})
		terraformDataModel.ID = types.StringValue(apiID)
		return terraformDataModel
	}

	secretManager := api_client.SecretManager(parts[0])
	terraformID := parts[1]

	terraformDataModel.ID = types.StringValue(terraformID)

	switch secretManager {
	case api_client.Vault:
		if api, ok := apiDataModel.(*api_client.VaultSecret); ok {
			terraformDataModel.Vault = SecretVaultFromApiToTerraform(ctx, api)
		}
	case api_client.AWS:
		if api, ok := apiDataModel.(*api_client.AWSSecret); ok {
			terraformDataModel.AWS = SecretAWSFromApiToTerraform(ctx, api)
		}
	case api_client.GCP:
		if api, ok := apiDataModel.(*api_client.GCPSecret); ok {
			terraformDataModel.GCP = SecretGCPFromApiToTerraform(ctx, api)
		}
	}

	tflog.Debug(ctx, "Result of SecretFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
