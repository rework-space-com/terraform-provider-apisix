package model

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SecretVaultType struct {
	Uri       types.String `tfsdk:"uri"`
	Prefix    types.String `tfsdk:"prefix"`
	Token     types.String `tfsdk:"token"`
	Namespace types.String `tfsdk:"namespace"`
}

var SecretVaultSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Set APISIX Secret Management Vault configuration.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"uri": schema.StringAttribute{
			Description: "URI of the vault server.",
			Required:    true,
		},
		"prefix": schema.StringAttribute{
			Description: "key prefix.",
			Required:    true,
		},
		"token": schema.StringAttribute{
			Description: "vault token.",
			Required:    true,
		},
		"namespace": schema.StringAttribute{
			Description: "Vault namespace, no default value.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("admin"),
		},
	},
}

func SecretVaultFromTerraformToApi(ctx context.Context, id *string, terraformDataModel *SecretVaultType) (apiDataModel *api_client.VaultSecret) {
	if terraformDataModel == nil || id == nil {
		tflog.Debug(ctx, "Can't transform Secret Vault to api model")
		return
	}

	result := api_client.VaultSecret{
		BaseSecret: api_client.BaseSecret{ID: fmt.Sprintf("%s/%s", api_client.Vault, *id)},
		Uri:        terraformDataModel.Uri.ValueStringPointer(),
		Prefix:     terraformDataModel.Prefix.ValueStringPointer(),
		Token:      terraformDataModel.Token.ValueStringPointer(),
		Namespace:  terraformDataModel.Namespace.ValueStringPointer(),
	}

	return &result
}

func SecretVaultFromApiToTerraform(ctx context.Context, apiDataModel *api_client.VaultSecret) (terraformDataModel *SecretVaultType) {
	if apiDataModel == nil {
		tflog.Debug(ctx, "Can't transform Secret Vault from api model")
		return
	}

	result := SecretVaultType{
		Uri:       types.StringPointerValue(apiDataModel.Uri),
		Prefix:    types.StringPointerValue(apiDataModel.Prefix),
		Token:     types.StringPointerValue(apiDataModel.Token),
		Namespace: types.StringPointerValue(apiDataModel.Namespace),
	}

	return &result
}
