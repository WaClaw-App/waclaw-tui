// Package util provides shared utility functions for TUI screen packages.
// These functions consolidate common patterns that were previously duplicated
// across screen sub-packages, following the DRY principle.
package util

import "strings"

// Slugify converts a category name to a filesystem-safe slug.
// This is the shared implementation previously duplicated as slugify() in the
// niche screen package and slugifyName() in the backend scenario engine.
//
// NOTE: This is a display-only FALLBACK for showing the folder path in the
// TUI when the backend has not yet provided an authoritative folder_slug.
// The backend owns the actual slug generation for file system operations;
// it sends folder_slug via HandleNavigate/HandleUpdate params. The TUI
// should use folder_slug preferentially and only fall back to Slugify()
// if the backend hasn't sent one yet.
func Slugify(name string) string {
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, "&", "and")
	return result
}
