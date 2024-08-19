package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSsl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ssl.ssl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_ssl.ssl", "status", "ON"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ssl.ssl", "url"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_ssl.ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMysqlSslUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ssl.ssl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_ssl.ssl", "status", "OFF"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ssl.ssl", "url"),
				),
			},
		},
	})
}

const testAccMysqlSsl = testAccMysql + `

resource "tencentcloud_mysql_ssl" "ssl" {
  instance_id = tencentcloud_mysql_instance.mysql.id
  status      = "ON"
}

`

const testAccMysqlSslUp = testAccMysql + `

resource "tencentcloud_mysql_ssl" "ssl" {
  instance_id = tencentcloud_mysql_instance.mysql.id
  status      = "OFF"
}

`
