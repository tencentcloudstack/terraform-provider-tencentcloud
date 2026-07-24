package ga2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

// ResourceTencentCloudGa2GlobalAccelerator manages a Tencent Cloud GA2 global accelerator instance.
//
// All write APIs (CreateGlobalAccelerator / ModifyGlobalAccelerator / DeleteGlobalAccelerator)
// are asynchronous: each returns a TaskId that must be polled via DescribeTaskResult until
// Status == "SUCCESS" before this resource considers the operation complete.
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Global accelerator instance name. Maximum length is 60 bytes.",
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "Billing mode. `PREPAID` for monthly subscription, `POSTPAID` for pay-as-you-go. " +
					"Default: `POSTPAID`. Currently only `POSTPAID` is supported. " +
					"Cannot be changed after creation; modifying this attribute forces a new resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Global accelerator instance description. Maximum length is 100 bytes.",
			},
			"cross_border_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Cross-border type. Valid values: `HighQuality` (premium BGP-IP cross-border), " +
					"`Unicom` (China Unicom dedicated-line cross-border).",
			},
			"cross_border_promise_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Whether the cross-border service commitment letter has been signed. Must be set to `true` " +
					"when `cross_border_type` is specified.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag key-value pairs to attach to the instance.",
			},

			// Computed
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global accelerator instance ID.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global accelerator instance state.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global accelerator instance status.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration domain (CNAME) assigned to the instance.",
			},
			"ddos_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DDoS protection instance ID associated with the global accelerator instance.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the global accelerator instance.",
			},
			"listener_counts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of listeners under this global accelerator instance.",
			},
			"accelerator_area_counts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of acceleration regions under this global accelerator instance.",
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
		gaId     string
		taskId   string
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

	//nolint:staticcheck
	if v, ok := d.GetOkExists("cross_border_promise_flag"); ok {
		request.CrossBorderPromiseFlag = helper.Bool(v.(bool))
	}

	// Tags are forwarded directly via CreateGlobalAccelerator to avoid an extra ModifyTags round-trip on Create.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		request.Tags = buildGa2GlobalAcceleratorTags(tags)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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
	gaId = *response.Response.GlobalAcceleratorId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(gaId)
	return resourceTencentCloudGa2GlobalAcceleratorRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		gaId    = d.Id()
	)

	respData, err := service.DescribeGa2GlobalAcceleratorById(ctx, gaId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_global_accelerator` [%s] not found, please check if it has been deleted.\n", logId, gaId)
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorId != nil {
		_ = d.Set("global_accelerator_id", respData.GlobalAcceleratorId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", respData.InstanceChargeType)
	}

	if respData.CrossBorderType != nil {
		_ = d.Set("cross_border_type", respData.CrossBorderType)
	}

	if respData.State != nil {
		_ = d.Set("state", respData.State)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Cname != nil {
		_ = d.Set("cname", respData.Cname)
	}

	if respData.DdosId != nil {
		_ = d.Set("ddos_id", respData.DdosId)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ListenerCounts != nil {
		_ = d.Set("listener_counts", int(*respData.ListenerCounts))
	}

	if respData.AcceleratorAreaCounts != nil {
		_ = d.Set("accelerator_area_counts", int(*respData.AcceleratorAreaCounts))
	}

	// DescribeGlobalAccelerators already returns TagSet on the instance, so we
	// hydrate `tags` directly from the response without a second round-trip.
	if len(respData.TagSet) > 0 {
		tags := make(map[string]string, len(respData.TagSet))
		for _, t := range respData.TagSet {
			if t == nil || t.Key == nil {
				continue
			}
			value := ""
			if t.Value != nil {
				value = *t.Value
			}
			tags[*t.Key] = value
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
		gaId  = d.Id()
	)

	// ModifyGlobalAccelerator only accepts: Name / Description / CrossBorderType / CrossBorderPromiseFlag.
	// Tags are managed through the standalone Tag service.
	modifyFields := []string{"name", "description", "cross_border_type", "cross_border_promise_flag"}
	needModify := false
	for _, f := range modifyFields {
		if d.HasChange(f) {
			needModify = true
			break
		}
	}

	if needModify {
		request := ga2v20250115.NewModifyGlobalAcceleratorRequest()
		request.GlobalAcceleratorId = helper.String(gaId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cross_border_type"); ok {
			if v.(string) != "NotAvailable" {
				request.CrossBorderType = helper.String(v.(string))
			}
		}

		//nolint:staticcheck
		if v, ok := d.GetOkExists("cross_border_promise_flag"); ok {
			request.CrossBorderPromiseFlag = helper.Bool(v.(bool))
		}

		var taskId string
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyGlobalAcceleratorWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

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
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("ga2", "ga", tcClient.Region, gaId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
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

	request.GlobalAcceleratorId = helper.String(d.Id())

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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

// buildGa2GlobalAcceleratorTags converts a map[string]string of tags into the SDK Tag slice.
func buildGa2GlobalAcceleratorTags(tags map[string]string) []*ga2v20250115.Tag {
	result := make([]*ga2v20250115.Tag, 0, len(tags))
	for k, v := range tags {
		key := k
		value := v
		result = append(result, &ga2v20250115.Tag{
			Key:   &key,
			Value: &value,
		})
	}
	return result
}
