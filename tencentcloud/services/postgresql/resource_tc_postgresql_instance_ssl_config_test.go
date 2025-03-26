package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlInstanceSslConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccPostgresqlInstanceSslConfig,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_ssl_config.postgresql_instance_ssl_config", "id")),
		}, {
			ResourceName:      "tencentcloud_postgresql_instance_ssl_config.postgresql_instance_ssl_config",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccPostgresqlInstanceSslConfig = `

resource "tencentcloud_postgresql_instance_ssl_config" "postgresql_instance_ssl_config" {
}
`
