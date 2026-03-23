package mongodb

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// mongodbAvailabilityZoneListDiffSuppress suppresses the diff when availability_zone_list
// has the same elements but in different order.
// This is necessary because:
// - The API requires strict ordering when creating resources (must use TypeList)
// - The API returns unordered list when reading resources
// - Without this suppression, users would see unnecessary diffs due to order differences
func mongodbAvailabilityZoneListDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// This function is called for every key related to availability_zone_list:
	// - "availability_zone_list" or "availability_zone_list.#" (list level)
	// - "availability_zone_list.0", "availability_zone_list.1", etc (element level)
	//
	// For TypeList, Terraform calls DiffSuppressFunc at BOTH list and element levels.
	// We need to compare the FULL lists for all these calls and return the same result.

	// Only process keys related to availability_zone_list
	if !strings.Contains(k, "availability_zone_list") {
		return false
	}

	// Get the complete lists from ResourceData
	oldList, newList := d.GetChange("availability_zone_list")

	// Handle nil cases
	if oldList == nil && newList == nil {
		return true
	}
	if oldList == nil || newList == nil {
		return false
	}

	// Convert to string slices
	oldZones := helper.InterfacesStrings(oldList.([]interface{}))
	newZones := helper.InterfacesStrings(newList.([]interface{}))

	// If lengths are different, there's a real change
	if len(oldZones) != len(newZones) {
		return false
	}

	// If both are empty, they're the same
	if len(oldZones) == 0 {
		return true
	}

	// Sort both lists and compare
	// Create copies to avoid modifying the original slices
	oldSorted := make([]string, len(oldZones))
	newSorted := make([]string, len(newZones))
	copy(oldSorted, oldZones)
	copy(newSorted, newZones)

	sort.Strings(oldSorted)
	sort.Strings(newSorted)

	// Compare element by element
	for i := range oldSorted {
		if oldSorted[i] != newSorted[i] {
			return false // Content is different
		}
	}

	// Content is the same, only order differs - suppress the diff
	return true
}
