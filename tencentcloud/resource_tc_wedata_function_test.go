package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataFunctionResource_basic -v
func TestAccTencentCloudNeedFixWedataFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataFunction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "kind"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "cluster_identifier"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "db_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "class_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "usage"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "param_desc"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "return_desc"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "example"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "comment"),
				),
			},
			{
				Config: testAccWedataFunctionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "kind"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "cluster_identifier"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "db_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "class_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "usage"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "param_desc"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "return_desc"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "example"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_function.example", "comment"),
				),
			},
		},
	})
}

const testAccWedataFunction = `
resource "tencentcloud_wedata_function" "example" {
  type               = "HIVE"
  kind               = "ANALYSIS"
  name               = "tf_example"
  cluster_identifier = "emr-m6u3qgk0"
  db_name            = "tf_db_example"
  project_id         = "1612982498218618880"
  class_name         = "tf_class_example"
  resource_list {
    path = "/wedata-demo-1314991481/untitled3-1.0-SNAPSHOT.jar"
    name = "untitled3-1.0-SNAPSHOT.jar"
    id   = "5b28bcdf-a0e6-4022-927d-927d399c4593"
    type = "cos"
  }
  description = "description."
  usage       = "usage info."
  param_desc  = "param info."
  return_desc = "return value info."
  example     = "example info."
  comment     = "V1"
}
`

const testAccWedataFunctionUpdate = `
resource "tencentcloud_wedata_function" "example" {
  type               = "HIVE"
  kind               = "ENCRYPTION"
  name               = "tf_example"
  cluster_identifier = "emr-m6u3qgk0"
  db_name            = "tf_db_example"
  project_id         = "1612982498218618880"
  class_name         = "tf_class_example_update"
  resource_list {
    path = "/wedata-demo-1314991481/untitled3-1.1-SNAPSHOT.jar"
    name = "untitled3-1.1-SNAPSHOT.jar"
    id   = "5b28bcdf-a0e6-4022-927d-927d399c4594"
    type = "cos"
  }
  description = "description update."
  usage       = "usage info update."
  param_desc  = "param info update."
  return_desc = "return value info update."
  example     = "example info update."
  comment     = "V2"
}
`
