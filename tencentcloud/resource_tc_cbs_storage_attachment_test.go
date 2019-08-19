package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCbsStorageAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsbStorageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCbsStorageAttachmentExists("tencentcloud_cbs_storage_attachment.my_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_storage_attachment.my_attachment", "storage_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_storage_attachment.my_attachment", "instance_id"),
				),
			},
		},
	})
}

func testAccCheckCsbStorageAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	cbsService := CbsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cbs_storage_attachment" {
			continue
		}

		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err.Error() == "disk id is not found" {
			return nil
		}
		if err != nil {
			return err
		}
		if *storage.Attached {
			return fmt.Errorf("cbs storage attchment still exists")
		}

	}

	return nil
}

func testAccCheckCbsStorageAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cbs storage attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cbs storage attachment id is not set")
		}
		cbsService := CbsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		storage, err := cbsService.DescribeDiskById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if *storage.Attached == false {
			return fmt.Errorf("cbs storage attchment not exists")
		}
		return nil
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
  storage_type      = "CLOUD_PREMIUM"
  storage_name      = "tf-test-storage"
}

resource "tencentcloud_instance" "my_instance" {
  instance_name 	= "tf-test-instance"
  availability_zone = "ap-guangzhou-3"
  image_id      	= "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type 	= "S3.SMALL1"
  system_disk_type  = "CLOUD_PREMIUM"
}

resource "tencentcloud_cbs_storage_attachment" "my_attachment" {
  storage_id  = "${tencentcloud_cbs_storage.my_storage.id}"
  instance_id = "${tencentcloud_instance.my_instance.id}"
}
`
