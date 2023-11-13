package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbModifyRealServerAccessStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbModifyRealServerAccessStrategy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_modify_real_server_access_strategy.modify_real_server_access_strategy", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_modify_real_server_access_strategy.modify_real_server_access_strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbModifyRealServerAccessStrategy = `

resource "tencentcloud_dcdb_modify_real_server_access_strategy" "modify_real_server_access_strategy" {
  instance_id = ""
  rs_access_strategy = 
}

`
