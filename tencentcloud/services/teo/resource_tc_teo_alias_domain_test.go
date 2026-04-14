package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoAliasDomainResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoAliasDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoAliasDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoAliasDomainExists("tencentcloud_teo_alias_domain.alias_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.alias_domain", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.alias_domain", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.alias_domain", "alias_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.alias_domain", "target_name"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.alias_domain", "paused", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_alias_domain.alias_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoAliasDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoAliasDomainExists("tencentcloud_teo_alias_domain.alias_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.alias_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.alias_domain", "target_name", "updated.tf-teo.xyz"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.alias_domain", "paused", "true"),
				),
			},
		},
	})
}

func testAccCheckTeoAliasDomainDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_alias_domain" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		aliasName := idSplit[1]

		aliasDomain, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if aliasDomain != nil {
			return fmt.Errorf("AliasDomain %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTeoAliasDomainExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		aliasName := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		aliasDomain, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if aliasDomain == nil {
			return fmt.Errorf("AliasDomain %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoAliasDomain = testAccTeoZone + `

resource "tencentcloud_teo_alias_domain" "alias_domain" {
    zone_id     = tencentcloud_teo_zone.basic.id
    alias_name  = "alias.test.tf-teo.xyz"
    target_name = "test.tf-teo.xyz"
}

`

const testAccTeoAliasDomainUpdate = testAccTeoZone + `

resource "tencentcloud_teo_alias_domain" "alias_domain" {
    zone_id     = tencentcloud_teo_zone.basic.id
    alias_name  = "alias.test.tf-teo.xyz"
    target_name = "updated.tf-teo.xyz"
    paused      = true
}

`
