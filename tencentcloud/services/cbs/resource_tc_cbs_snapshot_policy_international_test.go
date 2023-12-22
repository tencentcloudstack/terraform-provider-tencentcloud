package cbs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
)

func TestAccTencentCloudInternationalCbsResource_snapshotPolicy(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckInternationalCbsSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotPolicyExists("tencentcloud_cbs_snapshot_policy.snapshot_policy"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "snapshot_policy_name", "tf-test-snapshot-policy"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.0", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_weekdays.1", "3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "repeat_hours.0", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_snapshot_policy.snapshot_policy", "retention_days", "30"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_snapshot_policy.snapshot_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckInternationalCbsSnapshotPolicyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_snapshot_policy" {
			continue
		}

		policy, err := cbsService.DescribeSnapshotPolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if policy != nil {
			return fmt.Errorf("cbs snapshot policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckInternationalSnapshotPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs snapshot policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs snapshot policy id is not set")
		}
		cbsService := localcbs.NewCbsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		policy, err := cbsService.DescribeSnapshotPolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if policy == nil {
			return fmt.Errorf("cbs snapshot policy is not exist")
		}
		return nil
	}
}

const testAccInternationalCbsSnapshotPolicy = `
resource "tencentcloud_cbs_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "tf-test-snapshot-policy"
  repeat_weekdays      = [0, 3]
  repeat_hours         = [0]
  retention_days       = 30
}
`
