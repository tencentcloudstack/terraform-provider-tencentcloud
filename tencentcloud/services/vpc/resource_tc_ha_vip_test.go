package vpc_test

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
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func init() {
	resource.AddTestSweepers("tencentcloud_ha_vip", &resource.Sweeper{
		Name: "tencentcloud_ha_vip",
		F:    testSweepHaVipInstance,
	})
}

func testSweepHaVipInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	vpcService := svcvpc.NewVpcService(client.GetAPIV3Conn())

	instances, err := vpcService.DescribeHaVipByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := *v.HaVipId
		instanceName := *v.HaVipName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteHaVip(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudHaVip_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ha_vip" {
			continue
		}
		request := vpc.NewDescribeHaVipsRequest()
		request.HaVipIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeHaVips(request)
			if e != nil {
				ee, ok := e.(*sdkErrors.TencentCloudSDKError)
				if !ok {
					return tccommon.RetryError(errors.WithStack(e))
				}
				if ee.Code == svcvpc.VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e)
					return resource.NonRetryableError(e)
				} else {
					return tccommon.RetryError(errors.WithStack(e))
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
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("HA VIP instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HA VIP id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeHaVipsRequest()
		request.HaVipIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeHaVips(request)
			if e != nil {
				return tccommon.RetryError(errors.WithStack(e))
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

const testAccHaVipConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
`
const testAccHaVipConfigUpdate = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_update"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
`

const testAccHaVipConfigAssigned = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
  vip       = "172.16.0.137"
}
`
