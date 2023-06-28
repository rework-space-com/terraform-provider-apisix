package model

import (
	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TimeoutType struct {
	Connect types.Int64 `tfsdk:"connect"`
	Send    types.Int64 `tfsdk:"send"`
	Read    types.Int64 `tfsdk:"read"`
}

var TimeoutSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Sets the timeout (in seconds) for connecting to, and sending and receiving messages to and from the Upstream.",
	Optional:            true,
	Computed:            false,
	Attributes: map[string]schema.Attribute{
		"connect": schema.Int64Attribute{
			Required: true,
		},
		"send": schema.Int64Attribute{
			Required: true,
		},
		"read": schema.Int64Attribute{
			Required: true,
		},
	},
}

func TimeoutFromTerraformToAPI(terraformDataModel *TimeoutType) (apiDataModel *api_client.TimeoutType) {
	if terraformDataModel == nil {
		return
	}

	result := api_client.TimeoutType{
		Connect: terraformDataModel.Connect.ValueInt64(),
		Send:    terraformDataModel.Send.ValueInt64(),
		Read:    terraformDataModel.Read.ValueInt64(),
	}

	return &result
}

func TimeoutFromAPIToTerraform(apiDataModel *api_client.TimeoutType) (terraformDataModel *TimeoutType) {
	if apiDataModel == nil {
		return
	}

	result := TimeoutType{
		Connect: types.Int64Value(int64(apiDataModel.Connect)),
		Send:    types.Int64Value(int64(apiDataModel.Send)),
		Read:    types.Int64Value(int64(apiDataModel.Read)),
	}

	return &result
}
