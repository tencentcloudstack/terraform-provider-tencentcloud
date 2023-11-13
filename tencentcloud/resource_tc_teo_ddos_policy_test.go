package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoDdosPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDdosPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_policy.ddos_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_ddos_policy.ddos_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoDdosPolicy = `

resource "tencentcloud_teo_ddos_policy" "ddos_policy" {
  zone_id = &lt;nil&gt;
  policy_id = &lt;nil&gt;
  d_do_s_rule {
		switch = &lt;nil&gt;
		status_info {
			ply_level = &lt;nil&gt;
		}
		geo_ip {
			region_ids = &lt;nil&gt;
			switch = &lt;nil&gt;
		}
		allow_block {
			switch = &lt;nil&gt;
			allow_block_ips {
				ip = &lt;nil&gt;
				type = &lt;nil&gt;
			}
		}
		anti_ply {
			drop_tcp = &lt;nil&gt;
			drop_udp = &lt;nil&gt;
			drop_icmp = &lt;nil&gt;
			drop_other = &lt;nil&gt;
			source_create_limit = &lt;nil&gt;
			source_connect_limit = &lt;nil&gt;
			destination_create_limit = &lt;nil&gt;
			destination_connect_limit = &lt;nil&gt;
			abnormal_connect_num = &lt;nil&gt;
			abnormal_syn_ratio = &lt;nil&gt;
			abnormal_syn_num = &lt;nil&gt;
			connect_timeout = &lt;nil&gt;
			empty_connect_protect = &lt;nil&gt;
			udp_shard = &lt;nil&gt;
		}
		packet_filter {
			switch = &lt;nil&gt;
			packet_filters {
				action = &lt;nil&gt;
				protocol = &lt;nil&gt;
				dport_start = &lt;nil&gt;
				dport_end = &lt;nil&gt;
				packet_min = &lt;nil&gt;
				packet_max = &lt;nil&gt;
				sport_start = &lt;nil&gt;
				sport_end = &lt;nil&gt;
				match_type = &lt;nil&gt;
				is_not = &lt;nil&gt;
				offset = &lt;nil&gt;
				depth = &lt;nil&gt;
				match_begin = &lt;nil&gt;
				str = &lt;nil&gt;
				match_type2 = &lt;nil&gt;
				is_not2 = &lt;nil&gt;
				offset2 = &lt;nil&gt;
				depth2 = &lt;nil&gt;
				match_begin2 = &lt;nil&gt;
				str2 = &lt;nil&gt;
				match_logic = &lt;nil&gt;
			}
		}
		acl {
			switch = &lt;nil&gt;
			acls {
				dport_end = &lt;nil&gt;
				dport_start = &lt;nil&gt;
				sport_end = &lt;nil&gt;
				sport_start = &lt;nil&gt;
				protocol = &lt;nil&gt;
				action = &lt;nil&gt;
			}
		}
		speed_limit {
			package_limit = &lt;nil&gt;
			flux_limit = &lt;nil&gt;
		}

  }
}

`
