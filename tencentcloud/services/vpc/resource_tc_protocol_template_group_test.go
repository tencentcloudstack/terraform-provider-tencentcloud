package vpc_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudProtocolTemplateGroup_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckProtocolTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtocolTemplateGroup_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_protocol_template_group.group", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template_group.group", "template_ids.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_protocol_template_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProtocolTemplateGroup_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateGroupExists("tencentcloud_protocol_template_group.group"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template_group.group", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template_group.group", "template_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckProtocolTemplateGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_protocol_template_group" {
			continue
		}

		_, has, err := vpcService.DescribeServiceTemplateGroupById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("Service template group still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckProtocolTemplateGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Service template group %s is not found", n)
		}

		vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := vpcService.DescribeServiceTemplateGroupById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Service template group %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccProtocolTemplateGroup_basic = `
resource "tencentcloud_protocol_template" "template" {
  name = "test"
  protocols = ["tcp:all"]
}

resource "tencentcloud_protocol_template_group" "group"{
	name = "test"
	template_ids = [tencentcloud_protocol_template.template.id]
}
`

const testAccProtocolTemplateGroup_basic_update_remark = `
resource "tencentcloud_protocol_template" "template" {
  name = "test"
  protocols = ["tcp:all"]
}

resource "tencentcloud_protocol_template" "templateB" {
  name = "testB"
  protocols = ["tcp:80", "udp:90,111"]
}

resource "tencentcloud_protocol_template_group" "group"{
	name = "test_update"
	template_ids = [tencentcloud_protocol_template.templateB.id]
}
`
