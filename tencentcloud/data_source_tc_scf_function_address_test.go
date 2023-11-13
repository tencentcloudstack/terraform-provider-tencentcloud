package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionAddressDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_function_address.function_address")),
			},
		},
	})
}

const testAccScfFunctionAddressDataSource = `

data "tencentcloud_scf_function_address" "function_address" {
  function_name = ""
  qualifier = ""
  namespace = ""
    }

`
