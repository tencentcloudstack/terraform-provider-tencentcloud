package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbProjectSecurityGroupsDataSource_basic -v
func TestAccTencentCloudMariadbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbProjectSecurityGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_project_security_groups.project_security_groups"),
				),
			},
		},
	})
}

const testAccMariadbProjectSecurityGroupsDataSource = `
data "tencentcloud_mariadb_project_security_groups" "project_security_groups" {
  product    = "mariadb"
  project_id = 0
}
`
