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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverPublishSubscribeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverPublishSubscribe_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "publish_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "subscribe_instance_id"),
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
				ResourceName:            "tencentcloud_sqlserver_publish_subscribe.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_subscribe_db"},
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

const testAccSqlserverPublishSubscribe_basic = testAccSqlserverInstanceBasic + `
resource "tencentcloud_security_group" "foo" {
  name = "test-sg-tf"
}

resource "tencentcloud_sqlserver_instance" "publish_instance" {
  name                      = "tf_sqlserver_publish_instance"
  availability_zone         = var.availability_zone
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = "` + defaultVpcId + `"
  subnet_id                 = "` + defaultSubnetId + `"
  project_id                = 0
  memory                    = 2
  storage                   = 10
  maintenance_week_set      = [1,2,3]
  maintenance_start_time    = "09:00"
  maintenance_time_span     = 3
  security_groups           = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_instance" "subscribe_instance" {
  name                      = "tf_sqlserver_subscribe_instance"
  availability_zone         = var.availability_zone
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = "` + defaultVpcId + `"
  subnet_id                 = "` + defaultSubnetId + `"
  project_id                = 0
  memory                    = 2
  storage                   = 10
  maintenance_week_set      = [1,2,3]
  maintenance_start_time    = "09:00"
  maintenance_time_span     = 3
  security_groups           = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_db" "test_publish_subscribe" {
  instance_id = tencentcloud_sqlserver_instance.publish_instance.id
  name        = "test111"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.publish_instance.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.test_publish_subscribe.name
	}
}`

const testAccSqlserverPublishSubscribe_basic_update_name = testAccSqlserverInstanceBasic + `
resource "tencentcloud_sqlserver_instance" "publish_instance" {
  name                          = "tf_sqlserver_publish_instance"
  availability_zone             = var.availability_zone
  charge_type                   = "POSTPAID_BY_HOUR"
  vpc_id                        = "` + defaultVpcId + `"
  subnet_id                     = "` + defaultSubnetId + `"
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  security_groups               = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_instance" "subscribe_instance" {
  name                          = "tf_sqlserver_subscribe_instance"
  availability_zone             = var.availability_zone
  charge_type                   = "POSTPAID_BY_HOUR"
  vpc_id                        = "` + defaultVpcId + `"
  subnet_id                     = "` + defaultSubnetId + `"
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  security_groups               = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_db" "test_publish_subscribe" {
  instance_id = tencentcloud_sqlserver_instance.publish_instance.id
  name        = "test111"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.publish_instance.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example1"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.test_publish_subscribe.name
	}
}`
