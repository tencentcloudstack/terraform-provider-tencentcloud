package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsLoadBalancerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsLoadBalancer,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_load_balancer.load_balancer", "id")),
			},
			{
				ResourceName:      "tencentcloud_as_load_balancer.load_balancer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAsLoadBalancer = `

resource "tencentcloud_as_load_balancer" "load_balancer" {
  auto_scaling_group_id = "asg-12wjuh0s"
  forward_load_balancers {
    load_balancer_id = "lb-d8u76te5"
    listener_id      = "lbl-s8dh4y75"
    target_attributes {
      port   = 8080
      weight = 20
    }
    location_id = "loc-fsa87u6d"
    region      = "ap-guangzhou"
  }
}

`
