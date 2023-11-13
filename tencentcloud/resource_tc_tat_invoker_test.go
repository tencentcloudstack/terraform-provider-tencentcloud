package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTatInvokerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvoker,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tat_invoker.invoker", "id")),
			},
			{
				ResourceName:      "tencentcloud_tat_invoker.invoker",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatInvoker = `

resource "tencentcloud_tat_invoker" "invoker" {
  invoker_id = ""
}

`
