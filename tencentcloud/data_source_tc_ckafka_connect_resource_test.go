package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaConnectResourceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaConnectResourceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_connect_resource.connect_resource")),
			},
		},
	})
}

const testAccCkafkaConnectResourceDataSource = `

data "tencentcloud_ckafka_connect_resource" "connect_resource" {
  type = "DTS"
  search_word = "resourceName"
  offset = 0
  limit = 20
  resource_region = "region"
  }

`
