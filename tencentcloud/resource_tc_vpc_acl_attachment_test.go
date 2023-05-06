package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudVpcAclAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testVpcAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAclAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testVpcAclAttachmentExists("tencentcloud_vpc_acl_attachment.attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_acl_attachment.attachment", "acl_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_acl_attachment.attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testVpcAclAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc_acl_attachment" {
			continue
		}
		has, err := service.DescribeByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("[CHECK][ACL attachment][Destroy] check: ACL attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testVpcAclAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][ACL attachment][Exists] check:  %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][ACL attachment][Exists] check: id is not set")
		}
		has, err := service.DescribeByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if !has {
			return fmt.Errorf("[CHECK][ACL attachment][Exists] check: not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAclAttachment_basic = `
data "tencentcloud_vpc_instances" "id_instances" {
	is_default = true
}
resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id  = data.tencentcloud_vpc_instances.id_instances.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
		"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
    	"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
}
resource "tencentcloud_vpc_acl_attachment" "attachment"{
		acl_id = tencentcloud_vpc_acl.foo.id
		subnet_id = data.tencentcloud_vpc_instances.id_instances.instance_list[0].subnet_ids[0]
}
`
