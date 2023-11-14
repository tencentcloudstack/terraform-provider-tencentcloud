package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudChdfsFileSystemResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsFileSystem,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_file_system.file_system", "id")),
			},
			{
				ResourceName:      "tencentcloud_chdfs_file_system.file_system",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsFileSystem = `

resource "tencentcloud_chdfs_file_system" "file_system" {
  file_system_name = "test file system name"
  capacity_quota = 1073741824
  posix_acl = true
  description = "test desc"
  super_users = &lt;nil&gt;
  enable_ranger = false
  ranger_service_addresses = &lt;nil&gt;
}

`
