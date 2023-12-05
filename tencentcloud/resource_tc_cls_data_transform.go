package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsDataTransform() *schema.Resource {
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
				Description: "task type.",
			},

			"src_topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "src topic id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "task name.",
			},

			"etl_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "data transform content.",
			},

			"task_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "task type.",
			},

			"enable_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "task enable flag.",
			},

			"dst_resources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "data transform des resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "dst topic id.",
						},
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "alias.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsDataTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_data_transform.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateDataTransform(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cls_data_transform.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	dataTransformTaskId := d.Id()

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
	defer logElapsed("resource.tencentcloud_cls_data_transform.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyDataTransformRequest()

	dataTransformTaskId := d.Id()

	request.TaskId = &dataTransformTaskId

	immutableArgs := []string{"src_topic_id", "name", "etl_content", "enable_flag", "preview_log_statistics"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyDataTransform(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cls_data_transform.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	dataTransformTaskId := d.Id()

	if err := service.DeleteClsDataTransformById(ctx, dataTransformTaskId); err != nil {
		return err
	}

	return nil
}
