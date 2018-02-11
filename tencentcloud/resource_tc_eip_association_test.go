package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEipAssociationWithInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociationWithInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_eip_association.foo"),
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

func TestAccTencentCloudEipAssociationWithNetworkInterface(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEipAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEipAssociationWithNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_eip_association.foo"),
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

func testAccCheckEipAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*TencentCloudClient).commonConn
	cvmConn := testAccProvider.Meta().(*TencentCloudClient).cvmConn

	var assId string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eip_association" {
			continue
		}
		assId = rs.Primary.ID
	}

	association, err := parseAssociationId(assId)
	if err != nil {
		return err
	}
	eipId := association.eipId

	// make sure eip is deleted
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, _, err := findEipById(cvmConn, eipId)
		if err == nil {
			err = fmt.Errorf("eip can still be found after deleted")
			return resource.RetryableError(err)
		}
		if err == errEIPNotFound {
			return nil
		}
		return resource.RetryableError(errEIPStillDeleting)
	})
	if err != nil {
		return err
	}

	// make sure instance is deleted
	instanceId := association.instanceId
	if len(instanceId) > 0 {
		instanceIds := []string{instanceId}

		_, err = waitInstanceReachTargetStatus(client, instanceIds, "STOPPED")
		if err != nil {
			if err.Error() == instanceNotFoundErrorMsg(instanceIds) {
				return nil
			}
			return err
		}
		return nil
	}

	return nil
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

data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_instance" "my_instance" {
  instance_name     = "terraform_automation_test_kuruk"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
}

resource "tencentcloud_eip" "my_eip" {
  name = "tf_auto_test"
}

resource "tencentcloud_eip_association" "foo" {
  eip_id      = "${tencentcloud_eip.my_eip.id}"
  instance_id = "${tencentcloud_instance.my_instance.id}"
}
`

// TODO remove hard code and make netwrok_interface_id as a resource
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
