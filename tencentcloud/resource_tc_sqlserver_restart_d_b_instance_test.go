package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverRestartDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRestartDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_restart_d_b_instance.restart_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_restart_d_b_instance.restart_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRestartDBInstance = `

resource "tencentcloud_sqlserver_restart_d_b_instance" "restart_d_b_instance" {
  instance_id = "mssql-i1z41iwd"
}

`
