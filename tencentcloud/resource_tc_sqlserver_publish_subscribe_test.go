package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudSqlserverPublishSubscribeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverPublishSubscribeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverPublishSubscribe_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_instance_id", "mssql-82lhybgn"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "subscribe_instance_id", "mssql-12a60qdd"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_subscribe_name", "example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "database_tuples.#", "1"),
				),
				Destroy: false,
			},
			{
				Config: testAccSqlserverPublishSubscribe_basic_update_name,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_subscribe_name", "example1"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_publish_subscribe.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverPublishSubscribeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_publish_subscribe" {
			continue
		}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			continue
		}
		_, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, split[0], split[1])
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("SQL Server Publish Subscribe %s  still exists", split[0]+FILED_SP+split[1])
		}
	}
	return nil
}

func testAccCheckSqlserverPublishSubscribeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SQL Server Publish Subscribe %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SQL Server Publish Subscribe id is not set")
		}

		sqlserverService := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("SQL Server Publish Subscribe is not set: %s", rs.Primary.ID)
		}
		_, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, split[0], split[1])
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("SQL Server Publish Subscribe %s is not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccSqlserverPublishSubscribe_basic = `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = "mssql-82lhybgn"
	subscribe_instance_id           = "mssql-12a60qdd"
	publish_subscribe_name          = "example"
	database_tuples {
		publish_database            = "db_test_name"
		subscribe_database          = "db_test_name"
	}
}`

const testAccSqlserverPublishSubscribe_basic_update_name = `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = "mssql-82lhybgn"
	subscribe_instance_id           = "mssql-12a60qdd"
	publish_subscribe_name          = "example1"
	database_tuples {
		publish_database            = "db_test_name"
		subscribe_database          = "db_test_name"
	}
}`
