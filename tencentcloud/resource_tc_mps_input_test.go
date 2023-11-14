package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsInputResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsInput,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_input.input", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_input.input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsInput = `

resource "tencentcloud_mps_input" "input" {
  flow_id = ""
  input_group {
		input_name = ""
		protocol = ""
		description = ""
		allow_ip_list = 
		s_r_t_settings {
			mode = ""
			stream_id = ""
			latency = 
			recv_latency = 
			peer_latency = 
			peer_idle_timeout = 
			passphrase = ""
			pb_key_len = 
			source_addresses {
				ip = ""
				port = 
			}
		}
		r_t_p_settings {
			f_e_c = ""
			idle_timeout = 
		}
		fail_over = ""
		r_t_m_p_pull_settings {
			source_addresses {
				tc_url = ""
				stream_key = ""
			}
		}
		r_t_s_p_pull_settings {
			source_addresses {
				url = ""
			}
		}
		h_l_s_pull_settings {
			source_addresses {
				url = ""
			}
		}
		resilient_stream {
			enable = 
			buffer_time = 
		}

  }
}

`
