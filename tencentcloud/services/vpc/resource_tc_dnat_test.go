package vpc_test

import (
	"fmt"
	"log"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudNatGatewayTransRule_basic(t *testing.T) {
	t.Parallel()
	var dnatId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("tencentcloud_dnat.dev_dnat", &dnatId),
				),
			},
			{
				Config: testAccDnatConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("tencentcloud_dnat.dev_dnat", &dnatId),
				),
			},
		},
	})
}

func testAccCheckDnatExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DNAT instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DNAT id is not set")
		}
		_, params, e := svcvpc.ParseDnatId(rs.Primary.ID)
		if e != nil {
			return fmt.Errorf("[CRITAL]parse DNAT id fail, reason[%s]\n", e.Error())
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
		request.Filters = make([]*vpc.Filter, 0, len(params))
		for k, v := range params {
			filter := &vpc.Filter{
				Name:   helper.String(k),
				Values: []*string{helper.String(v)},
			}
			request.Filters = append(request.Filters, filter)
		}
		var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) != 1 {
			return fmt.Errorf("DNAT is not exists")
		}
		return nil
	}
}

func testAccCheckDnatDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dnat" {
			continue
		}
		_, params, e := svcvpc.ParseDnatId(rs.Primary.ID)
		if e != nil {
			log.Printf("[CRITAL]parse DNAT id fail, reason[%s]\n", e.Error())
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeNatGatewayDestinationIpPortTranslationNatRulesRequest()
		request.Filters = make([]*vpc.Filter, 0, len(params))
		for k, v := range params {
			filter := &vpc.Filter{
				Name:   helper.String(k),
				Values: []*string{helper.String(v)},
			}
			request.Filters = append(request.Filters, filter)
		}
		var response *vpc.DescribeNatGatewayDestinationIpPortTranslationNatRulesResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGatewayDestinationIpPortTranslationNatRules(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read DNAT failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewayDestinationIpPortTranslationNatRuleSet) != 0 {
			return fmt.Errorf("DNAT is still exists")
		}
	}
	return nil

}

const testAccDnatConfig = tcacctest.InstanceCommonTestCase + `
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = var.instance_name
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = var.instance_name
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = var.cvm_vpc_id
  name           = var.instance_name
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = tencentcloud_nat_gateway.my_nat.vpc_id
  nat_id       = tencentcloud_nat_gateway.my_nat.id
  protocol     = "TCP"
  elastic_ip   = tencentcloud_eip.eip_dev_dnat.public_ip
  elastic_port = "80"
  private_ip   = tencentcloud_instance.default.private_ip
  private_port = "9001"
}
`

const testAccDnatConfigUpdate = tcacctest.InstanceCommonTestCase + `
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = var.instance_name
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = var.instance_name
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = var.cvm_vpc_id
  name           = var.instance_name
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = tencentcloud_nat_gateway.my_nat.vpc_id
  nat_id       = tencentcloud_nat_gateway.my_nat.id
  protocol     = "TCP"
  elastic_ip   = tencentcloud_eip.eip_dev_dnat.public_ip
  elastic_port = "80"
  private_ip   = tencentcloud_instance.default.private_ip
  private_port = "9001"
  description  = var.instance_name
}
`
