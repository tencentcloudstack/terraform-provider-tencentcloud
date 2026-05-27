package ga2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2GlobalAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2GlobalAcceleratorCreate,
		Read:   resourceTencentCloudGa2GlobalAcceleratorRead,
		Update: resourceTencentCloudGa2GlobalAcceleratorUpdate,
		Delete: resourceTencentCloudGa2GlobalAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the global accelerator instance. Maximum length is 60 bytes.",
			},
			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "POSTPAID",
				Description: "Billing mode. Valid values: `PREPAID` (prepaid, monthly subscription), `POSTPAID` (postpaid, pay-as-you-go). Default: `POSTPAID`. Currently only postpaid is supported.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the global accelerator instance. Maximum length is 100 bytes.",
			},
			"cross_border_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cross-border type. Valid values: `HighQuality` (premium BGP-IP cross-border), `Unicom` (China Unicom dedicated line cross-border).",
			},
			"cross_border_promise_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag indicating acceptance of cross-border service agreement. Must be set to `true` when using cross-border service.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Tag information.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the global accelerator instance.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the global accelerator instance.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the global accelerator instance.",
			},
			"ddos_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DDoS protection ID of the global accelerator instance.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME domain of the global accelerator instance.",
			},
		},
	}
}

func resourceTencentCloudGa2GlobalAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = ga2v20250115.NewCreateGlobalAcceleratorRequest()
		response = ga2v20250115.NewCreateGlobalAcceleratorResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cross_border_type"); ok {
		request.CrossBorderType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cross_border_promise_flag"); ok {
		request.CrossBorderPromiseFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		tagList := make([]*ga2v20250115.Tag, 0, len(tags))
		for key, val := range tags {
			tagList = append(tagList, &ga2v20250115.Tag{
				Key:   helper.String(key),
				Value: helper.String(val.(string)),
			})
		}
		request.Tags = tagList
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 global accelerator failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 global accelerator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.GlobalAcceleratorId == nil {
		return fmt.Errorf("GlobalAcceleratorId is nil.")
	}
	globalAcceleratorId := *response.Response.GlobalAcceleratorId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId := *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(globalAcceleratorId)
	return resourceTencentCloudGa2GlobalAcceleratorRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	globalAcceleratorId := d.Id()

	respData, err := service.DescribeGa2GlobalAcceleratorById(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_global_accelerator` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", respData.InstanceChargeType)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.CrossBorderType != nil {
		_ = d.Set("cross_border_type", respData.CrossBorderType)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.State != nil {
		_ = d.Set("state", respData.State)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.DdosId != nil {
		_ = d.Set("ddos_id", respData.DdosId)
	}

	if respData.Cname != nil {
		_ = d.Set("cname", respData.Cname)
	}

	if respData.TagSet != nil {
		tags := make(map[string]string, len(respData.TagSet))
		for _, tag := range respData.TagSet {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudGa2GlobalAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	globalAcceleratorId := d.Id()

	needChange := false
	request := ga2v20250115.NewModifyGlobalAcceleratorRequest()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)

	if d.HasChange("name") {
		needChange = true
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		needChange = true
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("cross_border_type") {
		needChange = true
		if v, ok := d.GetOk("cross_border_type"); ok {
			request.CrossBorderType = helper.String(v.(string))
		}
	}

	if d.HasChange("cross_border_promise_flag") {
		needChange = true
		if v, ok := d.GetOk("cross_border_promise_flag"); ok {
			request.CrossBorderPromiseFlag = helper.Bool(v.(bool))
		}
	}

	if !needChange {
		return resourceTencentCloudGa2GlobalAcceleratorRead(d, meta)
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyGlobalAcceleratorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 global accelerator failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 global accelerator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2GlobalAcceleratorRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteGlobalAcceleratorRequest()
	)

	globalAcceleratorId := d.Id()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 global accelerator failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 global accelerator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}
