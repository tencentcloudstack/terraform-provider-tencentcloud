/*
Provides a resource to create a waf modify_access_period

Example Usage

```hcl
resource "tencentcloud_waf_modify_access_period" "example" {
  topic_id = "1ae37c76-df99-4e2b-998c-20f39eba6226"
  period   = 30
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafModifyAccessPeriod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafModifyAccessPeriodCreate,
		Read:   resourceTencentCloudWafModifyAccessPeriodRead,
		Delete: resourceTencentCloudWafModifyAccessPeriodDelete,

		Schema: map[string]*schema.Schema{
			"period": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 180),
				Description:  "Access log retention period, range is [1, 180].",
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

	var (
		logId   = getLogId(contextNil)
		request = waf.NewModifyAccessPeriodRequest()
		topicId string
	)

	if v, _ := d.GetOkExists("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		topicId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyAccessPeriod(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate waf ModifyAccessPeriod failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicId)

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
