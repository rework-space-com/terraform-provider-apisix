package apisix

import (
	"context"
	"fmt"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/holubovskyi/apisix-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consumerGroupResource{}
	_ resource.ResourceWithConfigure   = &consumerGroupResource{}
	_ resource.ResourceWithImportState = &consumerGroupResource{}
)

// NewConsumerGroupResource is a helper function to simplify the provider implementation.
func NewConsumerGroupResource() resource.Resource {
	return &consumerGroupResource{}
}

// consumerGroupResource is the resource implementation.
type consumerGroupResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *consumerGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer_group"
}

// Schema defines the schema for the resource.
func (r *consumerGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.ConsumerGroupSchema
}

// Configure adds the provider configured client to the resource.
func (r *consumerGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *consumerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the consumer group resource creation")
	// Retrieve values from plan
	var plan model.ConsumerGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newConsumerGroupRequest := model.ConsumerGroupFromTerraformToApi(ctx, &plan)

	// Create new consumer group
	newConsumerGroupResponse, err := r.client.CreateConsumerGroup(plan.ID.ValueString(), newConsumerGroupRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Consumer Group",
			"Could not create Consumer Group, unexpected error: "+err.Error(),
		)
		return
	}

	// APISIX can return plugins that differ from what was sent — re-fetch
	// when that happens, the same way Update() already does after
	// UpdateConsumerGroup(). See route_resource.go Create() for the full
	// explanation.
	createdConsumerGroup := newConsumerGroupResponse
	if model.PluginsChanged(newConsumerGroupRequest.Plugins, newConsumerGroupResponse.Plugins) {
		createdConsumerGroup, err = r.client.GetConsumerGroup(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading APISIX Consumer Group",
				"Could not read APISIX Consumer Group by ID "+plan.ID.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.ConsumerGroupFromApiToTerraform(ctx, createdConsumerGroup)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, plan.Plugins)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *consumerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the consumer group resource read")
	// Get current state
	var state model.ConsumerGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed consumer group from the APISIX
	consumerGroupStateResponse, err := r.client.GetConsumerGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Consumer Group",
			"Could not read APISIX Consumer Group by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.ConsumerGroupFromApiToTerraform(ctx, consumerGroupStateResponse)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, state.Plugins)

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the resource.
func (r *consumerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the consumer group resource update")
	// Retrieve values from plan
	var plan model.ConsumerGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateConsumerGroupRequest := model.ConsumerGroupFromTerraformToApi(ctx, &plan)

	// Update existing consumer group
	_, err := r.client.UpdateConsumerGroup(plan.ID.ValueString(), updateConsumerGroupRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Consumer Group",
			"Could not update Consumer Group, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated consumer group
	updatedConsumerGroup, err := r.client.GetConsumerGroup(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Consumer Group",
			"Could not read APISIX Consumer Group by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.ConsumerGroupFromApiToTerraform(ctx, updatedConsumerGroup)
	newState.Plugins = model.PreferPluginsValue(newState.Plugins, plan.Plugins)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *consumerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the consumer group resource delete")
	// Get current state
	var state model.ConsumerGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the consumer group
	err := r.client.DeleteConsumerGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Consumer Group",
			"Could not delete Consumer Group, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *consumerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the consumer group importing")
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
