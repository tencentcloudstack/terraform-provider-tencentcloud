package cdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMysqlClsLogAttachmentResource -v
func TestAccTencentCloudNeedFixMysqlClsLogAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_cls_log_attachment.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "instance_id", "cdb-8fk7id2l"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "log_type", "slowlog"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "log_set", "cb4ca863-13b4-489a-b3ef-aae3b0f1b172"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "log_topic", "286cb422-c09d-4d9c-96ab-107b6bba1561"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "period", "30"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "create_index", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_cls_log_attachment.example", "cls_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_cls_log_attachment.example", "allow_disk_redirect"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_cls_log_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsLogAttachment = `
resource "tencentcloud_mysql_cls_log_attachment" "example" {
  instance_id  = "cdb-8fk7id2l"
  log_type     = "slowlog"
  log_set      = "cb4ca863-13b4-489a-b3ef-aae3b0f1b172"
  log_topic    = "286cb422-c09d-4d9c-96ab-107b6bba1561"
  period       = 30
  create_index = true
  cls_region   = "ap-guangzhou"
}
`
