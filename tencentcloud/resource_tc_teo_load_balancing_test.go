package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoLoadBalancing_basic -v
func TestAccTencentCloudTeoLoadBalancing_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancingDestroy,
		Steps:        []resource.TestStep{
			//{
			//	Config: testAccTeoLoadBalancing,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckLoadBalancingExists("tencentcloud_teo_load_balancing.basic"),
			//		resource.TestCheckResourceAttr("tencentcloud_teo_load_balancing.basic", "zone_id", defaultZoneId),
			//		resource.TestCheckResourceAttr("tencentcloud_teo_load_balancing.basic", "host", "aaa."+defaultZoneName),
			//		resource.TestCheckResourceAttr("tencentcloud_teo_load_balancing.basic", "origin_group_id", "origin-8a6e424e-47b4-11ed-8422-5254006e4802"),
			//		resource.TestCheckResourceAttr("tencentcloud_teo_load_balancing.basic", "type", "proxied"),
			//	),
			//},
			//{
			//	ResourceName:      "tencentcloud_teo_load_balancing.basic",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

func testAccCheckLoadBalancingDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_load_balancing" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		loadBalancingId := idSplit[1]

		agents, err := service.DescribeTeoLoadBalancing(ctx, zoneId, loadBalancingId)
		if agents != nil {
			return fmt.Errorf("zone loadBalancing %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//const testAccTeoLoadBalancingVar = `
//variable "zone_id" {
//  default = "` + defaultZoneId + `"
//}
//
//variable "zone_name" {
//  default = "aaa.` + defaultZoneName + `"
//}
//`

//const testAccTeoLoadBalancing = testAccTeoLoadBalancingVar + `
//
//resource "tencentcloud_teo_load_balancing" "basic" {
//  host                   = var.zone_name
//  origin_group_id        = "origin-8a6e424e-47b4-11ed-8422-5254006e4802"
//  type                   = "proxied"
//  zone_id                = var.zone_id
//  backup_origin_group_id = ""
//}
//
//`
