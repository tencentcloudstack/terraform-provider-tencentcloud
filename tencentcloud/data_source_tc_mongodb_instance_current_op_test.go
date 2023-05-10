package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  instance_id = "cmgo-gwqk8669"
  op = "command"
  order_by_type = "desc"
}

`
