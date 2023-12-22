package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchRestartLogstashInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchRestartLogstashInstanceOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_restart_logstash_instance_operation.restart_logstash_instance_operation", "id"),
				),
			},
		},
	})
}

const testAccElasticsearchRestartLogstashInstanceOperation = tcacctest.DefaultEsVariables + `

resource "tencentcloud_elasticsearch_restart_logstash_instance_operation" "restart_logstash_instance_operation" {
  instance_id = var.logstash_id
  type = 0
}

`
