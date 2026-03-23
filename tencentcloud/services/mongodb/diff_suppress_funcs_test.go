package mongodb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMongodbAvailabilityZoneListDiffSuppress(t *testing.T) {
	testCases := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		{
			name:     "same order",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"},
			new:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"},
			expected: true,
		},
		{
			name:     "different order same content",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"},
			new:      []interface{}{"ap-guangzhou-6", "ap-guangzhou-3", "ap-guangzhou-4"},
			expected: true,
		},
		{
			name:     "different order same content - reverse",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"},
			new:      []interface{}{"ap-guangzhou-6", "ap-guangzhou-4", "ap-guangzhou-3"},
			expected: true,
		},
		{
			name:     "different content",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4"},
			new:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-5"},
			expected: false,
		},
		{
			name:     "different length - old longer",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4"},
			new:      []interface{}{"ap-guangzhou-3"},
			expected: false,
		},
		{
			name:     "different length - new longer",
			old:      []interface{}{"ap-guangzhou-3"},
			new:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4"},
			expected: false,
		},
		{
			name:     "both empty",
			old:      []interface{}{},
			new:      []interface{}{},
			expected: true,
		},
		{
			name:     "single element",
			old:      []interface{}{"ap-guangzhou-3"},
			new:      []interface{}{"ap-guangzhou-3"},
			expected: true,
		},
		{
			name:     "two elements swapped",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4"},
			new:      []interface{}{"ap-guangzhou-4", "ap-guangzhou-3"},
			expected: true,
		},
		{
			name:     "completely different zones",
			old:      []interface{}{"ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"},
			new:      []interface{}{"ap-shanghai-1", "ap-shanghai-2", "ap-shanghai-3"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test ResourceData with the old value
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"availability_zone_list": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			}, map[string]interface{}{
				"availability_zone_list": tc.old,
			})

			// Set the new value
			_ = d.Set("availability_zone_list", tc.new)

			// Test the diff suppress function
			result := mongodbAvailabilityZoneListDiffSuppress(
				"availability_zone_list",
				"",
				"",
				d,
			)

			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
				t.Errorf("old: %v, new: %v", tc.old, tc.new)
			}
		})
	}
}

func TestMongodbAvailabilityZoneListDiffSuppress_SubElements(t *testing.T) {
	// Test that the function returns false for sub-element keys
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"availability_zone_list": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}, map[string]interface{}{
		"availability_zone_list": []interface{}{"ap-guangzhou-3", "ap-guangzhou-4"},
	})

	_ = d.Set("availability_zone_list", []interface{}{"ap-guangzhou-4", "ap-guangzhou-3"})

	testCases := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "list root key",
			key:      "availability_zone_list",
			expected: true,
		},
		{
			name:     "list count key",
			key:      "availability_zone_list.#",
			expected: true,
		},
		{
			name:     "first element",
			key:      "availability_zone_list.0",
			expected: false,
		},
		{
			name:     "second element",
			key:      "availability_zone_list.1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mongodbAvailabilityZoneListDiffSuppress(tc.key, "", "", d)
			if result != tc.expected {
				t.Errorf("key %s: expected %v, got %v", tc.key, tc.expected, result)
			}
		})
	}
}
