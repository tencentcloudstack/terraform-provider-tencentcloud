package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvokerConfigResource_basic -v
func TestAccTencentCloudTatInvokerConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvokerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tat_invoker_config.invoker_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invoker_config.invoker_config", "invoker_status", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_tat_invoker_config.invoker_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatInvokerConfig = testAccTatInvoker + `

resource "tencentcloud_tat_invoker_config" "invoker_config" {
	invoker_id = tencentcloud_tat_invoker.invoker.id
	invoker_status = "off"
}

`
