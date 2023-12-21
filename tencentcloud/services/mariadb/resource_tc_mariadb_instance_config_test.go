package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbInstanceConfigResource_basic -v
func TestAccTencentCloudMariadbInstanceConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_instance_config.test", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_instance_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbInstanceConfig = `
resource "tencentcloud_mariadb_instance_config" "test" {
  instance_id        = "tdsql-9vqvls95"
  rs_access_strategy = 1
  extranet_access    = 0
}
`
