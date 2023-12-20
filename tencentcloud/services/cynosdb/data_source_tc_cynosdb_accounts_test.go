package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbAccountsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAccountsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_accounts.accounts"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_accounts.accounts", "account_set.#", "1"),
				),
			},
		},
	})
}

const testAccCynosdbAccountsDataSource = tcacctest.CommonCynosdb + `

data "tencentcloud_cynosdb_accounts" "accounts" {
	cluster_id = var.cynosdb_cluster_id
	account_names = ["root"]
}

`
