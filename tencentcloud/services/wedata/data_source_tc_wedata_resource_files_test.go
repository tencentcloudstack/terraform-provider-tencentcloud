package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataResourceFilesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataResourceFilesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_resource_files.wedata_resource_files"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_resource_files.wedata_resource_files", "data.#"),
			),
		}},
	})
}

const testAccWedataResourceFilesDataSource = `
data "tencentcloud_wedata_resource_files" "wedata_resource_files" {
  project_id         = 2905622749543821312
  resource_name      = "tftest.txt"
}
`
