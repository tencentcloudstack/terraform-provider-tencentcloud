package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaDatahubTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubTaskDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_datahub_task.datahub_task")),
			},
		},
	})
}

const testAccCkafkaDatahubTaskDataSource = `

data "tencentcloud_ckafka_datahub_task" "datahub_task" {
  limit = 20
  offset = 0
  search_word = "SearchWord"
  target_type = "CKafka"
  task_type = "SOURCE"
  source_type = "CKafka"
  resource = "Resource"
  }

`
