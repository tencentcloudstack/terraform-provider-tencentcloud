package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssWatermarksDataSource_basic -v
func TestAccTencentCloudCssWatermarksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermarksDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_watermarks.watermarks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.height"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.picture_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.watermark_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.watermark_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.width"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.x_position"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_watermarks.watermarks", "watermark_list.0.y_position"),
				),
			},
		},
	})
}

const testAccCssWatermarksDataSource = `

data "tencentcloud_css_watermarks" "watermarks" {
}

`
