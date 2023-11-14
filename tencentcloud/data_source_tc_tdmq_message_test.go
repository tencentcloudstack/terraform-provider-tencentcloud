package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqMessageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqMessageDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_message.message")),
			},
		},
	})
}

const testAccTdmqMessageDataSource = `

data "tencentcloud_tdmq_message" "message" {
  cluster_id = ""
  environment_id = ""
  topic_name = ""
  msg_id = ""
  query_dlq_msg = 
            }

`
