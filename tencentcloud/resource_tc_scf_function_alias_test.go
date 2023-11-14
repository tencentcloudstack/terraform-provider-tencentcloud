package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionAliasResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionAlias,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_function_alias.function_alias", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_function_alias.function_alias",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfFunctionAlias = `

resource "tencentcloud_scf_function_alias" "function_alias" {
  name = "test_func_alais"
  function_name = "test_function"
  function_version = "$LATEST"
  namespace = "test_namespace"
  routing_config {
		additional_version_weights {
			version = "1"
			weight = 
		}
		addtion_version_matchs {
			version = "1"
			key = "invoke.headers.User"
			method = "range"
			expression = "[1,2]"
		}

  }
  description = "test routing"
}

`
