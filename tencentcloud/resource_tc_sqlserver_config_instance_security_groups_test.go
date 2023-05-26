package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceSecurityGroupsResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceSecurityGroupsResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_security_groups.config_instance_security_groups", "id"),
				),
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
  instance_id = "mssql-qelbzgwf"
  security_group_id_set = ["sg-mayqdlt1", "sg-5aubsf8n"]
}
`
