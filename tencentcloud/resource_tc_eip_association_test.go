package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudEipAssociationWithInstance(t *testing.T) {
	id := "tencentcloud_eip_association.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudEipAssociationWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "instance_id"),
					resource.TestCheckNoResourceAttr(id, "network_interface_id"),
					resource.TestCheckNoResourceAttr(id, "private_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
				),
			},
		},
	})
}

func TestAccTencentCloudEipAssociationWithNetworkInterface(t *testing.T) {
	id := "tencentcloud_eip_association.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudEipAssociationWithNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists(id),
					resource.TestCheckResourceAttrSet(id, "eip_id"),
					resource.TestCheckResourceAttrSet(id, "network_interface_id"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckNoResourceAttr(id, "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.foo", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_eip.foo", "status", "UNBIND"),
				),
			},
		},
	})
}

func testAccCheckEipAssociationDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip_association" {
			continue
		}

		associationId, err := parseEipAssociationId(rs.Primary.ID)

		if err != nil {
			return err
		}

		eip, err := vpcService.DescribeEipById(ctx, associationId.EipId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, associationId.EipId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if eip == nil {
			return nil
		}
		if eip.InstanceId != nil && *eip.InstanceId == associationId.InstanceId {
			return fmt.Errorf("eip association still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEipAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("eip association %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("eip association id is not set")
		}
		vpcService := VpcService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		associationId, err := parseEipAssociationId(rs.Primary.ID)
		eip, err := vpcService.DescribeEipById(ctx, associationId.EipId)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				eip, err = vpcService.DescribeEipById(ctx, associationId.EipId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		if eip == nil || (eip.InstanceId != nil && *eip.InstanceId != associationId.InstanceId) {
			return fmt.Errorf("eip %s is not found", associationId.EipId)
		}
		return nil
	}
}

const testAccTencentCloudEipAssociationWithInstance = defaultInstanceVariable + `
resource "tencentcloud_eip" "foo" {
  name = var.instance_name
}

resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = data.tencentcloud_availability_zones.default.zones.0.name
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = tencentcloud_eip.foo.id
  instance_id = tencentcloud_instance.foo.id
}
`

const testAccTencentCloudEipAssociationWithNetworkInterface = defaultVpcVariable + `
resource "tencentcloud_eip" "foo" {
  name = var.instance_name
}

resource "tencentcloud_eni" "foo" {
  name        = var.instance_name
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = var.instance_name
  ipv4_count  = 1
}

resource "tencentcloud_eip_association" "foo" {
  eip_id               = tencentcloud_eip.foo.id
  network_interface_id = tencentcloud_eni.foo.id
  private_ip           = tencentcloud_eni.foo.ipv4_info.0.ip
}
`
