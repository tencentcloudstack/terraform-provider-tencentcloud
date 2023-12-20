package dcdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcdb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDCDBSecurityGroupAttachmentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDcdbSecurityGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbSecurityGroupAttachment, tcacctest.DefaultDcdbSGName, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbSecurityGroupAttachmentExists("tencentcloud_dcdb_security_group_attachment.security_group_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_security_group_attachment.security_group_attachment", "security_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_security_group_attachment.security_group_attachment", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_security_group_attachment.security_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDcdbSecurityGroupAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dcdb_security_group_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		instanceId := idSplit[0]
		securityGroupId := idSplit[1]

		ret, err := dcdbService.DescribeDcdbSecurityGroup(ctx, instanceId)
		if err != nil {
			return err
		}

		for _, sg := range ret.Groups {
			if securityGroupId == *sg.SecurityGroupId {
				return fmt.Errorf("dcdb sg attachment instance still exist, instanceId: %v", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckDcdbSecurityGroupAttachmentExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb sg attachment instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb sg attachment instance id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		instanceId := idSplit[0]
		securityGroupId := idSplit[1]

		ret, err := dcdbService.DescribeDcdbSecurityGroup(ctx, instanceId)
		if err != nil {
			return err
		}

		for _, sg := range ret.Groups {
			if securityGroupId == *sg.SecurityGroupId {
				return nil
			}
		}
		return fmt.Errorf("dcdb sg attachment instance %v not found", rs.Primary.ID)
	}
}

const testAcc_sg_vpc_config = `
data "tencentcloud_security_groups" "internal" {
	name = "%s"
}
	
locals {
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}
`

const testAccDcdbSecurityGroupAttachment = testAcc_sg_vpc_config + `

resource "tencentcloud_dcdb_security_group_attachment" "security_group_attachment" {
  security_group_id = local.sg_id
  instance_id = "%s"
}

`
