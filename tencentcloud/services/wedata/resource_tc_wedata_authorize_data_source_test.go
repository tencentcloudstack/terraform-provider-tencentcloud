package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataAuthorizeDataSourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataAuthorizeDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_authorize_data_source.example", "id"),
				),
			},
			{
				Config: testAccWedataAuthorizeDataSourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_authorize_data_source.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_authorize_data_source.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataAuthorizeDataSource = `
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_project_ids = [
    "1857740139240632320",
  ]

  auth_users = [
    "1857740139240632320_100028448903",
    "1857740139240632320_100028578751",
    "3108707295180644352_100028448903",
    "3108707295180644352_100032159948",
    "3108707295180644352_100044349576",
  ]
}
`

const testAccWedataAuthorizeDataSourceUpdate = `
resource "tencentcloud_wedata_authorize_data_source" "example" {
  data_source_id = "116203"
  auth_users = [
    "1857740139240632320_100028448903",
    "1857740139240632320_100028578751",
    "3108707295180644352_100028448903",
    "3108707295180644352_100032159948",
    "3108707295180644352_100044349576",
  ]
}
`
