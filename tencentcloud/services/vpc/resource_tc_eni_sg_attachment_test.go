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

func TestAccTencentCloudVpcEniSgAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEniSgAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEniSgAttachment,
				Check: resource.ComposeTestCheckFunc(testAccCheckEniSgAttachmentExists("tencentcloud_eni_sg_attachment.eni_sg_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_sg_attachment.eni_sg_attachment", "network_interface_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_sg_attachment.eni_sg_attachment", "security_group_ids.#")),
			},
			{
				ResourceName:      "tencentcloud_eni_sg_attachment.eni_sg_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEniSgAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		enis, err := service.DescribeEniById(ctx, []string{rs.Primary.ID})
		if err != nil {
			return err
		}
		tmpMap := make(map[string]struct{})

		for _, v := range enis[0].GroupSet {
			tmpMap[*v] = struct{}{}
		}
		value1, exists1 := rs.Primary.Attributes["security_group_ids.0"]
		value2, exists2 := rs.Primary.Attributes["security_group_ids.1"]
		if exists1 && exists2 {
			_, exists3 := tmpMap[value1]
			_, exists4 := tmpMap[value2]
			if exists3 && exists4 {
				return nil
			}
		}

		return fmt.Errorf("EniSgAttachment %s not found on server", rs.Primary.ID)
	}
}

func testAccCheckEniSgAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eni_sg_attachment" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		enis, err := service.DescribeEniById(ctx, []string{rs.Primary.ID})
		if err != nil {
			return err
		}
		if enis == nil || len(enis) < 1 {
			return nil
		}
		for _, v := range enis[0].GroupSet {
			if rs.Primary.ID == *v {
				return fmt.Errorf("delete EniSgAttachment %s fail, still on server", rs.Primary.ID)
			}
		}
		return nil
	}
	return nil
}

const testAccVpcEniSgAttachment = `

resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = ["eni-p0hkgx8p"]
  security_group_ids    = ["sg-902tl7t7", "sg-edmur627"]
}

`
