package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSqlScriptResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataSqlScript,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "access_scope"),
				),
			},
			{
				Config: testAccWedataSqlScriptUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "script_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_script.example", "access_scope"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_sql_script.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataSqlScript = `
resource "tencentcloud_wedata_sql_script" "example" {
  script_name        = "tf_example_script"
  project_id         = "2983848457986924544"
  parent_folder_path = "/tf_example"
  script_config {
    datasource_id    = "108826"
    compute_resource = "svmgao_stability"
  }

  script_content = "SHOW DATABASES;"
  access_scope   = "SHARED"
}
`

const testAccWedataSqlScriptUpdate = `
resource "tencentcloud_wedata_sql_script" "example" {
  script_name        = "tf_example_script"
  project_id         = "2983848457986924544"
  parent_folder_path = "/tf_example"
  script_config {
    datasource_id    = "108826"
    compute_resource = "svmgao_stability"
  }

  script_content = "SHOW DATABASES;"
  access_scope   = "PRIVATE"
}
`
