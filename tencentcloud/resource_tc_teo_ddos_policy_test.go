package tencentcloud

//func TestAccTencentCloudTeoDdosPolicy_basic(t *testing.T) {
//	t.Parallel()
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:  func() { testAccPreCheck(t) },
//		Providers: testAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccTeoDdosPolicy,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttrSet("tencentcloud_teo_ddos_policy.ddosPolicy", "id"),
//				),
//			},
//			{
//				ResourceName:      "tencentcloud_teo_ddos_policy.ddosPolicy",
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

const testAccTeoDdosPolicy = `

resource "tencentcloud_teo_ddos_policy" "ddosPolicy" {
  zone_id = ""
  policy_id = ""
  ddos_rule {
			switch = ""
			udp_shard_open = ""
		ddos_status_info {
				ply_level = ""
		}
		ddos_geo_ip {
				region_id = ""
				switch = ""
		}
		ddos_allow_block {
				switch = ""
			user_allow_block_ip {
					ip = ""
					mask = ""
					type = ""
					ip2 = ""
					mask2 = ""
			}
		}
		ddos_anti_ply {
				drop_tcp = ""
				drop_udp = ""
				drop_icmp = ""
				drop_other = ""
				source_create_limit = ""
				source_connect_limit = ""
				destination_create_limit = ""
				destination_connect_limit = ""
				abnormal_connect_num = ""
				abnormal_syn_ratio = ""
				abnormal_syn_num = ""
				connect_timeout = ""
				empty_connect_protect = ""
				udp_shard = ""
		}
		ddos_packet_filter {
				switch = ""
			packet_filter {
					action = ""
					protocol = ""
					dport_start = ""
					dport_end = ""
					packet_min = ""
					packet_max = ""
					sport_start = ""
					sport_end = ""
					match_type = ""
					is_not = ""
					offset = ""
					depth = ""
					match_begin = ""
					str = ""
					match_type2 = ""
					is_not2 = ""
					offset2 = ""
					depth2 = ""
					match_begin2 = ""
					str2 = ""
					match_logic = ""
			}
		}
		ddos_acl {
				switch = ""
			acl {
					dport_end = ""
					dport_start = ""
					sport_end = ""
					sport_start = ""
					protocol = ""
					action = ""
					default = ""
			}
		}

  }
}

`
