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

// GlobalRuleResourceModel maps the resource schema data.
type GlobalRuleResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Plugins types.String `tfsdk:"plugins"`
}

var GlobalRuleSchema = schema.Schema{
	Description: "Sets APISIX Plugins which run globally.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the global rule.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Required:    true,
		},
	},
}

func GlobalRuleFromTerraformToApi(ctx context.Context, terraformDataModel *GlobalRuleResourceModel) (apiDataModel api_client.GlobalRule) {
	apiDataModel.ID = terraformDataModel.ID.ValueStringPointer()
	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)

	tflog.Debug(ctx, "Result of the GlobalRuleFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func GlobalRuleFromApiToTerraform(ctx context.Context, apiDataModel *api_client.GlobalRule) (terraformDataModel GlobalRuleResourceModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)

	tflog.Debug(ctx, "Result of the GlobalRuleFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
