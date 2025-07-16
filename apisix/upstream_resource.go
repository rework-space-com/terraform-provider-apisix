package apisix

import (
	"context"
	"fmt"

	"github.com/holubovskyi/apisix-client-go"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	//	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &upstreamResource{}
	_ resource.ResourceWithConfigure        = &upstreamResource{}
	_ resource.ResourceWithImportState      = &upstreamResource{}
	_ resource.ResourceWithConfigValidators = &upstreamResource{}
)

// NewUpstreamResource is a helper function to simplify the provider implementation.
func NewUpstreamResource() resource.Resource {
	return &upstreamResource{}
}

// upstreamResource is the resource implementation.
type upstreamResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *upstreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_upstream"
}

// Schema defines the schema for the resource.
func (r *upstreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.UpstreamSchema
}

// Validate Config
func (r *upstreamResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("service_name"),
			path.MatchRoot("nodes"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("service_name"),
			path.MatchRoot("discovery_type"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("tls").AtName("client_cert_id"),
			path.MatchRoot("tls").AtName("client_cert"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("tls").AtName("client_cert_id"),
			path.MatchRoot("tls").AtName("client_key"),
		),
	}
}

// Configure adds the provider configured client to the resource.
func (r *upstreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *upstreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the upstream resource creation")
	// Retrieve values from plan
	var plan model.UpstreamResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newUpstreamRequest, labelsDiag := model.UpstreamFromTerraformToAPI(ctx, &plan)

	resp.Diagnostics.Append(labelsDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new upstream
	newUpstreamResponse, err := r.client.CreateUpstream(newUpstreamRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Upstream",
			"Could not create Upstream, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState, labelsDiag := model.UpstreamFromApiToTerraform(ctx, newUpstreamResponse)

	resp.Diagnostics.Append(labelsDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *upstreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the upstream resource read")
	// Get current state
	var state model.UpstreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed upstream from the APISIX
	upsreamResponse, err := r.client.GetUpstream(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Upstream",
			"Could not read APISIX Upstream by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState, labelsDiag := model.UpstreamFromApiToTerraform(ctx, upsreamResponse)

	resp.Diagnostics.Append(labelsDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the upstream resource.
func (r *upstreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the upstream update")
	// Retrieve values from plan
	var plan model.UpstreamResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateUpstreamRequest, labelsDiag := model.UpstreamFromTerraformToAPI(ctx, &plan)

	resp.Diagnostics.Append(labelsDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing upstream
	_, err := r.client.UpdateUpstream(plan.ID.ValueString(), updateUpstreamRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Upstream",
			"Could not update upstream, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated upstream from APISIX
	updatedUpstream, err := r.client.GetUpstream(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Upstream",
			"Could not read APISIX Upstream by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState, labelsDiag := model.UpstreamFromApiToTerraform(ctx, updatedUpstream)
	resp.Diagnostics.Append(labelsDiag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *upstreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the upstream delete")
	// Get current state
	var state model.UpstreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing certificate
	err := r.client.DeleteUpstream(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Upstream",
			"Could not delete upstream by ID "+state.ID.ValueString()+" unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *upstreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the upstream importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
