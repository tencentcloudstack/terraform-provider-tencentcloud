package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfFunctionAliasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionAliasesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_function_aliases.function_aliases")),
			},
		},
	})
}

const testAccScfFunctionAliasesDataSource = `

data "tencentcloud_scf_function_aliases" "function_aliases" {
  function_name = "keep-1676351130"
  namespace     = "default"
}

`
