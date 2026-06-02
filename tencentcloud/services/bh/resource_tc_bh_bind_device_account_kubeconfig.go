package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhBindDeviceAccountKubeconfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhBindDeviceAccountKubeconfigCreate,
		Read:   resourceTencentCloudBhBindDeviceAccountKubeconfigRead,
		Update: resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate,
		Delete: resourceTencentCloudBhBindDeviceAccountKubeconfigDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Container account Id. Maps to the SDK request field `Id`. Renamed in HCL because `id` is reserved by the Terraform Plugin SDK as the resource's internal identifier.",
			},

			"kubeconfig": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Container account kubeconfig credential.",
			},

			"manage_dimension": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Manage dimension. 1 means cluster.",
			},
		},
	}
}

func resourceTencentCloudBhBindDeviceAccountKubeconfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_account_kubeconfig.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var accountId int
	if v, ok := d.GetOkExists("account_id"); ok {
		accountId = v.(int)
	} else {
		return fmt.Errorf("`account_id` is required for resource tencentcloud_bh_bind_device_account_kubeconfig.")
	}

	d.SetId(fmt.Sprintf("%d", accountId))
	return resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate(d, meta)
}

func resourceTencentCloudBhBindDeviceAccountKubeconfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_account_kubeconfig.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	// The BH service does not currently expose a query API for kubeconfig
	// bindings, so Read is a no-op. State is authoritative; external drift
	// (e.g. credential rotation in the web console) is invisible to Terraform.
	return nil
}

func resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_account_kubeconfig.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewBindDeviceAccountKubeconfigRequest()
	)

	if v, ok := d.GetOkExists("account_id"); ok {
		request.Id = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("kubeconfig"); ok {
		request.Kubeconfig = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("manage_dimension"); ok {
		request.ManageDimension = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().BindDeviceAccountKubeconfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Bind bh device account kubeconfig failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s bind bh device account kubeconfig failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBhBindDeviceAccountKubeconfigRead(d, meta)
}

func resourceTencentCloudBhBindDeviceAccountKubeconfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_bind_device_account_kubeconfig.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	// The BH service does not provide an unbind API. terraform destroy only
	// removes the resource from local state; the kubeconfig binding on the
	// backend is preserved.
	return nil
}
