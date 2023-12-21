package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbRestartInstanceResource_basic -v
func TestAccTencentCloudMariadbRestartInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbRestartInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_restart_instance.restart_instance", "id"),
				),
			},
		},
	})
}

const testAccMariadbRestartInstance = `
resource "tencentcloud_mariadb_restart_instance" "restart_instance" {
  instance_id = "tdsql-9vqvls95"
}
`
