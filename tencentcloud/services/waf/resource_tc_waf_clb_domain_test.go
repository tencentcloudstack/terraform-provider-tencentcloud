package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafClbDomainResource_basic -v
func TestAccTencentCloudWafClbDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafClbDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_domain.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_clb_domain.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafClbDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_domain.example", "id"),
				),
			},
		},
	})
}

const testAccWafClbDomain = `
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id = "waf_2kxtlbky00b2v1fn"
  domain      = "xxx.com"
  is_cdn      = 3
  status      = 1
  engine      = 21

  load_balancer_set {
    load_balancer_id   = "lb-5dnrkgry"
    load_balancer_name = "keep-listener-clb"
    listener_id        = "lbl-nonkgvc2"
    listener_name      = "dsadasd"
    vip                = "106.55.220.8"
    vport              = "80"
    region             = "gz"
    protocol           = "HTTP"
    zone               = "ap-guangzhou-6"
    numerical_vpc_id   = "5232945"
    load_balancer_type = "OPEN"
  }

  region          = "gz"
  flow_mode       = 1
  alb_type        = "clb"
  bot_status      = 1
  api_safe_status = 1
  ip_headers      = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]
}
`

const testAccWafClbDomainUpdate = `
resource "tencentcloud_waf_clb_domain" "example" {
  instance_id = "waf_2kxtlbky00b2v1fn"
  domain      = "xxx.com"
  is_cdn      = 0
  status      = 1
  engine      = 12

  load_balancer_set {
    load_balancer_id   = "lb-5dnrkgry"
    load_balancer_name = "keep-listener-clb"
    listener_id        = "lbl-nonkgvc2"
    listener_name      = "dsadasd"
    vip                = "106.55.220.8"
    vport              = "80"
    region             = "gz"
    protocol           = "HTTP"
    zone               = "ap-guangzhou-6"
    numerical_vpc_id   = "5232945"
    load_balancer_type = "OPEN"
  }

  region          = "gz"
  flow_mode       = 0
  alb_type        = "clb"
  bot_status      = 0
  api_safe_status = 0
}
`
