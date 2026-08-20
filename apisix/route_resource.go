package apisix

import (
	"context"
	"fmt"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/holubovskyi/apisix-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &routeResource{}
	_ resource.ResourceWithConfigure        = &routeResource{}
	_ resource.ResourceWithImportState      = &routeResource{}
	_ resource.ResourceWithConfigValidators = &routeResource{}
)

// NewRouteResource is a helper function to simplify the provider implementation.
func NewRouteResource() resource.Resource {
	return &routeResource{}
}

// routeResource is the resource implementation.
type routeResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *routeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

// Schema defines the schema for the resource.
func (r *routeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.RouteSchema
}

// Validate Config
func (r *routeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("uri"),
			path.MatchRoot("uris"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("host"),
			path.MatchRoot("hosts"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("remote_addr"),
			path.MatchRoot("remote_addrs"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("plugin_config_id"),
			path.MatchRoot("script"),
		),
	}
}

// Configure adds the provider configured client to the resource.
func (r *routeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *routeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the route resource creation")
	// Retrieve values from plan
	var plan model.RouteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newRouteRequest := model.RouteFromTerraformToApi(ctx, &plan)

	// Create a new route
	newRouteResponse, err := r.client.CreateRoute(newRouteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Route",
			"Could not create Route, unexpected error: "+err.Error(),
		)
		return
	}

	// APISIX can return plugins that differ from what was sent (e.g.
	// encrypt_fields like openid-connect's client_secret come back as
	// ciphertext, decrypted only on a subsequent GET) — re-fetch when that
	// happens, the same way Update() already does after UpdateRoute(),
	// otherwise the resource is permanently tainted.
	createdRoute := newRouteResponse
	if model.PluginsChanged(newRouteRequest.Plugins, newRouteResponse.Plugins) {
		createdRoute, err = r.client.GetRoute(*newRouteResponse.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading APISIX Route",
				"Could not read APISIX Route by ID "+*newRouteResponse.ID+": "+err.Error(),
			)
			return
		}
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.RouteFromApiToTerraform(ctx, createdRoute)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, plan.Plugins)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *routeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the route resource read")
	// Get current state
	var state model.RouteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed route from the APISIX
	routeStateResponse, err := r.client.GetRoute(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Route",
			"Could not read APISIX Route by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.RouteFromApiToTerraform(ctx, routeStateResponse)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, state.Plugins)

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the resource.
func (r *routeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the route resource update")
	// Retrieve values from plan
	var plan model.RouteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateRouteRequest := model.RouteFromTerraformToApi(ctx, &plan)

	// Update existing route
	_, err := r.client.UpdateRoute(plan.ID.ValueString(), updateRouteRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Route",
			"Could not update Route, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated route
	updatedRoute, err := r.client.GetRoute(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Route",
			"Could not read APISIX Route by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.RouteFromApiToTerraform(ctx, updatedRoute)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, plan.Plugins)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *routeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the route resource delete")
	// Get current state
	var state model.RouteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the route
	err := r.client.DeleteRoute(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Route",
			"Could not delete rotue, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *routeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the route importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
