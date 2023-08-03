package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcIpv6EniAddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		//CheckDestroy: testAccCheckVpcIpv6EniAddressDestroy,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6EniAddress,
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckBandwidthPackageAttachmentExists("tencentcloud_vpc_ipv6_eni_address"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "network_interface_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.address"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.primary"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.address_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.description"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.is_wan_ip_blocked"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.state"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//
//func testAccCheckVpcIpv6EniAddressDestroy(s *terraform.State) error {
//	logId := getLogId(contextNil)
//	ctx := context.WithValue(context.TODO(), logIdKey, logId)
//	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type != "tencentcloud_vpc_bandwidth_package_attachment" {
//			continue
//		}
//		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
//		if len(idSplit) != 2 {
//			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
//		}
//		bandwidthPackageId := idSplit[0]
//		resourceId := idSplit[1]
//
//		bandwidthPackageResources, err := service.DescribeVpcBandwidthPackageAttachment(ctx, bandwidthPackageId, resourceId)
//		if err != nil {
//			log.Printf("[CRITAL]%s read VPN connection failed, reason:%s\n", logId, err.Error())
//			ee, ok := err.(*errors.TencentCloudSDKError)
//			if !ok {
//				return err
//			}
//			fmt.Print(ee)
//			if ee.Code == "InvalidParameterValue.BandwidthPackageNotFound" {
//				return nil
//			} else {
//				return err
//			}
//		} else {
//			if bandwidthPackageResources != nil {
//				return fmt.Errorf("vpc bandwidthPackageResources %s still exists", rs.Primary.ID)
//			}
//		}
//	}
//	return nil
//}
//
//func testAccCheckVpcIpv6EniAddressExists(r string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		logId := getLogId(contextNil)
//		ctx := context.WithValue(context.TODO(), logIdKey, logId)
//
//		rs, ok := s.RootModule().Resources[r]
//		if !ok {
//			return fmt.Errorf("resource %s is not found", r)
//		}
//
//		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
//		if len(idSplit) != 2 {
//			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
//		}
//		bandwidthPackageId := idSplit[0]
//		resourceId := idSplit[1]
//
//		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
//		bandwidthPackageResources, err := service.DescribeVpcBandwidthPackageAttachment(ctx, bandwidthPackageId, resourceId)
//		if bandwidthPackageResources == nil {
//			return fmt.Errorf("vpc bandwidthPackageResources %s is not found", rs.Primary.ID)
//		}
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//}

const testAccVpcIpv6EniAddress = `

resource "tencentcloud_vpc_ipv6_eni_address" "ipv6_eni_address" {
  vpc_id = ""
  network_interface_id = ""
  ipv6_addresses {
		address = ""
		primary = False
		address_id = ""
		description = ""
		is_wan_ip_blocked = 
		state = ""

  }
}

`
