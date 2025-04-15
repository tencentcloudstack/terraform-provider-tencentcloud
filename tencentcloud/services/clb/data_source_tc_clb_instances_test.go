package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbInstancesDataSource_internal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstancesDataSource_internal,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_name", "tf-clb-data-internal"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.network_type", "INTERNAL"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_vips.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.tags.test", "tf"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.numerical_vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstancesDataSource_open(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstancesDataSource_open,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.network_type", "OPEN"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_vips.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.target_region_info_vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.tags.test", "tf"),
				),
			},
		},
	})
}

const testAccClbInstancesDataSource_internal = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua-ci-temp-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_clb_instance" "clb" {
  network_type = "INTERNAL"
  clb_name     = "tf-clb-data-internal"
  vpc_id       = tencentcloud_vpc.foo.id
  subnet_id    = tencentcloud_subnet.subnet.id
  project_id   = 0

  tags = {
    test = "tf"
  }
}

data "tencentcloud_clb_instances" "clbs" {
  clb_id = tencentcloud_clb_instance.clb.id
}
`

const testAccClbInstancesDataSource_open = `
resource "tencentcloud_security_group" "foo" {
  name = "clb-instance-datasource-sg"
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb" {
  network_type              = "OPEN"
  clb_name                  = "tf-clb-data-open"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]

  tags = {
    test = "tf"
  }
}

data "tencentcloud_clb_instances" "clbs" {
  clb_id = tencentcloud_clb_instance.clb.id
}
`
