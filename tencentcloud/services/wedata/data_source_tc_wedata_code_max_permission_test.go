package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataCodeMaxPermissionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataCodeMaxPermissionDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_code_max_permission.wedata_code_max_permission"),
			),
		}},
	})
}

const testAccWedataCodeMaxPermissionDataSource = `
data "tencentcloud_wedata_code_max_permission" "example" {
  project_id  = "3108707295180644352"
  resource_id = "f0c14b9d-003e-4325-8830-d1a9fa934ed6"
}
`
