package igtm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIgtmPackageTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIgtmPackageTaskCreate,
		Read:   resourceTencentCloudIgtmPackageTaskRead,
		Update: resourceTencentCloudIgtmPackageTaskUpdate,
		Delete: resourceTencentCloudIgtmPackageTaskDelete,
		Schema: map[string]*schema.Schema{
			"task_detection_quantity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Value range: 1~10000.",
			},

			"auto_renew": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "Auto renewal: 1 enable auto renewal; 2 disable auto renewal.",
			},

			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Package duration in months, required for creation and renewal. Value range: 1~120.",
			},

			"auto_voucher": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Whether to automatically select vouchers, 1 yes; 0 no, default is 0.",
			},

			// computed
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task ID.",
			},
		},
	}
}

func resourceTencentCloudIgtmPackageTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_package_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = igtmv20231024.NewCreatePackageAndPayRequest()
		response = igtmv20231024.NewCreatePackageAndPayResponse()
		taskId   string
	)

	if v, ok := d.GetOkExists("task_detection_quantity"); ok {
		request.GoodsNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_voucher"); ok {
		request.AutoVoucher = helper.IntUint64(v.(int))
	}

	request.DealType = helper.String("CREATE")
	request.GoodsType = helper.String("TASK")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreatePackageAndPayWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ResourceIds == nil {
			return resource.NonRetryableError(fmt.Errorf("Create igtm package create failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create igtm package failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.ResourceIds) == 0 {
		return fmt.Errorf("ResourceIds is nil.")
	}

	taskId = *response.Response.ResourceIds[0]
	d.SetId(taskId)

	return resourceTencentCloudIgtmPackageTaskRead(d, meta)
}

func resourceTencentCloudIgtmPackageTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_package_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  = d.Id()
	)

	respData, err := service.DescribeIgtmPackageById(ctx, taskId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_igtm_package_task` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.PackageType != nil {
		_ = d.Set("goods_type", respData.PackageType)
	}

	if respData.AutoRenewFlag != nil {
		if *respData.AutoRenewFlag == 0 {
			_ = d.Set("auto_renew", 2)
		} else {
			_ = d.Set("auto_renew", 1)
		}
	}

	if respData.ResourceId != nil {
		_ = d.Set("task_id", respData.ResourceId)
	}

	return nil
}

func resourceTencentCloudIgtmPackageTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_package_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		taskId = d.Id()
	)

	if d.HasChange("task_detection_quantity") {
		request := igtmv20231024.NewCreatePackageAndPayRequest()
		if v, ok := d.GetOkExists("task_detection_quantity"); ok {
			request.GoodsNum = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_renew"); ok {
			request.AutoRenew = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntUint64(v.(int))
		}

		request.DealType = helper.String("MODIFY")
		request.GoodsType = helper.String("TASK")
		request.ResourceId = helper.String(taskId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreatePackageAndPayWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify igtm package task failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("time_span") {
		request := igtmv20231024.NewCreatePackageAndPayRequest()
		if v, ok := d.GetOkExists("task_detection_quantity"); ok {
			request.GoodsNum = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_renew"); ok {
			request.AutoRenew = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("time_span"); ok {
			request.TimeSpan = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("auto_voucher"); ok {
			request.AutoVoucher = helper.IntUint64(v.(int))
		}

		request.DealType = helper.String("RENEW")
		request.GoodsType = helper.String("TASK")
		request.ResourceId = &taskId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreatePackageAndPayWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify igtm package task failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("auto_renew") {
		request := igtmv20231024.NewModifyPackageAutoRenewRequest()
		if v, ok := d.GetOkExists("auto_renew"); ok {
			request.AutoRenew = helper.IntUint64(v.(int))
		}

		request.ResourceId = &taskId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().ModifyPackageAutoRenewWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update igtm package task failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudIgtmPackageTaskRead(d, meta)
}

func resourceTencentCloudIgtmPackageTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_package_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
