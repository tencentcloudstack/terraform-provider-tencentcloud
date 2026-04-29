package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoLoadBalancerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoLoadBalancer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_load_balancer.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_load_balancer.example", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_load_balancer.example", "name", "tf-example-lb"),
					resource.TestCheckResourceAttr("tencentcloud_teo_load_balancer.example", "type", "GENERAL"),
					resource.TestCheckResourceAttr("tencentcloud_teo_load_balancer.example", "steering_policy", "Pritory"),
					resource.TestCheckResourceAttr("tencentcloud_teo_load_balancer.example", "failover_policy", "OtherRecordInOriginGroup"),
				),
			},
			{
				Config: testAccTeoLoadBalancerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_load_balancer.example", "name", "tf-example-lb-updated"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_load_balancer.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoLoadBalancerBase = `
variable "zone_id" {
  default = "zone-2ju9lrnpaxol"
}

variable "origin_group_id" {
  default = "og-30l5kv5z2bse"
}
`

const testAccTeoLoadBalancer = testAccTeoLoadBalancerBase + `
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = var.zone_id
  name    = "tf-example-lb"
  type    = "GENERAL"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = var.origin_group_id
  }

  health_checker {
    type               = "HTTP"
    port               = 80
    interval           = 60
    timeout            = 5
    health_threshold   = 3
    critical_threshold = 2
    path               = "example.com/health"
    method             = "HEAD"
    expected_codes     = ["200"]
    follow_redirect    = "false"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherRecordInOriginGroup"
}
`

const testAccTeoLoadBalancerUpdate = testAccTeoLoadBalancerBase + `
resource "tencentcloud_teo_load_balancer" "example" {
  zone_id = var.zone_id
  name    = "tf-example-lb-updated"
  type    = "GENERAL"

  origin_groups {
    priority        = "priority_1"
    origin_group_id = var.origin_group_id
  }

  health_checker {
    type               = "HTTP"
    port               = 80
    interval           = 60
    timeout            = 5
    health_threshold   = 3
    critical_threshold = 2
    path               = "example.com/health"
    method             = "HEAD"
    expected_codes     = ["200", "301"]
    follow_redirect    = "true"
  }

  steering_policy = "Pritory"
  failover_policy = "OtherOriginGroup"
}
`
