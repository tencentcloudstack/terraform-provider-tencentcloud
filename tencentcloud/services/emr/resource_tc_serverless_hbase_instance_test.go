package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudServerlessHbaseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-shanghai")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessHbaseInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "instance_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9yg0eh05"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.node_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "tags.0.tag_value", "test"),
				),
			},
			{
				Config: testAccServerlessHbaseInstanceBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "instance_name", "tf-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9yg0eh05"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "zone_settings.0.node_num", "5"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "tags.0.tag_key", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance", "tags.0.tag_value", "test-update"),
				),
			},
			{
				ResourceName:            "tencentcloud_serverless_hbase_instance.serverless_hbase_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span", "time_unit"},
			},
		},
	})
}

func TestAccTencentCloudServerlessHbaseInstanceResource_prepay(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-shanghai")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessHbaseInstancePrePay,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "id"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "instance_name", "tf-test-prepay"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "pay_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "time_span", "1"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "time_unit", "m"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay", "auto_renew_flag", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_serverless_hbase_instance.serverless_hbase_instance_prepay",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span", "time_unit"},
			},
		},
	})
}

func TestAccTencentCloudServerlessHbaseInstanceResource_multiZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-shanghai")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessHbaseInstanceMultiZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "id"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "instance_name", "tf-test-multi-zone"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9yg0eh05"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.zone", "ap-shanghai-5"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.subnet_id", "subnet-0ucpc5fh"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.zone", "ap-shanghai-8"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.subnet_id", "subnet-koue2ynx"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "tags.0.tag_value", "test"),
				),
			},
			{
				Config: testAccServerlessHbaseInstanceMultiZoneUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "id"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "instance_name", "tf-test-multi-zone"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9yg0eh05"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.0.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.zone", "ap-shanghai-5"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.subnet_id", "subnet-0ucpc5fh"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.1.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.zone", "ap-shanghai-8"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.vpc_id", "vpc-oxyg5hxg"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.subnet_id", "subnet-koue2ynx"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "zone_settings.2.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_serverless_hbase_instance.serverless_hbase_instance_multi_zone", "tags.0.tag_value", "test"),
				),
			},
		},
	})
}

const testAccServerlessHbaseInstanceBasic = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance" {
  instance_name = "tf-test"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-9yg0eh05"
    }
    node_num = 3
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`

const testAccServerlessHbaseInstanceBasicUpdate = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance" {
  instance_name = "tf-test-update"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-9yg0eh05"
    }
    node_num = 5
  }
  tags {
    tag_key = "test-update"
    tag_value = "test-update"
  }
}
`

const testAccServerlessHbaseInstancePrePay = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance_prepay" {
  instance_name = "tf-test-prepay"
  pay_mode = 1
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-9yg0eh05"
    }
    node_num = 3
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
  time_span = 1
  time_unit = "m"
  auto_renew_flag = 1
}
`

const testAccServerlessHbaseInstanceMultiZone = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance_multi_zone" {
  instance_name = "tf-test-multi-zone"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-9yg0eh05"
    }
    node_num = 1
  }
  zone_settings {
    zone = "ap-shanghai-5"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-0ucpc5fh"
    }
    node_num = 1
  }
  zone_settings {
    zone = "ap-shanghai-8"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-koue2ynx"
    }
    node_num = 1
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`

const testAccServerlessHbaseInstanceMultiZoneUpdate = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance_multi_zone" {
  instance_name = "tf-test-multi-zone"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-9yg0eh05"
    }
    node_num = 2
  }
  zone_settings {
    zone = "ap-shanghai-5"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-0ucpc5fh"
    }
    node_num = 2
  }
  zone_settings {
    zone = "ap-shanghai-8"
    vpc_settings {
      vpc_id = "vpc-oxyg5hxg"
      subnet_id = "subnet-koue2ynx"
    }
    node_num = 2
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`
