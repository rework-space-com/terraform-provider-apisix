package model

import (
	"context"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type UpstreamChecksType struct {
	Active  *UpstreamChecksActiveType  `tfsdk:"active"`
	Passive *UpstreamChecksPassiveType `tfsdk:"passive"`
}

var UpstreamChecksSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Configures the parameters for the health check.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"active":  UpstreamChecksActiveSchemaAttribute,
		"passive": UpstreamChecksPassiveSchemaAttribute,
	},
}

func UpstreamChecksFromTerraformToAPI(ctx context.Context, terraformDataModel *UpstreamChecksType) (apiDataModel *api_client.UpstreamChecksType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamChecksType{
		Active:  UpstreamChecksActiveFromTerraformToApi(ctx, terraformDataModel.Active),
		Passive: UpstreamChecksPassiveFromTerraformToApi(ctx, terraformDataModel.Passive),
	}

	return &result
}

func UpstreamChecksFromApiToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamChecksType) (terraformDataModel *UpstreamChecksType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamChecksType{
		Active:  UpstreamChecksActiveFromApiToTerraform(ctx, apiDataModel.Active),
		Passive: UpstreamChecksPassiveFromApiToTerraform(ctx, apiDataModel.Passive),
	}

	return &result
}
