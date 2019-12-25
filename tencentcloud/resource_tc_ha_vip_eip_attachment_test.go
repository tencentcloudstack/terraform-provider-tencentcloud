package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudHaVipEipAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHaVipEipAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipEipAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipEipAttachmentExists("tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment_basic", "havip_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment_basic", "address_ip"),
				),
			},
			{
				ResourceName:      "tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckHaVipEipAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	vpcService := VpcService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ha_vip_eip_attachment" {
			continue
		}

		_, _, err := vpcService.DescribeHaVipEipById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("HA VIP EIP attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckHaVipEipAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("HA VIP EIP attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HA VIP EIP attachment id is not set")
		}
		vpcService := VpcService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, _, err := vpcService.DescribeHaVipEipById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccHaVipEipAttachment_basic = defaultVpcVariable + `
#Create EIP
resource "tencentcloud_eip" "eip" {
  name = "havip_eip"
}
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
resource "tencentcloud_ha_vip_eip_attachment" "ha_vip_eip_attachment_basic"{
  havip_id = tencentcloud_ha_vip.havip.id
  address_ip = tencentcloud_eip.eip.public_ip
}
`
