package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRumTawAreaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumTawAreaDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_taw_area.taw_area")),
			},
		},
	})
}

const testAccRumTawAreaDataSource = `

data "tencentcloud_rum_taw_area" "taw_area" {
  area_ids = 
  area_keys = 
  area_statuses = 
}

`
