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

func TestAccTencentCloudVpcAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
				ResourceName:      "tencentcloud_vpc_acl.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccTencentCloudVpcAclRulesResource_Update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigUpdateReduceAllRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigUpdateNoEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigUpdateNoIngress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigAllRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#0.0.0.0/0#ALL#ALL"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#0.0.0.0/0#ALL#ALL"),
				),
			},
		},
	})
}

func testAccCheckVpcACLExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

func testAccCheckVpcACLDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccVpcACLConfig = `
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

const testAccVpcACLConfigUpdate = `
data "tencentcloud_vpc_instances" "default" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_update"
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

const testAccVpcACLConfigUpdateReduceAllRule = `
data "tencentcloud_vpc_instances" "default" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_update"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
	]
} 
`

const testAccVpcACLConfigUpdateNoIngress = `
data "tencentcloud_vpc_instances" "default" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_update"
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
    	"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
} 
`
const testAccVpcACLConfigUpdateNoEgress = `
data "tencentcloud_vpc_instances" "default" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_update"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
		"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
} 
`
const testAccVpcACLConfigAllRules = `
data "tencentcloud_vpc_instances" "default" {
  is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {
  vpc_id  = data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
  name    = "test_acl_update"
  ingress = [
    "ACCEPT#0.0.0.0/0#ALL#ALL"
  ]
  egress  = [
    "ACCEPT#0.0.0.0/0#ALL#ALL"
  ]
}
`
