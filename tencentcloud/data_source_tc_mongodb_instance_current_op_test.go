package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstanceCurrentOpDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceCurrentOpDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_current_op.instance_current_op")),
			},
		},
	})
}

const testAccMongodbInstanceCurrentOpDataSource = `

data "tencentcloud_mongodb_instance_current_op" "instance_current_op" {
  instance_id = "cmgo-9d0p6umb"
  ns = ""
  millisecond_running = 10
  op = "update"
  replica_set_name = ""
  state = "secondary"
  limit = 10
  offset = 0
  order_by = ""
  order_by_type = "desc"
  }

`
