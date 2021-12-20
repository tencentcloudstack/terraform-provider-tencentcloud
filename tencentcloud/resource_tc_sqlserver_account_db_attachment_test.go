package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testSqlserverAccountDBAttachmentResourceName = "tencentcloud_sqlserver_account_db_attachment"
var testSqlserverAccountDBAttachmentResourceKey = testSqlserverAccountDBAttachmentResourceName + ".test"

func TestAccTencentCloudSqlserverAccountDBAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverAccountDBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverAccountDBAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountDBAttachmentExists(testSqlserverAccountDBAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountDBAttachmentResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", "test111"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "privilege", "ReadOnly"),
				),
			},
			{
				ResourceName:      testSqlserverAccountDBAttachmentResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccSqlserverAccountDBAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountDBAttachmentExists(testSqlserverAccountDBAttachmentResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountDBAttachmentResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", "test111"),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "privilege", "ReadWrite"),
				),
			},
		},
	})
}

func testAccCheckSqlserverAccountDBAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverAccountDBAttachmentResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 3 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		accountName := idStrs[1]
		dbName := idStrs[2]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete SQL Server account %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverAccountDBAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 3 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		accountName := idStrs[1]
		dbName := idStrs[2]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		if err != nil {
			_, has, err = service.DescribeAccountDBAttachmentById(ctx, instanceId, accountName, dbName)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("SQL Server account %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverAccountDBAttachment string = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_account" "test" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name = "tf_sqlserver_account"
  password = "testt123"
}

resource "tencentcloud_sqlserver_db" "test" {
  instance_id = "mssql-3cdq7kx5"
  name        = "test111"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "test" {
  instance_id = "mssql-3cdq7kx5"
  account_name = tencentcloud_sqlserver_account.test.name
  db_name = tencentcloud_sqlserver_db.test.name
  privilege = "ReadOnly"
}
`

const testAccSqlserverAccountDBAttachmentUpdate string = testAccSqlserverInstance + `
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
`
