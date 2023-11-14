package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqTopic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_topic.topic", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqTopic = `

resource "tencentcloud_tdmq_rocketmq_topic" "topic" {
  topic = &lt;nil&gt;
  namespaces = &lt;nil&gt;
  type = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
  remark = &lt;nil&gt;
  partition_num = &lt;nil&gt;
      }

`
