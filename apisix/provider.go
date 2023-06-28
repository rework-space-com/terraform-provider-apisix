package apisix

import (
	"context"
	"os"

	"github.com/holubovskyi/apisix-client-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &apisixProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &apisixProvider{
			version: version,
		}
	}
}

// apisixProvider is the provider implementation.
type apisixProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// apisixProviderModel maps provider schema data to a Go type.
type apisixProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

// Metadata returns the provider type name.
func (p *apisixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "apisix"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *apisixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with APISIX API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "Endpoint for APISIX API. May also be provided via APISIX_ENDPOINT environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "API Key for APISIX API. May also be provided via APISIX_APIKEY environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *apisixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring APISIX client")

	// Retrieve provider data from configuration
	var config apisixProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown APISIX API Endpoint",
			"The provider cannot create the APISIX API client as there is an unknown configuration value for the APISIX API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the APISIX_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown APISIX API Key",
			"The provider cannot create the APISIX API client as there is an unknown configuration value for the HashiCups API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the APISIX_APIKEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("APISIX_ENDPOINT")
	apiKey := os.Getenv("APISIX_APIKEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing APISIX API Endpoint",
			"The provider cannot create the APISIX API client as there is a missing or empty value for the APISIX API Endpoint. "+
				"Set the endpoint value in the configuration or use the APISIX_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing APISIX API Key",
			"The provider cannot create the APISIX API client as there is a missing or empty value for the HashiCups API Key. "+
				"Set the api_key value in the configuration or use the HASHICUPS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "apisix_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "apisix_apikey", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apisix_apikey")

	tflog.Debug(ctx, "Creating APISIX client")

	// Create a new APISIX client using the configuration values
	client, err := api_client.NewClient(&endpoint, &apiKey)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create APISIX API Client",
			"An unexpected error occurred when creating the APISIX API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"APISIX Client Error: "+err.Error(),
		)
		return
	}

	// Make the APISIX client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured APISIX client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *apisixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *apisixProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSLCertificateResource,
		NewUpstreamResource,
		NewServiceResource,
		NewConsumerResource,
		NewRouteResource,
		NewGlobalRuleResource,
		NewStreamRouteResource,
		NewConsumerGroupResource,
		NewPluginConfigResource,
	}
}
