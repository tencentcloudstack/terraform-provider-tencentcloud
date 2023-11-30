package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfInvokeFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfInvokeFunctionCreate,
		Read:   resourceTencentCloudScfInvokeFunctionRead,
		Delete: resourceTencentCloudScfInvokeFunctionDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"invocation_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Fill in RequestResponse for synchronized invocations (default and recommended) and Event for asychronized invocations. Note that for synchronized invocations, the max timeout period is 300s. Choose asychronized invocations if the required timeout period is longer than 300 seconds. You can also use InvokeFunction for synchronized invocations.",
			},

			"qualifier": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The version or alias of the triggered function. It defaults to $LATEST.",
			},

			"client_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function running parameter, which is in the JSON format. The maximum parameter size is 6 MB for synchronized invocations and 128KB for asynchronized invocations. This field corresponds to event input parameter.",
			},

			"log_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Null for async invocations.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace.",
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

func resourceTencentCloudScfInvokeFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_invoke_function.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = scf.NewInvokeRequest()
		response = scf.NewInvokeResponse()
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invocation_type"); ok {
		request.InvocationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("client_context"); ok {
		request.ClientContext = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().Invoke(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf InvokeFunction failed, reason:%+v", logId, err)
		return err
	}

	functionRequestId := *response.Response.Result.FunctionRequestId

	d.SetId(functionRequestId)

	return resourceTencentCloudScfInvokeFunctionRead(d, meta)
}

func resourceTencentCloudScfInvokeFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_invoke_function.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfInvokeFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_invoke_function.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
