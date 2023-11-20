package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingVpcAclRulesResource_Update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
				),
			},
			{
				Config: testAccTestingVpcACLConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
				),
			},
		},
	})
}

const testAccTestingVpcACLConfig = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}
resource "tencentcloud_vpc_acl" "foo" {
  vpc_id  = tencentcloud_vpc.foo.id
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

const testAccTestingVpcACLConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_vpc_acl" "foo" {
  vpc_id  = tencentcloud_vpc.foo.id
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
`
