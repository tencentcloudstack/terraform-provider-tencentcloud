package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbSecurityGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbSecurityGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_security_group.security_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_security_group.security_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbSecurityGroup = `

resource "tencentcloud_cynosdb_security_group" "security_group" {
  instance_ids = &lt;nil&gt;
  security_group_ids = &lt;nil&gt;
  zone = &lt;nil&gt;
}

`
