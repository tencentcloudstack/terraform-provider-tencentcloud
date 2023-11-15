package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDatasourceResource_basic -v
func TestAccTencentCloudNeedFixWedataDatasourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "category", "DB"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "type", "MYSQL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_name", "[zwcs]"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_ident", "体验项目"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "description", "description."),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "display", "tf_example_demo"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "cos_bucket", "wedata-agent-sh-1257305158"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "cos_region", "ap-shanghai"),
				),
			},
			{
				Config: testAccWedataDatasourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "category", "DB"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "type", "MYSQL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_name", "[zwcs]"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "owner_project_ident", "体验项目"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "description", "description update."),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "display", "tf_example_demo_update"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "cos_bucket", "wedata-agent-sh-1257305158"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_datasource.example", "cos_region", "ap-shanghai"),
				),
			},
		},
	})
}

const testAccWedataDatasource = `
resource "tencentcloud_wedata_datasource" "example" {
  name                = "tf_example"
  category            = "DB"
  type                = "MYSQL"
  owner_project_id    = "1612982498218618880"
  owner_project_name  = "[zwcs]"
  owner_project_ident = "体验项目"
  description         = "description."
  display             = "tf_example_demo"
  status              = 1
  cos_bucket          = "wedata-agent-sh-1257305158"
  cos_region          = "ap-shanghai"
  params              = jsonencode({
    "connectType" : "public",
    "authorityType" : "true",
    "deployType" : "CONNSTR_PUBLICDB",
    "url" : "jdbc:mysql://1.1.1.1:8080/database",
    "username" : "root",
    "password" : "password",
    "type" : "MYSQL"
  })
}
`

const testAccWedataDatasourceUpdate = `
resource "tencentcloud_wedata_datasource" "example" {
  name                = "tf_example"
  category            = "DB"
  type                = "MYSQL"
  owner_project_id    = "1612982498218618880"
  owner_project_name  = "[zwcs]"
  owner_project_ident = "体验项目"
  description         = "description update."
  display             = "tf_example_demo_update"
  status              = 1
  cos_bucket          = "wedata-agent-sh-1257305158"
  cos_region          = "ap-shanghai"
  params              = jsonencode({
    "connectType" : "public",
    "authorityType" : "true",
    "deployType" : "CONNSTR_PUBLICDB",
    "url" : "jdbc:mysql://2.2.2.2:8080/database",
    "username" : "admin",
    "password" : "new_password",
    "type" : "MYSQL"
  })
}
`
