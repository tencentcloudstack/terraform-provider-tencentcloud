package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbProjectSecurityGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbProjectSecurityGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_project_security_group.project_security_group")),
			},
		},
	})
}

const testAccCdbProjectSecurityGroupDataSource = `

data "tencentcloud_cdb_project_security_group" "project_security_group" {
  project_id = 
  }

`
