package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbReadOnlyInstanceExclusiveAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReadOnlyInstanceExclusiveAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_read_only_instance_exclusive_access.read_only_instance_exclusive_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbReadOnlyInstanceExclusiveAccess = `

resource "tencentcloud_cynosdb_read_only_instance_exclusive_access" "read_only_instance_exclusive_access" {
  cluster_id = "cynosdbmysql-12345678"
  instance_id = "cynosdbmysql-ins-12345678"
  vpc_id = "vpc-12345678"
  subnet_id = "subnet-12345678"
  port = 1234
  security_group_ids = 
}

`
