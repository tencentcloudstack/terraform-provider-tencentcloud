package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicFlowRankingDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicFlowRankingDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_flow_ranking.topic_flow_ranking")),
			},
		},
	})
}

const testAccCkafkaTopicFlowRankingDataSource = `
data "tencentcloud_ckafka_topic_flow_ranking" "topic_flow_ranking" {
	instance_id = "ckafka-vv7wpvae"
	ranking_type = "PRO"
	begin_date = "2023-05-29T00:00:00+08:00"
	end_date = "2021-05-29T23:59:59+08:00"
}
`
