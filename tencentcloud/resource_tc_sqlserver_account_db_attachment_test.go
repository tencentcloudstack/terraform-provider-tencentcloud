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

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_account_db_attachment
	resource.AddTestSweepers(testSqlserverAccountDBAttachmentResourceName, &resource.Sweeper{
		Name: testSqlserverAccountDBAttachmentResourceName,
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := SqlserverService{client}

			db, err := service.DescribeSqlserverInstances(ctx, "", defaultSQLServerName, -1, "", "", -1)

			if err != nil {
				return err
			}

			if len(db) == 0 {
				return fmt.Errorf("%s not exists", defaultSQLServerName)
			}

			instanceId := *db[0].InstanceId

			records, err := service.DescribeAccountDBAttachments(ctx, instanceId, defaultSQLServerAccount, defaultSQLServerDB)
			if err != nil {
				return err
			}

			if len(records) > 0 {
				err = service.ModifyAccountDBAttachment(ctx, instanceId, defaultSQLServerAccount, defaultSQLServerDB, "Delete")
			}

			if err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudSqlserverAccountDBAttachmentResource(t *testing.T) {
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
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", defaultSQLServerAccount),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", defaultSQLServerDB),
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
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "account_name", defaultSQLServerAccount),
					resource.TestCheckResourceAttr(testSqlserverAccountDBAttachmentResourceKey, "db_name", defaultSQLServerDB),
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

const testAccSqlserverAccountDBAttachment string = CommonPresetSQLServerAccount + `
resource "tencentcloud_sqlserver_account_db_attachment" "test" {
  instance_id = local.sqlserver_id
  account_name = local.sqlserver_account # "keep_sqlserver_account"
  db_name = local.sqlserver_db # "keep_sqlserver_db"
  privilege = "ReadOnly"
}
`

const testAccSqlserverAccountDBAttachmentUpdate string = CommonPresetSQLServerAccount + `
resource "tencentcloud_sqlserver_account_db_attachment" "test" {
  instance_id = local.sqlserver_id
  account_name = local.sqlserver_account
  db_name = local.sqlserver_db
  privilege = "ReadWrite"
}
`
