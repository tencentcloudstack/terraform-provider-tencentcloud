package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionVersionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_function_versions.function_versions")),
			},
		},
	})
}

const testAccScfFunctionVersionsDataSource = `

data "tencentcloud_scf_function_versions" "function_versions" {
  function_name = ""
  namespace = ""
  order = ""
  order_by = ""
    }

`
