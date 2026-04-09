package teo

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoEdgeKVGet_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTeoEdgeKVGet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "keys.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "data.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoEdgeKVGet_update(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTeoEdgeKVGet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "namespace"),
				),
			},
			{
				Config: testAccTencentCloudTeoEdgeKVGet_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "namespace"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoEdgeKVGet_disappears(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTeoEdgeKVGet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_kv_get.example", "id"),
				),
			},
			{
				// Simulate resource deletion from cloud
				// For query-only resource, deletion from state is immediate
				Destroy: true,
				Config:  testAccTencentCloudTeoEdgeKVGet_basic(),
			},
		},
	})
}

func TestAccTencentCloudTeoEdgeKVGet_invalidKey(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccTencentCloudTeoEdgeKVGet_invalidKey(),
				ExpectError: regexp.MustCompile("can only contain letters, numbers, hyphens, and underscores"),
			},
		},
	})
}

func TestAccTencentCloudTeoEdgeKVGet_tooManyKeys(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccTencentCloudTeoEdgeKVGet_tooManyKeys(),
				ExpectError: regexp.MustCompile("must contain at most 20 elements"),
			},
		},
	})
}

func TestAccTencentCloudTeoEdgeKVGet_emptyKeys(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccTencentCloudTeoEdgeKVGet_emptyKeys(),
				ExpectError: regexp.MustCompile("must contain at least 1 element"),
			},
		},
	})
}

func testAccTencentCloudTeoEdgeKVGet_basic() string {
	return `
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-xxxxx"
  namespace = "ns-xxxxx"
  keys      = ["key1", "key2"]
}
`
}

func testAccTencentCloudTeoEdgeKVGet_update() string {
	return `
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-xxxxx"
  namespace = "ns-xxxxx"
  keys      = ["key1", "key2", "key3"]
}
`
}

func testAccTencentCloudTeoEdgeKVGet_invalidKey() string {
	return `
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-xxxxx"
  namespace = "ns-xxxxx"
  keys      = ["key@invalid"]
}
`
}

func testAccTencentCloudTeoEdgeKVGet_tooManyKeys() string {
	keys := make([]string, 21)
	for i := 0; i < 21; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	keysStr := strings.Join(keys, `", "`)

	return fmt.Sprintf(`
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-xxxxx"
  namespace = "ns-xxxxx"
  keys      = ["%s"]
}
`, keysStr)
}

func testAccTencentCloudTeoEdgeKVGet_emptyKeys() string {
	return `
resource "tencentcloud_teo_edge_kv_get" "example" {
  zone_id   = "zone-xxxxx"
  namespace = "ns-xxxxx"
  keys      = []
}
`
}
