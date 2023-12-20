package cfs_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	localcfs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfs"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cfs_access_group
	resource.AddTestSweepers("tencentcloud_cfs_access_group", &resource.Sweeper{
		Name: "tencentcloud_cfs_access_group",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := localcfs.NewCfsService(client)

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
						err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
							if delErr := service.DeleteAccessRule(ctx, id, ruleId); delErr != nil {
								// retry when Pgroup is under deleting rule operation
								return tccommon.RetryError(delErr)
							}
							return nil
						})
						if err != nil {
							return err
						}
					}
				}

				if tcacctest.IsResourcePersist(name, nil) || !strings.HasPrefix(name, "test") {
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

// go test -i; go test -test.run TestAccTencentCloudCfsAccessGroup_basic -v
func TestAccTencentCloudCfsAccessGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCfsAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessGroupExists("tencentcloud_cfs_access_group.example"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.example", "name", "tx_example"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.example", "description", "desc."),
				),
			},
			{
				ResourceName:      "tencentcloud_cfs_access_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfsAccessGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessGroupExists("tencentcloud_cfs_access_group.example"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.example", "name", "tx_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_group.example", "description", "desc update."),
				),
			},
		},
	})
}

func testAccCheckCfsAccessGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := localcfs.NewCfsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cfs_access_group" {
			continue
		}

		accessGroups, err := cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				accessGroups, err = cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
				if err != nil {
					return tccommon.RetryError(err)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cfs access group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cfs access group id is not set")
		}
		cfsService := localcfs.NewCfsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		accessGroups, err := cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				accessGroups, err = cfsService.DescribeAccessGroup(ctx, rs.Primary.ID, "")
				if err != nil {
					return tccommon.RetryError(err)
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
resource "tencentcloud_cfs_access_group" "example" {
  name        = "tx_example"
  description = "desc."
}
`

const testAccCfsAccessGroupUpdate = `
resource "tencentcloud_cfs_access_group" "example" {
  name        = "tx_example_update"
  description = "desc update."
}
`
