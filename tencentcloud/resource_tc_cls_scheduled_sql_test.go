package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudClsScheduledSqlResource_basic -v
func TestAccTencentCloudClsScheduledSqlResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckClsScheduledSqlDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsScheduledSql,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsScheduledSqlExists("tencentcloud_cls_scheduled_sql.scheduled_sql"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "src_topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "enable_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "dst_resource.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "scheduled_sql_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_time_window"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "process_delay")),
			},
			{
				ResourceName:      "tencentcloud_cls_scheduled_sql.scheduled_sql",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsScheduledSqlDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clsService := ClsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_scheduled_sql" {
			continue
		}
		instance, err := clsService.DescribeClsScheduledSqlById(ctx, rs.Primary.ID)
		if err != nil {
			continue
		}
		if instance != nil {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Destroy] check: CLS ScheduledSql still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsScheduledSqlExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Exists] check: CLS ScheduledSql %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Create] check: CLS ScheduledSql id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		taskRes, err := clsService.DescribeClsScheduledSqlById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if taskRes == nil {
			return fmt.Errorf("[CHECK][CLS ScheduledSql][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsScheduledSql = `

resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}
resource "tencentcloud_cls_scheduled_sql" "scheduled_sql" {
  src_topic_id = tencentcloud_cls_topic.topic.id
  name = "tf-example"
  enable_flag = 1
  dst_resource {
    topic_id = tencentcloud_cls_topic.topic.id
    region = "ap-guangzhou"
    biz_type = 0
    metric_name = "test"

  }
  scheduled_sql_content = "xxx"
  process_start_time = 1723117637000
  process_type = 1
  process_period = 10
  process_time_window = "@m-15m,@m"
  process_delay = 5
  src_topic_region = "ap-guangzhou"
  syntax_rule = 0
}
`
