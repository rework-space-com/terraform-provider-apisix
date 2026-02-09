package model

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SecretGCPType struct {
	AuthConfig *SecretGCPAuthConfigType `tfsdk:"auth_config"`
	AuthFile   types.String             `tfsdk:"auth_file"`
	SslVerify  types.Bool               `tfsdk:"ssl_verify"`
}

var SecretGCPSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Set APISIX Secret Management GCP configuration.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"auth_config": SecretGCPAuthConfigSchemaAttribute,
		"auth_file": schema.StringAttribute{
			Description: "Path to the Google Cloud service account authentication JSON file. Either auth_config or auth_file must be provided..",
			Optional:    true,
		},
		"ssl_verify": schema.BoolAttribute{
			Description: "Path to the Google Cloud service account authentication JSON file. Either auth_config or auth_file must be provided..",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
	},
}

func SecretGCPFromTerraformToApi(ctx context.Context, id *string, terraformDataModel *SecretGCPType) (apiDataModel *api_client.GCPSecret) {
	if terraformDataModel == nil || id == nil {
		tflog.Debug(ctx, "Can't transform Secret GCP to api model")
		return
	}

	result := api_client.GCPSecret{
		BaseSecret: api_client.BaseSecret{ID: fmt.Sprintf("%s/%s", api_client.GCP, *id)},
		AuthConfig: SecretGCPAuthConfigFromTerraformToApi(ctx, terraformDataModel.AuthConfig),
		AuthFile:   terraformDataModel.AuthFile.ValueStringPointer(),
		SslVerify:  terraformDataModel.SslVerify.ValueBoolPointer(),
	}

	return &result
}

func SecretGCPFromApiToTerraform(ctx context.Context, apiDataModel *api_client.GCPSecret) (terraformDataModel *SecretGCPType) {
	if apiDataModel == nil {
		tflog.Debug(ctx, "Can't transform Secret GCP from api model")
		return
	}

	result := SecretGCPType{
		AuthConfig: SecretGCPAuthConfigFromApiToTerraform(ctx, apiDataModel.AuthConfig),
		AuthFile:   types.StringPointerValue(apiDataModel.AuthFile),
		SslVerify:  types.BoolPointerValue(apiDataModel.SslVerify),
	}

	return &result
}
