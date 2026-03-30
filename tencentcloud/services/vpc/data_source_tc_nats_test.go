package vpc_test

import (
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNatsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudNatsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.multi_nat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.0.name", "terraform_test_nats"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.1.bandwidth", "500"),
				),
			},
		},
	})
}

const testAccTencentCloudNatsDataSourceConfig_basic = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "dev_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

data "tencentcloud_nats" "multi_nat" {
  state          = 0
  name           = tencentcloud_nat_gateway.dev_nat.name
  vpc_id         = tencentcloud_vpc.main.id
  max_concurrent = tencentcloud_nat_gateway.test_nat.max_concurrent
  bandwidth      = tencentcloud_nat_gateway.test_nat.bandwidth
}
`

func TestAccTencentCloudNatsDataSource_VerboseLevel(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				// Test verbose_level with "DETAIL" value
				Config: testAccTencentCloudNatsDataSourceConfig_verboseLevelDetail,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.verbose_level_detail"),
				),
			},
			{
				// Test verbose_level with "COMPACT" value
				Config: testAccTencentCloudNatsDataSourceConfig_verboseLevelCompact,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.verbose_level_compact"),
				),
			},
			{
				// Test verbose_level with "SIMPLE" value
				Config: testAccTencentCloudNatsDataSourceConfig_verboseLevelSimple,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.verbose_level_simple"),
				),
			},
			{
				// Test omitted verbose_level parameter (default behavior)
				Config: testAccTencentCloudNatsDataSourceConfig_defaultBehavior,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.default_behavior"),
				),
			},
		},
	})
}

func TestAccTencentCloudNatsDataSource_InvalidVerboseLevel(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				// Test invalid verbose_level value (should fail validation)
				Config:      testAccTencentCloudNatsDataSourceConfig_invalidVerboseLevel,
				ExpectError: regexp.MustCompile(`verbose_level must be one of`),
			},
		},
	})
}

const testAccTencentCloudNatsDataSourceConfig_verboseLevelDetail = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_test" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test.public_ip,
  ]
}

data "tencentcloud_nats" "verbose_level_detail" {
  verbose_level = "DETAIL"
  vpc_id        = tencentcloud_vpc.main.id
}
`

const testAccTencentCloudNatsDataSourceConfig_verboseLevelCompact = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_test" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test.public_ip,
  ]
}

data "tencentcloud_nats" "verbose_level_compact" {
  verbose_level = "COMPACT"
  vpc_id        = tencentcloud_vpc.main.id
}
`

const testAccTencentCloudNatsDataSourceConfig_verboseLevelSimple = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_test" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test.public_ip,
  ]
}

data "tencentcloud_nats" "verbose_level_simple" {
  verbose_level = "SIMPLE"
  vpc_id        = tencentcloud_vpc.main.id
}
`

const testAccTencentCloudNatsDataSourceConfig_defaultBehavior = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_test" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test.public_ip,
  ]
}

data "tencentcloud_nats" "default_behavior" {
  vpc_id = tencentcloud_vpc.main.id
}
`

const testAccTencentCloudNatsDataSourceConfig_invalidVerboseLevel = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_test" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test.public_ip,
  ]
}

data "tencentcloud_nats" "invalid_verbose_level" {
  verbose_level = "INVALID"
  vpc_id        = tencentcloud_vpc.main.id
}
`
