package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorNoticeContentTmplsDataSource_basic -v
func TestAccTencentCloudMonitorNoticeContentTmplsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorNoticeContentTmplsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_notice_content_tmpls.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_notice_content_tmpls.example", "notice_content_tmpl_list.#"),
				),
			},
		},
	})
}

const testAccMonitorNoticeContentTmplsDataSource = `
data "tencentcloud_monitor_notice_content_tmpls" "example" {}
`
