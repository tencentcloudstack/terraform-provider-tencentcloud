package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_eip", &resource.Sweeper{
		Name: "tencentcloud_eip",
		F:    testSweepEipInstance,
	})
}

func testSweepEipInstance(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	vpcService := VpcService{
		client: client.apiV3Conn,
	}

	instances, err := vpcService.DescribeEipByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := *v.AddressId
		print(instanceId)
		instanceName := v.AddressName

		now := time.Now()

		createTime := stringTotime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()
		if instanceName != nil {
			if strings.HasPrefix(*instanceName, keepResource) || strings.HasPrefix(*instanceName, defaultResource) {
				continue
			}
		}

		// less than 30 minute, not delete
		if needProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteEip(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudEipResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBasicWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "gateway_eip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithNewTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckNoResourceAttr("tencentcloud_eip.foo", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "tags.abc", "abc"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
				),
			},
			{
				Config: testAccEipBasicWithoutName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.bar"),
					resource.TestCheckResourceAttr("tencentcloud_eip.bar", "status", "UNBIND"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.bar", "public_ip"),
				),
			},
			{
				ResourceName:      "tencentcloud_eip.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudEipResource_anycast(t *testing.T) {
	defer func() {
		os.Setenv(PROVIDER_REGION, "")
	}()
	os.Setenv(PROVIDER_REGION, "ap-hongkong")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAnycast,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "eip_anycast"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "type", "AnycastEIP"),
				),
			},
		},
	})
}

func TestAccTencentCloudEipResource_provider(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipProvider,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "eip_provider"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "type", "EIP"),
				),
			},
		},
	})
}

func TestAccTencentCloudEipResource_bandwidth(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", "eip_bandwidth"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "type", "EIP"),
				),
			},
		},
	})
}

func TestAccTencentCloudEipResource_chargetype(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipChargeType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "internet_charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
				),
			},
			{
				ResourceName:      "tencentcloud_eip.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudEipResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipPrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipExists("tencentcloud_eip.foo"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "internet_charge_type", "BANDWIDTH_PREPAID_BY_MONTH"),
				),
			},
			{
				PreConfig: func() { //sleep 1 min after update
					time.Sleep(10 * time.Second)
				},
				ResourceName:            "tencentcloud_eip.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"prepaid_period", "auto_renew_flag"},
			},
		},
	})
}

func testAccCheckEipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("eip %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("eip id is not set")
		}

		vpcService := VpcService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		eip, err := vpcService.DescribeEipById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if eip == nil {
			return fmt.Errorf("eip id is not found")
		}
		return nil
	}
}

func testAccCheckEipDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip" {
			continue
		}

		eip, err := vpcService.DescribeEipById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if eip != nil {
			return fmt.Errorf("eip still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccEipBasicWithName = `
resource "tencentcloud_eip" "foo" {
  name = "gateway_eip"
}
`
const testAccEipBasicWithNewName = `
resource "tencentcloud_eip" "foo" {
  name = "new_name"
}
`

const testAccEipBasicWithTags = `
resource "tencentcloud_eip" "foo" {
  name = "new_name"

  tags = {
    "test" = "test"
  }
}
`

const testAccEipBasicWithNewTags = `
resource "tencentcloud_eip" "foo" {
  name = "new_name"

  tags = {
    "abc" = "abc"
  }
}
`

const testAccEipBasicWithoutName = `
resource "tencentcloud_eip" "bar" {
}
`

const testAccEipAnycast = `
resource "tencentcloud_eip" "foo" {
  name = "eip_anycast"
  type = "AnycastEIP"
  anycast_zone = "ANYCAST_ZONE_OVERSEAS"
}
`

const testAccEipProvider = `
resource "tencentcloud_eip" "foo" {
  name = "eip_provider"
  internet_service_provider = "CMCC"
}
`

const testAccEipBandwidth = `
resource "tencentcloud_eip" "foo" {
	name = "eip_bandwidth"
	internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
	internet_max_bandwidth_out = 2
  }
`

const testAccEipChargeType = `
resource "tencentcloud_eip" "foo" {
	name = "eip_charge_type"
	internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  }
`

const testAccEipPrepaid = `
resource "tencentcloud_eip" "foo" {
  name = "eip_prepaid"
  internet_charge_type = "BANDWIDTH_PREPAID_BY_MONTH"
  prepaid_period = 6
  auto_renew_flag = 1
  internet_max_bandwidth_out = 2
}
`
