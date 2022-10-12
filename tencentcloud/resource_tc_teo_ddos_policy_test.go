package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTeoDdosPolicy_basic -v
func TestAccTencentCloudTeoDdosPolicy_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDdosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDdosPolicyExists("tencentcloud_teo_ddos_policy.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.0.action", "forward"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.0.dport_end", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.0.dport_start", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.0.sport_start", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.action", "drop"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.dport_end", "65535"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.dport_start", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.protocol", "all"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.sport_end", "65535"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.acl.0.acls.1.sport_start", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.allow_block.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.allow_block.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.abnormal_connect_num", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.abnormal_syn_num", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.abnormal_syn_ratio", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.destination_connect_limit", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.destination_create_limit", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.drop_icmp", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.drop_other", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.drop_tcp", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.drop_udp", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.empty_connect_protect", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.source_connect_limit", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.source_create_limit", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.anti_ply.0.udp_shard", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.geo_ip.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.geo_ip.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.packet_filter.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.packet_filter.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.speed_limit.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.speed_limit.0.flux_limit", "0 bps"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.speed_limit.0.package_limit", "0 pps"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.status_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_ddos_policy.basic", "ddos_rule.0.status_info.0.ply_level", "middle"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_ddos_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDdosPolicyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		policyId := idSplit[1]

		policyId64, err := strconv.ParseInt(policyId, 10, 64)
		if err != nil {
			log.Printf("[READ]%s read teo ddosPolicy parseInt[%v] failed, reason:%+v", logId, policyId, err)
			return err
		}

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoDdosPolicy(ctx, zoneId, policyId64)
		if agents == nil {
			return fmt.Errorf("zone DdosPolicy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoDdosPolicyVar = `
variable "default_zone_id" {
  default = "` + defaultZoneId + `"
}

variable "default_policy_id" {
  default = ` + defaultPolicyId + `
}`

const testAccTeoDdosPolicy = testAccTeoDdosPolicyVar + `

resource "tencentcloud_teo_ddos_policy" "basic" {
    policy_id = var.default_policy_id
    zone_id   = var.default_zone_id

    ddos_rule {
        switch = "on"

        acl {
            switch = "on"

            acls {
                action      = "forward"
                dport_end   = 8080
                dport_start = 8080
                protocol    = "tcp"
                sport_end   = 65535
                sport_start = 1
            }
            acls {
                action      = "drop"
                dport_end   = 65535
                dport_start = 1
                protocol    = "all"
                sport_end   = 65535
                sport_start = 1
            }
        }

        allow_block {
            switch = "on"
        }

        anti_ply {
            abnormal_connect_num      = 0
            abnormal_syn_num          = 0
            abnormal_syn_ratio        = 0
            connect_timeout           = 0
            destination_connect_limit = 0
            destination_create_limit  = 0
            drop_icmp                 = "off"
            drop_other                = "off"
            drop_tcp                  = "off"
            drop_udp                  = "off"
            empty_connect_protect     = "off"
            source_connect_limit      = 0
            source_create_limit       = 0
            udp_shard                 = "off"
        }

        geo_ip {
            region_ids = []
            switch     = "on"
        }

        packet_filter {
            switch = "on"
        }

        speed_limit {
            flux_limit    = "0 bps"
            package_limit = "0 pps"
        }

        status_info {
            ply_level = "middle"
        }
    }
}

`
