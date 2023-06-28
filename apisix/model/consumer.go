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

// ConsumerResourceModel maps the resource schema data.
type ConsumerResourceModel struct {
	Username    types.String `tfsdk:"username"`
	Description types.String `tfsdk:"desc"`
	Labels      types.Map    `tfsdk:"labels"`
	Plugins     types.String `tfsdk:"plugins"`
	GroupId     types.String `tfsdk:"group_id"`
}

var ConsumerSchema = schema.Schema{
	Description: "Manages APISIX Consumers.",
	Attributes: map[string]schema.Attribute{
		"username": schema.StringAttribute{
			Description: "Name and Identifier of the Consumer.",
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
			Description: "Attributes of the Consumer specified as key-value pairs.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Optional:    true,
		},
		"group_id": schema.StringAttribute{
			Description: "Group of the Consumer.",
			Optional:    true,
		},
	},
}

func ConsumerFromTerraformToApi(ctx context.Context, terraformDataModel *ConsumerResourceModel) (apiDataModel api_client.Consumer) {
	apiDataModel.Username = terraformDataModel.Username.ValueStringPointer()
	apiDataModel.Description = terraformDataModel.Description.ValueStringPointer()
	apiDataModel.GroupId = terraformDataModel.GroupId.ValueStringPointer()

	_ = terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, true)

	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)

	tflog.Debug(ctx, "Result of ConsumerFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func ConsumerFromApiToTerraform(ctx context.Context, apiDataModel *api_client.Consumer) (terraformDataModel ConsumerResourceModel) {
	terraformDataModel.Username = types.StringPointerValue(apiDataModel.Username)
	terraformDataModel.Description = types.StringPointerValue(apiDataModel.Description)
	terraformDataModel.GroupId = types.StringPointerValue(apiDataModel.GroupId)

	terraformDataModel.Labels, _ = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)

	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)

	tflog.Debug(ctx, "Result of ConsumerFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
