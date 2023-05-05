package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoOriginGroup_basic -v
func TestAccTencentCloudTeoOriginGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "configuration_type", "weight"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_group_name", "keep-group-1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_type", "self"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_records.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_records.0.port", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_records.0.private", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_records.0.record", defaultZoneName),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "origin_records.0.weight", "100"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOriginGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_origin_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup != nil {
			return fmt.Errorf("zone originGroup %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckOriginGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup == nil {
			return fmt.Errorf("zone originGroup %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoOriginGroupVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "zone_name" {
  default = "` + defaultZoneName + `"
}
`
const testAccTeoOriginGroup = testAccTeoOriginGroupVar + `

resource "tencentcloud_teo_origin_group" "basic" {
  configuration_type = "weight"
  origin_group_name  = "keep-group-1"
  origin_type        = "self"
  zone_id            = var.zone_id

  origin_records {
    area      = []
    port      = 8080
    private   = false
    record    = var.zone_name
    weight    = 100
  }
}

`
