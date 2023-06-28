package model

import (
	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpstreamKeepAlivePoolType struct {
	Size        types.Int64 `tfsdk:"size"`
	IdleTimeout types.Int64 `tfsdk:"idle_timeout"`
	Requests    types.Int64 `tfsdk:"requests"`
}

var UpstreamKeepAlivePoolSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Sets the `keepalive_pool`.",
	Optional:            true,
	Computed:            false,
	Attributes: map[string]schema.Attribute{
		"size": schema.Int64Attribute{
			Required: true,
		},
		"idle_timeout": schema.Int64Attribute{
			Required: true,
		},
		"requests": schema.Int64Attribute{
			Required: true,
		},
	},
}

func UpstreamKeepAlivePoolFromTerraformToAPI(terraformDataModel *UpstreamKeepAlivePoolType) (apiDataModel *api_client.UpstreamKeepAlivePoolType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.UpstreamKeepAlivePoolType{
		Size:        terraformDataModel.Size.ValueInt64(),
		IdleTimeout: terraformDataModel.IdleTimeout.ValueInt64(),
		Requests:    terraformDataModel.Requests.ValueInt64(),
	}

	return &result
}

func UpstreamKeepAlivePoolFromAPIToTerraform(apiDataModel *api_client.UpstreamKeepAlivePoolType) (terraformDataModel *UpstreamKeepAlivePoolType) {
	if apiDataModel == nil {
		return
	}

	result := UpstreamKeepAlivePoolType{
		Size:        types.Int64Value(int64(apiDataModel.Size)),
		IdleTimeout: types.Int64Value(int64(apiDataModel.IdleTimeout)),
		Requests:    types.Int64Value(int64(apiDataModel.Requests)),
	}

	return &result
}
