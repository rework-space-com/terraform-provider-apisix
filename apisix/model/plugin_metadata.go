package model

import (
	"context"

	api_client "github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PluginMetadataResourceModel maps the resource schema data.
type PluginMetadataResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Metadata types.String `tfsdk:"metadata"`
}

var PluginMetadataSchema = schema.Schema{
	Description: "Manages APISIX Plugin Metadata resource.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The name of the plugin.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"metadata": schema.StringAttribute{
			Description: "Metadata associated with the plugin.",
			Required:    true,
		},
	},
}

func PluginMetadataFromTerraformToApi(ctx context.Context, terraformDataModel *PluginMetadataResourceModel) (apiDataModel api_client.PluginMetadata) {
	apiDataModel.Id = terraformDataModel.Id.ValueStringPointer()
	apiDataModel.Metadata = PluginsStringToJson(ctx, terraformDataModel.Metadata)

	tflog.Debug(ctx, "Result of the PluginMetadataFromTerraformToApi", map[string]any{
		"id":       apiDataModel.Id,
		"metadata": apiDataModel.Metadata,
	})

	return apiDataModel
}

func PluginMetadataFromApiToTerraform(ctx context.Context, apiDataModel *api_client.PluginMetadata) (terraformDataModel PluginMetadataResourceModel) {
	terraformDataModel.Id = types.StringPointerValue(apiDataModel.Id)
	terraformDataModel.Metadata = PluginsFromJsonToString(ctx, apiDataModel.Metadata)

	tflog.Debug(ctx, "Result of the PluginMetadataFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
