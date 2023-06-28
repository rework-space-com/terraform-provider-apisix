package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ServiceResourceModel maps the resource schema data.
type ServiceResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"desc"`
	EnableWebsocket types.Bool   `tfsdk:"enable_websocket"`
	Hosts           types.List   `tfsdk:"hosts"`
	Labels          types.Map    `tfsdk:"labels"`
	Plugins         types.String `tfsdk:"plugins"`
	UpstreamId      types.String `tfsdk:"upstream_id"`
}

var ServiceSchema = schema.Schema{
	Description: "Manages APISIX services.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the service.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "Identifier for the service.",
			Optional:    true,
		},
		"desc": schema.StringAttribute{
			Description: "Description of usage scenarios.",
			Optional:    true,
		},
		"enable_websocket": schema.BoolAttribute{
			MarkdownDescription: "Enables a websocket. Set to `false` by default.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"hosts": schema.ListAttribute{
			MarkdownDescription: "Matches with any one of the multiple `hosts` specified in the form of a non-empty list.",
			ElementType:         types.StringType,
			Optional:            true,
		},
		"labels": schema.MapAttribute{
			Description: "Attributes of the Service specified as key-value pairs.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Optional:    true,
		},
		"upstream_id": schema.StringAttribute{
			Description: "Id of the Upstream service.",
			Optional:    true,
		},
	},
}

func ServiceFromTerraformToApi(ctx context.Context, terraformDataModel *ServiceResourceModel) (apiDataModel api_client.Service) {
	apiDataModel.Name = terraformDataModel.Name.ValueStringPointer()
	apiDataModel.Description = terraformDataModel.Description.ValueStringPointer()
	apiDataModel.EnableWebsocket = terraformDataModel.EnableWebsocket.ValueBoolPointer()
	apiDataModel.UpstreamId = terraformDataModel.UpstreamId.ValueStringPointer()

	_ = terraformDataModel.Hosts.ElementsAs(ctx, &apiDataModel.Hosts, true)
	_ = terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, true)

	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)

	tflog.Debug(ctx, "Result of ServiceFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func ServiceFromApiToTerraform(ctx context.Context, apiDataModel *api_client.Service) (terraformDataModel ServiceResourceModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Name = types.StringPointerValue(apiDataModel.Name)
	terraformDataModel.Description = types.StringPointerValue(apiDataModel.Description)
	terraformDataModel.EnableWebsocket = types.BoolPointerValue(apiDataModel.EnableWebsocket)
	terraformDataModel.UpstreamId = types.StringPointerValue(apiDataModel.UpstreamId)

	terraformDataModel.Hosts, _ = types.ListValueFrom(ctx, types.StringType, apiDataModel.Hosts)
	terraformDataModel.Labels, _ = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)

	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)

	tflog.Debug(ctx, "Result of the ServiceFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
