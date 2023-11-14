package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigTerminateDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigTerminateDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_terminate_d_b_instance.config_terminate_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_terminate_d_b_instance.config_terminate_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigTerminateDBInstance = `

resource "tencentcloud_sqlserver_config_terminate_d_b_instance" "config_terminate_d_b_instance" {
  instance_id_set = 
}

`
