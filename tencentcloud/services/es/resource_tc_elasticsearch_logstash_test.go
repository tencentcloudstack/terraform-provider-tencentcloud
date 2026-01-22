package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchLogstashResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchLogstash,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash.logstash", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "instance_name", "logstash-test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "node_type", "LOGSTASH.SA2.MEDIUM4"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "disk_size", "20"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "license_type", "xpack"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "node_num", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash.logstash", "operation_duration.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_elasticsearch_logstash.logstash",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudElasticsearchLogstashResource_MultiZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchLogstashMultiZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash.logstash", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "instance_name", "logstash-test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "node_type", "LOGSTASH.SA2.MEDIUM4"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "disk_size", "20"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "license_type", "xpack"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "node_num", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash.logstash", "operation_duration.#"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "deploy_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash.logstash", "multi_zone_infos.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_elasticsearch_logstash.logstash",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccElasticsearchLogstash = `
resource "tencentcloud_elasticsearch_logstash" "logstash" {
  instance_name    = "logstash-test"
  zone             = "ap-guangzhou-6"
  logstash_version = "7.14.2"
  vpc_id           = "vpc-4owdpnwr"
  subnet_id        = "subnet-4o0zd840"
  node_num         = 1
  charge_type      = "POSTPAID_BY_HOUR"
  node_type        = "LOGSTASH.SA2.MEDIUM4"
  disk_type        = "CLOUD_SSD"
  disk_size        = 20
  license_type     = "xpack"
  operation_duration {
    periods    = [1, 2, 3, 4, 5, 6, 0]
    time_start = "02:00"
    time_end   = "06:00"
    time_zone  = "UTC+8"
  }
  tags = {
    tagKey = "tagValue"
  }
}
`

const testAccElasticsearchLogstashMultiZone = `
resource "tencentcloud_elasticsearch_logstash" "logstash" {
  instance_name    = "logstash-test"
  zone             = "-"
  logstash_version = "7.14.2"
  vpc_id           = "vpc-axrsmmrv"
  subnet_id        = "-"
  node_num         = 2
  charge_type      = "POSTPAID_BY_HOUR"
  node_type        = "LOGSTASH.SA2.MEDIUM4"
  disk_type        = "CLOUD_SSD"
  disk_size        = 20
  license_type     = "xpack"
  operation_duration {
    periods    = [1, 2, 3, 4, 5, 6, 0]
    time_start = "02:00"
    time_end   = "06:00"
    time_zone  = "UTC+8"
  }
  tags = {
    tagKey = "tagValue"
  }
  deploy_mode = 1
  multi_zone_infos {
    availability_zone = "ap-guangzhou-3"
    subnet_id         = "subnet-j5vja918"
  }
  multi_zone_infos {
    availability_zone = "ap-guangzhou-4"
    subnet_id         = "subnet-oi7ya2j6"
  }
}
`
