package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverAccountDBAttachmentsName = "data.tencentcloud_sqlserver_account_db_attachments.test"

func TestAccTencentCloudDataSqlserverAccountDBAttachments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSqlserverAccountDBAttachmentsName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataSqlserverAccountDBAttachmentsName, "list.0.instance_id", "mssql-3cdq7kx5"),
					resource.TestCheckResourceAttr(testDataSqlserverAccountDBAttachmentsName, "list.0.account_name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.db_name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverAccountDBAttachmentsName, "list.0.privilege"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverAccountDBAttachmentsBasic = testAccSqlserverDB_basic + `
data "tencentcloud_sqlserver_account_db_attachments" "test"{
  instance_id = tencentcloud_sqlserver_instance.test.id
  account_name = "tf_sqlserver_account"
}
`
