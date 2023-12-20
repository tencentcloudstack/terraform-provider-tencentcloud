package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssTimeShiftStreamListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssTimeShiftStreamListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_css_time_shift_stream_list.time_shift_stream_list"),
				),
			},
		},
	})
}

const testAccCssTimeShiftStreamListDataSource = `

data "tencentcloud_css_time_shift_stream_list" "time_shift_stream_list" {
  start_time   = 1698768000
  end_time     = 1698820641
  stream_name  = "live"
  domain       = "177154.push.tlivecloud.com"
  domain_group = "tf-test"
}

`
