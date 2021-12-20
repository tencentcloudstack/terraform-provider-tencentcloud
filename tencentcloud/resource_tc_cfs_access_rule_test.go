package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCfsAccessRule(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCfsAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessRuleExists("tencentcloud_cfs_access_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_rule.foo", "auth_client_ip", "10.10.1.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_rule.foo", "priority", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfs_access_rule.foo", "access_group_id"),
				),
			},
		},
	})
}

func testAccCheckCfsAccessRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cfsService := CfsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cfs_access_rule" {
			continue
		}

		accessGroupId := rs.Primary.Attributes["access_group_id"]
		accessRules, err := cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				accessRules, err = cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if len(accessRules) > 0 {
			return fmt.Errorf("cfs access rule still exist: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCfsAccessRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cfs access rule %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cfs access rule id is not set")
		}
		accessGroupId := rs.Primary.Attributes["access_group_id"]
		cfsService := CfsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		accessRules, err := cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				accessRules, err = cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if len(accessRules) < 1 {
			return fmt.Errorf("cfs access rule is not found")
		}
		return nil
	}
}

const testAccCfsAccessRule = `
resource "tencentcloud_cfs_access_group" "foo" {
  name = "test_cfs_access_rule"
}

resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = tencentcloud_cfs_access_group.foo.id
  auth_client_ip = "10.10.1.0/24"
  priority = 1
}
`
