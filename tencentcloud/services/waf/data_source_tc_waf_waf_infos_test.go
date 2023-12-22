package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafWafInfosDataSource_basic -v
func TestAccTencentCloudNeedFixWafWafInfosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWafInfosDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_waf_infos.example"),
				),
			},
		},
	})
}

const testAccWafWafInfosDataSource = `
data "tencentcloud_waf_waf_infos" "example" {
  params {
    load_balancer_id = "lb-A8VF445"
    listener_id      = "lbl-nonkgvc2"
    domain_id        = "waf-MPtWPK5Q"
  }
}
`
