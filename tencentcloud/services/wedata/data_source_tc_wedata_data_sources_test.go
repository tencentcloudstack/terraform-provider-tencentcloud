package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataDataSourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataDataSourcesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_sources.example"),
			),
		}},
	})
}

const testAccWedataDataSourcesDataSource = `
data "tencentcloud_wedata_data_sources" "example" {
  project_id   = "2982667120655491072"
  name         = "tf_example"
  display_name = "display_name"
  type         = ["MYSQL", "ORACLE"]
  creator      = "user"
}
`
