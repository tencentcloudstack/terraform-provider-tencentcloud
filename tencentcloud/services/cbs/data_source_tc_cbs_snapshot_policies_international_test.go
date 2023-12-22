package cbs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalCbsDataSource_snapshotPolicies(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCbsSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalCbsSnapshotPoliciesDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSnapshotPolicyExists("tencentcloud_cbs_snapshot_policy.policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.0.snapshot_policy_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.0.snapshot_policy_name", "tf-test-snapshot-policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.0.repeat_weekdays.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.0.repeat_hours.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_snapshot_policies.policies", "snapshot_policy_list.0.retention_days", "30"),
				),
			},
		},
	})
}

const testAccInternationalCbsSnapshotPoliciesDataSource = `
resource "tencentcloud_cbs_snapshot_policy" "policy" {
  snapshot_policy_name = "tf-test-snapshot-policy"
  repeat_weekdays      = [0, 3]
  repeat_hours         = [0]
  retention_days       = 30
}

data "tencentcloud_cbs_snapshot_policies" "policies" {
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.policy.id
  snapshot_policy_name = tencentcloud_cbs_snapshot_policy.policy.snapshot_policy_name
}
`
