package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudLiteHbaseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiteHbaseInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lite_hbase_instance.lite_hbase_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "instance_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-muytmxhk"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9ye3xm5v"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.node_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "tags.0.tag_value", "test"),
				),
			},
			{
				Config: testAccLiteHbaseInstanceBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lite_hbase_instance.lite_hbase_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "instance_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-muytmxhk"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9ye3xm5v"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "zone_settings.0.node_num", "5"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance", "tags.0.tag_value", "test"),
				),
			},
			{
				ResourceName:      "tencentcloud_lite_hbase_instance.lite_hbase_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudLiteHbaseInstanceResource_multiZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiteHbaseInstanceMultiZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "instance_name", "tf-test-multi-zone"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "disk_type", "CLOUD_HSSD"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "node_type", "4C16G"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.0.zone", "ap-shanghai-2"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.vpc_id", "vpc-muytmxhk"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.0.vpc_settings.0.subnet_id", "subnet-9ye3xm5v"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.0.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.1.zone", "ap-shanghai-5"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.vpc_id", "vpc-muytmxhk"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.1.vpc_settings.0.subnet_id", "subnet-1ppkfg6t"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.1.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.2.zone", "ap-shanghai-8"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.vpc_id", "vpc-muytmxhk"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.2.vpc_settings.0.subnet_id", "subnet-1tup7mn1"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "zone_settings.2.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "tags.0.tag_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_lite_hbase_instance.lite_hbase_instance_multi_zone", "tags.0.tag_value", "test"),
				),
			},
		},
	})
}

const testAccLiteHbaseInstanceBasic = `
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance" {
  instance_name = "tf-test"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-9ye3xm5v"
    }
    node_num = 3
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`

const testAccLiteHbaseInstanceBasicUpdate = `
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance" {
  instance_name = "tf-test"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-9ye3xm5v"
    }
    node_num = 5
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`

const testAccLiteHbaseInstanceMultiZone = `
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance_multi_zone" {
  instance_name = "tf-test-multi-zone"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-9ye3xm5v"
    }
    node_num = 1
  }
  zone_settings {
    zone = "ap-shanghai-5"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-1ppkfg6t"
    }
    node_num = 1
  }
  zone_settings {
    zone = "ap-shanghai-8"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-1tup7mn1"
    }
    node_num = 1
  }
  tags {
    tag_key = "test"
    tag_value = "test"
  }
}
`
