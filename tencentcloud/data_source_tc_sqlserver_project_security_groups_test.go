package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverProjectSecurityGroupsDataSource_basic -v
func TestAccTencentCloudSqlserverProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverProjectSecurityGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_project_security_groups.example")),
			},
		},
	})
}

const testAccSqlserverProjectSecurityGroupsDataSource = `
data "tencentcloud_sqlserver_project_security_groups" "example" {
  project_id = 0
}
`
