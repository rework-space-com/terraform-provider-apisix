package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpstreamChecksPassiveHealthyType struct {
	HTTPStatuses types.List  `tfsdk:"http_statuses"`
	Successes    types.Int64 `tfsdk:"successes"`
}

var UpstreamChecksPassiveHealthySchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Passive health check refers to judging whether the corresponding upstream node is healthy by judging the response status of the request forwarded from APISIX to the upstream node.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"http_statuses": schema.ListAttribute{
			MarkdownDescription: "Passive check (healthy node) HTTP or HTTPS type check, the HTTP status code of the healthy node.",
			ElementType:         types.Int64Type,
			Optional:            true,
			Computed:            true,
			Validators: []validator.List{
				listvalidator.ValueInt64sAre(int64validator.Between(200, 599)),
			},
			Default: listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{
				types.Int64Value(200),
				types.Int64Value(201),
				types.Int64Value(202),
				types.Int64Value(203),
				types.Int64Value(204),
				types.Int64Value(205),
				types.Int64Value(206),
				types.Int64Value(207),
				types.Int64Value(208),
				types.Int64Value(226),
				types.Int64Value(300),
				types.Int64Value(301),
				types.Int64Value(302),
				types.Int64Value(303),
				types.Int64Value(304),
				types.Int64Value(305),
				types.Int64Value(306),
				types.Int64Value(307),
				types.Int64Value(308),
			})),
		},
		"successes": schema.Int64Attribute{
			MarkdownDescription: "Passive checks (healthy node) determine the number of times a node is healthy.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.Between(1, 254),
			},
			Default: int64default.StaticInt64(5),
		},
	},
}

func UpstreamChecksPassiveHealthyFromTerraformToApi(ctx context.Context, terraformDataModel *UpstreamChecksPassiveHealthyType) (apiDataModel *api_client.UpstreamChecksPassiveHealthyType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamChecksPassiveHealthyType{
		Successes: terraformDataModel.Successes.ValueInt64(),
	}
	_ = terraformDataModel.HTTPStatuses.ElementsAs(ctx, &result.HTTPStatuses, false)

	return &result
}

func UpstreamChecksPassiveHealthyFromApiToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamChecksPassiveHealthyType) (terraformDataModel *UpstreamChecksPassiveHealthyType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamChecksPassiveHealthyType{
		Successes: types.Int64Value(int64(apiDataModel.Successes)),
	}
	result.HTTPStatuses, _ = types.ListValueFrom(ctx, types.Int64Type, apiDataModel.HTTPStatuses)

	return &result
}
