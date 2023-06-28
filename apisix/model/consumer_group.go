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

// ConsumerGroupResourceModel maps the resource schema data.
type ConsumerGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"desc"`
	Labels      types.Map    `tfsdk:"labels"`
	Plugins     types.String `tfsdk:"plugins"`
}

var ConsumerGroupSchema = schema.Schema{
	Description: "Manages APISIX Group of Plugins which can be reused across Consumers.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the consumer group.",
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
			Description: "Attributes of the Consumer group specified as key-value pairs.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Required:    true,
		},
	},
}

func ConsumerGroupFromTerraformToApi(ctx context.Context, terraformDataModel *ConsumerGroupResourceModel) (apiDataModel api_client.ConsumerGroup) {
	apiDataModel.ID = terraformDataModel.ID.ValueStringPointer()
	apiDataModel.Description = terraformDataModel.Description.ValueStringPointer()
	terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, true)
	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)

	tflog.Debug(ctx, "Result of the ConsumerGroupFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func ConsumerGroupFromApiToTerraform(ctx context.Context, apiDataModel *api_client.ConsumerGroup) (terraformDataModel ConsumerGroupResourceModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Description = types.StringPointerValue(apiDataModel.Description)
	terraformDataModel.Labels, _ = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)
	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)

	tflog.Debug(ctx, "Result of the ConsumerGroupFromApiToTerraform", map[string]any{
		"Values": apiDataModel,
	})

	return terraformDataModel
}
