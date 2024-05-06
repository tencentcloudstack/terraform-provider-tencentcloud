package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoOriginGroup_basic -v
func TestAccTencentCloudTeoOriginGroup_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "name", "keep-group-1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "type", "GENERAL"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.record"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.type", "IP_DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.weight", "100"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckOriginGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_origin_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup != nil {
			return fmt.Errorf("zone originGroup %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckOriginGroupExists(r string) resource.TestCheckFunc {
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
		originGroupId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup == nil {
			return fmt.Errorf("zone originGroup %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoOriginGroup = testAccTeoZone + `

resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = tencentcloud_teo_zone.basic.id

  records {
    record  = var.zone_name
    type    = "IP_DOMAIN"
    weight  = 100
    private = false
  }
}

`
