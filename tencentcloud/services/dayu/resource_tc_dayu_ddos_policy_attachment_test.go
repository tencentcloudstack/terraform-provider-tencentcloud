package dayu_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDayuDdosPolicyAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuDdosPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosPolicyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyAttachmentExists("tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment_basic", "resource_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment_basic", "policy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment_basic", "resource_type")),
			},
		},
	})
}

func testAccCheckDayuDdosPolicyAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dayuService := svcdayu.NewDayuService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_ddos_policy_attachment" {
			continue
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy attachment")
		}
		resourceId := items[0]
		resourceType := items[1]
		policyId := items[2]
		_, has, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("DDoS policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuDdosPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DDoS policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DDoS policy attachment id is not set")
		}
		dayuService := svcdayu.NewDayuService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy attachment")
		}
		resourceId := items[0]
		resourceType := items[1]
		policyId := items[2]
		_, has, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("DDoS policy attachment does not exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDayuDdosPolicyAttachment_basic = `
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

resource "tencentcloud_dayu_ddos_policy_attachment" "dayu_ddos_policy_attachment_basic" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  resource_id = "bgpip-00000294"
  policy_id = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
`
