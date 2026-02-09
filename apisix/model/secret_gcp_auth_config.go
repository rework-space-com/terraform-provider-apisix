package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SecretGCPAuthConfigType struct {
	ClientEmail types.String `tfsdk:"client_email"`
	PrivateKey  types.String `tfsdk:"private_key"`
	ProjectId   types.String `tfsdk:"project_id"`
	TokenUri    types.String `tfsdk:"token_uri"`
	EntriesUri  types.String `tfsdk:"entries_uri"`
	Scope       types.List   `tfsdk:"scope"`
}

var SecretGCPAuthConfigSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Set APISIX Secret Management GCP Auth Config configuration.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"client_email": schema.StringAttribute{
			Description: "Email address of the Google Cloud service account.",
			Required:    true,
		},
		"private_key": schema.StringAttribute{
			Description: "Private key of the Google Cloud service account.",
			Required:    true,
			Sensitive:   true,
		},
		"project_id": schema.StringAttribute{
			Description: "Project ID in the Google Cloud service account.",
			Required:    true,
		},
		"token_uri": schema.StringAttribute{
			Description: "Token URI of the Google Cloud service account.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("https://oauth2.googleapis.com/token"),
		},
		"entries_uri": schema.StringAttribute{
			Description: "The API access endpoint for the Google Secrets Manager.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("https://secretmanager.googleapis.com/v1"),
		},
		"scope": schema.ListAttribute{
			Description: "Access scopes of the Google Cloud service account.",
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			Default: listdefault.StaticValue(
				types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("https://www.googleapis.com/auth/cloud-platform"),
				}),
			),
		},
	},
}

func SecretGCPAuthConfigFromTerraformToApi(ctx context.Context, terraformDataModel *SecretGCPAuthConfigType) (apiDataModel *api_client.AuthConfigType) {
	if terraformDataModel == nil {
		tflog.Debug(ctx, "Can't transform Secret GCP AuthConfig to api model")
		return
	}

	var scopes []string
	_ = terraformDataModel.Scope.ElementsAs(ctx, &scopes, false)

	result := api_client.AuthConfigType{
		ClientEmail: terraformDataModel.ClientEmail.ValueStringPointer(),
		PrivateKey:  terraformDataModel.PrivateKey.ValueStringPointer(),
		ProjectId:   terraformDataModel.ProjectId.ValueStringPointer(),
		TokenUri:    terraformDataModel.TokenUri.ValueStringPointer(),
		EntriesUri:  terraformDataModel.EntriesUri.ValueStringPointer(),
		Scope:       &scopes,
	}

	return &result
}

func SecretGCPAuthConfigFromApiToTerraform(ctx context.Context, apiDataModel *api_client.AuthConfigType) (terraformDataModel *SecretGCPAuthConfigType) {
	if apiDataModel == nil {
		tflog.Debug(ctx, "Can't transform Secret GCP AuthConfig from api model")
		return nil
	}

	scopes, _ := types.ListValueFrom(ctx, types.StringType, apiDataModel.Scope)

	result := SecretGCPAuthConfigType{
		ClientEmail: types.StringPointerValue(apiDataModel.ClientEmail),
		PrivateKey:  types.StringPointerValue(apiDataModel.PrivateKey),
		ProjectId:   types.StringPointerValue(apiDataModel.ProjectId),
		TokenUri:    types.StringPointerValue(apiDataModel.TokenUri),
		EntriesUri:  types.StringPointerValue(apiDataModel.EntriesUri),
		Scope:       scopes,
	}

	return &result
}
