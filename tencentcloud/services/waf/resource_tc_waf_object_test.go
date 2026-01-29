package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafObjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWafObject,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_object.waf_object", "id")),
		}, {
			ResourceName:      "tencentcloud_waf_object.waf_object",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWafObject = `

resource "tencentcloud_waf_object" "waf_object" {
}
`
