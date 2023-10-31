package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixBiDatasourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiDatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bi_datasource.datasource", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "charset", "utf8"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_host", "bj-cdb-1lxqg5r6.sql.tencentcdb.com"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_port", "63694"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_type", "MYSQL"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_pwd", "ABc123,,,"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "db_user", "root"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "project_id", "11015030"),
					resource.TestCheckResourceAttr("tencentcloud_bi_datasource.datasource", "source_name", "tf-source-name"),
				),
			},
			{
				ResourceName:            "tencentcloud_bi_datasource.datasource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_pwd"},
			},
		},
	})
}

const testAccBiDatasource = `

resource "tencentcloud_bi_datasource" "datasource" {
  charset     = "utf8"
  db_host     = "bj-cdb-1lxqg5r6.sql.tencentcdb.com"
  db_name     = "tf-test"
  db_port     = 63694
  db_type     = "MYSQL"
  db_pwd      = "ABc123,,,"
  db_user     = "root"
  project_id  = 11015030
  source_name = "tf-source-name"
}

`
