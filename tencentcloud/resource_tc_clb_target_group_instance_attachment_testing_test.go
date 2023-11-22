package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTestingClbTGAttachmentInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTGAttachmentInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbTGAttachmentInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip", "172.16.0.17"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "port", "99"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "weight", "3"),
				),
			},
			{
				Config: testAccTestingClbTGAttachmentInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip", "172.16.0.17"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "port", "99"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group_instance_attachment.test", "weight", "5"),
				),
			},
		},
	})
}

const testAccTestingClbTGAttachmentInstance_basic = `

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = "vpc-humgpppd"
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = "172.16.0.17"
    port            = 99
    weight          = 3
}
`

const testAccTestingClbTGAttachmentInstance_update = `

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = "vpc-humgpppd"
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = "172.16.0.17"
    port            = 99
    weight          = 5
}
`
