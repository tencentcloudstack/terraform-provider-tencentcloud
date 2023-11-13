package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbSecurityGroupsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSecurityGroups,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_security_groups.security_groups", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_security_groups.security_groups",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbSecurityGroups = `

resource "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_ids = &lt;nil&gt;
  security_group_id = &lt;nil&gt;
  product = &lt;nil&gt;
}

`
