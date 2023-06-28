package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/holubovskyi/apisix-client-go"
)

type UpstreamChecksPassiveUnhealthyType struct {
	HTTPStatuses types.List  `tfsdk:"http_statuses"`
	TCPFailures  types.Int64 `tfsdk:"tcp_failures"`
	Timeouts     types.Int64 `tfsdk:"timeouts"`
	HTTPFailures types.Int64 `tfsdk:"http_failures"`
}

var UpstreamChecksPassiveUnhealthySchemaAttribute = schema.SingleNestedAttribute{
	Optional: true,
	Attributes: map[string]schema.Attribute{
		"http_statuses": schema.ListAttribute{
			MarkdownDescription: "Passive check (unhealthy node) HTTP or HTTPS type check, the HTTP status code of the non-healthy node.",
			ElementType:         types.Int64Type,
			Optional:            true,
			Computed:            true,
			Validators: []validator.List{
				listvalidator.ValueInt64sAre(int64validator.Between(200, 599)),
			},
			Default: listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{
				types.Int64Value(429),
				types.Int64Value(500),
				types.Int64Value(503),
			})),
		},
		"http_failures": schema.Int64Attribute{
			MarkdownDescription: "Passive check (unhealthy node) The number of times that the node is not healthy during HTTP or HTTPS type checking.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(0, 254),
			},
			Default: int64default.StaticInt64(5),
		},
		"tcp_failures": schema.Int64Attribute{
			MarkdownDescription: "Passive check (unhealthy node) When TCP type is checked, determine the number of times that the node is not healthy.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(0, 254),
			},
			Default: int64default.StaticInt64(2),
		},
		"timeouts": schema.Int64Attribute{
			MarkdownDescription: "Passive checks (unhealthy node) determine the number of timeouts for unhealthy nodes.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(0, 254),
			},
			Default: int64default.StaticInt64(7),
		},
	},
}

func UpstreamChecksPassiveUnhealthyFromTerraformToApi(ctx context.Context, terraformDataModel *UpstreamChecksPassiveUnhealthyType) (apiDataModel *api_client.UpstreamChecksPassiveUnhealthyType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamChecksPassiveUnhealthyType{
		TCPFailures:  terraformDataModel.TCPFailures.ValueInt64(),
		Timeouts:     terraformDataModel.Timeouts.ValueInt64(),
		HTTPFailures: terraformDataModel.HTTPFailures.ValueInt64(),
	}

	_ = terraformDataModel.HTTPStatuses.ElementsAs(ctx, &result.HTTPStatuses, false)

	return &result
}

func UpstreamChecksPassiveUnhealthyFromApiToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamChecksPassiveUnhealthyType) (terraformDataModel *UpstreamChecksPassiveUnhealthyType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamChecksPassiveUnhealthyType{
		TCPFailures:  types.Int64Value(int64(apiDataModel.TCPFailures)),
		Timeouts:     types.Int64Value(int64(apiDataModel.Timeouts)),
		HTTPFailures: types.Int64Value(int64(apiDataModel.HTTPFailures)),
	}
	result.HTTPStatuses, _ = types.ListValueFrom(ctx, types.Int64Type, apiDataModel.HTTPStatuses)

	return &result
}
