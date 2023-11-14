package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbAccountsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccountsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_accounts.accounts")),
			},
		},
	})
}

const testAccCynosdbAccountsDataSource = `

data "tencentcloud_cynosdb_accounts" "accounts" {
  cluster_id = "cynosdbmysql-on5xw0ni"
  account_names = 
  db_type = "MYSQL"
  hosts = 
  }

`
