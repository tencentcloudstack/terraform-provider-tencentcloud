package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDayuDdosIpAttachmentV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuDdosIpAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuDdosIpAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuDdosIpAttachmentExists("tencentcloud_dayu_ddos_ip_attachment_v2.boundip"),
					resource.TestCheckResourceAttrSet("tencentcloud_dayu_ddos_ip_attachment_v2.boundip", "bgp_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dayu_ddos_ip_attachment_v2.boundip", "bound_ip_list.#", "1"),
				),
			},
		},
	})
}

func testAccCheckDayuDdosIpAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dayu_ddos_ip_attachment_v2" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bgpInstanceId := idSplit[0]
		boundIps := idSplit[1]
		boundIpMap := make(map[string]bool)
		for _, boundIp := range strings.Split(boundIps, COMMA_SP) {
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DDoS ip attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DDoS ip attachment id is not set")
		}
		service := AntiddosService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "tencentcloud_dayu_ddos_ip_attachment_v2" {
				continue
			}
			idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
			for _, item := range strings.Split(boundIps, COMMA_SP) {
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
	bgp_instance_id = "bgp-000001co"
	bound_ip_list {
		ip = "43.136.81.73"
		biz_type = "public"
		instance_id = "ins-eukucmzm"
		device_type = "cvm"
	}
  }
`
