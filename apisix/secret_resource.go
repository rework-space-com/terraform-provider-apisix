package apisix

import (
	"context"
	"fmt"
	"strings"

	"github.com/holubovskyi/apisix-client-go"

	"terraform-provider-apisix/apisix/model"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &secretResource{}
	_ resource.ResourceWithConfigure        = &secretResource{}
	_ resource.ResourceWithImportState      = &secretResource{}
	_ resource.ResourceWithConfigValidators = &secretResource{}
)

// NewSecretResource is a helper function to simplify the provider implementation.
func NewSecretResource() resource.Resource {
	return &secretResource{}
}

// secretResource is the resource implementation.
type secretResource struct {
	client *api_client.ApiClient
}

// Metadata returns the resource type name.
func (r *secretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

// Schema defines the schema for the resource.
func (r *secretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = model.SecretSchema
}

// Validate Config
func (r *secretResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("vault"),
			path.MatchRoot("aws"),
			path.MatchRoot("gcp"),
		),
	}
}

// Configure adds the provider configured client to the resource.
func (r *secretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *secretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Start of the secret resource creation")
	// Retrieve values from plan
	var plan model.SecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	secretManager, newSecretRequest := model.SecretFromTerraformToApi(ctx, &plan)

	// Create new secret
	newSecretReponse, err := r.client.CreateSecret(secretManager, plan.ID.ValueString(), newSecretRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secret",
			"Could not create Secret, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	newState := model.SecretFromApiToTerraform(ctx, newSecretReponse)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *secretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Start of the secret resource read")
	// Get current state
	var state model.SecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretManager, secretManagerError := GetSecretManagerFromState(state)
	if secretManagerError != nil {
		resp.Diagnostics.AddError("Provider Selection Error", secretManagerError.Error())
		return
	}

	// Get refreshed secret from the APISIX
	secretStateResponse, err := r.client.GetSecret(secretManager, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Secret",
			"Could not read APISIX Secret by ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite with refreshed state
	newState := model.SecretFromApiToTerraform(ctx, secretStateResponse)

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update the resource.
func (r *secretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Start of the secret resource update")
	// Retrieve values from plan
	var plan model.SecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	secretManager, updateSecretRequest := model.SecretFromTerraformToApi(ctx, &plan)

	// Update existing rule
	_, err := r.client.UpdateSecret(secretManager, plan.ID.ValueString(), updateSecretRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating APISIX Secret",
			"Could not update Secret, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated rule
	updatedSecret, err := r.client.GetSecret(secretManager, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading APISIX Secret",
			"Could not read APISIX Secret by ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := model.SecretFromApiToTerraform(ctx, updatedSecret)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *secretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Start of the secret resource delete")
	// Get current state
	var state model.SecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretManager, secretManagerError := GetSecretManagerFromState(state)
	if secretManagerError != nil {
		resp.Diagnostics.AddError("Provider Selection Error", secretManagerError.Error())
		return
	}

	// Delete the secret
	err := r.client.DeleteSecret(secretManager, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting APISIX Secret",
			"Could not delete Secret, unexpected error: "+err.Error(),
		)
		return
	}
}

// Import resource into state
func (r *secretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Start of the rule importing")
	// Retrieve import ID and save to id attribute

	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in the format <manager>/<id>, got: %s", req.ID),
		)
		return
	}

	manager := api_client.SecretManager(parts[0])
	terraformID := parts[1]

	// 1. We manually set the ID to the "clean" version (toto)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), terraformID)...)
	switch manager {
	case api_client.Vault:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vault"), &model.SecretVaultType{})...)
	case api_client.AWS:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("aws"), &model.SecretAWSType{})...)
	case api_client.GCP:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("gcp"), &model.SecretGCPType{})...)
	}
}

func GetSecretManagerFromState(state model.SecretResourceModel) (secretManager api_client.SecretManager, err error) {
	if state.Vault != nil {
		return api_client.Vault, nil
	} else if state.AWS != nil {
		return api_client.AWS, nil
	} else if state.GCP != nil {
		return api_client.GCP, nil
	} else {
		return "", fmt.Errorf("unknow secret manager: %s", secretManager)
	}
}
