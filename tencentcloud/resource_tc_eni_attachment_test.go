package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEniAttachment_basic(t *testing.T) {
	var (
		eniId string
		cvmId string
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEniAttachmentDestroy(&eniId),
		Steps: []resource.TestStep{
			{
				Config: testAccEniAttachmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniAttachmentExists("tencentcloud_eni_attachment.foo", &eniId, &cvmId),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_attachment.foo", "eni_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni_attachment.foo", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_eni_attachment.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEniAttachmentExists(n string, eniId, cvmId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no eni attachment id is set")
		}

		split := strings.Split(rs.Primary.ID, "+")
		*eniId, *cvmId = split[0], split[1]

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		enis, err := service.DescribeEniById(context.TODO(), []string{*eniId})
		if err != nil {
			return err
		}

		for _, e := range enis {
			if e.NetworkInterfaceId == nil {
				return errors.New("eni id is nil")
			}

			if *e.NetworkInterfaceId == *eniId {
				if e.Attachment == nil {
					return errors.New("eni attachment is nil")
				}

				if e.Attachment.InstanceId == nil {
					return errors.New("eni attach instance id is nil")
				}

				if *e.Attachment.InstanceId != *cvmId {
					return errors.New("eni attach instance id is not right")
				}

				return nil
			}
		}

		return fmt.Errorf("eni attachment not found: %s", rs.Primary.ID)
	}
}

func testAccCheckEniAttachmentDestroy(eniId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := VpcService{client: client}

		enis, err := service.DescribeEniById(context.TODO(), []string{*eniId})
		if err != nil {
			return err
		}

		if len(enis) > 0 {
			return errors.New("eni still exists")
		}

		return nil
	}
}

const testAccEniAttachmentBasic = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = "${var.availability_zone}"
  name              = "ci-test-eni-subnet"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 1
}

data "tencentcloud_image" "my_favorite_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "foo" {
  instance_name            = "ci-test-eni-attach"
  availability_zone        = "ap-guangzhou-3"
  image_id                 = "${data.tencentcloud_image.my_favorite_image.image_id}"
  instance_type            = "${data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type}"
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = "${tencentcloud_vpc.foo.id}"
  subnet_id                = "${tencentcloud_subnet.foo.id}"
}

resource "tencentcloud_eni_attachment" "foo" {
  eni_id      = "${tencentcloud_eni.foo.id}"
  instance_id = "${tencentcloud_instance.foo.id}"
}
`
