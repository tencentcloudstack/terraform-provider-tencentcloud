package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testDayuEipResourceName = "tencentcloud_dayu_eip"
var testDayuEipResourceKey = testDayuEipResourceName + ".test"

func TestAccTencentCloudDayuEipResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuEipResource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuEipExists(testDayuEipResourceKey),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "bind_resource_id", "ins-4m0jvxic"),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "bind_resource_region", "hk"),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "eip", "162.62.163.50"),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "eip_address_status", "BINDING"),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "resource_id", "bgpip-000004xg"),
					resource.TestCheckResourceAttr(testDayuEipResourceKey, "resource_region", "ap-hongkong"),
				),
			},
		},
	})
}

func testAccCheckDayuEipDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuL4RuleResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of eip")
		}
		resourceId := items[0]
		eip := items[1]

		service := AntiddosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		result, err := service.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)
		if err != nil {
			return err
		}
		if len(result) > 0 {
			err := service.DisassociateDDoSEipAddress(ctx, resourceId, eip)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("delete L4 rule %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuEipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of eip")
		}
		resourceId := items[0]

		service := AntiddosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		result, err := service.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)

		if err != nil {
			return err
		}
		if len(result) > 0 {
			return nil
		} else {
			return fmt.Errorf("eip rule %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuEipResource string = `
resource "tencentcloud_dayu_eip" "test" {
	resource_id = "bgpip-000004xg"
	eip = "162.62.163.50"
	bind_resource_id = "ins-4m0jvxic"
	bind_resource_region = "hk"
	bind_resource_type = "cvm"
  }
`
