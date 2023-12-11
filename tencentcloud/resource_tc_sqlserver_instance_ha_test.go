package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverInstanceHaResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceHa,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_ha.instance_ha", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_instance_ha.instance_ha",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverInstanceHa = `

resource "tencentcloud_sqlserver_instance_ha" "instance_ha" {
  instance_id = "mssql-i1z41iwd"
  wait_switch = 0
}

`
