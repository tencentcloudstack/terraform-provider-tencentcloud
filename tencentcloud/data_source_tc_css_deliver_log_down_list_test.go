package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssDeliverLogDownListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssDeliverLogDownListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_deliver_log_down_list.deliver_log_down_list"),
				),
			},
		},
	})
}

const testAccCssDeliverLogDownListDataSource = `

data "tencentcloud_css_deliver_log_down_list" "deliver_log_down_list" {
}

`
