package ssm

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSsmRotateProductSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmRotateProductSecretCreate,
		Read:   resourceTencentCloudSsmRotateProductSecretRead,
		Delete: resourceTencentCloudSsmRotateProductSecretDelete,

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},
		},
	}
}

func resourceTencentCloudSsmRotateProductSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_rotate_product_secret.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = ssm.NewRotateProductSecretRequest()
		response     = ssm.NewRotateProductSecretResponse()
		asyncRequest = ssm.NewDescribeAsyncRequestInfoRequest()
		secretName   string
		flowId       int64
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
		secretName = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().RotateProductSecret(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssm rotateProductSecret failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(secretName)

	// wait
	flowId = *response.Response.FlowID
	asyncRequest.FlowID = &flowId
	err = resource.Retry(tccommon.WriteRetryTimeout*3, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().DescribeAsyncRequestInfo(asyncRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.TaskStatus == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("operate ssm describeAsyncRequestInfo status is %d", *result.Response.TaskStatus))
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssm describeAsyncRequestInfo failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSsmRotateProductSecretRead(d, meta)
}

func resourceTencentCloudSsmRotateProductSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_rotate_product_secret.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSsmRotateProductSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_rotate_product_secret.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
