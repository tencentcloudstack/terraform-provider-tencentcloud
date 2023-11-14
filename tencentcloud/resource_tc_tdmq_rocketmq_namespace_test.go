package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqNamespaceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqNamespace,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_namespace.namespace", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_namespace.namespace",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqNamespace = `

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id = &lt;nil&gt;
  namespace_id = &lt;nil&gt;
  ttl = &lt;nil&gt;
  retention_time = &lt;nil&gt;
  remark = &lt;nil&gt;
    }

`
