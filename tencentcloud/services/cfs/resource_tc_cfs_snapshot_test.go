package cfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsSnapshotResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsSnapshot,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfs_snapshot.snapshot", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfs_snapshot.snapshot",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfsSnapshot = `

resource "tencentcloud_cfs_snapshot" "snapshot" {
  file_system_id = "cfs-iobiaxtj"
  snapshot_name = "test"
  tags = {
    "createdBy" = "terraform"
  }
}

`
