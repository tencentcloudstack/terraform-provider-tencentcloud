package vpn_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func init() {
	resource.AddTestSweepers("tencentcloud_vpn_customer_gateway", &resource.Sweeper{
		Name: "tencentcloud_vpn_customer_gateway",
		F:    testSweepVpnCustomerGateway,
	})
}

func testSweepVpnCustomerGateway(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	vpcService := svcvpc.NewVpcService(client)

	instances, err := vpcService.DescribeCustomerGatewayByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.CustomerGatewayName, *v.CustomerGatewayName) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.CustomerGatewayId,
				Name: *v.CustomerGatewayName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.CustomerGatewayId,
			Name:       *v.CustomerGatewayName,
			CreateTime: *v.CustomerGatewayName,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateCustomerGateway")

	for _, v := range instances {
		customerGwId := *v.CustomerGatewayId
		customerName := *v.CustomerGatewayName

		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(customerName, tcacctest.KeepResource) || strings.HasPrefix(customerName, tcacctest.DefaultResource) {
			continue
		}
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteCustomerGateway(ctx, customerGwId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", customerGwId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudVpnCustomerGateway_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnCustomerGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("tencentcloud_vpn_customer_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "public_ip_address", "1.1.1.2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "tags.test", "tf"),
				),
			},
			{
				Config: testAccVpnCustomerGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnCustomerGatewayExists("tencentcloud_vpn_customer_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_customer_gateway.my_cgw", "public_ip_address", "1.1.1.2"),
				),
			},
		},
	})
}

func testAccCheckVpnCustomerGatewayDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_customer_gateway" {
			continue
		}
		request := vpc.NewDescribeCustomerGatewaysRequest()
		request.CustomerGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeCustomerGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeCustomerGateways(request)
			if e != nil {
				ee, ok := e.(*errors.TencentCloudSDKError)
				if !ok {
					return tccommon.RetryError(e)
				}
				if ee.Code == svcvpc.VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return tccommon.RetryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
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
			if len(response.Response.CustomerGatewaySet) != 0 {
				return fmt.Errorf("VPN customer gateway id is still exists")
			}
		}

	}
	return nil
}

func testAccCheckVpnCustomerGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN customer gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN customer gateway id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeCustomerGatewaysRequest()
		request.CustomerGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeCustomerGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeCustomerGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.CustomerGatewaySet) != 1 {
			return fmt.Errorf("VPN customer gateway id is not found")
		}
		return nil
	}
}

const testAccVpnCustomerGatewayConfig = `
resource "tencentcloud_vpn_customer_gateway" "my_cgw" {
  name              = "terraform_test"
  public_ip_address = "1.5.5.5"

  tags = {
    test = "tf"
  }
}
`
const testAccVpnCustomerGatewayConfigUpdate = `
resource "tencentcloud_vpn_customer_gateway" "my_cgw" {
  name              = "terraform_update"
  public_ip_address = "1.5.5.5"

  tags = {
    test = "test"
  }
}
`
