package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsExecuteFunctionOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsExecuteFunctionOperationCreate,
		Read:   resourceTencentCloudMpsExecuteFunctionOperationRead,
		Delete: resourceTencentCloudMpsExecuteFunctionOperationDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of called backend API.",
			},

			"function_arg": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "API parameter. Parameter format will depend on the actual function definition.",
			},
		},
	}
}

func resourceTencentCloudMpsExecuteFunctionOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_execute_function_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = mps.NewExecuteFunctionRequest()
		functionName string
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
		functionName = v.(string)
	}

	if v, ok := d.GetOk("function_arg"); ok {
		request.FunctionArg = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ExecuteFunction(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps executeFunctionOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(functionName)

	return resourceTencentCloudMpsExecuteFunctionOperationRead(d, meta)
}

func resourceTencentCloudMpsExecuteFunctionOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_execute_function_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsExecuteFunctionOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_execute_function_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
