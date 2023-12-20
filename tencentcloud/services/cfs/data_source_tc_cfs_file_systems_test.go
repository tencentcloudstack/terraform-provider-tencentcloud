package cfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsFileSystemsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCfsFileSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsFileSystemsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsFileSystemExists("tencentcloud_cfs_file_system.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.file_system_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.protocol", "NFS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.access_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_file_systems.file_systems", "file_system_list.0.mount_ip"),
				),
			},
		},
	})
}

const testAccCfsFileSystemsDataSource = DefaultCfsAccessGroup + `
resource "tencentcloud_vpc" "vpc" {
  name       = "test-cfs-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "test-cfs-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_cfs_file_system" "foo" {
  name = "test_cfs_file_system"
  availability_zone = "ap-guangzhou-3"
  access_group_id = local.cfs_access_group_id
  protocol = "NFS"
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_cfs_file_systems" "file_systems" {
  file_system_id = tencentcloud_cfs_file_system.foo.id
  name = tencentcloud_cfs_file_system.foo.name
}
`
