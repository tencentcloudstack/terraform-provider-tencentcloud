package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCfsFileSystemClientsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsFileSystemClientsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfs_file_system_clients.file_system_clients")),
			},
		},
	})
}

const testAccCfsFileSystemClientsDataSource = `

data "tencentcloud_cfs_file_system_clients" "file_system_clients" {
  file_system_id = "cfs-iobiaxtj"
}

`
