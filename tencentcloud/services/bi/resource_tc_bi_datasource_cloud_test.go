package bi_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixBiDatasourceCloudResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  project_id = "11015030"
  db_pwd     = "zxcvb12345"
  service_type {
    instance_id = "cdb-1ub45mjx"
    region     = "ap-guangzhou"
    type       = "Cloud"
  }
  source_name = "bi_test"
  vip         = "172.16.64.9"
  vport       = "3306"
  region_id   = "gz"
  vpc_id      = 5232945
}
`
