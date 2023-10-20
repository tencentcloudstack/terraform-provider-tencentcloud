package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudBiDatasourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiDatasource,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_datasource.datasource", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_datasource.datasource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiDatasource = `

resource "tencentcloud_bi_datasource" "datasource" {
  db_host = "1.2.3.4"
  db_port = 3306
  service_type = "Own"
  db_type = "Database type."
  charset = "utf8"
  db_user = "root"
  db_pwd = "abc"
  db_name = "abc"
  source_name = "abc"
  project_id = 123
  catalog = "presto"
  data_origin = "abc"
  data_origin_project_id = "abc"
  data_origin_datasource_id = "abc"
  extra_param = ""
  uniq_vpc_id = ""
  vip = ""
  vport = ""
  vpc_id = ""
}

`
