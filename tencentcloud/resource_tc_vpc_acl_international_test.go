package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudInternationalVpcAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.#", "2"),
				),
			},
			{
				ResourceName: "tencentcloud_vpc_acl.foo",
				ImportState:  true,
			},
		},
	})
}

func testAccInternationalCheckVpcACLExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := service.DescribeNetWorkByACLID(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("vpc network acl %s not exists", rs.Primary.ID)
	}
}

func testAccInternationalCheckVpcACLDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc_acl" {
			continue
		}
		_, has, err := service.DescribeNetWorkByACLID(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}

		return fmt.Errorf("vpc acl %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccInternationalVpcACLConfig = `
data "tencentcloud_vpc_instances" "default" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id  = data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
    	"ACCEPT#192.168.1.0/24#80-90#TCP",
	]
}  
`
