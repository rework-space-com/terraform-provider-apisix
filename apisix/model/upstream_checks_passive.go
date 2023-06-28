package model

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	api_client "github.com/holubovskyi/apisix-client-go"
)

type UpstreamChecksPassiveType struct {
	Healthy   *UpstreamChecksPassiveHealthyType   `tfsdk:"healthy"`
	Unhealthy *UpstreamChecksPassiveUnhealthyType `tfsdk:"unhealthy"`
}

var UpstreamChecksPassiveSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Passive health check refers to judging whether the corresponding upstream node is healthy by judging the response status of the request forwarded from APISIX to the upstream node.",
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"healthy":   UpstreamChecksPassiveHealthySchemaAttribute,
		"unhealthy": UpstreamChecksPassiveUnhealthySchemaAttribute,
	},
}

func UpstreamChecksPassiveFromTerraformToApi(ctx context.Context, terraformDataModel *UpstreamChecksPassiveType) (apiDataModel *api_client.UpstreamChecksPassiveType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamChecksPassiveType{
		Healthy:   UpstreamChecksPassiveHealthyFromTerraformToApi(ctx, terraformDataModel.Healthy),
		Unhealthy: UpstreamChecksPassiveUnhealthyFromTerraformToApi(ctx, terraformDataModel.Unhealthy),
	}

	return &result
}

func UpstreamChecksPassiveFromApiToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamChecksPassiveType) (terraformDataModel *UpstreamChecksPassiveType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamChecksPassiveType{
		Healthy:   UpstreamChecksPassiveHealthyFromApiToTerraform(ctx, apiDataModel.Healthy),
		Unhealthy: UpstreamChecksPassiveUnhealthyFromApiToTerraform(ctx, apiDataModel.Unhealthy),
	}

	return &result
}
