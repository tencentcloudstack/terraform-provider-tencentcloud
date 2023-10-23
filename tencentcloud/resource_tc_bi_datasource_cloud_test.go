package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixBiDatasourceCloudResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiDatasourceCloud,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_datasource_cloud.datasource_cloud", "id"),
				),
			},
		},
	})
}

const testAccBiDatasourceCloud = `

resource "tencentcloud_bi_datasource_cloud" "datasource_cloud" {
  charset    = "utf8"
  db_name    = "bi_dev"
  db_type    = "MYSQL"
  db_user    = "root"
  project_id = "11015056"
  db_pwd     = "xxxxxx"
  service_type {
    instance_id = "cdb-12viotu5"
    region     = "ap-guangzhou"
    type       = "Cloud"
  }
  source_name = "tf-test1"
  vip         = "10.0.0.4"
  vport       = "3306"
  region_id   = "gz"
  vpc_id      = 5292713
}
`
