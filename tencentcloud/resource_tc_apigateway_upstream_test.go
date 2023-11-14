package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayUpstreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUpstream,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_upstream.upstream", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_upstream.upstream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayUpstream = `

resource "tencentcloud_apigateway_upstream" "upstream" {
  scheme = ""
  algorithm = ""
  uniq_vpc_id = ""
  upstream_name = ""
  upstream_description = ""
  upstream_type = ""
  retries = 
  upstream_host = ""
  nodes {
		host = ""
		port = 
		weight = 
		vm_instance_id = ""
		tags = 
		healthy = ""
		service_name = ""
		name_space = ""
		cluster_id = ""
		source = ""
		unique_service_name = ""

  }
  health_checker {
		enable_active_check = 
		enable_passive_check = 
		healthy_http_status = ""
		unhealthy_http_status = ""
		tcp_failure_threshold = 
		timeout_threshold = 
		http_failure_threshold = 
		active_check_http_path = ""
		active_check_timeout = 
		active_check_interval = 
		active_request_header = 
		unhealthy_timeout = 

  }
  k8s_service {
		weight = 
		cluster_id = ""
		namespace = ""
		service_name = ""
		port = 
		extra_labels {
			key = ""
			value = ""
		}
		name = ""

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
