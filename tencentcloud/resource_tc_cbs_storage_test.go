package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCbsStorage_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_name", "tf-storage-basic"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "storage_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_basic", "availability_zone", "ap-guangzhou-3"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_storage.storage_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudCbsStorage_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCbsStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorage_full,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_name", "tf-storage-full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "tags.test", "tf"),
				),
			},
			{
				Config: testAccCbsStorage_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageExists("tencentcloud_cbs_storage.storage_full"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_name", "tf-storage-update"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "storage_size", "60"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_storage.storage_full", "tags.test", "tf"),
				),
			},
		},
	})
}

func testAccCheckCbsStorageDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cbsService := CbsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_storage" {
			continue
		}

		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if storage != nil {
			return fmt.Errorf("cbs storage still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckStorageExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs storage %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs storage id is not set")
		}
		cbsService := CbsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if storage == nil {
			return fmt.Errorf("cbs storage is not exist")
		}
		return nil
	}
}

const testAccCbsStorage_basic = `
resource "tencentcloud_cbs_storage" "storage_basic" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-basic"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
}
`

const testAccCbsStorage_full = `
resource "tencentcloud_cbs_storage" "storage_full" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-full"
	storage_size      = 50
	availability_zone = "ap-guangzhou-3"
	project_id = 0
	encrypt = false
	tags = {
		test = "tf"
	}
}
`
const testAccCbsStorage_update = `
resource "tencentcloud_cbs_storage" "storage_full" {
	storage_type      = "CLOUD_PREMIUM"
	storage_name      = "tf-storage-update"
	storage_size      = 60
	availability_zone = "ap-guangzhou-3"
	project_id = 0
	encrypt = false
	tags = {
		test = "tf"
	}
}
`
