package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingCbsStorageResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingCbsStorage_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_name", "tf-storage-basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "availability_zone", "ap-guangzhou-2"),
				),
			},
			{
				Config: testAccTestingCbsStorageUp_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_name", "tf-storage-basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "availability_zone", "ap-guangzhou-2"),
				),
			},
			{
				ResourceName:            "tencentcloud_cbs_storage.storage_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
		},
	})
}

const testAccTestingCbsStorage_basic = `
resource "tencentcloud_cbs_storage" "storage_basic" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-basic"
	storage_size      = 50
	availability_zone = "ap-guangzhou-2"
}
`
const testAccTestingCbsStorageUp_basic = `
resource "tencentcloud_cbs_storage" "storage_basic" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-basic"
	storage_size      = 100
	availability_zone = "ap-guangzhou-2"
}
`
