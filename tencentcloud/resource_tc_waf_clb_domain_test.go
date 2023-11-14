package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafClbDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafClbDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_domain.clb_domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_clb_domain.clb_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafClbDomain = `

resource "tencentcloud_waf_clb_domain" "clb_domain" {
  host {
		domain = ""
		domain_id = ""
		main_domain = ""
		mode = 
		status = 
		state = 
		engine = 
		is_cdn = 
		load_balancer_set {
			load_balancer_id = ""
			load_balancer_name = ""
			listener_id = ""
			listener_name = ""
			vip = ""
			vport = 
			region = ""
			protocol = ""
			zone = ""
			numerical_vpc_id = 
			load_balancer_type = ""
		}
		region = ""
		edition = ""
		flow_mode = 
		cls_status = 
		level = 
		cdc_clusters = 
		alb_type = ""
		ip_headers = 
		engine_type = 

  }
  instance_i_d = ""
}

`
