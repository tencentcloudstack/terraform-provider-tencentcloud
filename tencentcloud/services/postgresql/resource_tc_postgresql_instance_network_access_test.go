package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlInstanceNetworkAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceNetworkAccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.example", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.example", "vip"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_instance_network_access.example",
				ImportState:       true,
				ImportStateVerify: true,
			}},
	})
}

const testAccPostgresqlInstanceNetworkAccess = `
resource "tencentcloud_postgresql_instance_network_access" "example" {
  db_instance_id = "postgres-ai46555b"
  vpc_id         = "vpc-i5yyodl9"
  subnet_id      = "subnet-d4umunpy"
  vip            = "10.0.10.11"
}
`
