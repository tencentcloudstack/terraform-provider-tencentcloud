package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataDatasourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDatasource,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_datasource.datasource", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_datasource.datasource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataDatasource = `

resource "tencentcloud_wedata_datasource" "datasource" {
  name = "name"
  category = "DB"
  type = "MYSQL"
  owner_project_id = "110111121"
  owner_project_name = "ownerprojectname"
  owner_project_ident = "OwnerProjectIdent"
  biz_params = "{}"
  params = "{}"
  description = "descr"
  display = "Display"
  database_name = "db"
  instance = "instance"
  status = 1
  cluster_id = "cid"
  collect = "false"
  c_o_s_bucket = "aaaa"
  c_o_s_region = "ap-beijing"
}

`
