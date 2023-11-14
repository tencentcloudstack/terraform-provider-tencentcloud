package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbProjectSecurityGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_project_security_groups.project_security_groups")),
			},
		},
	})
}

const testAccDcdbProjectSecurityGroupsDataSource = `

data "tencentcloud_dcdb_project_security_groups" "project_security_groups" {
  product = ""
  project_id = 
  }

`
