package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCbsStorageAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCbsStorageAttached("tencentcloud_cbs_storage.my_storage", "tencentcloud_instance.my_instance"),
				),
			},
		},
	})
}

func testAccCheckCbsStorageAttached(storageTag string, instanceTag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cbsRs, ok := s.RootModule().Resources[storageTag]
		if !ok {
			return fmt.Errorf("Storage not found: %s", storageTag)
		}

		insRs, ok := s.RootModule().Resources[instanceTag]
		if !ok {
			return fmt.Errorf("Instance not found: %s", instanceTag)
		}

		if cbsRs.Primary.ID == "" {
			return fmt.Errorf("Storage no ID is set")
		}

		if insRs.Primary.ID == "" {
			return fmt.Errorf("Instance no ID is set")
		}

		provider := testAccProvider
		// Ignore if Meta is empty, this can happen for validation providers
		if provider.Meta() == nil {
			return fmt.Errorf("Provider Meta is nil")
		}
		client := provider.Meta().(*TencentCloudClient).commonConn

		if _, err := waitInstanceReachTargetStatus(client, []string{insRs.Primary.ID}, "PENDING"); err != nil {
			return err
		}
		if _, err := waitInstanceReachOneOfTargetStatusList(client, []string{insRs.Primary.ID}, []string{"RUNNING", "STOPPED"}); err != nil {
			return err
		}

		storage, _, err := describeCbsStorage(cbsRs.Primary.ID, client)
		if err != nil {
			return err
		}

		if storage.Attached == 1 {
			if storage.InstanceId != insRs.Primary.ID {
				return fmt.Errorf("disk(%s) is attached in %s, not %s", cbsRs.Primary.ID, storage.InstanceId, insRs.Primary.ID)
			}
			return nil
		}
		return fmt.Errorf("disk(%s) not attached!", cbsRs.Primary.ID)
	}

}

const testAccCbsStorageAttachmentConfig = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

resource "tencentcloud_cbs_storage" "my_storage" {
  availability_zone = "ap-guangzhou-3"
  storage_size      = 100
  storage_type      = "cloudSSD"
  period            = 1
  storage_name      = "testAccCbsStorageTest"
}

resource "tencentcloud_instance" "my_instance" {
  instance_name 	= "testAccCbsAttachmentTest"
  availability_zone = "ap-guangzhou-3"
  image_id      	= "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type 	= "S3.SMALL1"
  system_disk_type  = "CLOUD_SSD"
}

resource "tencentcloud_cbs_storage_attachment" "my_attachment" {
  storage_id  = "${tencentcloud_cbs_storage.my_storage.id}"
  instance_id = "${tencentcloud_instance.my_instance.id}"
}
`
