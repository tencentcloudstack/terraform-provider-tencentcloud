package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoLoadBalancingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoLoadBalancing,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_load_balancing.load_balancing", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_load_balancing.load_balancing",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoLoadBalancing = `

resource "tencentcloud_teo_load_balancing" "load_balancing" {
  zone_id = &lt;nil&gt;
    host = &lt;nil&gt;
  type = &lt;nil&gt;
  origin_group_id = &lt;nil&gt;
  backup_origin_group_id = &lt;nil&gt;
  t_t_l = &lt;nil&gt;
  status = &lt;nil&gt;
    }

`
