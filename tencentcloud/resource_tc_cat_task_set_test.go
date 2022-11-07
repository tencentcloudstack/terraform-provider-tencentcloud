package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCatTaskSet_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatTaskSet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cat_task_set.taskSet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCatTaskSet = `

resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
			name = "demo"
			target_address = "http://www.baidu.com"

  }
  task_type = 1
  nodes = 
    interval = 5
  parameters = "{&quot;ipType&quot;:0,&quot;netIcmpOn&quot;:1,&quot;netIcmpActivex&quot;:0,&quot;netIcmpTimeout&quot;:20,&quot;netIcmpInterval&quot;:0.5,&quot;netIcmpNum&quot;:4,&quot;netIcmpSize&quot;:32,&quot;netIcmpDataCut&quot;:1,&quot;netDnsOn&quot;:1,&quot;netDnsTimeout&quot;:20,&quot;netDnsQuerymethod&quot;:1,&quot;netDnsNs&quot;:&quot;&quot;,&quot;netDigOn&quot;:0,&quot;netDnsServer&quot;:0,&quot;netTracertOn&quot;:1,&quot;netTracertTimeout&quot;:20,&quot;netTracertNum&quot;:20,&quot;whiteList&quot;:&quot;&quot;,&quot;blackList&quot;:&quot;&quot;,&quot;netIcmpActivexStr&quot;:&quot;&quot;}"
  task_category = 1
  probe_type = 1
  plugin_source = "CDN"
  cron = "* 0-5 * * *"
  client_num = "3198058"
    tags = {
    "createdBy" = "terraform"
  }
}

`
