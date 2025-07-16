package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UpstreamResourceModel struct {
	ID            types.String               `tfsdk:"id"`
	Type          types.String               `tfsdk:"type"`
	ServiceName   types.String               `tfsdk:"service_name"`
	DiscoveryType types.String               `tfsdk:"discovery_type"`
	Timeout       *TimeoutType               `tfsdk:"timeout"`
	Name          types.String               `tfsdk:"name"`
	Desc          types.String               `tfsdk:"desc"`
	PassHost      types.String               `tfsdk:"pass_host"`
	Scheme        types.String               `tfsdk:"scheme"`
	Retries       types.Int64                `tfsdk:"retries"`
	RetryTimeout  types.Int64                `tfsdk:"retry_timeout"`
	Labels        types.Map                  `tfsdk:"labels"`
	UpstreamHost  types.String               `tfsdk:"upstream_host"`
	HashOn        types.String               `tfsdk:"hash_on"`
	Key           types.String               `tfsdk:"key"`
	KeepalivePool *UpstreamKeepAlivePoolType `tfsdk:"keepalive_pool"`
	TLS           *UpstreamTLSType           `tfsdk:"tls"`
	Checks        *UpstreamChecksType        `tfsdk:"checks"`
	Nodes         *[]UpstreamNodeType        `tfsdk:"nodes"`
}

var UpstreamSchema = schema.Schema{
	Description: "Manages APISIX Upstreams.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the upstream.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "Load balancing algorithm to be used, and the default value is `roundrobin`.\n" +
				"Can be one of the following: `roundrobin`, `chash`, `ewma` or `least_conn`",
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("roundrobin"),
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"roundrobin", "chash", "ewma", "least_conn"}...),
			},
		},
		"service_name": schema.StringAttribute{
			MarkdownDescription: "Service name used for service discovery. Can't be used with `nodes`",
			Optional:            true,
		},
		"discovery_type": schema.StringAttribute{
			MarkdownDescription: "The type of service discovery. Required, if `service_name` is used",
			Optional:            true,
		},
		"timeout": TimeoutSchemaAttribute,
		"name": schema.StringAttribute{
			MarkdownDescription: "Identifier for the Upstream.",
			Optional:            true,
		},
		"desc": schema.StringAttribute{
			MarkdownDescription: "Description of usage scenarios.",
			Optional:            true,
		},
		"pass_host": schema.StringAttribute{
			MarkdownDescription: "Configures the `host` when the request is forwarded to the upstream. " +
				"Can be one of `pass`, `node` or `rewrite`. Defaults to `pass` if not specified.",
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("pass"),
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"pass", "node", "rewrite"}...),
			},
		},
		"scheme": schema.StringAttribute{
			MarkdownDescription: "The scheme used when communicating with the Upstream. " +
				"For an L7 proxy, this value can be one of `http`, `https`, `grpc`, `grpcs`. " +
				"For an L4 proxy, this value could be one of `tcp`, `udp`, `tls`. Defaults to `http`.",
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("http"),
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"http", "https", "grpc", "grpcs", "tcp", "udp", "tls"}...),
			},
		},
		"retries": schema.Int64Attribute{
			MarkdownDescription: "Sets the number of retries while passing the request to Upstream using the underlying Nginx mechanism. " +
				"Setting this to `0` disables retry.",
			Optional: true,
		},
		"retry_timeout": schema.Int64Attribute{
			MarkdownDescription: "Timeout to continue with retries. Setting this to `0` disables the retry timeout.",
			Optional:            true,
		},
		"upstream_host": schema.StringAttribute{
			MarkdownDescription: "Specifies the host of the Upstream request. This is only valid if the `pass_host` is set to `rewrite`.",
			Optional:            true,
		},
		"hash_on": schema.StringAttribute{
			MarkdownDescription: "Only valid if the type is chash. Supports Nginx variables (vars), custom headers (header), cookie and consumer. " +
				"Defaults to vars.",
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("vars"),
		},
		"key": schema.StringAttribute{
			MarkdownDescription: "Nginx var",
			Optional:            true,
		},
		"labels": schema.MapAttribute{
			MarkdownDescription: "Attributes of the Upstream specified as `key-value` pairs.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"keepalive_pool": UpstreamKeepAlivePoolSchemaAttribute,
		"tls":            UpstreamTLSSchemaAttribute,
		"checks":         UpstreamChecksSchemaAttribute,
		"nodes":          UpstreamNodesSchemaAttribute,
	},
}

