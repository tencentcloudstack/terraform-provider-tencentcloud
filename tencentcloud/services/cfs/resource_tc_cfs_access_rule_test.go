package cfs_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcfs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCfsAccessRule(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCfsAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessRuleExists("tencentcloud_cfs_access_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_rule.foo", "auth_client_ip", "10.11.1.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_cfs_access_rule.foo", "priority", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfs_access_rule.foo", "access_group_id"),
				),
			},
		},
	})
}

func testAccCheckCfsAccessRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := localcfs.NewCfsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cfs_access_rule" {
			continue
		}

		accessGroupId := rs.Primary.Attributes["access_group_id"]
		accessRules, err := cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				accessRules, err = cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cfs access rule %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cfs access rule id is not set")
		}
		accessGroupId := rs.Primary.Attributes["access_group_id"]
		cfsService := localcfs.NewCfsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		accessRules, err := cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				accessRules, err = cfsService.DescribeAccessRule(ctx, accessGroupId, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
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

const testAccCfsAccessRule = DefaultCfsAccessGroup + `

resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = local.cfs_access_group_id
  auth_client_ip = "10.11.1.0/24"
  priority = 1
}
`
