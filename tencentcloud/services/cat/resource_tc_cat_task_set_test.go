package cat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCatTaskSet_basic -v
func TestAccTencentCloudCatTaskSet_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatTaskSet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "batch_tasks.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "batch_tasks.0.name", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "batch_tasks.0.target_address", "http://www.baidu.com"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "task_type", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "nodes.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "interval", "6"),
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "parameters"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "task_category", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "cron", "* 0-6 * * *"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_cat_task_set.task_set",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCatTaskSetUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "batch_tasks.0.name", "demo_test"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "batch_tasks.0.target_address", "http://www.baidu.com"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "task_type", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "nodes.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "interval", "6"),
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "parameters"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "task_category", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "cron", "* 0-6 * * *"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "tags.createdBy", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "status", "4"),
				),
			},
			{
				Config: testAccCatTaskSetResume,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cat_task_set.task_set", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cat_task_set.task_set", "status", "10"),
				),
			},
		},
	})
}

const testAccCatTaskSet = `

resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
    name           = "demo"
    target_address = "http://www.baidu.com"
  }
  task_type     = 5
  nodes         = ["12136", "12137", "12138", "12141", "12144"]
  interval      = 6
  parameters    = jsonencode(
  {
    "ipType"            = 0,
    "grabBag"           = 0,
    "filterIp"          = 0,
    "netIcmpOn"         = 1,
    "netIcmpActivex"    = 0,
    "netIcmpTimeout"    = 20,
    "netIcmpInterval"   = 0.5,
    "netIcmpNum"        = 20,
    "netIcmpSize"       = 32,
    "netIcmpDataCut"    = 1,
    "netDnsOn"          = 1,
    "netDnsTimeout"     = 5,
    "netDnsQuerymethod" = 1,
    "netDnsNs"          = "",
    "netDigOn"          = 1,
    "netDnsServer"      = 2,
    "netTracertOn"      = 1,
    "netTracertTimeout" = 60,
    "netTracertNum"     = 30,
    "whiteList"         = "",
    "blackList"         = "",
    "netIcmpActivexStr" = ""
  }
  )
  task_category = 1
  cron          = "* 0-6 * * *"
  tags          = {
    "createdBy" = "terraform"
  }
}

`

const testAccCatTaskSetUp = `

resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
    name           = "demo_test"
    target_address = "http://www.baidu.com"
  }
  task_type     = 5
  nodes         = ["12136", "12137", "12138", "12141", "12144"]
  interval      = 6
  parameters    = jsonencode(
  {
    "ipType"            = 0,
    "grabBag"           = 0,
    "filterIp"          = 0,
    "netIcmpOn"         = 1,
    "netIcmpActivex"    = 0,
    "netIcmpTimeout"    = 20,
    "netIcmpInterval"   = 0.5,
    "netIcmpNum"        = 20,
    "netIcmpSize"       = 32,
    "netIcmpDataCut"    = 1,
    "netDnsOn"          = 1,
    "netDnsTimeout"     = 5,
    "netDnsQuerymethod" = 1,
    "netDnsNs"          = "",
    "netDigOn"          = 1,
    "netDnsServer"      = 2,
    "netTracertOn"      = 1,
    "netTracertTimeout" = 60,
    "netTracertNum"     = 30,
    "whiteList"         = "",
    "blackList"         = "",
    "netIcmpActivexStr" = ""
  }
  )
  task_category = 1
  cron          = "* 0-6 * * *"
  tags          = {
    "createdBy" = "terraform"
  }
  operate       = "suspend"
}

`

const testAccCatTaskSetResume = `

resource "tencentcloud_cat_task_set" "task_set" {
  batch_tasks {
    name           = "demo_test"
    target_address = "http://www.baidu.com"
  }
  task_type     = 5
  nodes         = ["12136", "12137", "12138", "12141", "12144"]
  interval      = 6
  parameters    = jsonencode(
  {
    "ipType"            = 0,
    "grabBag"           = 0,
    "filterIp"          = 0,
    "netIcmpOn"         = 1,
    "netIcmpActivex"    = 0,
    "netIcmpTimeout"    = 20,
    "netIcmpInterval"   = 0.5,
    "netIcmpNum"        = 20,
    "netIcmpSize"       = 32,
    "netIcmpDataCut"    = 1,
    "netDnsOn"          = 1,
    "netDnsTimeout"     = 5,
    "netDnsQuerymethod" = 1,
    "netDnsNs"          = "",
    "netDigOn"          = 1,
    "netDnsServer"      = 2,
    "netTracertOn"      = 1,
    "netTracertTimeout" = 60,
    "netTracertNum"     = 30,
    "whiteList"         = "",
    "blackList"         = "",
    "netIcmpActivexStr" = ""
  }
  )
  task_category = 1
  cron          = "* 0-6 * * *"
  tags          = {
    "createdBy" = "terraform"
  }
  operate       = "resume"
}

`
