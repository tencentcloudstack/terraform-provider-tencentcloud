package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudNeedFixMariadbCloseDBExtranetAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCloseDBExtranetAccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_close_db_extranet_access.close_db_extranet_access", "id"),
				),
			},
		},
	})
}

const testAccMariadbCloseDBExtranetAccess = `
resource "tencentcloud_mariadb_close_db_extranet_access" "close_db_extranet_access" {
  instance_id = "tdsql-9vqvls95"
}
`
