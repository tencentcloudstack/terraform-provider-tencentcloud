package teo_test

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
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionRuntimeEnvironment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.key", "test-a"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.type", "string"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.value", "AAA"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.key", "test-b"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.type", "string"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.value", "BBB"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionRuntimeEnvironmentUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.key", "test-b"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.type", "string"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.0.value", "BBB"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.key", "test-a"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.type", "string"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment", "environment_variables.1.value", "AAA"),
				),
			},
		},
	})
}

const testAccTeoFunctionRuntimeEnvironment = `

resource "tencentcloud_teo_function_runtime_environment" "teo_function_runtime_environment" {
    function_id = "ef-txx7fnua"
    zone_id     = "zone-2qtuhspy7cr6"

    environment_variables {
        key   = "test-a"
        type  = "string"
        value = "AAA"
    }
    environment_variables {
        key   = "test-b"
        type  = "string"
        value = "BBB"
    }
}
`
const testAccTeoFunctionRuntimeEnvironmentUp = `

resource "tencentcloud_teo_function_runtime_environment" "teo_function_runtime_environment" {
    function_id = "ef-txx7fnua"
    zone_id     = "zone-2qtuhspy7cr6"

    environment_variables {
        key   = "test-b"
        type  = "string"
        value = "BBB"
    }
    environment_variables {
        key   = "test-a"
        type  = "string"
        value = "AAA"
    }
}
`
