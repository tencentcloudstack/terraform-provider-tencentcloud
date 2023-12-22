package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbRenewInstanceResource_basic -v
func TestAccTencentCloudMariadbRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		CheckDestroy: testAccCheckMariadbInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbRenewInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_renew_instance.renew_instance", "id"),
				),
			},
		},
	})
}

const testAccMariadbRenewInstance = testAccMariadbInstance + `
resource "tencentcloud_mariadb_renew_instance" "renew_instance" {
  instance_id = tencentcloud_mariadb_instance.instance.id
  period      = 1
}
`
