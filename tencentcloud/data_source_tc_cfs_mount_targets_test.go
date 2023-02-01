package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCfsMountTargetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsMountTargetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfs_mount_targets.mount_targets")),
			},
		},
	})
}

const testAccCfsMountTargetsDataSource = `

data "tencentcloud_cfs_mount_targets" "mount_targets" {
  file_system_id = "cfs-iobiaxtj"
}

`
