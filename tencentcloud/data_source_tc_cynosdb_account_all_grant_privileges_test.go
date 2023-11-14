package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbAccountAllGrantPrivilegesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccountAllGrantPrivilegesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_account_all_grant_privileges.account_all_grant_privileges")),
			},
		},
	})
}

const testAccCynosdbAccountAllGrantPrivilegesDataSource = `

data "tencentcloud_cynosdb_account_all_grant_privileges" "account_all_grant_privileges" {
  cluster_id = "xxx"
  account {
		account_name = ""
		host = ""

  }
        }

`
