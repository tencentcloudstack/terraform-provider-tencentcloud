package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseAccountPermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseAccountPermission,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_account_permission.account_permission", "id")),
			},
			{
				ResourceName:      "tencentcloud_clickhouse_account_permission.account_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClickhouseAccountPermission = tcacctest.DefaultClickhouseVariables + `

resource "tencentcloud_clickhouse_account_permission" "account_permission" {
	instance_id = var.instance_id
	cluster = "default_cluster"
	user_name = "keep_test"
	all_database = true
	global_privileges = ["SELECT", "ALTER"]
}
`
