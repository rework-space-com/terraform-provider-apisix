package model

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type SecretAWSType struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	SessionToken    types.String `tfsdk:"session_token"`
	Region          types.String `tfsdk:"region"`
	EndpointUrl     types.String `tfsdk:"endpoint_url"`
}

var SecretAWSSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Set APISIX Secret Management AWS configuration.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"access_key_id": schema.StringAttribute{
			Description: "AWS Access Key ID.",
			Required:    true,
		},
		"secret_access_key": schema.StringAttribute{
			Description: "AWS Secret Access Key.",
			Required:    true,
		},
		"session_token": schema.StringAttribute{
			Description: "Temporary access credential information.",
			Optional:    true,
		},
		"region": schema.StringAttribute{
			Description: "AWS Region.",
			Optional:    true,
		},
		"endpoint_url": schema.StringAttribute{
			Description: "AWS Secret Manager URL.",
			Optional:    true,
		},
	},
}

func SecretAWSFromTerraformToApi(ctx context.Context, id *string, terraformDataModel *SecretAWSType) (apiDataModel *api_client.AWSSecret) {
	if terraformDataModel == nil || id == nil {
		tflog.Debug(ctx, "Can't transform Secret AWS to api model")
		return
	}

	result := api_client.AWSSecret{
		BaseSecret:      api_client.BaseSecret{ID: fmt.Sprintf("%s/%s", api_client.AWS, *id)},
		AccessKeyId:     terraformDataModel.AccessKeyId.ValueStringPointer(),
		SecretAccessKey: terraformDataModel.SecretAccessKey.ValueStringPointer(),
		SessionToken:    terraformDataModel.SessionToken.ValueStringPointer(),
		Region:          terraformDataModel.Region.ValueStringPointer(),
		EndpointUrl:     terraformDataModel.EndpointUrl.ValueStringPointer(),
	}

	return &result
}

func SecretAWSFromApiToTerraform(ctx context.Context, apiDataModel *api_client.AWSSecret) (terraformDataModel *SecretAWSType) {
	if apiDataModel == nil {
		tflog.Debug(ctx, "Can't transform Secret AWS from api model")
		return
	}

	result := SecretAWSType{
		AccessKeyId:     types.StringPointerValue(apiDataModel.AccessKeyId),
		SecretAccessKey: types.StringPointerValue(apiDataModel.SecretAccessKey),
		SessionToken:    types.StringPointerValue(apiDataModel.SessionToken),
		Region:          types.StringPointerValue(apiDataModel.Region),
		EndpointUrl:     types.StringPointerValue(apiDataModel.EndpointUrl),
	}

	return &result
}
