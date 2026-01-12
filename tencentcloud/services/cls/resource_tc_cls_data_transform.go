package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsDataTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsDataTransformCreate,
		Read:   resourceTencentCloudClsDataTransformRead,
		Update: resourceTencentCloudClsDataTransformUpdate,
		Delete: resourceTencentCloudClsDataTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"func_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type. `1`: Specify the theme; `2`: Dynamic creation.",
			},

			"src_topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Source topic ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Task name.",
			},

			"etl_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Data transform content. If `func_type` is `2`, must use `log_auto_output`.",
			},

			"task_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Task type. `1`: Use random data from the source log theme for processing preview; `2`: Use user-defined test data for processing preview; `3`: Create real machining tasks.",
			},

			"enable_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task enable flag. `1`: enable, `2`: disable, Default is `1`.",
			},

			"dst_resources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Data transform des resources. If `func_type` is `1`, this parameter is required. If `func_type` is `2`, this parameter does not need to be filled in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dst topic ID.",
						},
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alias.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsDataTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateDataTransformRequest()
		response = cls.NewCreateDataTransformResponse()
		taskId   string
	)

	if v, ok := d.GetOkExists("func_type"); ok {
		request.FuncType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("src_topic_id"); ok {
		request.SrcTopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("etl_content"); ok {
		request.EtlContent = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_flag"); ok {
		request.EnableFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dst_resources"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataTransformResouceInfo := cls.DataTransformResouceInfo{}
			if v, ok := dMap["topic_id"]; ok {
				dataTransformResouceInfo.TopicId = helper.String(v.(string))
			}

			if v, ok := dMap["alias"]; ok {
				dataTransformResouceInfo.Alias = helper.String(v.(string))
			}

			request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	dataTransform, err := service.DescribeClsDataTransformById(ctx, dataTransformTaskId)
	if err != nil {
		return err
	}

	if dataTransform == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsDataTransform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dataTransform.SrcTopicId != nil {
		_ = d.Set("src_topic_id", dataTransform.SrcTopicId)
	}

	if dataTransform.Name != nil {
		_ = d.Set("name", dataTransform.Name)
	}

	if dataTransform.EtlContent != nil {
		_ = d.Set("etl_content", dataTransform.EtlContent)
	}

	if dataTransform.EnableFlag != nil {
		_ = d.Set("enable_flag", dataTransform.EnableFlag)
	}

	if dataTransform.DstResources != nil {
		var dstResourcesList []interface{}
		for _, dstResources := range dataTransform.DstResources {
			dstResourcesMap := map[string]interface{}{}

			if dstResources.TopicId != nil {
				dstResourcesMap["topic_id"] = dstResources.TopicId
			}

			if dstResources.Alias != nil {
				dstResourcesMap["alias"] = dstResources.Alias
			}

			dstResourcesList = append(dstResourcesList, dstResourcesMap)
		}

		_ = d.Set("dst_resources", dstResourcesList)
	}

	return nil
}

func resourceTencentCloudClsDataTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		request             = cls.NewModifyDataTransformRequest()
		dataTransformTaskId = d.Id()
	)

	immutableArgs := []string{"src_topic_id", "preview_log_statistics"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TaskId = &dataTransformTaskId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("etl_content") {
		if v, ok := d.GetOk("etl_content"); ok {
			request.EtlContent = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_flag") {
		if v, ok := d.GetOkExists("enable_flag"); ok {
			request.EnableFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("dst_resources") {
		if v, ok := d.GetOk("dst_resources"); ok {
			for _, item := range v.([]interface{}) {
				dataTransformResouceInfo := cls.DataTransformResouceInfo{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["topic_id"]; ok {
					dataTransformResouceInfo.TopicId = helper.String(v.(string))
				}

				if v, ok := dMap["alias"]; ok {
					dataTransformResouceInfo.Alias = helper.String(v.(string))
				}

				request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	if err := service.DeleteClsDataTransformById(ctx, dataTransformTaskId); err != nil {
		return err
	}

	return nil
}
