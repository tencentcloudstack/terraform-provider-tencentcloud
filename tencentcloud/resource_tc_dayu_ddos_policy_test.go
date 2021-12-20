package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDayuDdosPolicyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyExists("tencentcloud_dayu_ddos_policy.test_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_policy.test_policy", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "name", "tf_test_policy"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "white_ips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "black_ips.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_tcp", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_udp", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_icmp", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_other", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_abroad", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.check_sync_conn", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.null_conn_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.s_new_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.d_new_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.s_conn_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.d_conn_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.tcp_mbps_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.udp_mbps_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.icmp_mbps_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.bad_conn_threshold", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.conn_timeout", "500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.syn_rate", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.syn_limit", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.start_port", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.end_port", "2500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.protocol", "all"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.kind", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.action", "drop"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.d_start_port", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.d_end_port", "1500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.s_start_port", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.s_end_port", "2500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.pkt_length_max", "1400"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.pkt_length_min", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.is_include", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.match_begin", "begin_l5"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.match_type", "pcre"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.depth", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.offset", "500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.tcp_port_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.tcp_port_list.0", "2000-3000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.udp_port_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.udp_port_list.0", "5000-6000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.offset", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.auto_remove", "true"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.open_switch", "true"),
				),
			},
			{
				Config: testAccDayuDdosPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyExists("tencentcloud_dayu_ddos_policy.test_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_policy.test_policy", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "name", "tf_test_policy_update"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "black_ips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "white_ips.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_tcp", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_udp", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_icmp", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_other", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.drop_abroad", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.check_sync_conn", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.null_conn_enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.s_new_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.d_new_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.s_conn_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.d_conn_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.tcp_mbps_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.udp_mbps_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.icmp_mbps_limit", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.bad_conn_threshold", "50"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.conn_timeout", "3000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.syn_rate", "80"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "drop_options.0.syn_limit", "80"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.start_port", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.end_port", "150"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.kind", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "port_filters.0.action", "transmit"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.d_start_port", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.d_end_port", "150"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.s_start_port", "200"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.s_end_port", "250"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.pkt_length_max", "140"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.pkt_length_min", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.is_include", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.match_begin", "no_match"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.match_type", "sunday"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.depth", "500"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "packet_filters.0.offset", "100"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.tcp_port_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.tcp_port_list.0", "5000-6000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.udp_port_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.udp_port_list.0", "2000-3000"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.offset", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.auto_remove", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy.test_policy", "watermark_filters.0.open_switch", "false"),
				),
			},
		},
	})
}

func testAccCheckDayuDdosPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_ddos_policy" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDos policy")
		}
		resourceType := items[0]
		policyId := items[1]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeDdosPolicy(ctx, resourceType, policyId)
		if err != nil {
			_, has, err = service.DescribeDdosPolicy(ctx, resourceType, policyId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuDdosPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDos policy")
		}
		resourceType := items[0]
		policyId := items[1]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeDdosPolicy(ctx, resourceType, policyId)
		if err != nil {
			_, has, err = service.DescribeDdosPolicy(ctx, resourceType, policyId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("DDoS policy %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuDdosPolicy string = `
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
`
const testAccDayuDdosPolicyUpdate string = `
resource "tencentcloud_dayu_ddos_policy" "test_policy" {
  resource_type         = "bgpip"
  name                  = "tf_test_policy_update"
  white_ips = ["2.2.2.2"]
  
  drop_options {
    drop_tcp  = false 
	drop_udp  = false
	drop_icmp  = false
	drop_other  = false
	drop_abroad  = false
	check_sync_conn = false
	s_new_limit = 50
	d_new_limit = 50
	s_conn_limit = 50
	d_conn_limit = 50
	tcp_mbps_limit = 50
	udp_mbps_limit = 50
	icmp_mbps_limit = 50
	other_mbps_limit = 50
	bad_conn_threshold = 50
	null_conn_enable = false
	conn_timeout = 3000
	syn_rate = 80
	syn_limit = 80
  }

  port_filters {
	start_port = "100"
	end_port = "150"
	protocol = "tcp"
  	action = "transmit"
	kind = 0
  }

  packet_filters {
	protocol = "icmp"
	action = "drop_black"
	d_start_port = 100
	d_end_port = 150
	s_start_port = 200
	s_end_port = 250
	pkt_length_max = 140
	pkt_length_min = 100
	is_include = false
	match_begin = "no_match"
	match_type = "sunday"
	depth = 500
	offset = 100
  }

  watermark_filters {
  	tcp_port_list = ["5000-6000"]
	udp_port_list = ["2000-3000", "3500-4000"]
	offset = 0
	auto_remove = false
	open_switch = false
  }
}
`
