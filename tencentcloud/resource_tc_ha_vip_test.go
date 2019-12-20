package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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
				ee, ok := e.(*sdkErrors.TencentCloudSDKError)
				if !ok {
					return retryError(errors.WithStack(e))
				}
				if ee.Code == VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e)
					return resource.NonRetryableError(e)
				} else {
					return retryError(errors.WithStack(e))
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
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
				return retryError(errors.WithStack(e))
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%s\n", logId, err)
			return err
		}
		if len(response.Response.HaVipSet) != 1 {
			return fmt.Errorf("HA VIP id is not found")
		}
		return nil
	}
}

const testAccHaVipConfig = defaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = "${var.vpc_id}"
  subnet_id = "${var.subnet_id}"
}
`
const testAccHaVipConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_update"
  vpc_id    = "${var.vpc_id}"
  subnet_id = "${var.subnet_id}"
}
`

const testAccHaVipConfigAssigned = defaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = "${var.vpc_id}"
  subnet_id = "${var.subnet_id}"
  vip       = "172.16.0.255"
}
`
