package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchLogstashResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
