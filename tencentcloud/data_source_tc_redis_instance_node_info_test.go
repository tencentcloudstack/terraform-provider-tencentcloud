package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceNodeInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_node_info.instance_node_info")),
			},
		},
	})
}

const testAccRedisInstanceNodeInfoDataSource = `

data "tencentcloud_redis_instance_node_info" "instance_node_info" {
  instance_id = "crs-c1nl9rpv"
        }

`
