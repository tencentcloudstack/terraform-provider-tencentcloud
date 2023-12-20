package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaDatahubTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubTask,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_datahub_task.datahub_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_ckafka_datahub_task.datahub_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCkafkaDatahubTask = `

resource "tencentcloud_ckafka_datahub_task" "datahub_task" {
	task_name = "test-task123321"
	task_type = "SOURCE"
	source_resource {
		  type = "POSTGRESQL"
		  postgre_sql_param {
			  database = "postgres"
			  table = "*"
			  resource = "resource-y9nxnw46"
			  plugin_name = "decoderbufs"
			  snapshot_mode = "never"
			  is_table_regular = false
			  key_columns = ""
			  record_with_schema = false
		  }
	}
	target_resource {
		  type = "TOPIC"
		  topic_param {
			  compression_type = "none"
			  resource = "1308726196-keep-topic"
			  use_auto_create_topic = false
		  }
	}
  }
`
