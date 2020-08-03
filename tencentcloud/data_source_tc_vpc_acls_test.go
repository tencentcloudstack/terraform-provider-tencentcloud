package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVpcACLBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcACLInstances,

				Check: resource.ComposeTestCheckFunc(
					// id filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_acl.default"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_acl.default", "name", "test_acl"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_acl.default", "egress.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_acl.default", "ingress.#", "1"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcACLInstances = `
data "tencentcloud_vpc_instances" "test" {
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.test.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
	]
}  

data "tencentcloud_vpc_instances" "default" {
	name = "test_acl"
	result_output_file="data_source_tc_vpc_acls.txt"
}
`
