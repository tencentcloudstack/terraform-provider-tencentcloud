package dayuv2_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdayuv2 "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayuv2"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDayuDdosIpAttachmentV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuDdosIpAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosIpAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosIpAttachmentExists("tencentcloud_dayu_ddos_ip_attachment_v2.boundip"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_ip_attachment_v2.boundip", "bgp_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_ip_attachment_v2.boundip", "bound_ip_list.#", "2"),
				),
			},
		},
	})
}

func testAccCheckDayuDdosIpAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcdayuv2.NewAntiddosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_ddos_ip_attachment_v2" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bgpInstanceId := idSplit[0]
		boundIps := idSplit[1]
		boundIpMap := make(map[string]bool)
		for _, boundIp := range strings.Split(boundIps, tccommon.COMMA_SP) {
			boundIpMap[boundIp] = true
		}
		boundip, err := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
		if err != nil {
			return err
		}
		if boundip.EipProductInfos == nil {
			return nil
		}
		for _, item := range boundip.EipProductInfos {
			if _, ok := boundIpMap[*item.Ip]; ok {
				return fmt.Errorf("DDoS ip attachment still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckDayuDdosIpAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DDoS ip attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DDoS ip attachment id is not set")
		}
		service := svcdayuv2.NewAntiddosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "tencentcloud_dayu_ddos_ip_attachment_v2" {
				continue
			}
			idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
			if len(idSplit) != 2 {
				return fmt.Errorf("id is broken,%s", rs.Primary.ID)
			}
			bgpInstanceId := idSplit[0]
			boundIps := idSplit[1]

			boundip, err := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
			if err != nil {
				return err
			}
			if boundip.EipProductInfos == nil {
				return nil
			}
			boundIpMap := make(map[string]bool)
			for _, item := range boundip.EipProductInfos {
				boundIpMap[*item.Ip] = true
			}
			for _, item := range strings.Split(boundIps, tccommon.COMMA_SP) {
				if _, ok := boundIpMap[item]; !ok {
					return fmt.Errorf("DDoS ip attachment not exists: %s, %s", rs.Primary.ID, item)
				}
			}
		}
		return nil
	}
}

const testAccDayuDdosIpAttachment_basic = `
resource "tencentcloud_dayu_ddos_ip_attachment_v2" "boundip" {
	bgp_instance_id = "bgp-00000fih"
	bound_ip_list {
		ip = "43.136.81.73"
		biz_type = "public"
		instance_id = "ins-eukucmzm"
		device_type = "cvm"
	}
	bound_ip_list {
		ip = "43.139.245.210"
		biz_type = "public"
		instance_id = "ins-c6vwi48a"
		device_type = "cvm"
	}
  }
`
