package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverInstanceSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceSsl,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_ssl.instance_ssl", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_instance_ssl.instance_ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverInstanceSsl = `

resource "tencentcloud_sqlserver_instance_ssl" "instance_ssl" {
  instance_id = "mssql-i1z41iwd"
  type = "enable"
  wait_switch = 0
}

`
