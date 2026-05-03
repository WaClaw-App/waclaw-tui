// Package util provides shared utility functions for TUI screen packages.
// These functions consolidate common patterns that were previously duplicated
// across screen sub-packages, following the DRY principle.
package util

// ToInt extracts an int from an any value, handling the float64 type that
// JSON deserialization produces. Returns def if the value cannot be converted.
// This is the shared implementation previously duplicated in onboarding,
// niche, and license screen packages.
func ToInt(v any, def int) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case float32:
		return int(n)
	case int64:
		return int(n)
	case int32:
		return int(n)
	case int:
		return n
	}
	return def
}

// ToStringSlice converts a []any to []string, skipping non-string elements.
// This is the shared implementation previously duplicated in onboarding and
// niche screen packages.
func ToStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// ToString extracts a string from an any value, returning def if not a string.
func ToString(v any, def string) string {
	if s, ok := v.(string); ok {
		return s
	}
	return def
}

// ToBool extracts a bool from an any value, returning def if not a bool.
func ToBool(v any, def bool) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return def
}

// ToFloat64 extracts a float64 from an any value, returning def if not a float64.
func ToFloat64(v any, def float64) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return def
}

// ToMap extracts a map[string]any from an any value, returning nil if not a map.
func ToMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}

// ToSlice extracts a []any from an any value, returning nil if not a slice.
func ToSlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}
