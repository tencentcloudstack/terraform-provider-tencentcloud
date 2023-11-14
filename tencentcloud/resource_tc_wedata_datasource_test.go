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
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "category"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "owner_project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "owner_project_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "owner_project_ident"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "biz_params"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "params"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "display"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "database_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "collect"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "cos_bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.example", "cos_region"),
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
