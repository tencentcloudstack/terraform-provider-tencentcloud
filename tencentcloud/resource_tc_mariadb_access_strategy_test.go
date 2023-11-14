package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbAccessStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbAccessStrategy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_access_strategy.access_strategy", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_access_strategy.access_strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbAccessStrategy = `

resource "tencentcloud_mariadb_access_strategy" "access_strategy" {
  instance_id = ""
  rs_access_strategy = 
}

`
