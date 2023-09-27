package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudVpcBandwidthPackageResource_basic -v
func TestAccTencentCloudVpcBandwidthPackageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwidthPackageExists("tencentcloud_vpc_bandwidth_package.bandwidth_package"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "bandwidth_package_name", "iac-test-001"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "charge_type", "TOP5_POSTPAID_BY_MONTH"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "network_type", "BGP"),
				),
			},
			{
				Config: testAccVpcBandwidthPackageUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwidthPackageExists("tencentcloud_vpc_bandwidth_package.bandwidth_package"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "bandwidth_package_name", "iac-test-002"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "charge_type", "TOP5_POSTPAID_BY_MONTH"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "network_type", "BGP"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_bandwidth_package.bandwidth_package",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudVpcBandwidthPackageResource_Egress(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackageEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwidthPackageExists("tencentcloud_vpc_bandwidth_package.bandwidth_package"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "bandwidth_package_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "charge_type", "ENHANCED95_POSTPAID_BY_MONTH"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "network_type", "SINGLEISP_CMCC"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "internet_max_bandwidth", "400"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_bandwidth_package.bandwidth_package", "egress", "center_egress2"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_bandwidth_package.bandwidth_package",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBandwidthPackageDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc_bandwidth_package" {
			continue
		}

		bandwidthPackage, err := service.DescribeVpcBandwidthPackage(ctx, rs.Primary.ID)
		if bandwidthPackage != nil {
			return fmt.Errorf("vpc bandwidthPackage %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckBandwidthPackageExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		bandwidthPackage, err := service.DescribeVpcBandwidthPackage(ctx, rs.Primary.ID)
		if bandwidthPackage == nil {
			return fmt.Errorf("vpc bandwidthPackage %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccVpcBandwidthPackage = `

resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type            = "BGP"
  charge_type             = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name  = "iac-test-001"
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccVpcBandwidthPackageUpdate = `

resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type            = "BGP"
  charge_type             = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name  = "iac-test-002"
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccVpcBandwidthPackageEgress = `

resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type           = "SINGLEISP_CMCC"
  charge_type            = "ENHANCED95_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  internet_max_bandwidth = 400
  egress                 = "center_egress2"
  tags                   = {
    "createdBy" = "terraform"
  }
}

`
