package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCbsStorageAttachment(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsbStorageAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsStorageAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCbsStorageAttachmentExists("tencentcloud_cbs_storage_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_storage_attachment.foo", "storage_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_storage_attachment.foo", "instance_id"),
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
		if storage == nil {
			continue
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
		if storage == nil || *storage.Attached == false {
			return fmt.Errorf("cbs storage attchment not exists")
		}
		return nil
	}
}

const testAccCbsStorageAttachmentConfig = instanceCommonTestCase + `
resource "tencentcloud_cbs_storage" "foo" {
  availability_zone = var.availability_zone
  storage_size      = 100
  storage_type      = "CLOUD_PREMIUM"
  storage_name      = var.instance_name
}

resource "tencentcloud_cbs_storage_attachment" "foo" {
  storage_id  = tencentcloud_cbs_storage.foo.id
  instance_id = tencentcloud_instance.default.id
}
`
