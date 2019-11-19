package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	errors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func TestAccTencentCloudHaVip_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHaVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("tencentcloud_ha_vip.havip"),
					resource.TestCheckResourceAttr("tencentcloud_ha_vip.havip", "name", "terraform_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "create_time"),
				),
			},
			{
				Config: testAccHaVipConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("tencentcloud_ha_vip.havip"),
					resource.TestCheckResourceAttr("tencentcloud_ha_vip.havip", "name", "terraform_update"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudHaVip_assigned(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHaVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipConfigAssigned,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHaVipExists("tencentcloud_ha_vip.havip"),
					resource.TestCheckResourceAttr("tencentcloud_ha_vip.havip", "name", "terraform_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_ha_vip.havip", "create_time"),
				),
			},
		},
	})
}

func testAccCheckHaVipDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ha_vip" {
			continue
		}
		request := vpc.NewDescribeHaVipsRequest()
		request.HaVipIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeHaVips(request)
			if e != nil {
				ee, ok := e.(*errors.TencentCloudSDKError)
				if !ok {
					return retryError(e)
				}
				if ee.Code == VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return retryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*errors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound" {
				return nil
			} else {
				return err
			}
		} else {
			if len(response.Response.HaVipSet) != 0 {
				return fmt.Errorf("HA VIP id is still exists")
			}
		}
	}
	return nil
}

func testAccCheckHaVipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("HA VIP instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HA VIP id is not set")
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeHaVipsRequest()
		request.HaVipIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeHaVips(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.HaVipSet) != 1 {
			return fmt.Errorf("HA VIP id is not found")
		}
		return nil
	}
}

const testAccHaVipConfig = `
# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
data "tencentcloud_vpc_subnets" "subnet" {
  name              = "Default-Subnet-Terraform-勿删"
}
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
  subnet_id = "${data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id}"
}
`
const testAccHaVipConfigUpdate = `
# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
data "tencentcloud_vpc_subnets" "subnet" {
  name              = "Default-Subnet-Terraform-勿删"
}
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_update"
  vpc_id    = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
  subnet_id = "${data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id}"
}
`

const testAccHaVipConfigAssigned = `
# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
data "tencentcloud_vpc_subnets" "subnet" {
  name              = "Default-Subnet-Terraform-勿删"
}
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
  subnet_id = "${data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id}"
  vip       = "172.16.16.255"
}
`
