package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDataDayuDdosPolicyAttachmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuDdosPolicyAttachmentsDataSource_basic, defaultDayuBgpIp),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuDdosPolicyAttachmentExists("tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment"),
					resource.TestCheckResourceAttr("data.tencentcloud_dayu_ddos_policy_attachments.dayu_ddos_policy_attachments", "dayu_ddos_policy_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dayu_ddos_policy_attachments.dayu_ddos_policy_attachments", "dayu_ddos_policy_attachment_list.0.resource_id"),
				),
			},
		},
	})
}

const testAccDayuDdosPolicyAttachmentsDataSource_basic = `
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

resource "tencentcloud_dayu_ddos_policy_attachment" "dayu_ddos_policy_attachment" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  resource_id = "%s"
  policy_id = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}

data "tencentcloud_dayu_ddos_policy_attachments" "dayu_ddos_policy_attachments" {
  resource_id = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_id
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
`
