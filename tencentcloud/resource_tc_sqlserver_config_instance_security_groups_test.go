package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigInstanceSecurityGroupsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceSecurityGroups,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceSecurityGroups = `

resource "tencentcloud_sqlserver_config_instance_security_groups" "config_instance_security_groups" {
  instance_id = "mssql-i1z41iwd"
  security_group_id_set = 
}

`
