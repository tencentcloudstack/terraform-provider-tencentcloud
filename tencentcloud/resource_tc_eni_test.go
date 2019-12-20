package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
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
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
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
			{
				Config: testAccEniUpdateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "tags.test", "test"),
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
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
				),
			},
			{
				Config: testAccEniUpdateCountAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "30"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "30"),
				),
			},
			{
				Config: testAccEniUpdateCountSub,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_count", "20"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "20"),
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
					resource.TestCheckNoResourceAttr("tencentcloud_eni.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "mac"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "primary", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eni.foo", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
				),
			},
			{
				Config: testAccEniManuallyUpdatePrimaryDesc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "1"),
				),
			},
			{
				Config: testAccEniManuallyUpdateAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4s.#", "30"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "30"),
				),
			},
			{
				Config: testAccEniManuallyUpdateSub,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEniExists("tencentcloud_eni.foo", &eniId),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4s.#", "15"),
					resource.TestCheckResourceAttr("tencentcloud_eni.foo", "ipv4_info.#", "15"),
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

		enis, err := service.DescribeEniById(context.TODO(), []string{rs.Primary.ID})
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

		enis, err := service.DescribeEniById(context.TODO(), []string{*id})
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
  cidr_block        = "10.0.0.0/16"
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

const testAccEniUpdateTags = testAccEniVpc + `

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

  tags = {
    "test" = "test"
  }
}
`

const testAccEniUpdateCountAdd = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 30
}
`

const testAccEniUpdateCountSub = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 20
}
`

const testAccEniManually = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  
  ipv4s {
    ip      = "10.0.0.10"
    primary = true
    description = "desc"
  }
}
`

const testAccEniManuallyUpdatePrimaryDesc = testAccEniVpc + `

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  
  ipv4s {
    ip          = "10.0.0.10"
    primary     = true
    description = ""
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
    ip          = "10.0.0.10"
    primary     = true
    description = ""
  }

  ipv4s {
    ip      = "10.0.0.11"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.12"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.13"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.14"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.15"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.16"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.17"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.18"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.19"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.21"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.22"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.23"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.24"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.25"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.26"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.27"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.28"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.29"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.30"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.31"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.32"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.33"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.34"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.35"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.36"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.37"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.38"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.39"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.40"
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
    ip          = "10.0.0.10"
    primary     = true
    description = "" // set empty desc to test if SDK can set private IP desc empty or not
  }

  ipv4s {
    ip      = "10.0.0.11"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.12"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.13"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.14"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.15"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.16"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.17"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.18"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.19"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.21"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.22"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.23"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.24"
    primary = false
  }

  ipv4s {
    ip      = "10.0.0.25"
    primary = false
  }
}
`
