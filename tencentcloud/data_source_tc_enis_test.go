package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccDataSourceTencentCloudEnis_basic -v
func TestAccDataSourceTencentCloudEnis_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudEnisBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.name", "ci-test-eni"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.description", "eni desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.subnet_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.primary", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.ipv4s.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.foo", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.ipv4s.0.primary", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.foo", "enis.0.ipv4s.0.description", "eni desc"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudEnis_filter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudEnisFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.vpc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.name", "ci-test-eni"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.description", "eni desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.subnet_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.primary", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.tags.test", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.ipv4s.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.vpc", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.ipv4s.0.primary", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.vpc", "enis.0.ipv4s.0.description", "eni desc"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.subnet"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "security_group"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.name", "ci-test-eni"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.description", "eni desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.subnet_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.primary", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.tags.test", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.ipv4s.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.subnet", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.ipv4s.0.primary", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.subnet", "enis.0.ipv4s.0.description", "eni desc"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.name", "name", "ci-test-eni"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.name", "enis.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.name", "enis.0.name", "ci-test-eni"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.primary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.name", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.create_time"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.name", "enis.0.ipv4s.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.name", "enis.0.ipv4s.0.primary"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "description"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.description", "enis.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.description", "enis.0.description", "eni desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.primary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.description", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.create_time"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.description", "enis.0.ipv4s.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.description", "enis.0.ipv4s.0.primary"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.ipv4"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "ipv4"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.ipv4", "enis.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.primary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.ipv4", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.create_time"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.ipv4", "enis.0.ipv4s.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.ipv4", "enis.0.ipv4s.0.primary"),

					testAccCheckTencentCloudDataSourceID("data.tencentcloud_enis.tags"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.tags", "enis.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.primary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.mac"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.tags", "enis.0.state", ENI_STATE_AVAILABLE),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_enis.tags", "enis.0.tags.test", "test"),
					resource.TestMatchResourceAttr("data.tencentcloud_enis.tags", "enis.0.ipv4s.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.ipv4s.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_enis.tags", "enis.0.ipv4s.0.primary"),
				),
			},
		},
	})
}

const testAccEnisVpc = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = var.availability_zone
  name              = "ci-test-eni-subnet"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}
`

const TestAccDataSourceTencentCloudEnisBasic = testAccEnisVpc + `

resource "tencentcloud_security_group" "foo" {
  name = "test-ci-eni-sg1"
}

resource "tencentcloud_eni" "foo" {
  name            = "ci-test-eni"
  vpc_id          = tencentcloud_vpc.foo.id
  subnet_id       = tencentcloud_subnet.foo.id
  description     = "eni desc"
  security_groups = [tencentcloud_security_group.foo.id]
  ipv4_count      = 1
}

data "tencentcloud_enis" "foo" {
  ids = [tencentcloud_eni.foo.id]
}
`

const TestAccDataSourceTencentCloudEnisFilter = testAccEnisVpc + `

resource "tencentcloud_security_group" "foo" {
  name = "test-ci-eni-sg1"
}

resource "tencentcloud_eni" "foo" {
  name            = "ci-test-eni"
  vpc_id          = tencentcloud_vpc.foo.id
  subnet_id       = tencentcloud_subnet.foo.id
  description     = "eni desc"
  security_groups = [tencentcloud_security_group.foo.id]
  ipv4_count      = 1

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_enis" "vpc" {
  vpc_id = tencentcloud_eni.foo.vpc_id
}

data "tencentcloud_enis" "subnet" {
  subnet_id      = tencentcloud_eni.foo.subnet_id
  security_group = tencentcloud_security_group.foo.id
}

data "tencentcloud_enis" "name" {
  name = tencentcloud_eni.foo.name
}

data "tencentcloud_enis" "description" {
  description = tencentcloud_eni.foo.description
}

data "tencentcloud_enis" "ipv4" {
  ipv4 = tencentcloud_eni.foo.ipv4_info.0.ip
}

data "tencentcloud_enis" "tags" {
  tags = tencentcloud_eni.foo.tags
}
`
