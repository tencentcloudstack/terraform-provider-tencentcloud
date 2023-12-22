package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSnapshotPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSnapshotPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_snapshot_policy.snapshot_policy", "id")),
			},
			{
				Config: testAccVpcSnapshotPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_snapshot_policy.snapshot_policy", "snapshot_policy_name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_snapshot_policy.snapshot_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcSnapshotPolicy = `

resource "tencentcloud_vpc_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "terraform-test"
  backup_type          = "time"
  cos_bucket           = "cos-lock-1308919341"
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "02:03:03"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "04:13:23"
  }
}

`

const testAccVpcSnapshotPolicyUpdate = `

resource "tencentcloud_vpc_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "terraform-for-test"
  backup_type          = "time"
  cos_bucket           = "cos-lock-1308919341"
  cos_region           = "ap-guangzhou"
  create_new_cos       = false
  keep_time            = 2

  backup_policies {
    backup_day  = "monday"
    backup_time = "00:00:00"
  }
  backup_policies {
    backup_day  = "tuesday"
    backup_time = "02:03:03"
  }
  backup_policies {
    backup_day  = "wednesday"
    backup_time = "04:13:23"
  }
}

`
