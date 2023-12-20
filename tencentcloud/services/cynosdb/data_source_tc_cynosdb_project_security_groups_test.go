package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbProjectSecurityGroupsDataSource_basic -v
func TestAccTencentCloudCynosdbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProjectSecurityGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_project_security_groups.project_security_groups"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.inbound.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.outbound.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.security_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.security_group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_project_security_groups.project_security_groups", "groups.0.security_group_remark"),
				),
			},
		},
	})
}

const testAccCynosdbProjectSecurityGroupsDataSource = `
data "tencentcloud_cynosdb_project_security_groups" "project_security_groups" {
  project_id = 1250480
  search_key = "自定义模版"
}
`
