package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Config: testAccChdfsFileSystemUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_chdfs_file_system.file_system", "id"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_file_system.file_system", "file_system_name", "terraform-test-for"),
				),
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
  capacity_quota           = 10995116277760
  description              = "file system for terraform test"
  enable_ranger            = true
  file_system_name         = "terraform-test"
  posix_acl                = false
  ranger_service_addresses = [
    "127.0.0.1:80",
    "127.0.0.1:8000",
  ]
  super_users              = [
    "terraform",
    "iac",
  ]
}

`

const testAccChdfsFileSystemUpdate = `

resource "tencentcloud_chdfs_file_system" "file_system" {
  capacity_quota           = 10995116277760
  description              = "file system for terraform test"
  enable_ranger            = true
  file_system_name         = "terraform-test-for"
  posix_acl                = false
  ranger_service_addresses = [
    "127.0.0.1:80",
    "127.0.0.1:8000",
  ]
  super_users              = [
    "terraform",
    "iac",
  ]
}

`
