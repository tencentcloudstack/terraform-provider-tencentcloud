package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBParametersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbParameters_basic, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_parameters.parameters"),
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
