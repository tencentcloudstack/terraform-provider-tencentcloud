package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionRuntimeEnvironmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoFunctionRuntimeEnvironment,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoFunctionRuntimeEnvironment = `

resource "tencentcloud_teo_function_runtime_environment" "teo_function_runtime_environment" {
  environment_variables = {
  }
}
`
