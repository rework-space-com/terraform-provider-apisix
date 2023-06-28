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
	_ resource.Resource                = &consumerResource{}
	_ resource.ResourceWithConfigure   = &consumerResource{}
	_ resource.ResourceWithImportState = &consumerResource{}
)

// NewConsumerResource is a helper function to simplify the provider implementation.
func NewConsumerResource() resource.Resource {
	return &consumerResource{}
}

// consumerResource is the resource implementation.
type consumerResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *consumerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer"
}

// Schema defines the schema for the resource.
func (r *consumerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.ConsumerSchema
}

// Configure adds the provider configured client to the resource.
func (r *consumerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *consumerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the Consumer resource creation")
	// Retrieve values from plan
	var plan model.ConsumerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	newConsumerRequest := model.ConsumerFromTerraformToApi(ctx, &plan)

	// Create new consumer
	newConsumerResponse, err := r.client.CreateConsumer(newConsumerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Consumer",
			"Could not create Consumer, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.ConsumerFromApiToTerraform(ctx, newConsumerResponse)
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
func (r *consumerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the consumer resource read")
	// Get current state
	var state model.ConsumerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed service from the APISIX
	consumerStateResponse, err := r.client.GetConsumer(state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Consumer",
			"Could not read APISIX Consumer by name "+state.Username.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.ConsumerFromApiToTerraform(ctx, consumerStateResponse)
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
func (r *consumerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the consumer resource update")
	// Retrieve values from plan
	var plan model.ConsumerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateConsumerRequest := model.ConsumerFromTerraformToApi(ctx, &plan)

	// Update existing consumer
	_, err := r.client.UpdateConsumer(updateConsumerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Consumer",
			"Could not update Consumer, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated consumer
	updatedConsumer, err := r.client.GetConsumer(plan.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Consumer",
			"Could not read APISIX Consumer by name "+plan.Username.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.ConsumerFromApiToTerraform(ctx, updatedConsumer)
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
func (r *consumerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the consumer resource delete")
	// Get current state
	var state model.ConsumerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the consumer
	err := r.client.DeleteConsumer(state.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Consumer",
			"Could not delete consumer, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *consumerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the consumer importing")
	// Retrieve import user name and save to 'username' attribute
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
