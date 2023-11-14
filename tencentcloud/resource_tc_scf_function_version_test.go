package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_function_version.function_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_function_version.function_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfFunctionVersion = `

resource "tencentcloud_scf_function_version" "function_version" {
  function_name = "test_function"
  description = "test function"
  namespace = "test_namespace"
}

`
