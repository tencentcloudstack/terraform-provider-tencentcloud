package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDayuCCPolicyV2Resource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATION) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuCCPolicyV2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosPolicyV2Exists("tencentcloud_dayu_cc_policy_v2.demo"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_cc_policy_v2.demo", "cc_black_white_ips.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_cc_policy_v2.demo", "cc_geo_ip_policys.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_cc_policy_v2.demo", "cc_precision_policys.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_cc_policy_v2.demo", "cc_precision_req_limits.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_cc_policy_v2.demo", "thresholds.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDayuCCPolicyV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_cc_policy_v2" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 2 {
			return fmt.Errorf("broken ID of DDoS policy")
		}
		instanceId := items[0]
		business := items[1]
		antiddosService := AntiddosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		thresholdList, err := antiddosService.DescribeCCThresholdList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(thresholdList) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}

		ccGeoIpPolicys, err := antiddosService.DescribeCcGeoIPBlockConfigList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(ccGeoIpPolicys) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}

		ccBlackWhiteIpList, err := antiddosService.DescribeCcBlackWhiteIpList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(ccBlackWhiteIpList) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}

		ccPrecisionPlyList, err := antiddosService.DescribeCCPrecisionPlyList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(ccPrecisionPlyList) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}

		ccLevelList, err := antiddosService.DescribeCCLevelList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(ccLevelList) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}
		ccReqLimitPolicyList, err := antiddosService.DescribeCCReqLimitPolicyList(ctx, business, instanceId)
		if err != nil {
			return err
		}
		if len(ccReqLimitPolicyList) != 0 {
			return fmt.Errorf("delete cc policy %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuCCPolicyV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		return nil
	}
}

const testAccDayuCCPolicyV2 string = `
resource "tencentcloud_dayu_cc_policy_v2" "demo" {
	resource_id="bgpip-000004xf"
	business="bgpip"
	thresholds {
	  domain="12.com"
	  threshold=0
	}
	cc_geo_ip_policys {
	  action="drop"
	  region_type="china"
	  domain="12.com"
	  protocol="http"
	}
  
	cc_black_white_ips {
	  protocol="http"
	  domain="12.com"
	  black_white_ip="1.2.3.4"
	  type="black"
	}
	cc_precision_policys{
	  policy_action="drop"
	  domain="1.com"
	  protocol="http"
	  ip="162.62.163.34"
	  policys {
		field_name="cgi"
		field_type="value"
		value="12123.com"
		value_operator="equal"
	  }
	}
	cc_precision_req_limits {
	  domain="11.com"
	  protocol="http"
	  level="loose"
	  policys {
		action="alg"
		execute_duration=2
		mode="equal"
		period=5
		request_num=12
		uri="15.com"
	  }
	}
  }`
