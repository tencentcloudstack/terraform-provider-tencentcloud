package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoSecurityPolicy_basic -v
func TestAccTencentCloudTeoSecurityPolicy_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityPolicyExists("tencentcloud_teo_security_policy.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_policy.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_policy.basic", "entity", "www."+defaultZoneName),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_policy.basic", "config.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSecurityPolicyExists(r string) resource.TestCheckFunc {
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
		entity := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoSecurityPolicy(ctx, zoneId, entity)
		if agents == nil {
			return fmt.Errorf("zone DnsSec %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoSecurityPolicyVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}

variable "zone_name" {
  default = "www.` + defaultZoneName + `"
}`

const testAccTeoSecurityPolicy = testAccTeoSecurityPolicyVar + `

resource "tencentcloud_teo_security_policy" "basic" {
  entity  = var.zone_name
  zone_id = var.zone_id

  config {
    switch_config {
      web_switch = "on"
    }

    bot_config {
      switch = "on"

      intelligence_rule {
        switch = "off"

        items {
          action = "drop"
          label  = "evil_bot"
        }
        items {
          action = "alg"
          label  = "suspect_bot"
        }
        items {
          action = "monitor"
          label  = "good_bot"
        }
        items {
          action = "trans"
          label  = "normal"
        }
      }

      managed_rule {
        action           = "monitor"
        alg_managed_ids  = []
        cap_managed_ids  = []
        drop_managed_ids = []
        mon_managed_ids  = [
          100000, 100001, 100002, 100003, 100006, 100007, 100008, 100009, 100010, 100011, 100012, 100013, 100014,
          100015, 100016, 100017, 100018, 100019, 100020, 100021, 100022, 10000003, 10000004, 10000005, 10000006,
          10000007, 10000008, 10000009,
        ]
        page_id           = 0
        punish_time       = 0
        response_code     = 0
        trans_managed_ids = []
      }

      portrait_rule {
        alg_managed_ids  = []
        cap_managed_ids  = []
        drop_managed_ids = []
        mon_managed_ids  = []
        switch           = "off"
      }
    }
  }
}`
