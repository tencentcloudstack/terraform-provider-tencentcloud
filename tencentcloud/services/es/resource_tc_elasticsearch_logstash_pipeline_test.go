package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchLogstashPipelineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchLogstashPipeline,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.config"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.pipeline_id", "logstash-pipeline-test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.queue_type", "memory"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.queue_check_point_writes", "0"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.batch_delay", "50"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.batch_size", "125"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.workers", "1"),
				),
			},
			{
				Config: testAccElasticsearchLogstashPipelineUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.config"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline", "pipeline.0.workers", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_elasticsearch_logstash_pipeline.logstash_pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"op_type"},
			},
		},
	})
}

const testAccElasticsearchLogstashPipeline = tcacctest.DefaultEsVariables + `

resource "tencentcloud_elasticsearch_logstash_pipeline" "logstash_pipeline" {
  instance_id = var.logstash_id
  pipeline {
  pipeline_id = "logstash-pipeline-test"
  pipeline_desc = ""
  config =<<EOF
input{

}
filter{

}
output{

}
EOF
  queue_type = "memory"
  queue_check_point_writes = 0
  queue_max_bytes = ""
  batch_delay = 50
  batch_size = 125
  workers = 1
  }
  op_type = 2
}
`

const testAccElasticsearchLogstashPipelineUpdate = tcacctest.DefaultEsVariables + `

resource "tencentcloud_elasticsearch_logstash_pipeline" "logstash_pipeline" {
  instance_id = var.logstash_id
  pipeline {
  pipeline_id = "logstash-pipeline-test"
  pipeline_desc = ""
  config =<<EOF
input{


}
filter{

}
output{

}
EOF
  queue_type = "memory"
  queue_check_point_writes = 0
  queue_max_bytes = ""
  batch_delay = 50
  batch_size = 125
  workers = 2
  }
  op_type = 2
}
`
