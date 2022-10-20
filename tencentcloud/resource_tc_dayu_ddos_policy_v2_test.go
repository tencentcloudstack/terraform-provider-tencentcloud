package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDayuDdosPolicyV2Resource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosPolicyV2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyV2Exists("tencentcloud_dayu_ddos_policy_v2.test_policy"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "black_white_ips.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "ddos_connect_limit.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "ddos_geo_ip_block_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "ddos_speed_limit_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "packet_filters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_policy_v2.test_policy", "protocol_block_config.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDayuDdosPolicyV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_ddos_policy_v2" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDoS policy")
		}
		resourceId := items[0]

		antiddosService := AntiddosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		blackWhiteIpRelationList, err := antiddosService.DescribeListBlackWhiteIpList(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(blackWhiteIpRelationList) != 0 {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}

		aclConfigRelationList, err := antiddosService.DescribeListPortAclList(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(aclConfigRelationList) != 0 {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}

		ddosGeoIPBlockConfigRelations, err := antiddosService.DescribeListDDoSGeoIPBlockConfig(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(ddosGeoIPBlockConfigRelations) != 0 {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}

		ddosSpeedLimitConfigRelations, err := antiddosService.DescribeListDDoSSpeedLimitConfig(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(ddosSpeedLimitConfigRelations) != 0 {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}

		packetFilterRelationList, err := antiddosService.DescribeListPacketFilterConfig(ctx, resourceId)
		if err != nil {
			return err
		}
		if len(packetFilterRelationList) != 0 {
			return fmt.Errorf("delete DDoS policy %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuDdosPolicyV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		return nil
	}
}

const testAccDayuDdosPolicyV2 string = `
resource "tencentcloud_dayu_ddos_policy_v2" "test_policy" {
	resource_id = "bgpip-000004xf"
	business = "bgpip"
	ddos_threshold="100"
	ddos_level="low"
	black_white_ips {
		ip = "1.2.3.4"
		ip_type = "black"
	}
	acls {
		action = "transmit"
		d_port_start = 1
		d_port_end = 10
		s_port_start=10
		s_port_end=20
		priority=9
		forward_protocol="all"
	}
	protocol_block_config {
		drop_icmp=1
		drop_tcp=0
		drop_udp=0
		drop_other=0
	}
	ddos_connect_limit {
		sd_new_limit=10
		sd_conn_limit=11
		dst_new_limit=20
		dst_conn_limit=21
		bad_conn_threshold=30
		syn_rate=10
		syn_limit=20
		conn_timeout=30
		null_conn_enable=1
	}
	ddos_ai="on"
	ddos_geo_ip_block_config {
		action="drop"
		area_list=["100001"]
		region_type="customized"
	}
	ddos_speed_limit_config {
		protocol_list="TCP"
		dst_port_list="10"
		mode=1
		packet_rate=10
		bandwidth=20
	}
	packet_filters {
		action="drop"
		protocol="all"
		s_port_start=10
		s_port_end=10
		d_port_start=20
		d_port_end=20
		pktlen_min=30
		pktlen_max=30
		str="12"
		str2="30"
		match_logic="and"
		match_type="pcre"
		match_type2="pcre"
		match_begin="begin_l3"
		match_begin2="begin_l3"
		depth=2
		depth2=3
		offset=1
		offset2=2
		is_not=0
		is_not2=0
	}
}`
