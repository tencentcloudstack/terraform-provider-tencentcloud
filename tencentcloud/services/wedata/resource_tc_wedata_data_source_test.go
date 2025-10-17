package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataDataSourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_data_source.example", "id"),
				),
			},
			{
				Config: testAccWedataDataSourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_data_source.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_data_source.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataDataSource = `
resource "tencentcloud_wedata_data_source" "example" {
  project_id          = "2982667120655491072"
  name                = "tf_example"
  type                = "MYSQL"
  prod_con_properties = ""
  dev_con_properties  = ""
  display_name        = "display_name"
  description         = "description"
}
`

const testAccWedataDataSourceUpdate = `
resource "tencentcloud_wedata_data_source" "example" {
  project_id          = "2982667120655491072"
  name                = "tf_example"
  type                = "MYSQL"
  prod_con_properties = ""
  dev_con_properties  = ""
  display_name        = "display_name"
  description         = "description"
}
`
