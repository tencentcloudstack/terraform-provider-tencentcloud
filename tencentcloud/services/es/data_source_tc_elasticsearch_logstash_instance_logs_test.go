package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchLogstashInstanceLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchLogstashInstanceLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_elasticsearch_logstash_instance_logs.logstash_instance_logs")),
			},
		},
	})
}

const testAccElasticsearchLogstashInstanceLogsDataSource = `
data "tencentcloud_elasticsearch_logstash_instance_logs" "logstash_instance_logs" {
	instance_id = "ls-kru90fkz"
	log_type = 1
	start_time = "2023-10-31 10:30:00"
	end_time = "2023-10-31 10:30:10"
}
`
