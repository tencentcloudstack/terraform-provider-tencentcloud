package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalClbResource_instance(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalClbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "clb_name", "tf-clb"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "tags.test", "tf"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "tags.test1", "tf1"),
				),
			},
			{
				ResourceName:            "tencentcloud_clb_instance.clb_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dynamic_vip"},
			},
		},
	})
}

const testAccInternationalClbInstance_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb"
  tags = {
    test = "tf"
    test1 = "tf1"
  }
}
`
