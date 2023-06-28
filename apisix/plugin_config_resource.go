package apisix

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pluginConfigResource{}
	_ resource.ResourceWithConfigure   = &pluginConfigResource{}
	_ resource.ResourceWithImportState = &pluginConfigResource{}
)

// NewPluginConfigResource is a helper function to simplify the provider implementation.
func NewPluginConfigResource() resource.Resource {
	return &pluginConfigResource{}
}

// pluginConfigResource is the resource implementation.
type pluginConfigResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *pluginConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin_config"
}

// Schema defines the schema for the resource.
func (r *pluginConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.PluginConfigSchema
}

// Configure adds the provider configured client to the resource.
func (r *pluginConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_client.ApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api_client.ApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create a new resource.
func (r *pluginConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the plugin config resource creation")
	// Retrieve values from plan
	var plan model.PluginConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newPluginConfigRequest := model.PluginConfigFromTerraformToApi(ctx, &plan)

	// Create new plugin config
	newPluginConfigResponse, err := r.client.CreatePluginConfig(plan.ID.ValueString(), newPluginConfigRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Plugin Config",
			"Could not create Plugin Config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.PluginConfigFromApiToTerraform(ctx, newPluginConfigResponse)
	if !newState.Plugins.IsNull() {
		newState.Plugins = types.StringValue(plan.Plugins.ValueString())
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *pluginConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the plugin config resource read")
	// Get current state
	var state model.PluginConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed plugin config from the APISIX
	pluginConfigStateResponse, err := r.client.GetPluginConfig(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Plugin Config",
			"Could not read APISIX Plugin Config by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.PluginConfigFromApiToTerraform(ctx, pluginConfigStateResponse)
	if !newState.Plugins.IsNull() {
		newState.Plugins = types.StringValue(state.Plugins.ValueString())
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the resource.
func (r *pluginConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the plugin config resource update")
	// Retrieve values from plan
	var plan model.PluginConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updatePluginConfigRequest := model.PluginConfigFromTerraformToApi(ctx, &plan)

	// Update existing plugin config
	_, err := r.client.UpdatePluginConfig(plan.ID.ValueString(), updatePluginConfigRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Plugin Config",
			"Could not update Plugin Config, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated rule
	updatedPluginConfig, err := r.client.GetPluginConfig(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Plugin Config",
			"Could not read APISIX Plugin Config by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.PluginConfigFromApiToTerraform(ctx, updatedPluginConfig)
	if !newState.Plugins.IsNull() {
		newState.Plugins = types.StringValue(plan.Plugins.ValueString())
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *pluginConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the plugin config resource delete")
	// Get current state
	var state model.PluginConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the plugin config
	err := r.client.DeletePluginConfig(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Plugin Config",
			"Could not delete Plugin Config, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *pluginConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the plugin config importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
