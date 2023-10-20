package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudBiDatasourceCloudResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiDatasourceCloud,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_bi_datasource_cloud.datasource_cloud", "id")),
			},
			{
				ResourceName:      "tencentcloud_bi_datasource_cloud.datasource_cloud",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBiDatasourceCloud = `

resource "tencentcloud_bi_datasource_cloud" "datasource_cloud" {
  service_type = "Cloud"
  db_type = "Database type."
  charset = "utf8"
  db_user = "root"
  db_pwd = "abc"
  db_name = "abc"
  source_name = "abc"
  project_id = "123"
  vip = "1.2.3.4"
  vport = "3306"
  vpc_id = ""
  uniq_vpc_id = ""
  region_id = ""
  extra_param = ""
  instance_id = ""
  prod_db_name = ""
  data_origin = "abc"
  data_origin_project_id = "abc"
  data_origin_datasource_id = "abc"
}

`
