package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTCRVPCAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRVPCAttachment_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRVPCAttachmentExists("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment", "status"),
					//this access ip will solve out with very long time
					//resource.TestCheckResourceAttrSet("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment", "access_ip"),
				),
				Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTCRVPCAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_vpc_attachment" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		vpcId := items[1]
		subnetId := items[2]

		_, has, err := tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
		if has {
			return fmt.Errorf("vpc attachment still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRVPCAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vpc attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc attachment id is not set")
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		vpcId := items[1]
		subnetId := items[2]

		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
		if !has {
			return fmt.Errorf("vpc attachment %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRVPCAttachment_basic = defaultVpcSubnets + `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "test_resource_attach"
  instance_type = "basic"
  delete_bucket = true
}

resource "tencentcloud_tcr_vpc_attachment" "mytcr_vpc_attachment" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
}`
