package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRumWhitelistResource_basic -v
func TestAccTencentCloudRumWhitelistResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRumWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRumWhitelist,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRumWhitelistExists("tencentcloud_rum_whitelist.whitelist"),
					resource.TestCheckResourceAttr("tencentcloud_rum_whitelist.whitelist", "instance_id", defaultRumInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_rum_whitelist.whitelist", "remark", "white list remark"),
					resource.TestCheckResourceAttr("tencentcloud_rum_whitelist.whitelist", "whitelist_uin", "20221122"),
				),
			},
			{
				ResourceName:      "tencentcloud_rum_whitelist.whitelist",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRumWhitelistDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_rum_project" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceID := idSplit[0]
		wid := idSplit[1]

		whitelist, err := service.DescribeRumWhitelist(ctx, instanceID, wid)
		if whitelist != nil {
			return fmt.Errorf("rum whitelist %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRumWhitelistExists(r string) resource.TestCheckFunc {
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
		instanceID := idSplit[0]
		wid := idSplit[1]

		service := RumService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		whitelist, err := service.DescribeRumWhitelist(ctx, instanceID, wid)
		if whitelist == nil {
			return fmt.Errorf("rum whitelist %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccRumWhitelistVar = `
variable "instance_id" {
  default = "` + defaultRumInstanceId + `"
}
`

const testAccRumWhitelist = testAccRumWhitelistVar + `

resource "tencentcloud_rum_whitelist" "whitelist" {
	instance_id = var.instance_id
	remark = "white list remark"
	whitelist_uin = "20221122"
	# aid = ""
  }
`
