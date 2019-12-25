package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCbsStoragesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStoragesDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_storages.storages", "storage_list.0.storage_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.storage_name", "tf-test-storage"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.storage_size", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.encrypt", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.attached", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_storages.storages", "storage_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cbs_storages.storages", "storage_list.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_cbs_storages.storages", "storage_list.0.tags.test", "tf"),
				),
			},
		},
	})
}

const testAccCbsStoragesDataSource = `
resource "tencentcloud_cbs_storage" "storage" {
  storage_type      = "CLOUD_PREMIUM"
  storage_name      = "tf-test-storage"
  storage_size      = 50
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false
  
  tags = {
    test = "tf"
  }
}

data "tencentcloud_cbs_storages" "storages" {
  storage_id = tencentcloud_cbs_storage.storage.id
}
`
