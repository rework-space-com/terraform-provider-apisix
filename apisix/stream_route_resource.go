package apisix

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &streamRouteResource{}
	_ resource.ResourceWithConfigure   = &streamRouteResource{}
	_ resource.ResourceWithImportState = &streamRouteResource{}
)

// NewStreamRouteResource is a helper function to simplify the provider implementation.
func NewStreamRouteResource() resource.Resource {
	return &streamRouteResource{}
}

// streamRouteResource is the resource implementation.
type streamRouteResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *streamRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stream_route"
}

// Schema defines the schema for the resource.
func (r *streamRouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.StreamRouteSchema
}

// Configure adds the provider configured client to the resource.
func (r *streamRouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *streamRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the stream route resource creation")
	// Retrieve values from plan
	var plan model.StreamRouteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newStreamRouteRequest := model.StreamRouteFromTerraformToApi(ctx, &plan)

	// Create new stream route
	newStreamRouteReponse, err := r.client.CreateStreamRoute(newStreamRouteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Stream Route",
			"Could not create Stream Route, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.StreamRouteFromApiToTerraform(ctx, newStreamRouteReponse)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *streamRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the stream route resource read")
	// Get current state
	var state model.StreamRouteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed stream route from the APISIX
	streamRouteStateResponse, err := r.client.GetStreamRoute(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Stream Route",
			"Could not read APISIX Stream Route by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.StreamRouteFromApiToTerraform(ctx, streamRouteStateResponse)

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the resource.
func (r *streamRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the stream route resource update")
	// Retrieve values from plan
	var plan model.StreamRouteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateStreamRouteRequest := model.StreamRouteFromTerraformToApi(ctx, &plan)

	// Update existing stream route
	_, err := r.client.UpdateStreamRoute(plan.ID.ValueString(), updateStreamRouteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Stream Route",
			"Could not update Stream Route, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated stream route
	updatedStreamRoute, err := r.client.GetStreamRoute(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Stream Route",
			"Could not read APISIX Stream Route by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.StreamRouteFromApiToTerraform(ctx, updatedStreamRoute)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *streamRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the stream route resource delete")
	// Get current state
	var state model.StreamRouteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the Stream Route
	err := r.client.DeleteStreamRoute(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Stream Route",
			"Could not delete stream route, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *streamRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the stream route importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
