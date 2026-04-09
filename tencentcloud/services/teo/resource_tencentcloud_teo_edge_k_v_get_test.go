package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoEdgeKVGet_basic -v
func TestAccTencentCloudTeoEdgeKVGet_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoEdgeKVGetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoEdgeKVGet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoEdgeKVGetExists("tencentcloud_teo_edge_k_v_get.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_k_v_get.basic", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_edge_k_v_get.basic", "namespace"),
					resource.TestCheckResourceAttr("tencentcloud_teo_edge_k_v_get.basic", "keys.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_edge_k_v_get.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudTeoEdgeKVGet_update -v
func TestAccTencentCloudTeoEdgeKVGet_update(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoEdgeKVGetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoEdgeKVGet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoEdgeKVGetExists("tencentcloud_teo_edge_k_v_get.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_edge_k_v_get.basic", "keys.#", "2"),
				),
			},
			{
				Config: testAccTeoEdgeKVGet_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoEdgeKVGetExists("tencentcloud_teo_edge_k_v_get.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_edge_k_v_get.basic", "keys.#", "3"),
				),
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudTeoEdgeKVGet_nonExistentKey -v
func TestAccTencentCloudTeoEdgeKVGet_nonExistentKey(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoEdgeKVGetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoEdgeKVGet_nonExistentKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoEdgeKVGetExists("tencentcloud_teo_edge_k_v_get.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_edge_k_v_get.basic", "keys.#", "1"),
				),
			},
		},
	})
}

func testAccCheckTeoEdgeKVGetDestroy(s *terraform.State) error {
	// For query resources, destroy only removes from state, no need to check cloud resources
	return nil
}

func testAccCheckTeoEdgeKVGetExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}

		// For query resources, we just verify the ID format is correct
		zoneId := idSplit[0]
		namespace := idSplit[1]
		keysHash := idSplit[2]

		if zoneId == "" {
			return fmt.Errorf("zone_id is empty")
		}
		if namespace == "" {
			return fmt.Errorf("namespace is empty")
		}
		if keysHash == "" {
			return fmt.Errorf("keys_hash is empty")
		}

		return nil
	}
}

const testAccTeoEdgeKVGet = testAccTeoZone + `
resource "tencentcloud_teo_edge_k_v_get" "basic" {
  zone_id   = tencentcloud_teo_zone.basic.id
  namespace = "default"
  keys      = ["test_key_1", "test_key_2"]
}
`

const testAccTeoEdgeKVGet_update = testAccTeoZone + `
resource "tencentcloud_teo_edge_k_v_get" "basic" {
  zone_id   = tencentcloud_teo_zone.basic.id
  namespace = "default"
  keys      = ["test_key_1", "test_key_2", "test_key_3"]
}
`

const testAccTeoEdgeKVGet_nonExistentKey = testAccTeoZone + `
resource "tencentcloud_teo_edge_k_v_get" "basic" {
  zone_id   = tencentcloud_teo_zone.basic.id
  namespace = "default"
  keys      = ["non_existent_key"]
}
`
