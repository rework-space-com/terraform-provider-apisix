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
	_ resource.Resource                = &globalRuleResource{}
	_ resource.ResourceWithConfigure   = &globalRuleResource{}
	_ resource.ResourceWithImportState = &globalRuleResource{}
)

// NewGlobalRuleResource is a helper function to simplify the provider implementation.
func NewGlobalRuleResource() resource.Resource {
	return &globalRuleResource{}
}

// globalRuleResource is the resource implementation.
type globalRuleResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *globalRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_rule"
}

// Schema defines the schema for the resource.
func (r *globalRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.GlobalRuleSchema
}

// Configure adds the provider configured client to the resource.
func (r *globalRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *globalRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the global rule resource creation")
	// Retrieve values from plan
	var plan model.GlobalRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newGlobalRuleRequest := model.GlobalRuleFromTerraformToApi(ctx, &plan)

	// Create new global rule
	newGlobalRuleReponse, err := r.client.CreateGlobalRule(plan.ID.ValueString(), newGlobalRuleRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Global Rule",
			"Could not create Global Rule, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.GlobalRuleFromApiToTerraform(ctx, newGlobalRuleReponse)
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
func (r *globalRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the global rule resource read")
	// Get current state
	var state model.GlobalRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed global rule from the APISIX
	globalRuleStateResponse, err := r.client.GetGlobalRule(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Global Rule",
			"Could not read APISIX Global Rule by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.GlobalRuleFromApiToTerraform(ctx, globalRuleStateResponse)
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
func (r *globalRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the global rule resource update")
	// Retrieve values from plan
	var plan model.GlobalRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateGlobalRuleRequest := model.GlobalRuleFromTerraformToApi(ctx, &plan)

	// Update existing rule
	_, err := r.client.UpdateGlobalRule(plan.ID.ValueString(), updateGlobalRuleRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Global Rule",
			"Could not update Global Rule, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated rule
	updatedGlobalRule, err := r.client.GetGlobalRule(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Global Rule",
			"Could not read APISIX Global Rule by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.GlobalRuleFromApiToTerraform(ctx, updatedGlobalRule)
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
func (r *globalRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the global rule resource delete")
	// Get current state
	var state model.GlobalRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the global rule
	err := r.client.DeleteGlobalRule(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Global Rule",
			"Could not delete Global Rule, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *globalRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the rule importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
