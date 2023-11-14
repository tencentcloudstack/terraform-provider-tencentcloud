package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDBInstanceSecurityGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDBInstanceSecurityGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_d_b_instance_security_group.d_b_instance_security_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_d_b_instance_security_group.d_b_instance_security_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresDBInstanceSecurityGroup = `

resource "tencentcloud_postgres_d_b_instance_security_group" "d_b_instance_security_group" {
  security_group_id_set = 
  d_b_instance_id = ""
  read_only_group_id = ""
}

`
