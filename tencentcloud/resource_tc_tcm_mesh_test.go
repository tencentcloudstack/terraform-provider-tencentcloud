package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmMeshResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmMesh,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcm_mesh.mesh", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcm_mesh.mesh",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmMesh = `

resource "tencentcloud_tcm_mesh" "mesh" {
  mesh_id = "mesh-xxxxxxxx"
  display_name = "test mesh"
  mesh_version = "1.8.1"
  type = "HOSTED"
  config {
		tracing {
			enable = true
			a_p_m {
				enable = true
				region = "ap-shanghai"
				instance_id = "apm-xxx"
			}
			sampling = 
			zipkin {
				address = "10.10.10.10:9411"
			}
		}
		prometheus {
			vpc_id = "vpc-xxx"
			subnet_id = "subnet-xxx"
			region = "sh"
			instance_id = "prom-xxx"
			custom_prom {
				is_public_addr = false
				vpc_id = "vpc-xxx"
				url = "http://x.x.x.x:9090"
				auth_type = "none, basic"
				username = "test"
				password = "test"
			}
		}
		istio {
			outbound_traffic_policy = "ALLOW_ANY"
			disable_policy_checks = true
			enable_pilot_h_t_t_p = true
			disable_h_t_t_p_retry = true
			smart_d_n_s {
				istio_meta_d_n_s_capture = true
				istio_meta_d_n_s_auto_allocate = true
			}
		}
		inject {
			exclude_i_p_ranges = 
			hold_application_until_proxy_starts = true
			hold_proxy_until_application_ends = true
		}
		sidecar_resources {
			limits {
				name = "cpu"
				quantity = "100m"
			}
			requests {
				name = "cpu"
				quantity = "100m"
			}
		}
		access_log {
			enable = true
			template = "istio"
			selected_range {
				items {
					namespace = "test"
					gateways = 
				}
				all = true
			}
			c_l_s {
				enable = true
				log_set = "f832fd4a-2b57-4573-ab6c-c3c69caf84c9"
				topic = "1ad05336-8afc-4e56-91e5-28d8a8511761"
			}
			encoding = "JSON"
			format = ""
			address = "10.10.10.4:3398"
			enable_server = true
			enable_stdout = true
		}

  }
  tag_list {
		key = "key"
		value = "value"
		passthrough = true

  }
}

`
