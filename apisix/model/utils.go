package model

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func PluginsStringToJson(ctx context.Context, pluginsString types.String) *map[string]interface{} {
	if pluginsString.IsNull() || pluginsString.IsUnknown() {
		return nil
	}

	pluginsStr := pluginsString.ValueString()
	if pluginsStr == "" {
		return nil
	}

	var pluginsMap map[string]interface{}
	if err := json.Unmarshal([]byte(pluginsStr), &pluginsMap); err != nil {
		tflog.Error(ctx, "Failed to parse plugins JSON", map[string]interface{}{
			"error":       err.Error(),
			"json_string": pluginsStr,
		})
		return nil
	}

	tflog.Debug(ctx, "Parsed metadata JSON", map[string]interface{}{
		"input_string": pluginsStr,
		"parsed_map":   pluginsMap,
	})

	return &pluginsMap
}

func PluginsFromJsonToString(ctx context.Context, metadataMap *map[string]interface{}) types.String {
	if metadataMap == nil || len(*metadataMap) == 0 {
		return types.StringNull()
	}

	metadataBytes, err := json.Marshal(*metadataMap)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal metadata to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return types.StringNull()
	}

	jsonString := string(metadataBytes)
	tflog.Debug(ctx, "Converted metadata to JSON string", map[string]interface{}{
		"input_map":     *metadataMap,
		"output_string": jsonString,
	})

	return types.StringValue(jsonString)
}

func VarsStringToJson(ctx context.Context, str types.String) (jsonPointer *[]interface{}) {

	if str.IsNull() {
		return nil
	}

	var result []interface{}
	bytes := []byte(str.ValueString())

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]interface{ any }{
			"Error converting vars to json": err,
			"Input string is":               str.ValueString(),
		})
		panic(err)
	}

	return &result
}

func VarsFromJsonToString(ctx context.Context, jsonPointer *[]interface{}) (str types.String) {
	if jsonPointer == nil {
		return types.StringNull()
	}

	data, err := json.Marshal(jsonPointer)
	if err != nil {
		tflog.Error(ctx, "Error converting vars to terraform values")
		panic(err)
	}

	jsonStr := string(data)

	return types.StringValue(jsonStr)

}
