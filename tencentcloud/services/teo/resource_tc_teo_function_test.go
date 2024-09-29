package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoFunction,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_function.teo_function",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoFunction = `

resource "tencentcloud_teo_function" "teo_function" {
}
`
