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
	api_client "github.com/holubovskyi/apisix-client-go"
)

type UpstreamChecksActiveUnhealthyType struct {
	Interval     types.Int64 `tfsdk:"interval"`
	HTTPStatuses types.List  `tfsdk:"http_statuses"`
	TCPFailures  types.Int64 `tfsdk:"tcp_failures"`
	Timeouts     types.Int64 `tfsdk:"timeouts"`
	HTTPFailures types.Int64 `tfsdk:"http_failures"`
}

var UpstreamChecksActiveUnhealthySchemaAttribute = schema.SingleNestedAttribute{
	Optional: true,
	Attributes: map[string]schema.Attribute{
		"interval": schema.Int64Attribute{
			MarkdownDescription: "Active check (unhealthy node) check interval (unit: second)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"http_statuses": schema.ListAttribute{
			MarkdownDescription: "Active check (unhealthy node) HTTP or HTTPS type check, the HTTP status code of the non-healthy node.",
			ElementType:         types.Int64Type,
			Optional:            true,
			Computed:            true,
			Validators: []validator.List{
				listvalidator.ValueInt64sAre(int64validator.Between(200, 599)),
			},
			Default: listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{
				types.Int64Value(429),
				types.Int64Value(404),
				types.Int64Value(500),
				types.Int64Value(501),
				types.Int64Value(502),
				types.Int64Value(503),
				types.Int64Value(504),
				types.Int64Value(505),
			})),
		},
		"http_failures": schema.Int64Attribute{
			MarkdownDescription: "Active check (unhealthy node) HTTP or HTTPS type check, determine the number of times that the node is not healthy.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 254),
			},
			Default: int64default.StaticInt64(5),
		},
		"tcp_failures": schema.Int64Attribute{
			MarkdownDescription: "Active check (unhealthy node) TCP type check, determine the number of times that the node is not healthy.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 254),
			},
			Default: int64default.StaticInt64(2),
		},
		"timeouts": schema.Int64Attribute{
			MarkdownDescription: "Active check (unhealthy node) to determine the number of timeouts for unhealthy nodes.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 254),
			},
			Default: int64default.StaticInt64(3),
		},
	},
}

func UpstreamChecksActiveUnhealthyFromTerraformToApi(ctx context.Context, terraformDataModel *UpstreamChecksActiveUnhealthyType) (apiDataModel *api_client.UpstreamChecksActiveUnhealthyType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamChecksActiveUnhealthyType{
		Interval:     terraformDataModel.Interval.ValueInt64(),
		TCPFailures:  terraformDataModel.TCPFailures.ValueInt64(),
		Timeouts:     terraformDataModel.Timeouts.ValueInt64(),
		HTTPFailures: terraformDataModel.HTTPFailures.ValueInt64(),
	}

	_ = terraformDataModel.HTTPStatuses.ElementsAs(ctx, &result.HTTPStatuses, false)

	return &result
}

func UpstreamChecksActiveUnhealthyFromApiToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamChecksActiveUnhealthyType) (terraformDataModel *UpstreamChecksActiveUnhealthyType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamChecksActiveUnhealthyType{
		Interval:     types.Int64Value(int64(apiDataModel.Interval)),
		TCPFailures:  types.Int64Value(int64(apiDataModel.TCPFailures)),
		Timeouts:     types.Int64Value(int64(apiDataModel.Timeouts)),
		HTTPFailures: types.Int64Value(int64(apiDataModel.HTTPFailures)),
	}
	result.HTTPStatuses, _ = types.ListValueFrom(ctx, types.Int64Type, apiDataModel.HTTPStatuses)

	return &result
}
