package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PluginConfigResourceModel maps the resource schema data.
type PluginConfigResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"desc"`
	Labels      types.Map    `tfsdk:"labels"`
	Plugins     types.String `tfsdk:"plugins"`
}

var PluginConfigSchema = schema.Schema{
	Description: "Manages APISIX Group of Plugins which can be reused across Routes.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the plugin config.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"desc": schema.StringAttribute{
			Description: "Description of usage scenarios.",
			Optional:    true,
		},
		"labels": schema.MapAttribute{
			Description: "Attributes of the Plugin config specified as key-value pairs.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Required:    true,
		},
	},
}

func PluginConfigFromTerraformToApi(ctx context.Context, terraformDataModel *PluginConfigResourceModel) (apiDataModel api_client.PluginConfig) {
	apiDataModel.ID = terraformDataModel.ID.ValueStringPointer()
	apiDataModel.Description = terraformDataModel.Description.ValueStringPointer()
	terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, true)
	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)

	tflog.Debug(ctx, "Result of the PluginConfigFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func PluginConfigFromApiToTerraform(ctx context.Context, apiDataModel *api_client.PluginConfig) (terraformDataModel PluginConfigResourceModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Description = types.StringPointerValue(apiDataModel.Description)
	terraformDataModel.Labels, _ = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)
	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)

	tflog.Debug(ctx, "Result of the PluginConfigFromApiToTerraform", map[string]any{
		"Values": apiDataModel,
	})

	return terraformDataModel
}
