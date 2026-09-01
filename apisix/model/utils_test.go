package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPreferPluginsValuePrefersFreshApiValue(t *testing.T) {
	fresh := types.StringValue(`{"limit-count":{"count":2000}}`)
	stale := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPluginsValue(fresh, stale); got.ValueString() != fresh.ValueString() {
		t.Fatalf("expected API value to win when both values are set, got %s", got.ValueString())
	}
}

func TestPreferPluginsValueFallsBackToKnownState(t *testing.T) {
	fresh := types.StringNull()
	stale := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPluginsValue(fresh, stale); got.ValueString() != stale.ValueString() {
		t.Fatalf("expected fallback state value when the API response is null, got %s", got.ValueString())
	}
}

func TestPreferPlanValuePrefersPlannedValue(t *testing.T) {
	plan := types.StringValue(`{"limit-count":{"count":2000}}`)
	api := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPlanValue(plan, api); got.ValueString() != plan.ValueString() {
		t.Fatalf("expected planned value to win when both values are set, got %s", got.ValueString())
	}
}

func TestPreferPlanValueFallsBackToApiValue(t *testing.T) {
	plan := types.StringNull()
	api := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPlanValue(plan, api); got.ValueString() != api.ValueString() {
		t.Fatalf("expected fallback API value when the plan is null, got %s", got.ValueString())
	}
}

func TestPreferPluginsValueReturnsNullWhenBothNull(t *testing.T) {
	if got := PreferPluginsValue(types.StringNull(), types.StringNull()); !got.IsNull() {
		t.Fatalf("expected null when both API and state are null, got %s", got.ValueString())
	}
}

func TestPreferPluginsValueSkipsUnknownApiValue(t *testing.T) {
	api := types.StringUnknown()
	state := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPluginsValue(api, state); got.ValueString() != state.ValueString() {
		t.Fatalf("expected fallback state value when the API response is unknown, got %s", got.ValueString())
	}
}

func TestPreferPlanValueReturnsNullWhenBothNull(t *testing.T) {
	if got := PreferPlanValue(types.StringNull(), types.StringNull()); !got.IsNull() {
		t.Fatalf("expected null when both plan and API are null, got %s", got.ValueString())
	}
}

func TestPreferPlanValueSkipsUnknownPlanValue(t *testing.T) {
	plan := types.StringUnknown()
	api := types.StringValue(`{"limit-count":{"count":5}}`)

	if got := PreferPlanValue(plan, api); got.ValueString() != api.ValueString() {
		t.Fatalf("expected fallback API value when the plan is unknown, got %s", got.ValueString())
	}
}
