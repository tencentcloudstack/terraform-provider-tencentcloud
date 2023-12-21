package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbSecurityGroups_basic -v
func TestAccTencentCloudMariadbSecurityGroups_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_security_groups.security_groups", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_security_groups.security_groups",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbSecurityGroupsVar = `
variable "security_group_id" {
  default = "` + tcacctest.DefaultMariadbSecurityGroupId + `"
}
`

const testAccMariadbSecurityGroups = testAccMariadbHourDbInstance + testAccMariadbSecurityGroupsVar + `

resource "tencentcloud_mariadb_security_groups" "security_groups" {
	instance_id       = tencentcloud_mariadb_hour_db_instance.basic.id
	security_group_id = var.security_group_id
	product           = "mariadb"
}

`
