package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsScheduledSqlResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsScheduledSql,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_scheduled_sql.scheduled_sql", "id")),
			},
		},
	})
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
