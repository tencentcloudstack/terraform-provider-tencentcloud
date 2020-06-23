package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudSqlserverDB_basic_and_update(t *testing.T) {
	var instanceId string
	var dbName string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverDBDestroy(&instanceId, &dbName),
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDB_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverDBExists("tencentcloud_sqlserver_db.mysqlserver_db", &instanceId, &dbName),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "instance_id", "mssql-3cdq7kx5"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "name", "testAccSqlserverDB"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "charset", "Chinese_PRC_BIN"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "remark", "testACC-remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.mysqlserver_db", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.mysqlserver_db", "status"),
				),
			},
			{
				Config: testAccSqlserverDB_basic_update_name(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverDBExists("tencentcloud_sqlserver_db.mysqlserver_db", &instanceId, &dbName),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "name", "testAccSqlserverDB_update"),
				),
			},
			{
				Config: testAccSqlserverDB_basic_update_remark(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverDBExists("tencentcloud_sqlserver_db.mysqlserver_db", &instanceId, &dbName),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "remark", "testACC-remark_update"),
				),
			},
		},
	})
}

func testAccCheckSqlserverDBDestroy(instanceId *string, dbName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "tencentcloud_sqlserver_db" {
				continue
			}
			_, has, err := sqlserverService.DescribeDBDetailsByName(ctx, *instanceId, *dbName)
			if has {
				return fmt.Errorf("SQLServer DB still exists")
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckSqlserverDBExists(n string, instanceId *string, dbName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SQLServer DB %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SQLServer DB id is not set")
		}

		sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := sqlserverService.DescribeDBDetailsByName(ctx, rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		if !has {
			return fmt.Errorf("SQLServer DB %s is not found", rs.Primary.Attributes["name"])
		}
		if err != nil {
			return err
		}

		*instanceId = rs.Primary.Attributes["instance_id"]
		*dbName = rs.Primary.Attributes["name"]
		return nil
	}
}

func testAccSqlserverDB_basic() string {
	return `
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-3cdq7kx5"
  name        = "testAccSqlserverDB"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}`
}

func testAccSqlserverDB_basic_update_name() string {
	return `
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-3cdq7kx5"
  name        = "testAccSqlserverDB_update"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}`
}

func testAccSqlserverDB_basic_update_remark() string {
	return `
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-3cdq7kx5"
  name        = "testAccSqlserverDB_update"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark_update"
}`
}
