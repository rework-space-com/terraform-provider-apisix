package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// StreamRouteModel maps the resource schema data.
type StreamRouteModel struct {
	ID         types.String `tfsdk:"id"`
	UpstreamId types.String `tfsdk:"upstream_id"`
	RemoteAddr types.String `tfsdk:"remote_addr"`
	ServerAddr types.String `tfsdk:"server_addr"`
	ServerPort types.Int64  `tfsdk:"server_port"`
	SNI        types.String `tfsdk:"sni"`
}

var StreamRouteSchema = schema.Schema{
	Description: "Manages APISIX Routes used in the Stream Proxy.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the stream route.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"upstream_id": schema.StringAttribute{
			Description: "Id of the Upstream service.",
			Required:    true,
		},
		"remote_addr": schema.StringAttribute{
			MarkdownDescription: "Filters Upstream forwards by matching with client IP. IPv4 (`127.0.0.1`) OR CIDR format (`127.0.0.1/32`).",
			Optional:            true,
		},
		"server_addr": schema.StringAttribute{
			MarkdownDescription: "Filters Upstream forwards by matching with APISIX Server IP. IPv4 (`127.0.0.1`) OR CIDR format (`127.0.0.1/32`).",
			Optional:            true,
		},
		"server_port": schema.Int64Attribute{
			Description: "Filters Upstream forwards by matching with APISIX Server port.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.Between(1, 65535),
			},
		},
		"sni": schema.StringAttribute{
			MarkdownDescription: "Server Name Indication. Matches with domain names such as `foo.com`",
			Optional:            true,
		},
	},
}

func StreamRouteFromTerraformToApi(ctx context.Context, terraformDataModel *StreamRouteModel) (apiDataModel api_client.StreamRoute) {
	apiDataModel.UpstreamId = terraformDataModel.UpstreamId.ValueStringPointer()
	apiDataModel.RemoteAddr = terraformDataModel.RemoteAddr.ValueStringPointer()
	apiDataModel.ServerAddr = terraformDataModel.ServerAddr.ValueStringPointer()
	apiDataModel.ServerPort = terraformDataModel.ServerPort.ValueInt64Pointer()
	apiDataModel.SNI = terraformDataModel.SNI.ValueStringPointer()

	tflog.Debug(ctx, "Result of the StreamRouteFromTerraformToApi", map[string]any{
		"Values": apiDataModel,
	})

	return apiDataModel
}

func StreamRouteFromApiToTerraform(ctx context.Context, apiDataModel *api_client.StreamRoute) (terraformDataModel StreamRouteModel) {
	terraformDataModel.ID = types.StringPointerValue(apiDataModel.ID)
	terraformDataModel.UpstreamId = types.StringPointerValue(apiDataModel.UpstreamId)
	terraformDataModel.RemoteAddr = types.StringPointerValue(apiDataModel.RemoteAddr)
	terraformDataModel.ServerAddr = types.StringPointerValue(apiDataModel.ServerAddr)
	terraformDataModel.ServerPort = types.Int64PointerValue(apiDataModel.ServerPort)
	terraformDataModel.SNI = types.StringPointerValue(apiDataModel.SNI)

	tflog.Debug(ctx, "Result of the StreamRouteFromApiToTerraform", map[string]any{
		"Values": terraformDataModel,
	})

	return terraformDataModel
}
