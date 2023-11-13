package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbProjectSecurityGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_project_security_groups.project_security_groups")),
			},
		},
	})
}

const testAccMariadbProjectSecurityGroupsDataSource = `

data "tencentcloud_mariadb_project_security_groups" "project_security_groups" {
  product = ""
  project_id = 
  }

`
