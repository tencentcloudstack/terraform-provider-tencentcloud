package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbSecurityGroupsDataSource -v
func TestAccTencentCloudMariadbSecurityGroupsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMariadbSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_security_groups.security_groups"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_security_groups.security_groups", "list.#"),
				),
			},
		},
	})
}

const testAccDataSourceMariadbSecurityGroups = testAccMariadbSecurityGroups + `

data "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  product = "mariadb"
}

`
