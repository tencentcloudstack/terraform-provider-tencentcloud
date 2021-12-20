package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataDayuDdosPoliciesName = "data.tencentcloud_dayu_ddos_policies.id_test"

func TestAccTencentCloudDataDayuDdosPolicies(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuDdosPoliciesBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuDdosPolicyExists("tencentcloud_dayu_ddos_policy.test_policy"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataDayuDdosPoliciesName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataDayuDdosPoliciesName, "list.0.policy_id"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.name", "tf_test_policy"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.drop_options.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.black_ips.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.port_filters.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.watermark_filters.#", "1"),
					resource.TestCheckResourceAttr(testDataDayuDdosPoliciesName, "list.0.packet_filters.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuDdosPoliciesBaic = `
resource "tencentcloud_dayu_ddos_policy" "test_policy" {
  resource_type         = "bgpip"
  name                  = "tf_test_policy"
  black_ips = ["1.1.1.1"]

  drop_options{
    drop_tcp  = true 
	drop_udp  = true
	drop_icmp  = true
	drop_other  = true
	drop_abroad  = true
	check_sync_conn = true
	s_new_limit = 100
	d_new_limit = 100
	s_conn_limit = 100
	d_conn_limit = 100
	tcp_mbps_limit = 100
	udp_mbps_limit = 100
	icmp_mbps_limit = 100
	other_mbps_limit = 100
	bad_conn_threshold = 100
	null_conn_enable = true
	conn_timeout = 500
	syn_rate = 50
	syn_limit = 100
  }

  port_filters{
	start_port = "2000"
	end_port = "2500"
	protocol = "all"
  	action = "drop"
	kind = 1
  }

  packet_filters{
	protocol = "tcp"
	action = "drop"
	d_start_port = 1000
	d_end_port = 1500
	s_start_port = 2000
	s_end_port = 2500
	pkt_length_max = 1400
	pkt_length_min = 1000
	is_include = true
	match_begin = "begin_l5"
	match_type = "pcre"
	depth = 1000
	offset = 500
  }

  watermark_filters{
  	tcp_port_list = ["2000-3000", "3500-4000"]
	udp_port_list = ["5000-6000"]
	offset = 50
	auto_remove = true
	open_switch = true
  }
}

data "tencentcloud_dayu_ddos_policies" "id_test" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  policy_id      = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
`
