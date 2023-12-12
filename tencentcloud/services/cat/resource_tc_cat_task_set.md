Provides a resource to create a cat task_set

Example Usage

```hcl
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
```
Import

cat task_set can be imported using the id, e.g.
```
$ terraform import tencentcloud_cat_task_set.task_set taskSet_id
```