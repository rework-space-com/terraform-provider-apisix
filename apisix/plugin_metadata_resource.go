package apisix

import (
	"context"
	"fmt"

	api_client "github.com/holubovskyi/apisix-client-go"

	"terraform-provider-apisix/apisix/model"

	//"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	//"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pluginMetadataResource{}
	_ resource.ResourceWithConfigure   = &pluginMetadataResource{}
	_ resource.ResourceWithImportState = &pluginMetadataResource{}
)

// NewPluginMetadataResource is a helper function to simplify the provider implementation.
func NewPluginMetadataResource() resource.Resource {
	return &pluginMetadataResource{}
}

// pluginMetadataResource is the resource implementation.
type pluginMetadataResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *pluginMetadataResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin_metadata"
}

// Schema defines the schema for the resource.
func (r *pluginMetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.PluginMetadataSchema
}

// Configure adds the provider configured client to the resource.
func (r *pluginMetadataResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *pluginMetadataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the plugin metadata resource creation")

	var plan model.PluginMetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newPluginMetadataRequest := model.PluginMetadataFromTerraformToApi(ctx, &plan)
	// Debug: Log what we're about to send
	tflog.Debug(ctx, "Create - Sending to API", map[string]interface{}{
		"plugin_id":     plan.Id.ValueString(),
		"request":       newPluginMetadataRequest,
		"plan_metadata": plan.Metadata.ValueString(),
	})

	// Create new plugin metadata
	newPluginMetadataResponse, err := r.client.CreatePluginMetadata(plan.Id.ValueString(), newPluginMetadataRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Plugin Metadata",
			"Could not create Plugin Metadata, unexpected error: "+err.Error(),
		)
		return
	}

	// Debug: Log the create response
	tflog.Debug(ctx, "Create - API response", map[string]interface{}{"response": newPluginMetadataResponse})

	// Map response body to schema
	newState := model.PluginMetadataFromApiToTerraform(ctx, newPluginMetadataResponse)

	// Debug: Log the converted state
	tflog.Debug(ctx, "Create - Converted state", map[string]interface{}{
		"metadata_is_null": newState.Metadata.IsNull(),
		"metadata_value":   newState.Metadata.ValueString(),
	})

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

// Read resource information.
func (r *pluginMetadataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.PluginMetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plugin metadata from API
	pluginMetadataResponse, err := r.client.GetPluginMetadata(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read plugin metadata, got error: %s", err))
		return
	}

	// Convert API response to Terraform state
	newState := model.PluginMetadataFromApiToTerraform(ctx, pluginMetadataResponse)

	// CRITICAL FIX: If API returns null metadata but we have metadata in state, preserve it
	if newState.Metadata.IsNull() && !state.Metadata.IsNull() {
		newState.Metadata = state.Metadata
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

// Update the resource.
func (r *pluginMetadataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the plugin metadata resource update")

	var plan model.PluginMetadataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePluginMetadataRequest := model.PluginMetadataFromTerraformToApi(ctx, &plan)
	// Debug: Log what we're about to send
	tflog.Debug(ctx, "Update - Sending to API", map[string]interface{}{
		"plugin_id":     plan.Id.ValueString(),
		"request":       updatePluginMetadataRequest,
		"plan_metadata": plan.Metadata.ValueString(),
	})

	// Update existing plugin metadata
	updateResponse, err := r.client.UpdatePluginMetadata(plan.Id.ValueString(), updatePluginMetadataRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Plugin Metadata",
			"Could not update plugin metadata, unexpected error: "+err.Error(),
		)
		return
	}

	// Debug: Log the update response
	tflog.Debug(ctx, "Update - API response", map[string]interface{}{"response": updateResponse})

	// Fetch updated metadata
	updatedPluginMetadata, err := r.client.GetPluginMetadata(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Plugin Metadata After Update",
			"Could not read plugin metadata after update: "+err.Error(),
		)
		return
	}

	// Debug: Log what we got back
	tflog.Debug(ctx, "Update - Get response", map[string]interface{}{"response": updatedPluginMetadata})

	// Convert to state
	newState := model.PluginMetadataFromApiToTerraform(ctx, updatedPluginMetadata)

	// Debug: Log the converted state
	tflog.Debug(ctx, "Update - Converted state", map[string]interface{}{
		"metadata_is_null": newState.Metadata.IsNull(),
		"metadata_value":   newState.Metadata.ValueString(),
	})

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

// Delete resource.
func (r *pluginMetadataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the plugin metadata resource delete")
	// Get current state
	var state model.PluginMetadataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the plugin metadata
	err := r.client.DeletePluginMetadata(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Plugin Metadata",
			"Could not delete Plugin Metadata, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *pluginMetadataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID should be the plugin name
	pluginName := req.ID

	// Read the plugin metadata from API
	pluginMetadataResponse, err := r.client.GetPluginMetadata(pluginName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Plugin Metadata",
			"Could not read APISIX Plugin Metadata during import: "+err.Error(),
		)
		return
	}

	// Debug: Log what we got from API
	tflog.Debug(ctx, "Import - API response", map[string]interface{}{"response": pluginMetadataResponse})

	// Convert API response to Terraform state
	state := model.PluginMetadataFromApiToTerraform(ctx, pluginMetadataResponse)

	// Debug: Log the converted state
	tflog.Debug(ctx, "Import - Converted state", map[string]interface{}{"metadata": state.Metadata.ValueString()})

	// Set the imported state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
