package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRedisSecurityGroupAttachmentResource_basic -v
func TestAccTencentCloudRedisSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSecurityGroupAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudRedisSecurityGroupExists("tencentcloud_redis_security_group_attachment.security_group_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_security_group_attachment.security_group_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_security_group_attachment.security_group_attachment", "instance_id", defaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_security_group_attachment.security_group_attachment", "security_group_id", defaultCrsSecurityGroups),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_security_group_attachment.security_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTencentCloudRedisSecurityGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := items[0]
		securityGroupId := items[1]

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		securityGroup, err := service.DescribeRedisSecurityGroupAttachmentById(ctx, PRODUCT, instanceId, securityGroupId)
		if err != nil {
			return err
		}
		if securityGroup == nil {
			return fmt.Errorf("redis securityGroup %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccTencentCloudRedisSecurityGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_security_group_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := items[0]
		securityGroupId := items[1]

		securityGroup, err := service.DescribeRedisSecurityGroupAttachmentById(ctx, PRODUCT, instanceId, securityGroupId)
		if err != nil {
			return err
		}
		if securityGroup != nil {
			return fmt.Errorf("redis securityGroup %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccRedisSecurityGroupVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}
variable "security_group_id" {
	default = "` + defaultCrsSecurityGroups + `"
}
`

const testAccRedisSecurityGroupAttachment = testAccRedisSecurityGroupVar + `

resource "tencentcloud_redis_security_group_attachment" "security_group_attachment" {
	instance_id       = var.instance_id
	security_group_id = var.security_group_id
}

`
