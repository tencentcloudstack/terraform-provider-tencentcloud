package tcmq_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSourceTcmqTopic = "data.tencentcloud_tcmq_topic.topic"

func TestAccTencentCloudTcmqTopicDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTcmqTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTcmqTopicDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceTcmqTopic, "topic_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudTcmqTopicDataSource_basic = `
resource "tencentcloud_tcmq_topic" "topic" {
	topic_name = "test_topic_datasource"
}
	
data "tencentcloud_tcmq_topic" "topic" {
	topic_name = tencentcloud_tcmq_topic.topic.topic_name
}
`
