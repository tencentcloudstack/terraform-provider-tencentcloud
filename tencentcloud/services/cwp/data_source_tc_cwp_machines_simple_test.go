package cwp_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCwpMachinesSimpleDataSource_basic -v
func TestAccTencentCloudNeedFixCwpMachinesSimpleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpMachinesSimpleDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cwp_machines_simple.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cwp_machines_simple.example", "machine_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cwp_machines_simple.example", "machine_region"),
				),
			},
		},
	})
}

const testAccCwpMachinesSimpleDataSource = `
data "tencentcloud_cwp_machines_simple" "example" {
  machine_type   = "CVM"
  machine_region = "ap-guangzhou"
  project_ids    = [0]
}
`
