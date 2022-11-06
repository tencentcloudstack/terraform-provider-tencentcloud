package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbParametersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbParameters_basic, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_parameters.parameters"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_parameters.parameters", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_parameters.parameters", "list.0.param"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_parameters.parameters", "list.0.value"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbParameters_basic = `

data "tencentcloud_dcdb_parameters" "parameters" {
  instance_id = "%s"
  }

`
