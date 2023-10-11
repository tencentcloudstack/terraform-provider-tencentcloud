package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsOutputResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsOutput,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_output.output", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_output.output",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsOutput = `

resource "tencentcloud_mps_output" "output" {
  flow_id = ""
  output {
		output_name = ""
		description = ""
		protocol = ""
		output_region = ""
		s_r_t_settings {
			destinations {
				ip = ""
				port = 
			}
			stream_id = ""
			latency = 
			recv_latency = 
			peer_latency = 
			peer_idle_timeout = 
			passphrase = ""
			pb_key_len = 
			mode = ""
		}
		r_t_m_p_settings {
			destinations {
				url = ""
				stream_key = ""
			}
			chunk_size = 
		}
		r_t_p_settings {
			destinations {
				ip = ""
				port = 
			}
			f_e_c = ""
			idle_timeout = 
		}
		allow_ip_list = 
		max_concurrent = 

  }
}

`
