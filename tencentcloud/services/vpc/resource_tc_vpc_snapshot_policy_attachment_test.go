package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSnapshotPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSnapshotPolicyAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_snapshot_policy_attachment.snapshot_policy_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_snapshot_policy_attachment.snapshot_policy_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcSnapshotPolicyAttachment = `

resource "tencentcloud_vpc_snapshot_policy" "snapshot_policy" {
  snapshot_policy_name = "terraform-attachment"
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

resource "tencentcloud_vpc_snapshot_policy_attachment" "snapshot_policy_attachment" {
  snapshot_policy_id = tencentcloud_vpc_snapshot_policy.snapshot_policy.id

  instances {
    instance_id        = "sg-r8ibzbd9"
    instance_name      = "cm-eks-cls-eizsc1iw-security-group"
    instance_region    = "ap-guangzhou"
    instance_type      = "securitygroup"
  }
  instances {
    instance_id        = "sg-k3tn70lh"
    instance_name      = "keep-ci-temp-test-sg"
    instance_region    = "ap-guangzhou"
    instance_type      = "securitygroup"
  }
}

`
