package tke

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesRollOutSequence() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesRollOutSequenceCreate,
		Read:   resourceTencentCloudKubernetesRollOutSequenceRead,
		Update: resourceTencentCloudKubernetesRollOutSequenceUpdate,
		Delete: resourceTencentCloudKubernetesRollOutSequenceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the roll-out sequence.",
			},
			"sequence_flows": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The sequence flow steps of the roll-out sequence.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The tags for the sequence flow step.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Tag values.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"soak_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Wait time in seconds between steps.",
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the roll-out sequence is enabled.",
			},
		},
	}
}

func resourceTencentCloudKubernetesRollOutSequenceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_roll_out_sequence.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = tkev20180525.NewCreateRollOutSequenceRequest()
		response = tkev20180525.NewCreateRollOutSequenceResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sequence_flows"); ok {
		for _, item := range v.([]interface{}) {
			flowMap := item.(map[string]interface{})
			sequenceFlow := tkev20180525.SequenceFlow{}
			if v, ok := flowMap["tags"]; ok {
				for _, tagItem := range v.([]interface{}) {
					tagMap := tagItem.(map[string]interface{})
					sequenceTag := tkev20180525.SequenceTag{}
					if v, ok := tagMap["key"].(string); ok {
						sequenceTag.Key = helper.String(v)
					}

					if v, ok := tagMap["value"]; ok {
						for _, val := range v.([]interface{}) {
							sequenceTag.Value = append(sequenceTag.Value, helper.String(val.(string)))
						}
					}

					sequenceFlow.Tags = append(sequenceFlow.Tags, &sequenceTag)
				}
			}

			if v, ok := flowMap["soak_time"].(int); ok {
				sequenceFlow.SoakTime = helper.Int64(int64(v))
			}

			request.SequenceFlows = append(request.SequenceFlows, &sequenceFlow)
		}
	}

	request.Enabled = helper.Bool(d.Get("enabled").(bool))

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().CreateRollOutSequenceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create kubernetes roll out sequence failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes roll out sequence failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ID == nil {
		return fmt.Errorf("ID is nil.")
	}

	d.SetId(strconv.FormatInt(*response.Response.ID, 10))
	return resourceTencentCloudKubernetesRollOutSequenceRead(d, meta)
}

func resourceTencentCloudKubernetesRollOutSequenceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_roll_out_sequence.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idStr := d.Id()
	sequenceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is invalid, %s: %v", idStr, err)
	}

	var sequence *tkev20180525.RollOutSequence
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var offset int64
		var limit int64 = 20
		for {
			request := tkev20180525.NewDescribeRollOutSequencesRequest()
			request.Offset = helper.Int64(offset)
			request.Limit = helper.Int64(limit)

			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeRollOutSequencesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe kubernetes roll out sequences failed, Response is nil."))
			}

			for _, seq := range result.Response.Sequences {
				if seq.ID != nil && *seq.ID == sequenceId {
					sequence = seq
					return nil
				}
			}

			if result.Response.TotalCount == nil || offset+limit >= *result.Response.TotalCount {
				break
			}

			offset += limit
		}

		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if sequence == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_roll_out_sequence` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if sequence.Name != nil {
		_ = d.Set("name", sequence.Name)
	}

	if sequence.SequenceFlows != nil && len(sequence.SequenceFlows) > 0 {
		sequenceFlowsList := make([]map[string]interface{}, 0, len(sequence.SequenceFlows))
		for _, flow := range sequence.SequenceFlows {
			flowMap := map[string]interface{}{}
			if flow.Tags != nil && len(flow.Tags) > 0 {
				tagsList := make([]map[string]interface{}, 0, len(flow.Tags))
				for _, tag := range flow.Tags {
					tagMap := map[string]interface{}{}
					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}

					if tag.Value != nil {
						valueList := make([]interface{}, 0, len(tag.Value))
						for _, v := range tag.Value {
							valueList = append(valueList, *v)
						}
						tagMap["value"] = valueList
					}

					tagsList = append(tagsList, tagMap)
				}

				flowMap["tags"] = tagsList
			}

			if flow.SoakTime != nil {
				flowMap["soak_time"] = flow.SoakTime
			}

			sequenceFlowsList = append(sequenceFlowsList, flowMap)
		}

		_ = d.Set("sequence_flows", sequenceFlowsList)
	}

	if sequence.Enabled != nil {
		_ = d.Set("enabled", sequence.Enabled)
	}

	return nil
}

func resourceTencentCloudKubernetesRollOutSequenceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_roll_out_sequence.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idStr := d.Id()
	sequenceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is invalid, %s: %v", idStr, err)
	}

	needChange := false
	mutableArgs := []string{"name", "sequence_flows", "enabled"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyRollOutSequenceRequest()
		request.ID = helper.Int64(sequenceId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("sequence_flows"); ok {
			for _, item := range v.([]interface{}) {
				flowMap := item.(map[string]interface{})
				sequenceFlow := tkev20180525.SequenceFlow{}
				if v, ok := flowMap["tags"]; ok {
					for _, tagItem := range v.([]interface{}) {
						tagMap := tagItem.(map[string]interface{})
						sequenceTag := tkev20180525.SequenceTag{}
						if v, ok := tagMap["key"].(string); ok {
							sequenceTag.Key = helper.String(v)
						}

						if v, ok := tagMap["value"]; ok {
							for _, val := range v.([]interface{}) {
								sequenceTag.Value = append(sequenceTag.Value, helper.String(val.(string)))
							}
						}

						sequenceFlow.Tags = append(sequenceFlow.Tags, &sequenceTag)
					}
				}

				if v, ok := flowMap["soak_time"].(int); ok {
					sequenceFlow.SoakTime = helper.Int64(int64(v))
				}

				request.SequenceFlows = append(request.SequenceFlows, &sequenceFlow)
			}
		}

		request.Enabled = helper.Bool(d.Get("enabled").(bool))

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().ModifyRollOutSequenceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update kubernetes roll out sequence failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudKubernetesRollOutSequenceRead(d, meta)
}

func resourceTencentCloudKubernetesRollOutSequenceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_roll_out_sequence.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewDeleteRollOutSequenceRequest()
	)

	idStr := d.Id()
	sequenceId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("id is invalid, %s: %v", idStr, err)
	}

	request.ID = helper.Int64(sequenceId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DeleteRollOutSequenceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes roll out sequence failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
