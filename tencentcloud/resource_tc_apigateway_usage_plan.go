/*
Provides a resource to create a apigateway usage_plan

Example Usage

```hcl
resource "tencentcloud_apigateway_usage_plan" "usage_plan" {
  usage_plan_name = ""
  usage_plan_desc = ""
  max_request_num =
  max_request_num_pre_sec =
}
```

Import

apigateway usage_plan can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_usage_plan.usage_plan usage_plan_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudApigatewayUsagePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayUsagePlanCreate,
		Read:   resourceTencentCloudApigatewayUsagePlanRead,
		Update: resourceTencentCloudApigatewayUsagePlanUpdate,
		Delete: resourceTencentCloudApigatewayUsagePlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"usage_plan_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User defined usage plan name.",
			},

			"usage_plan_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User defined usage plan description.",
			},

			"max_request_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The total number of requested quotas, with a value range of -1 or [1, 99999999], defaults to -1, indicating that it is not enabled.",
			},

			"max_request_num_pre_sec": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of request limits per second, with a value range of -1 or [1, 2000]. The default is -1, which means it is not enabled.",
			},
		},
	}
}

func resourceTencentCloudApigatewayUsagePlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_usage_plan.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = apigateway.NewCreateUsagePlanRequest()
		response    = apigateway.NewCreateUsagePlanResponse()
		usagePlanId string
	)
	if v, ok := d.GetOk("usage_plan_name"); ok {
		request.UsagePlanName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("usage_plan_desc"); ok {
		request.UsagePlanDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_request_num"); ok {
		request.MaxRequestNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_request_num_pre_sec"); ok {
		request.MaxRequestNumPreSec = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateUsagePlan(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway usagePlan failed, reason:%+v", logId, err)
		return err
	}

	usagePlanId = *response.Response.UsagePlanId
	d.SetId(usagePlanId)

	return resourceTencentCloudApigatewayUsagePlanRead(d, meta)
}

func resourceTencentCloudApigatewayUsagePlanRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_usage_plan.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	usagePlanId := d.Id()

	usagePlan, err := service.DescribeApigatewayUsagePlanById(ctx, usagePlanId)
	if err != nil {
		return err
	}

	if usagePlan == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayUsagePlan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if usagePlan.UsagePlanName != nil {
		_ = d.Set("usage_plan_name", usagePlan.UsagePlanName)
	}

	if usagePlan.UsagePlanDesc != nil {
		_ = d.Set("usage_plan_desc", usagePlan.UsagePlanDesc)
	}

	if usagePlan.MaxRequestNum != nil {
		_ = d.Set("max_request_num", usagePlan.MaxRequestNum)
	}

	if usagePlan.MaxRequestNumPreSec != nil {
		_ = d.Set("max_request_num_pre_sec", usagePlan.MaxRequestNumPreSec)
	}

	return nil
}

func resourceTencentCloudApigatewayUsagePlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_usage_plan.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := apigateway.NewModifyUsagePlanRequest()

	usagePlanId := d.Id()

	request.UsagePlanId = &usagePlanId

	immutableArgs := []string{"usage_plan_name", "usage_plan_desc", "max_request_num", "max_request_num_pre_sec"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("usage_plan_name") {
		if v, ok := d.GetOk("usage_plan_name"); ok {
			request.UsagePlanName = helper.String(v.(string))
		}
	}

	if d.HasChange("usage_plan_desc") {
		if v, ok := d.GetOk("usage_plan_desc"); ok {
			request.UsagePlanDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("max_request_num") {
		if v, ok := d.GetOkExists("max_request_num"); ok {
			request.MaxRequestNum = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("max_request_num_pre_sec") {
		if v, ok := d.GetOkExists("max_request_num_pre_sec"); ok {
			request.MaxRequestNumPreSec = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().ModifyUsagePlan(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway usagePlan failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudApigatewayUsagePlanRead(d, meta)
}

func resourceTencentCloudApigatewayUsagePlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_usage_plan.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	usagePlanId := d.Id()

	if err := service.DeleteApigatewayUsagePlanById(ctx, usagePlanId); err != nil {
		return err
	}

	return nil
}
