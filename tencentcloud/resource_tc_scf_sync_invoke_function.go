package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfSyncInvokeFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfSyncInvokeFunctionCreate,
		Read:   resourceTencentCloudScfSyncInvokeFunctionRead,
		Delete: resourceTencentCloudScfSyncInvokeFunctionDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"qualifier": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Version or alias of the function. It defaults to $DEFAULT.",
			},

			"event": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function running parameter, which is in the JSON format. Maximum parameter size is 6 MB. This field corresponds to event input parameter.",
			},

			"log_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Valid value: None (default) or Tail. If the value is Tail, log in the response will contain the corresponding function execution log (up to 4KB).",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace. default is used if it's left empty.",
			},

			"routing_key": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Traffic routing config in json format, e.g., {k:v}. Please note that both k and v must be strings. Up to 1024 bytes allowed.",
			},
		},
	}
}

func resourceTencentCloudScfSyncInvokeFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_sync_invoke_function.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request           = scf.NewInvokeFunctionRequest()
		response          = scf.NewInvokeFunctionResponse()
		functionRequestId string
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("event"); ok {
		request.Event = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("routing_key"); ok {
		request.RoutingKey = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().InvokeFunction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf syncInvokeFunction failed, reason:%+v", logId, err)
		return err
	}

	functionRequestId = *response.Response.Result.FunctionRequestId
	d.SetId(functionRequestId)

	return resourceTencentCloudScfSyncInvokeFunctionRead(d, meta)
}

func resourceTencentCloudScfSyncInvokeFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_sync_invoke_function.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfSyncInvokeFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_sync_invoke_function.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
