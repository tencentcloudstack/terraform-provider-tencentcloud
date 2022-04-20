package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverAccountDBAttachmentsName = "data.tencentcloud_sqlserver_account_db_attachments.test"

func TestAccDataSourceTencentCloudSqlserverAccountDBAttachments(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSqlserverAccountDBAttachmentsName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataSqlserverAccountDBAttachmentsName, "list.0.account_name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.db_name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.privilege"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = "tf_sqlserver_account"
  password = "testt123"
}
resource "tencentcloud_sqlserver_db" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name        = "test111"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  account_name = tencentcloud_sqlserver_account.test.name
  db_name = tencentcloud_sqlserver_db.test.name
  privilege = "ReadWrite"
}
data "tencentcloud_sqlserver_account_db_attachments" "test"{
  instance_id = tencentcloud_sqlserver_instance.test.id
  account_name = tencentcloud_sqlserver_account_db_attachment.test.account_name
}
`
