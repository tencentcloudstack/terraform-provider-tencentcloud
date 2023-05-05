package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccScfFunctionAliasUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function_alias.function_alias", "id"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function_alias.function_alias", "description", "weight test first"),
				),
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
  description      = "weight test"
  function_name    = "keep-1676351130"
  function_version = "$LATEST"
  name             = "weight"
  namespace        = "default"

  routing_config {
    additional_version_weights {
      version = "2"
      weight  = 0.4
    }
  }
}

`

const testAccScfFunctionAliasUpdate = `

resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "weight test first"
  function_name    = "keep-1676351130"
  function_version = "$LATEST"
  name             = "weight"
  namespace        = "default"

  routing_config {
    additional_version_weights {
      version = "2"
      weight  = 0.2
    }
  }
}

`
