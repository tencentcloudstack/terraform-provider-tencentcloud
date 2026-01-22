package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataAddCalcEnginesToProjectOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataAddCalcEnginesToProjectOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_add_calc_engines_to_project_operation.example", "id"),
				),
			},
		},
	})
}

const testAccWedataAddCalcEnginesToProjectOperation = `
resource "tencentcloud_wedata_add_calc_engines_to_project_operation" "example" {
  project_id = "20241107221758402"
  dlc_info {
    compute_resources = [
      "dlc_linau6d4bu8bd5u52ffu52a8"
    ]
    region           = "ap-guangzhou"
    default_database = "default_db"
  }
}
`
