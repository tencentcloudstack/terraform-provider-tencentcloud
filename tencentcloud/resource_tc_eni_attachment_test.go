package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudEniAttachmentBasic(t *testing.T) {
	t.Parallel()
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
		if *eniId == "" {
			return nil
		}
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

const testAccEniAttachmentBasic = instanceCommonTestCase + `
resource "tencentcloud_eni" "foo" {
  name        = var.instance_name
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = var.instance_name
  ipv4_count  = 1
}

resource "tencentcloud_eni_attachment" "foo" {
  eni_id      = tencentcloud_eni.foo.id
  instance_id = tencentcloud_instance.default.id
}
`
