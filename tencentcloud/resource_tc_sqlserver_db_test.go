package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudSqlserverDB_basic_and_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDB_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "name", "testAccSqlserverDB"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "charset", "Chinese_PRC_BIN"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "remark", "testACC-remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.mysqlserver_db", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.mysqlserver_db", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_db.mysqlserver_db", "instance_id"),
				),
				Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_sqlserver_db.mysqlserver_db",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverDB_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverDBExists("tencentcloud_sqlserver_db.mysqlserver_db"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_db.mysqlserver_db", "remark", "testACC-remark_update"),
				),
			},
		},
	})
}

func testAccCheckSqlserverDBDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_db" {
			continue
		}
		_, has, err := sqlserverService.DescribeDBDetailsById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("SQLServer DB still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckSqlserverDBExists(n string) resource.TestCheckFunc {
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
		_, has, err := sqlserverService.DescribeDBDetailsById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("SQLServer DB %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccSqlserverDB_basic = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name        = "testAccSqlserverDB"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}`

const testAccSqlserverDB_basic_update_remark = testAccSqlserverInstance + `
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  name        = "testAccSqlserverDB"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark_update"
}`
