package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVpcAclAttachment_basic(t *testing.T) {
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
		err := service.DescribeByAclId(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][ACL attachment][Destroy] check: acl attachment still exists: %s", rs.Primary.ID)
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
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][ACL attachment][Exists] check:  %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][ACL attachment][Exists] check: id is not set")
		}
		err := service.DescribeByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][ACL attachment][Exists] check:  still exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAclAttachment_basic = `
data "tencentcloud_vpc_instances" "id_instances" {
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
		subnet_ids = data.tencentcloud_vpc_instances.id_instances.instance_list[0].subnet_ids
}
`
