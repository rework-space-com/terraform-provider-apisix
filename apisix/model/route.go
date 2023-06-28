package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RouteResourceModel maps the resource schema data.
type RouteResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"desc"`
	URI             types.String `tfsdk:"uri"`
	URIS            types.List   `tfsdk:"uris"`
	Host            types.String `tfsdk:"host"`
	Hosts           types.List   `tfsdk:"hosts"`
	RemoteAddr      types.String `tfsdk:"remote_addr"`
	RemoteAddrs     types.List   `tfsdk:"remote_addrs"`
	Methods         types.List   `tfsdk:"methods"`
	Priority        types.Int64  `tfsdk:"priority"`
	Vars            types.String `tfsdk:"vars"`
	FilterFunc      types.String `tfsdk:"filter_func"`
	Plugins         types.String `tfsdk:"plugins"`
	Script          types.String `tfsdk:"script"`
	UpstreamId      types.String `tfsdk:"upstream_id"`
	ServiceId       types.String `tfsdk:"service_id"`
	PluginConfigId  types.String `tfsdk:"plugin_config_id"`
	Labels          types.Map    `tfsdk:"labels"`
	Timeout         *TimeoutType `tfsdk:"timeout"`
	EnableWebsocket types.Bool   `tfsdk:"enable_websocket"`
	Status          types.Int64  `tfsdk:"status"`
}

var RouteSchema = schema.Schema{
	Description: "Manages APISIX routes.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the route.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "Identifier for the route.",
			Optional:    true,
		},
		"desc": schema.StringAttribute{
			Description: "Description of usage scenarios.",
			Optional:    true,
		},
		"uri": schema.StringAttribute{
			Description: "Matches the uri.",
			Optional:    true,
		},
		"uris": schema.ListAttribute{
			MarkdownDescription: "Matches with any one of the multiple `uri`s specified in the form of a non-empty list.",
			ElementType:         types.StringType,
			Optional:            true,
		},
		"host": schema.StringAttribute{
			Description: "Matches with domain names such as `foo.com` or PAN domain names like `*.foo.com`.",
			Optional:    true,
		},
		"hosts": schema.ListAttribute{
			MarkdownDescription: "Matches with any one of the multiple `host`s specified in the form of a non-empty list.",
			ElementType:         types.StringType,
			Optional:            true,
		},
		"remote_addr": schema.StringAttribute{
			Description: "Matches with the specified IP address in standard IPv4 format (`192.168.1.101`), CIDR format (`192.168.1.0/24`), or in IPv6 format.",
			Optional:    true,
		},
		"remote_addrs": schema.ListAttribute{
			MarkdownDescription: "Matches with any one of the multiple `remote_addrs` specified in the form of a non-empty list.",
			ElementType:         types.StringType,
			Optional:            true,
		},
		"methods": schema.ListAttribute{
			MarkdownDescription: "Matches with the specified HTTP methods. Matches all methods if empty or unspecified.",
			ElementType:         types.StringType,
			Optional:            true,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf(HttpMethods...),
				),
			},
		},
		"priority": schema.Int64Attribute{
			MarkdownDescription: "If different Routes matches to the same `uri`, then the Route is matched based on its `priority`." +
				"A higher value corresponds to higher priority." +
				"It is set to `0` by default.",
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(0),
		},
		"vars": schema.StringAttribute{
			MarkdownDescription: "Matches based on the specified variables consistent with variables in Nginx. Takes the form `[[var, operator, val], [var, operator, val], ...]]`.",
			Optional:            true,
		},
		"filter_func": schema.StringAttribute{
			MarkdownDescription: "Matches based on a user-defined filtering function." +
				"Used in scenarios requiring complex matching. These functions can accept an input parameter `vars` which can be used to access the Nginx variables.",
			Optional: true,
		},
		"plugins": schema.StringAttribute{
			Description: "Plugins that are executed during the request/response cycle.",
			Optional:    true,
		},
		"plugin_config_id": schema.StringAttribute{
			Description: "Plugin config bound to the Route.",
			Optional:    true,
		},
		"script": schema.StringAttribute{
			Description: "Used for writing arbitrary Lua code or directly calling existing plugins to be executed.",
			Optional:    true,
		},
		"upstream_id": schema.StringAttribute{
			Description: "Id of the Upstream service.",
			Optional:    true,
		},
		"service_id": schema.StringAttribute{
			Description: "Configuration of the bound Service.",
			Optional:    true,
		},
		"labels": schema.MapAttribute{
			Description: "Attributes of the Service specified as key-value pairs.",
			ElementType: types.StringType,
			Optional:    true,
		},
		"timeout": TimeoutSchemaAttribute,
		"enable_websocket": schema.BoolAttribute{
			MarkdownDescription: "Enables a websocket. Set to `false` by default.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"status": schema.Int64Attribute{
			MarkdownDescription: "Enables the current Route. Set to `1` (enabled) by default. `1` to enable, `0` to disable",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
			Validators: []validator.Int64{
				int64validator.OneOf([]int64{0, 1}...),
			},
		},
	},
}

func RouteFromTerraformToApi(ctx context.Context, terraformDataModel *RouteResourceModel) (apiDataModel api_client.Route) {
	apiDataModel.Name = terraformDataModel.Name.ValueStringPointer()
	apiDataModel.Description = terraformDataModel.Description.ValueStringPointer()
	apiDataModel.URI = terraformDataModel.URI.ValueStringPointer()

	terraformDataModel.URIS.ElementsAs(ctx, &apiDataModel.URIS, false)

	apiDataModel.Host = terraformDataModel.Host.ValueStringPointer()

	terraformDataModel.Hosts.ElementsAs(ctx, &apiDataModel.Hosts, false)

	apiDataModel.RemoteAddr = terraformDataModel.RemoteAddr.ValueStringPointer()

	terraformDataModel.RemoteAddrs.ElementsAs(ctx, &apiDataModel.RemoteAddrs, false)
	terraformDataModel.Methods.ElementsAs(ctx, &apiDataModel.Methods, false)

	apiDataModel.Priority = terraformDataModel.Priority.ValueInt64Pointer()

	apiDataModel.Vars = VarsStringToJson(ctx, terraformDataModel.Vars)

	apiDataModel.FilterFunc = terraformDataModel.FilterFunc.ValueStringPointer()
	apiDataModel.Plugins = PluginsStringToJson(ctx, terraformDataModel.Plugins)
	apiDataModel.Script = terraformDataModel.Script.ValueStringPointer()
	apiDataModel.UpstreamId = terraformDataModel.UpstreamId.ValueStringPointer()
	apiDataModel.ServiceId = terraformDataModel.ServiceId.ValueStringPointer()
	apiDataModel.PluginConfigId = terraformDataModel.PluginConfigId.ValueStringPointer()

	terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, false)

	apiDataModel.Timeout = TimeoutFromTerraformToAPI(terraformDataModel.Timeout)

	apiDataModel.EnableWebsocket = terraformDataModel.EnableWebsocket.ValueBoolPointer()
	apiDataModel.Status = terraformDataModel.Status.ValueInt64Pointer()

	tflog.Debug(ctx, "Result of the RouteFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func RouteFromApiToTerraform(ctx context.Context, apiDataModel *api_client.Route) (terraformDataModel RouteResourceModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Name = types.StringPointerValue(apiDataModel.Name)
	terraformDataModel.Description = types.StringPointerValue(apiDataModel.Description)
	terraformDataModel.URI = types.StringPointerValue(apiDataModel.URI)

	terraformDataModel.URIS, _ = types.ListValueFrom(ctx, types.StringType, apiDataModel.URIS)

	terraformDataModel.Host = types.StringPointerValue(apiDataModel.Host)

	terraformDataModel.Hosts, _ = types.ListValueFrom(ctx, types.StringType, apiDataModel.Hosts)

	terraformDataModel.RemoteAddr = types.StringPointerValue(apiDataModel.RemoteAddr)

	terraformDataModel.RemoteAddrs, _ = types.ListValueFrom(ctx, types.StringType, apiDataModel.RemoteAddrs)

	terraformDataModel.Methods, _ = types.ListValueFrom(ctx, types.StringType, apiDataModel.Methods)
	terraformDataModel.Priority = types.Int64PointerValue(apiDataModel.Priority)

	terraformDataModel.Vars = VarsFromJsonToString(ctx, apiDataModel.Vars)

	terraformDataModel.FilterFunc = types.StringPointerValue(apiDataModel.FilterFunc)
	terraformDataModel.Plugins = PluginsFromJsonToString(ctx, apiDataModel.Plugins)
	terraformDataModel.Script = types.StringPointerValue(apiDataModel.Script)
	terraformDataModel.UpstreamId = types.StringPointerValue(apiDataModel.UpstreamId)
	terraformDataModel.ServiceId = types.StringPointerValue(apiDataModel.ServiceId)
	terraformDataModel.PluginConfigId = types.StringPointerValue(apiDataModel.PluginConfigId)

	terraformDataModel.Labels, _ = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)

	terraformDataModel.Timeout = TimeoutFromAPIToTerraform(apiDataModel.Timeout)

	terraformDataModel.EnableWebsocket = types.BoolPointerValue(apiDataModel.EnableWebsocket)
	terraformDataModel.Status = types.Int64PointerValue(apiDataModel.Status)

	tflog.Debug(ctx, "Result of the RouteFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
