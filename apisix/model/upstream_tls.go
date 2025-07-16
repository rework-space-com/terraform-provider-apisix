package model

import (
	"context"
	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UpstreamTLSType struct {
	ClientCertID types.String `tfsdk:"client_cert_id"`
	ClientCert   types.String `tfsdk:"client_cert"`
	ClientKey    types.String `tfsdk:"client_key"`
}

var UpstreamTLSSchemaAttribute = schema.SingleNestedAttribute{
	MarkdownDescription: "Configures the TLS client certificate for the upstream.",
	Optional:            true,

	Attributes: map[string]schema.Attribute{
		"client_cert_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the client certificate to use for TLS. Can't be used with `tls.client_cert` and `tls.client_key`.",
			Optional:            true,
		},
		"client_cert": schema.StringAttribute{
			MarkdownDescription: "Sets the client certificate while connecting to a TLS Upstream. Can't be used with `tls.client_cert_id`.",
			Optional:            true,
		},
		"client_key": schema.StringAttribute{
			MarkdownDescription: "Sets the client key while connecting to a TLS Upstream. Can't be used with `tls.client_cert_id`.",
			Optional:            true,
		},
	},
}

func UpstreamTLSFromTerraformToAPI(ctx context.Context, terraformDataModel *UpstreamTLSType) (apiDataModel *api_client.UpstreamTLSType) {
	if terraformDataModel == nil {
		tflog.Debug(ctx, "Can't transform upstream TLS to api model")
		return
	}

	result := api_client.UpstreamTLSType{
		ClientCertID: terraformDataModel.ClientCertID.ValueString(),
		ClientCert:   terraformDataModel.ClientCert.ValueString(),
		ClientKey:    terraformDataModel.ClientKey.ValueString(),
	}

	return &result
}

func UpstreamTLSFromAPIToTerraform(ctx context.Context, apiDataModel *api_client.UpstreamTLSType) (terraformDataModel *UpstreamTLSType) {
	if apiDataModel == nil {
		tflog.Debug(ctx, "Can't transform upstream TLS from api model")
		return
	}

	result := UpstreamTLSType{
		ClientCertID: types.StringValue(apiDataModel.ClientCertID),
		ClientCert:   types.StringValue(apiDataModel.ClientCert),
		ClientKey:    types.StringValue(apiDataModel.ClientKey),
	}

	return &result
}
