package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstanceBasic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_instance.cdwch_instance", "id")),
			},
			{
				ResourceName:            "tencentcloud_clickhouse_instance.cdwch_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span"},
			},
		},
	})
}

func TestAccTencentCloudClickhouseInstanceResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstancePrepaid,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_instance.cdwch_instance_prepaid", "id")),
			},
			{
				ResourceName:            "tencentcloud_clickhouse_instance.cdwch_instance_prepaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span"},
			},
		},
	})
}

func TestAccTencentCloudClickhouseInstanceResource_withSecondaryZoneInfo(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstanceWithSecondaryZoneInfo,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_instance.cdwch_instance_secondary_zone", "id")),
			},
			{
				ResourceName:            "tencentcloud_clickhouse_instance.cdwch_instance_secondary_zone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span"},
			},
		},
	})
}

const testAccClickhouseInstanceBasic = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "cdwch_instance" {
  zone            = var.availability_zone
  ha_flag         = true
  vpc_id          = tencentcloud_vpc.vpc.id
  subnet_id       = tencentcloud_subnet.subnet.id
  product_version = "21.8.12.29"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }
  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
  charge_type   = "POSTPAID_BY_HOUR"
  instance_name = "tf-test-clickhouse"
}
`

const testAccClickhouseInstancePrepaid = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "cdwch_instance_prepaid" {
  zone            = var.availability_zone
  ha_flag         = true
  vpc_id          = tencentcloud_vpc.vpc.id
  subnet_id       = tencentcloud_subnet.subnet.id
  product_version = "21.8.12.29"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }
  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
  charge_type   = "PREPAID"
  renew_flag    = 1
  time_span     = 1
  instance_name = "tf-test-clickhouse-prepaid"
}
`

const testAccClickhouseInstanceWithSecondaryZoneInfo = `
resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-6"
  is_multicast      = false
}
resource "tencentcloud_subnet" "subnet1" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet1"
  cidr_block        = "10.0.2.0/24"
  availability_zone = "ap-guangzhou-3"
  is_multicast      = false
}
resource "tencentcloud_subnet" "subnet2" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet2"
  cidr_block        = "10.0.3.0/24"
  availability_zone = "ap-guangzhou-4"
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "cdwch_instance_secondary_zone" {
  zone            = "ap-guangzhou-6"
  ha_flag         = true
  vpc_id          = tencentcloud_vpc.vpc.id
  subnet_id       = tencentcloud_subnet.subnet.id
  product_version = "21.8.12.29"
  data_spec {
    spec_name = "S_4_16_S"
    count     = 4
    disk_size = 300
  }
  common_spec {
    spec_name = "S_4_16_S"
    count     = 3
    disk_size = 300
  }
  charge_type   = "POSTPAID_BY_HOUR"
  instance_name = "tf-test"
  secondary_zone_info {
    secondary_zone   = "ap-guangzhou-4"
    secondary_subnet = tencentcloud_subnet.subnet2.id
  }
  secondary_zone_info {
    secondary_zone   = "ap-guangzhou-3"
    secondary_subnet = tencentcloud_subnet.subnet1.id
  }
  ha_zk = true
}
`
