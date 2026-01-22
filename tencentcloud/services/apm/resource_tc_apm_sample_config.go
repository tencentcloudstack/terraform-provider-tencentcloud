package apm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apmv20210622 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmSampleConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmSampleConfigCreate,
		Read:   resourceTencentCloudApmSampleConfigRead,
		Update: resourceTencentCloudApmSampleConfigUpdate,
		Delete: resourceTencentCloudApmSampleConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business system ID.",
			},

			"sample_rate": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Sampling rate.",
			},

			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application name.",
			},

			"sample_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Sampling rule name.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Sampling tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key value definition.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value definition.",
						},
					},
				},
			},

			"operation_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API name.",
			},

			"operation_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "0: exact match (default); 1: prefix match; 2: suffix match.",
			},
		},
	}
}

func resourceTencentCloudApmSampleConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_sample_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = apmv20210622.NewCreateApmSampleConfigRequest()
		instanceId string
		sampleName string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("sample_rate"); ok {
		request.SampleRate = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("service_name"); ok {
		request.ServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sample_name"); ok {
		request.SampleName = helper.String(v.(string))
		sampleName = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			aPMKVItem := apmv20210622.APMKVItem{}
			if v, ok := tagsMap["key"].(string); ok && v != "" {
				aPMKVItem.Key = helper.String(v)
			}

			if v, ok := tagsMap["value"].(string); ok && v != "" {
				aPMKVItem.Value = helper.String(v)
			}

			request.Tags = append(request.Tags, &aPMKVItem)
		}
	}

	if v, ok := d.GetOk("operation_name"); ok {
		request.OperationName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("operation_type"); ok {
		request.OperationType = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().CreateApmSampleConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create apm sample config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{instanceId, sampleName}, tccommon.FILED_SP))
	return resourceTencentCloudApmSampleConfigRead(d, meta)
}

func resourceTencentCloudApmSampleConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_sample_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	sampleName := idSplit[1]

	respData, err := service.DescribeApmSampleConfigById(ctx, instanceId, sampleName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_sample_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceKey != nil {
		_ = d.Set("instance_id", respData.InstanceKey)
	}

	if respData.SampleRate != nil {
		_ = d.Set("sample_rate", respData.SampleRate)
	}

	if respData.ServiceName != nil {
		_ = d.Set("service_name", respData.ServiceName)
	}

	if respData.SampleName != nil {
		_ = d.Set("sample_name", respData.SampleName)
	}

	if respData.Tags != nil && len(respData.Tags) > 0 {
		tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
		for _, tags := range respData.Tags {
			tagsMap := map[string]interface{}{}
			if tags.Key != nil {
				tagsMap["key"] = tags.Key
			}

			if tags.Value != nil {
				tagsMap["value"] = tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)
	}

	if respData.OperationName != nil {
		_ = d.Set("operation_name", respData.OperationName)
	}

	if respData.OperationType != nil {
		_ = d.Set("operation_type", respData.OperationType)
	}

	return nil
}

func resourceTencentCloudApmSampleConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_sample_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	sampleName := idSplit[1]

	needChange := false
	mutableArgs := []string{"sample_rate", "service_name", "tags", "operation_name", "operation_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := apmv20210622.NewModifyApmSampleConfigRequest()
		if v, ok := d.GetOk("sample_name"); ok {
			request.SampleName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sample_rate"); ok {
			request.SampleRate = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("service_name"); ok {
			request.ServiceName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("operation_name"); ok {
			request.OperationName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("tags"); ok {
			for _, item := range v.([]interface{}) {
				tagsMap := item.(map[string]interface{})
				aPMKVItem := apmv20210622.APMKVItem{}
				if v, ok := tagsMap["key"].(string); ok && v != "" {
					aPMKVItem.Key = helper.String(v)
				}

				if v, ok := tagsMap["value"].(string); ok && v != "" {
					aPMKVItem.Value = helper.String(v)
				}

				request.Tags = append(request.Tags, &aPMKVItem)
			}
		}

		if v, ok := d.GetOkExists("operation_type"); ok {
			request.OperationType = helper.IntInt64(v.(int))
		}

		request.InstanceId = &instanceId
		request.SampleName = &sampleName
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmSampleConfigWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update apm sample config failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudApmSampleConfigRead(d, meta)
}

func resourceTencentCloudApmSampleConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_sample_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = apmv20210622.NewDeleteApmSampleConfigRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	sampleName := idSplit[1]

	request.InstanceId = &instanceId
	request.SampleName = &sampleName
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().DeleteApmSampleConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete apm sample config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
