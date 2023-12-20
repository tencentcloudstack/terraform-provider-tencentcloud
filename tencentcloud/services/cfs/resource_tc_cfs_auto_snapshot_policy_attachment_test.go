package cfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsAutoSnapshotPolicyAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAutoSnapshotPolicyAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfs_auto_snapshot_policy_attachment.auto_snapshot_policy_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfs_auto_snapshot_policy_attachment.auto_snapshot_policy_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfsAutoSnapshotPolicyAttachment = `

resource "tencentcloud_cfs_auto_snapshot_policy_attachment" "auto_snapshot_policy_attachment" {
  auto_snapshot_policy_id = "asp-basic"
  file_system_ids         = "cfs-iobiaxtj"
}

`
