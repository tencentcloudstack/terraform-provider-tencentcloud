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
			{
				ResourceName:      "tencentcloud_cls_scheduled_sql.scheduled_sql",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsScheduledSql = `

resource "tencentcloud_cls_scheduled_sql" "scheduled_sql" {
  src_topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  name = "task"
  enable_flag = 1
  dst_resource {
		topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
		region = "ap-guangzhou"
		biz_type = 0
		metric_name = "test"

  }
  scheduled_sql_content = "xxx"
  process_start_time = 1690515360000
  process_type = 1
  process_period = 10
  process_time_window = "@m-15m,@m"
  process_delay = 5
  src_topic_region = "ap-guangzhou"
  process_end_time = 1690515360000
  syntax_rule = 0
}

`
