package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverAccountDBAttachmentsName = "data.tencentcloud_sqlserver_account_db_attachments.test"

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverAccountDBAttachments_basic -v
func TestAccDataSourceTencentCloudSqlserverAccountDBAttachments_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.#"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.account_name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.db_name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.privilege"),
				),
			},
		},
	})
}

const testAccSQLServerAttachDataDB = "test_db_attachment"

const testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic = CommonPresetSQLServerAccount + `
resource "tencentcloud_sqlserver_db" "test" {
  instance_id = local.sqlserver_id
  name        = "` + testAccSQLServerAttachDataDB + `"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "test" {
  instance_id  = local.sqlserver_id
  account_name = local.sqlserver_account
  db_name      = tencentcloud_sqlserver_db.test.name
  privilege    = "ReadWrite"
}
data "tencentcloud_sqlserver_account_db_attachments" "test"{
  instance_id  = local.sqlserver_id
  account_name = tencentcloud_sqlserver_account_db_attachment.test.account_name
}
`
