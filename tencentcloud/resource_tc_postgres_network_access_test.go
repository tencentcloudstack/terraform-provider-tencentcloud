package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresNetworkAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresNetworkAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_network_access.network_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_network_access.network_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresNetworkAccess = `

resource "tencentcloud_postgres_network_access" "network_access" {
  read_only_group_id = "pgro-xxxx"
  vpc_id = "vpc-xxx"
  subnet_id = "subnet-xxx"
  is_assign_vip = false
  vip = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
