package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaDatahubTopicDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubTopicDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_datahub_topic.datahub_topic")),
			},
		},
	})
}

const testAccCkafkaDatahubTopicDataSource = `

data "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  search_word = "topicName"
  offset = 0
  limit = 20
  }

`