func UpstreamFromTerraformToAPI(ctx context.Context, terraformDataModel *UpstreamResourceModel) (apiDataModel api_client.Upstream, labelsDiag diag.Diagnostics) {
	apiDataModel.Type = terraformDataModel.Type.ValueStringPointer()
	apiDataModel.ServiceName = terraformDataModel.ServiceName.ValueStringPointer()
	apiDataModel.DiscoveryType = terraformDataModel.DiscoveryType.ValueStringPointer()
	apiDataModel.Name = terraformDataModel.Name.ValueStringPointer()
	apiDataModel.Desc = terraformDataModel.Desc.ValueStringPointer()
	apiDataModel.PassHost = terraformDataModel.PassHost.ValueStringPointer()
	apiDataModel.Scheme = terraformDataModel.Scheme.ValueStringPointer()
	apiDataModel.Retries = terraformDataModel.Retries.ValueInt64Pointer()
	apiDataModel.RetryTimeout = terraformDataModel.RetryTimeout.ValueInt64Pointer()
	apiDataModel.UpstreamHost = terraformDataModel.UpstreamHost.ValueStringPointer()
	apiDataModel.HashOn = terraformDataModel.HashOn.ValueStringPointer()
	apiDataModel.Key = terraformDataModel.Key.ValueStringPointer()

	labelsDiag = terraformDataModel.Labels.ElementsAs(ctx, &apiDataModel.Labels, false)

	apiDataModel.Timeout = TimeoutFromTerraformToAPI(terraformDataModel.Timeout)
	apiDataModel.KeepalivePool = UpstreamKeepAlivePoolFromTerraformToAPI(terraformDataModel.KeepalivePool)
	apiDataModel.Checks = UpstreamChecksFromTerraformToAPI(ctx, terraformDataModel.Checks)
	apiDataModel.Nodes = UpstreamNodesFromTerraformToAPI(ctx, terraformDataModel.Nodes)
	apiDataModel.TLS = UpstreamTLSFromTerraformToAPI(ctx, terraformDataModel.TLS)

	tflog.Debug(ctx, "Result of UpstreamFromTerraformToAPI", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel, labelsDiag
}

func UpstreamFromApiToTerraform(ctx context.Context, apiDataModel *api_client.Upstream) (terraformDataModel UpstreamResourceModel, labelsDiag diag.Diagnostics) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.Type = types.StringPointerValue(apiDataModel.Type)
	terraformDataModel.ServiceName = types.StringPointerValue(apiDataModel.ServiceName)
	terraformDataModel.DiscoveryType = types.StringPointerValue(apiDataModel.DiscoveryType)
	terraformDataModel.Name = types.StringPointerValue(apiDataModel.Name)
	terraformDataModel.Desc = types.StringPointerValue(apiDataModel.Desc)
	terraformDataModel.PassHost = types.StringPointerValue(apiDataModel.PassHost)
	terraformDataModel.Scheme = types.StringPointerValue(apiDataModel.Scheme)
	terraformDataModel.Retries = types.Int64PointerValue(apiDataModel.Retries)
	terraformDataModel.RetryTimeout = types.Int64PointerValue(apiDataModel.RetryTimeout)
	terraformDataModel.UpstreamHost = types.StringPointerValue(apiDataModel.UpstreamHost)
	terraformDataModel.HashOn = types.StringPointerValue(apiDataModel.HashOn)
	terraformDataModel.Key = types.StringPointerValue(apiDataModel.Key)

	terraformDataModel.Labels, labelsDiag = types.MapValueFrom(ctx, types.StringType, apiDataModel.Labels)

	terraformDataModel.Timeout = TimeoutFromAPIToTerraform(apiDataModel.Timeout)
	terraformDataModel.KeepalivePool = UpstreamKeepAlivePoolFromAPIToTerraform(apiDataModel.KeepalivePool)
	terraformDataModel.Checks = UpstreamChecksFromApiToTerraform(ctx, apiDataModel.Checks)
	terraformDataModel.Nodes = UpstreamNodesFromApiToTerraform(ctx, apiDataModel.Nodes)
	terraformDataModel.TLS = UpstreamTLSFromAPIToTerraform(ctx, apiDataModel.TLS)

	tflog.Debug(ctx, "Result of UpstreamFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel, labelsDiag
}
