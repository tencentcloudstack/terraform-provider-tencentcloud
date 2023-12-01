package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoAccelerationDomainResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeoAccelerationDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoAccelerationDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoAccelerationDomainExists("tencentcloud_teo_acceleration_domain.acceleration_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_acceleration_domain.acceleration_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_acceleration_domain.acceleration_domain", "domain_name", "aaa.tf-teo.xyz"),
					resource.TestCheckResourceAttr("tencentcloud_teo_acceleration_domain.acceleration_domain", "origin_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_acceleration_domain.acceleration_domain", "origin_info.0.origin", "150.109.8.1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_acceleration_domain.acceleration_domain", "origin_info.0.origin_type", "IP_DOMAIN"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_acceleration_domain.acceleration_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTeoAccelerationDomainDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_acceleration_domain" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		domainName := idSplit[1]

		agents, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
		if agents != nil {
			return fmt.Errorf("AccelerationDomain %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTeoAccelerationDomainExists(r string) resource.TestCheckFunc {
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
		domainName := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
		if agents == nil {
			return fmt.Errorf("AccelerationDomain %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoAccelerationDomain = testAccTeoZone + `

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
    zone_id     = "zone-2o0i41pv2h8c"
    domain_name = "aaa.makn.cn"

    origin_info {
        origin      = "150.109.8.1"
        origin_type = "IP_DOMAIN"
    }
}

`
