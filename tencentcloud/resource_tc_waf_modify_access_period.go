/*
Provides a resource to create a waf modify_access_period

Example Usage

```hcl
resource "tencentcloud_waf_modify_access_period" "modify_access_period" {
  period =
  topic_id = ""
}
```

Import

waf modify_access_period can be imported using the id, e.g.

```
terraform import tencentcloud_waf_modify_access_period.modify_access_period modify_access_period_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWafModifyAccessPeriod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafModifyAccessPeriodCreate,
		Read:   resourceTencentCloudWafModifyAccessPeriodRead,
		Delete: resourceTencentCloudWafModifyAccessPeriodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Access log retention period, range is [1, 180].",
			},

			"topic_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log topic, new version does not need to be uploaded.",
			},
		},
	}
}

func resourceTencentCloudWafModifyAccessPeriodCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_modify_access_period.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = waf.NewModifyAccessPeriodRequest()
		response  = waf.NewModifyAccessPeriodResponse()
		requestId string
	)
	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyAccessPeriod(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate waf ModifyAccessPeriod failed, reason:%+v", logId, err)
		return err
	}

	requestId = *response.Response.RequestId
	d.SetId(requestId)

	return resourceTencentCloudWafModifyAccessPeriodRead(d, meta)
}

func resourceTencentCloudWafModifyAccessPeriodRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_modify_access_period.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWafModifyAccessPeriodDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_modify_access_period.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
