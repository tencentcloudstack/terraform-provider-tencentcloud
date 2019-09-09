package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudEni_basic(t *testing.T) {
	var eniId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEniDestroy(&eniId),
		Steps: []resource.TestStep{
			{
				Config: testAccEniBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "name", "ci-test-eni"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "description", "eni desc"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "security_groups.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "ipv6_count"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv6_info.#", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_eni.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudEni_updateAttr(t *testing.T) {
	var eniId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEniDestroy(&eniId),
		Steps: []resource.TestStep{
			{
				Config: testAccEniBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "name", "ci-test-eni"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "description", "eni desc"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "security_groups.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "ipv6_count"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv6_info.#", "0"),
				),
			},
			{
				Config: testAccEniUpdateAttr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "name", "ci-test-eni-new"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "description", "eni desc new"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "security_groups.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudEni_updateCount(t *testing.T) {
	var eniId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEniDestroy(&eniId),
		Steps: []resource.TestStep{
			{
				Config: testAccEniBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "name", "ci-test-eni"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "description", "eni desc"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "security_groups.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "ipv6_count"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv6_info.#", "0"),
				),
			},
			{
				Config: testAccEniUpdateCountAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "3"),
				),
			},
			{
				Config: testAccEniUpdateCountSub,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudEni_updateManually(t *testing.T) {
	var eniId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEniDestroy(&eniId),
		Steps: []resource.TestStep{
			{
				Config: testAccEniManually,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "name", "ci-test-eni"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "description", "eni desc"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "security_groups.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4s.#", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "ipv4_count"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "ipv6_count"),
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv6_info.#", "0"),
				),
			},
			{
				Config: testAccEniManuallyUpdateAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4s.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "3"),
				),
			},
			{
				Config: testAccEniManuallyUpdateSub,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4s.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "2"),
				),
			},
		},
	})
}

func testAccCheckEniExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no eni id is set")
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		enis, err := service.DescribeEniById(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		for _, e := range enis {
			if e.NetworkInterfaceId == nil {
				return errors.New("eni id is nil")
			}

			if *e.NetworkInterfaceId == rs.Primary.ID {
				*id = rs.Primary.ID
				return nil
			}
		}

		return fmt.Errorf("eni not found: %s", rs.Primary.ID)
	}
}

func testAccCheckEniDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := VpcService{client: client}

		enis, err := service.DescribeEniById(context.TODO(), *id)
		if err != nil {
			return err
		}

		if len(enis) > 0 {
			return errors.New("eni still exists")
		}

		return nil
	}
}

const testAccEniVpc = `
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
`

const testAccEniBasic = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 1
}
`

const testAccEniUpdateCountAdd = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 3
}
`

const testAccEniUpdateCountSub = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 2
}
`

const testAccEniManually = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  
  ipv4s {
    ip = "10.0.20.11"
    primary = true
  }
}
`

const testAccEniManuallyUpdateAdd = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  
  ipv4s {
    ip = "10.0.20.11"
    primary = true
  }

  ipv4s {
    ip = "10.0.20.12"
    primary = false
  }

  ipv4s {
    ip = "10.0.20.13"
    primary = false
  }
}
`

const testAccEniManuallyUpdateSub = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  
  ipv4s {
    ip = "10.0.20.11"
    primary = true
  }

  ipv4s {
    ip = "10.0.20.12"
    primary = false
  }
}
`

const testAccEniUpdateAttr = testAccEniVpc + `

resource "tencentcloud_security_group" "foo" {
  name = "test-ci-eni-sg1"
}

resource "tencentcloud_security_group" "bar" {
  name = "test-ci-eni-sg2"
}

resource "tencentcloud_eni" "foo" {
  name            = "ci-test-eni-new"
  vpc_id          = "${tencentcloud_vpc.foo.id}"
  subnet_id       = "${tencentcloud_subnet.foo.id}"
  description     = "eni desc new"
  security_groups = ["${tencentcloud_security_group.foo.id}", "${tencentcloud_security_group.bar.id}"]
  ipv4_count      = 1
}
`
