package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cfs_access_group
	resource.AddTestSweepers("tencentcloud_cfs_access_group", &resource.Sweeper{
		Name: "tencentcloud_cfs_access_group",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := CfsService{client}

			groups, err := service.DescribeAccessGroup(ctx, "", "")

			if err != nil {
				return err
			}

			for i := range groups {
				id := *groups[i].PGroupId
				name := *groups[i].Name

				rules, err := service.DescribeAccessRule(ctx, id, "")

				if err == nil { // ignore deleting the access rules when an error happened
					for _, item := range rules {
						ruleId := *item.RuleId
						err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
							if delErr := service.DeleteAccessRule(ctx, id, ruleId); delErr != nil {
								// retry when Pgroup is under deleting rule operation
								return retryError(delErr)
							}
							return nil
						})
						if err != nil {
							return err
						}
					}
				}

				if isResourcePersist(name, nil) || !strings.HasPrefix(name, "test") {
					continue
				}
				if err := service.DeleteAccessGroup(ctx, id); err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCfsAccessGroup(t *testing.T) {
	t.Parallel()
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
