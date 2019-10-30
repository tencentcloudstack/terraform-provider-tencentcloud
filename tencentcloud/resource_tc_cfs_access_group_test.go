package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCfsAccessGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCfsAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessGroupExists("tencentcloud_cfs_access_group.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.foo", "name", "test_cfs_access_group"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.foo", "description", "test"),
				),
			},
		},
	})
}

func testAccCheckCfsAccessGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cfsService := CfsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cfs_access_group" {
			continue
		}

		accessGroups, err := cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				accessGroups, err = cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if len(accessGroups) > 0 {
			return fmt.Errorf("cfs access group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCfsAccessGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cfs access group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cfs access group id is not set")
		}
		cfsService := CfsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		accessGroups, err := cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				accessGroups, err = cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if len(accessGroups) < 1 {
			return fmt.Errorf("cfs access group is not found")
		}
		return nil
	}
}

const testAccCfsAccessGroup = `
resource "tencentcloud_cfs_access_group" "foo" {
  name = "test_cfs_access_group"
  description = "test"
}
`
