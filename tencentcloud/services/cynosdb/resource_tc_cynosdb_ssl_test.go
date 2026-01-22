package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbSsl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_ssl.cynosdb_ssl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_ssl.cynosdb_ssl", "status", "ON"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_ssl.cynosdb_ssl", "download_url"),
				),
			},
			{
				Config: testAccCynosdbSsl_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_ssl.cynosdb_ssl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_ssl.cynosdb_ssl", "status", "OFF"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_ssl.cynosdb_ssl", "download_url"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_ssl.cynosdb_ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbSsl = `
resource "tencentcloud_cynosdb_ssl" "cynosdb_ssl" {
  cluster_id = "cynosdbmysql-7yr4dde5"
  instance_id = "cynosdbmysql-ins-4f62d5tq"
  status = "ON"
}
`

const testAccCynosdbSsl_update = `
resource "tencentcloud_cynosdb_ssl" "cynosdb_ssl" {
  cluster_id = "cynosdbmysql-7yr4dde5"
  instance_id = "cynosdbmysql-ins-4f62d5tq"
  status = "OFF"
}
`
