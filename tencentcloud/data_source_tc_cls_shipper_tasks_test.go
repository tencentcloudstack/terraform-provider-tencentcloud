package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClsShipperTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsShipperTasksDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cls_shipper_tasks.shipper_tasks")),
			},
		},
	})
}

const testAccClsShipperTasksDataSource = `

data "tencentcloud_cls_shipper_tasks" "shipper_tasks" {
  shipper_id = "dbde3c9b-ea16-4032-bc2a-d8fa65567a8e"
  start_time = 160749910700
  end_time = 160749910800
}

`
