package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaTopicFlowRankingDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicFlowRankingDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_flow_ranking.topic_flow_ranking")),
			},
		},
	})
}

const testAccCkafkaTopicFlowRankingDataSource = `

data "tencentcloud_ckafka_topic_flow_ranking" "topic_flow_ranking" {
  instance_id = "InstanceId"
  ranking_type = "PRO"
  begin_date = "2021-05-13T07:23:11+08:00"
  end_date = "2021-05-14T07:23:11+08:00"
  }

`
