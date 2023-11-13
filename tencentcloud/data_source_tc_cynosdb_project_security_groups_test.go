package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProjectSecurityGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_project_security_groups.project_security_groups")),
			},
		},
	})
}

const testAccCynosdbProjectSecurityGroupsDataSource = `

data "tencentcloud_cynosdb_project_security_groups" "project_security_groups" {
  project_id = 11954
  search_key = ""
  }

`
