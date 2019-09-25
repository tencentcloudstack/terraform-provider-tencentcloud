package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEipAssociationWithInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociationWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists("tencentcloud_eip_association.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.my_eip", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.my_eip", "name", "tf_auto_test"),
					resource.TestCheckResourceAttr("tencentcloud_eip.my_eip", "status", "UNBIND"),

					resource.TestCheckResourceAttrSet("tencentcloud_eip_association.foo", "eip_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip_association.foo", "instance_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_eip_association.foo", "network_interface_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_eip_association.foo", "private_ip"),
				),
			},
		},
	})
}

/*
func TestAccTencentCloudEipAssociationWithNetworkInterface(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociationWithNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEipAssociationExists("tencentcloud_eip_association.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip.my_eip", "public_ip"),
					resource.TestCheckResourceAttr("tencentcloud_eip.my_eip", "name", "tf_auto_test"),
					resource.TestCheckResourceAttr("tencentcloud_eip.my_eip", "status", "UNBIND"),

					resource.TestCheckResourceAttrSet("tencentcloud_eip_association.foo", "eip_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_eip_association.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip_association.foo", "network_interface_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip_association.foo", "private_ip"),
				),
			},
		},
	})
}
*/

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
		if eip == nil || eip.InstanceId == nil || *eip.InstanceId != associationId.InstanceId {
			return fmt.Errorf("eip %s is not found", associationId.EipId)
		}
		return nil
	}
}

const testAccEipAssociationWithInstance = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }

  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_instance" "my_instance" {
  instance_name     = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type  = "CLOUD_SSD"
}

resource "tencentcloud_eip" "my_eip" {
  name = "tf_auto_test"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = "${tencentcloud_eip.my_eip.id}"
  instance_id = "${tencentcloud_instance.my_instance.id}"
}
`

/*
// TODO remove hard code and make network_interface_id as a resource
const testAccEipAssociationWithNetworkInterface = `
resource "tencentcloud_eip" "my_eip" {
  name = "tf_auto_test"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = "${tencentcloud_eip.my_eip.id}"
  network_interface_id = "eni-auqmq7hp"
  private_ip = "10.0.1.6"
}
`
*/
