package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTestingClbInstanceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "clb_name", BasicClbName),
				),
			},
			{
				ResourceName:            "tencentcloud_clb_instance.clb_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dynamic_vip"},
			},
		},
	})
}

func TestAccTencentCloudTestingClbInstanceResource_open(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbInstance_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", OpenClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info_vpc_id"),
				),
			},
			{
				Config: testAccTestingClbInstance_update_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", OpenClbNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info_vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbTestingInstanceResource_multiple_instance(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbInstance_Multi_Instance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.multiple_instance"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "clb_name", MultiClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "master_zone_id", "100004"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "slave_zone_id", "100003"),
				),
			},
			{
				Config: testAccTestingClbInstance_Multi_Instance_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.multiple_instance"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "network_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "clb_name", MultiClbName),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "master_zone_id", "100004"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.multiple_instance", "slave_zone_id", "100003"),
				),
			},
		},
	})
}

const testAccTestingClbInstance_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "` + BasicClbName + `"
}
`

const testAccTestingClbInstance_open = `
resource "tencentcloud_security_group" "foo" {
  name = "keep-ci-temp-test-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "` + OpenClbName + `"
  project_id                = 0
  vpc_id                    = tencentcloud_vpc.foo.id
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]
}
`
const testAccTestingClbInstance_update_open = `
resource "tencentcloud_security_group" "foo" {
  name = "clb-instance-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "clb-instance-open-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
  network_type              = "OPEN"
  clb_name                  = "` + OpenClbNameUpdate + `"
  vpc_id                    = tencentcloud_vpc.foo.id
  project_id                = 0
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = tencentcloud_vpc.foo.id
  security_groups           = [tencentcloud_security_group.foo.id]
}
`
const testAccTestingClbInstance_Multi_Instance = `
resource "tencentcloud_clb_instance" "multiple_instance" {
  network_type              = "OPEN"
  clb_name                  = "` + MultiClbName + `"
  master_zone_id = "100004"
  slave_zone_id = "100003"
}
`

const testAccTestingClbInstance_Multi_Instance_Update = `
resource "tencentcloud_clb_instance" "multiple_instance" {
  network_type              = "OPEN"
  clb_name                  = "` + MultiClbName + `"
  master_zone_id = "100004"
  slave_zone_id = "100003"
}
`
