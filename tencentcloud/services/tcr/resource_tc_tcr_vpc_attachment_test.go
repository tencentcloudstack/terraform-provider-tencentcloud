package tcr_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTcrVPCAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRVPCAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRVPCAttachment_basic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRVPCAttachmentExists("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment_resource"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment_resource", "status"),
					// this access ip will solve out with very long time
					// resource.TestCheckResourceAttrSet("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment_resource", "access_ip"),
				),
				// Destroy: false,
			},
			{
				ResourceName:      "tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTCRVPCAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_vpc_attachment" {
			continue
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vpc attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc attachment id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		vpcId := items[1]
		subnetId := items[2]

		tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTCRVPCAttachment_basic = DefaultTcrVpcSubnets + `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "test-resource-attach"
  instance_type = "basic"
  delete_bucket = true
}

resource "tencentcloud_tcr_vpc_attachment" "mytcr_vpc_attachment_resource" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
}`
