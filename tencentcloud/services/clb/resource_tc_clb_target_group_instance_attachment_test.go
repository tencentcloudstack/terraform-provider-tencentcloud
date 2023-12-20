package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudClbTGAttachmentInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbTGAttachmentInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTGAttachmentInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "weight"),
				),
			},
			{
				Config: testAccClbTGAttachmentInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "weight"),
				),
			},
		},
	})
}

func testAccCheckClbTGAttachmentInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group_instance_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("target group instance attachment id is not set")
		}
		targetGroupId := idSplit[0]
		bindIp := idSplit[1]
		port, err := strconv.ParseUint(idSplit[2], 0, 64)
		if err != nil {
			return err
		}

		filters := make(map[string]string)
		filters["TargetGroupId"] = targetGroupId
		filters["BindIP"] = bindIp
		targetGroupInstances, err := clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return err
		}
		for _, tgInstance := range targetGroupInstances {
			if *tgInstance.Port == port {
				return fmt.Errorf("[CHECK][CLB target group instance attachment][Destroy] check: CLB target group instance attachment still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckClbTGAttachmentInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] check: CLB target group id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("target group instance attachment id is not set")
		}
		targetGroupId := idSplit[0]
		bindIp := idSplit[1]
		port, err := strconv.ParseUint(idSplit[2], 0, 64)
		if err != nil {
			return err
		}

		filters := make(map[string]string)
		filters["TargetGroupId"] = targetGroupId
		filters["BindIP"] = bindIp
		targetGroupInstances, err := clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return err
		}
		for _, tgInstance := range targetGroupInstances {
			if *tgInstance.Port == port {
				return nil
			}
		}
		return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] id %s is not exist", rs.Primary.ID)
	}
}

const testAccClbTGAttachmentInstance_basic = tcacctest.InstanceCommonTestCase + `

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.default.id
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = var.cvm_vpc_id
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip 
    port            = 88
    weight          = 3
}
`

const testAccClbTGAttachmentInstance_update = tcacctest.InstanceCommonTestCase + `

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.default.id
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = var.cvm_vpc_id
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip 
    port            = 88
    weight          = 5
}
`
