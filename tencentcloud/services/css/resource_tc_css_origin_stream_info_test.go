package css_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCssOriginStreamInfoResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssOriginStreamInfo,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_origin_stream_info.example", "id"),
				),
			},
			{
				Config: testAccCssOriginStreamInfoUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_origin_stream_info.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_origin_stream_info.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssOriginStreamInfo = `
resource "tencentcloud_css_origin_stream_info" "example" {
  domain_name             = "arunma.cn"
  origin_stream_play_type = "rtmp"
  cdn_stream_play_type    = ["rtmp"]
  origin_stream_type      = 1
  origin_address_type     = 1
  origin_address          = ["1.1.1.1:8080"]
  origin_timeout          = 10000
  origin_retry_times      = 10
}
`

const testAccCssOriginStreamInfoUpdate = `
resource "tencentcloud_css_origin_stream_info" "example" {
  domain_name             = "arunma.cn"
  origin_stream_play_type = "rtmp"
  cdn_stream_play_type    = ["rtmp"]
  origin_stream_type      = 1
  origin_address_type     = 1
  origin_address          = ["1.1.1.1:8090"]
  origin_timeout          = 10000
  origin_retry_times      = 10
}
`
