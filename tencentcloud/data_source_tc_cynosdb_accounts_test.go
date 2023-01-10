package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_accounts.accounts"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_accounts.accounts", "account_set.#", "1"),
				),
			},
		},
	})
}

const testAccCynosdbAccountsDataSource = CommonCynosdb + `

data "tencentcloud_cynosdb_accounts" "accounts" {
	cluster_id = var.cynosdb_cluster_id
	account_names = ["root"]
}

`
