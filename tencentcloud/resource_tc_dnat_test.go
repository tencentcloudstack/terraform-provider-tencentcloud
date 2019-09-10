package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func TestAccTencentCloudDnat_basic(t *testing.T) {
	var dnatId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("tencentcloud_dnat.dev_dnat", &dnatId),
				),
			},
		},
	})
}

func testAccCheckDnatExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("dnat instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dnat id is not set")
		}
		_, params, e := parseDnatId(rs.Primary.ID)
		if e != nil {
			log.Printf("[CRITAL]parse dnat id fail, reason[%s]\n", e.Error())
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
		request.Filters = make([]*vpc.Filter, 0, len(params))
		for k, v := range params {
			filter := &vpc.Filter{
				Name:   stringToPointer(k),
				Values: []*string{stringToPointer(v)},
			}
			request.Filters = append(request.Filters, filter)
		}
		var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) != 1 {
			return fmt.Errorf("dnat is not exists")
		}
		return nil
	}
}

func testAccCheckDnatDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dnat" {
			continue
		}
		_, params, e := parseDnatId(rs.Primary.ID)
		if e != nil {
			log.Printf("[CRITAL]parse dnat id fail, reason[%s]\n", e.Error())
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
		request.Filters = make([]*vpc.Filter, 0, len(params))
		for k, v := range params {
			filter := &vpc.Filter{
				Name:   stringToPointer(k),
				Values: []*string{stringToPointer(v)},
			}
			request.Filters = append(request.Filters, filter)
		}
		var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) != 0 {
			return fmt.Errorf("dnat is still exists")
		}
	}
	return nil

}

const testAccDnatConfig = `
data "tencentcloud_availability_zones" "my_favorate_zones" {
	name = "ap-guangzhou-3"
  }
  
  data "tencentcloud_image" "my_favorate_image" {
	filter {
	  name   = "image-type"
	  values = ["PUBLIC_IMAGE"]
	}
  }
  
  # Create VPC and Subnet
  data "tencentcloud_vpc_instances" "foo" {
	name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "foo" {
	subnet_id = "subnet-pqfek0t8"
  }
  
  # Create EIP 
  resource "tencentcloud_eip" "eip_dev_dnat" {
	name = "terraform_test"
  }
  resource "tencentcloud_eip" "eip_test_dnat" {
	name = "terraform_test"
  }
  
  # Create NAT Gateway
  resource "tencentcloud_nat_gateway" "my_nat" {
	vpc_id           = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
	name             = "terraform test"
	max_concurrent   = 3000000
	bandwidth        = 500
	assigned_eip_set = [
	  "${tencentcloud_eip.eip_dev_dnat.public_ip}",
	  "${tencentcloud_eip.eip_test_dnat.public_ip}",
	]
  }
  
  # Create CVM
  resource "tencentcloud_instance" "foo" {
	availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
	image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
	vpc_id           = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
	subnet_id         = "${data.tencentcloud_vpc_subnets.foo.instance_list.0.subnet_id}"
	system_disk_type  = "CLOUD_SSD"
  }
  
  # Add DNAT Entry
  resource "tencentcloud_dnat" "dev_dnat" {
	vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
	nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
	protocol     = "TCP"
	elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
	elastic_port = "80"
	private_ip   = "${tencentcloud_instance.foo.private_ip}"
	private_port = "9001"
  }
`
