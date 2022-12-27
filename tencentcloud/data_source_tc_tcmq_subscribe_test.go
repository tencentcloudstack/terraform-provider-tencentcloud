package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSourceTcmqSubscribe = "data.tencentcloud_tcmq_subscribe.subscribe"

func TestAccTencentCloudTcmqSubscribeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTcmqSubscribeDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceTcmqSubscribe, "subscription_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudTcmqSubscribeDataSource_basic = `
resource "tencentcloud_tcmq_topic" "topic" {
	topic_name = "test_subscribe_datasource_topic"
}
	
resource "tencentcloud_tcmq_subscribe" "subscribe" {
	topic_name = tencentcloud_tcmq_topic.topic.topic_name
	subscription_name = "test_subscribe"
	protocol = "http"
	endpoint = "http://mikatong.com"
  }
  
  data "tencentcloud_tcmq_subscribe" "subscribe" {
	topic_name = tencentcloud_tcmq_topic.topic.topic_name
	subscription_name = tencentcloud_tcmq_subscribe.subscribe.subscription_name
  }
`
