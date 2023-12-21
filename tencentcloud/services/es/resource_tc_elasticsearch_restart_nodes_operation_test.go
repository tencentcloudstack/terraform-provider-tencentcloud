package es_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudElasticsearchRestartNodesOperationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchRestartNodesOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_restart_nodes_operation.restart_nodes_operation", "id")),
			},
		},
	})
}

const testAccElasticsearchRestartNodesOperation = `

resource "tencentcloud_elasticsearch_restart_nodes_operation" "restart_nodes_operation" {
	instance_id = "es-5wn36he6"
	node_names = ["1648026612002990732"]
}

`
