package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRedisSecurityGroupResource_basic -v
func TestAccTencentCloudRedisSecurityGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccTencentCloudRedisSecurityGroupExists("tencentcloud_redis_security_group.security_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_security_group.security_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_security_group.security_group", "instance_id", defaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_security_group.security_group", "security_group_id", defaultCrsSecurityGroup),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_security_group.security_group",
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
		// securityGroup := items[1]

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		securityGroupId, err := service.DescribeRedisSecurityGroupById(ctx, instanceId)
		if err != nil {
			return err
		}

		if securityGroupId == "" {
			return fmt.Errorf("security group %s not found", rs.Primary.ID)
		}
		return nil
	}
}

func testAccTencentCloudRedisSecurityGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_security_group" {
			continue
		}
		time.Sleep(5 * time.Second)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		// securityGroup := items[1]

		securityGroupId, err := service.DescribeRedisSecurityGroupById(ctx, instanceId)
		if err != nil {
			return err
		}
		if securityGroupId != "" {
			return nil
		}
		return fmt.Errorf("security group %s still exists", rs.Primary.ID)
	}
	return nil
}

const testAccRedisSecurityGroupVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}

variable "security_group_id" {
	default = "` + defaultCrsSecurityGroup + `"
}
`

const testAccRedisSecurityGroup = testAccRedisSecurityGroupVar + `

resource "tencentcloud_redis_security_group" "security_group" {
    instance_id       = var.instance_id
    security_group_id = var.security_group_id
}

`
