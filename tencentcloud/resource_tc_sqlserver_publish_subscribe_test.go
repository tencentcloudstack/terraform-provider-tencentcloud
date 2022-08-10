package tencentcloud

import (
	"context"
	"fmt"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_publish_subscribe
	resource.AddTestSweepers("tencentcloud_sqlserver_publish_subscribe", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_publish_subscribe",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := SqlserverService{client}
			instance, err := service.DescribeSqlserverInstances(ctx, "", defaultPubSQLServerName, -1, "", "", 1)
			if err != nil {
				return err
			}
			subInstances, err := service.DescribeSqlserverInstances(ctx, "", defaultSubSQLServerName, -1, "", "", 1)

			if err != nil {
				return err
			}

			pubInstanceId := *instance[0].InstanceId
			subInstanceId := *subInstances[0].InstanceId

			testAccUnsubscribePubDB(ctx, &service, pubInstanceId)

			database, err := service.DescribeDBsOfInstance(ctx, subInstanceId)
			if err != nil {
				return err
			}

			if len(database) == 0 {
				log.Printf("no DBs in %s", subInstanceId)
				return nil
			}

			for i := range database {
				item := database[i]
				created := time.Time{}
				name := *item.Name
				if item.CreateTime != nil {
					created = stringTotime(*item.CreateTime)
				}
				if name != defaultSQLServerPubSubDB || isResourcePersist("", &created) {
					continue
				}
				if err := service.DeleteSqlserverDB(ctx, subInstanceId, []*string{item.Name}); err != nil {
					log.Printf("err: %s", err.Error())
				}
			}
			return nil
		},
	})
}

func testAccUnsubscribePubDB(ctx context.Context, service *SqlserverService, instanceId string) {

	pubsubs, err := service.DescribeSqlserverPublishSubscribes(ctx, map[string]interface{}{
		"instanceId": instanceId,
	})

	if len(pubsubs) == 0 {
		log.Printf("NO pubsub result")
		return
	}

	pubSubId := *pubsubs[0].Id

	pubSub := &sqlserver.PublishSubscribe{
		Id: &pubSubId,
	}
	tuples := []*sqlserver.DatabaseTuple{
		{
			helper.String(defaultSQLServerPubSubDB),
			helper.String(defaultSQLServerPubSubDB),
		},
	}
	err = service.DeletePublishSubscribe(ctx, pubSub, tuples)
	if err != nil {
		fmt.Printf("[ERROR] %s", err.Error())
	}
}

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
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "publish_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "subscribe_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_subscribe_name", "example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "database_tuples.#", "1"),
				),
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

const testAccSqlserverPublishSubscribe_basic = CommonPubSubSQLServer + `

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = local.pub_sqlserver_id
	subscribe_instance_id           = local.sub_sqlserver_id
	publish_subscribe_name          = "example"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = local.sqlserver_pubsub_db
	}
}`

const testAccSqlserverPublishSubscribe_basic_update_name = CommonPubSubSQLServer + `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = local.pub_sqlserver_id
	subscribe_instance_id           = local.sub_sqlserver_id
	publish_subscribe_name          = "example1"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = local.sqlserver_pubsub_db
	}
}`
