package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchLogstashInstanceOperationsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchLogstashInstanceOperationsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_logstash_instance_operations.logstash_instance_operations")),
			},
		},
	})
}

const testAccElasticsearchLogstashInstanceOperationsDataSource = `
data "tencentcloud_elasticsearch_logstash_instance_operations" "logstash_instance_operations" {
	instance_id = "ls-kru90fkz"
	start_time = "2018-01-01 00:00:00"
	end_time = "2023-10-31 10:12:45"
}
`
