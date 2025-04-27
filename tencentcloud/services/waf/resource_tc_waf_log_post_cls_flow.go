package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafLogPostClsFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafLogPostClsFlowCreate,
		Read:   resourceTencentCloudWafLogPostClsFlowRead,
		Update: resourceTencentCloudWafLogPostClsFlowUpdate,
		Delete: resourceTencentCloudWafLogPostClsFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cls_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the CLS is delivered. The default value is ap-shanghai.",
			},

			"logset_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the log set where the delivered CLS is located. The default value is waf_post_logset.",
			},

			"log_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "1- Access log, 2- Attack log, the default is access log.",
			},

			"log_topic_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the log subject where the submitted CLS is located. The default value is waf_post_logtopic.",
			},

			"flow_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID for post cls flow.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CLS logset ID.",
			},

			"log_topic_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CLS log topic ID.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status 0- Off 1- On.",
			},
		},
	}
}

func resourceTencentCloudWafLogPostClsFlowCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_cls_flow.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = wafv20180125.NewCreatePostCLSFlowRequest()
		flowId  int64
		logType int64 = 1
		hasFlow bool
	)

	if v, ok := d.GetOk("cls_region"); ok {
		request.CLSRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_topic_name"); ok {
		request.LogTopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_type"); ok {
		request.LogType = helper.IntInt64(v.(int))
		logType = int64(v.(int))
	}

	// check unique first
	respData, err := service.DescribeWafLogPostClsFlowById(ctx, logType)
	if err != nil {
		return err
	}

	if respData != nil && len(respData.PostCLSFlows) != 0 {
		return fmt.Errorf("In the case of `log_type` is %d, only one resource can be created and cannot be created multiple times.", logType)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreatePostCLSFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf log post cls flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// get flowId
	respData, err = service.DescribeWafLogPostClsFlowById(ctx, logType)
	if err != nil {
		return err
	}

	if respData == nil || len(respData.PostCLSFlows) == 0 {
		return fmt.Errorf("resource `waf_log_post_cls_flow` not found, please check if it has been deleted.")
	}

	for _, item := range respData.PostCLSFlows {
		if *item.LogType == logType {
			flowId = *item.FlowId
			hasFlow = true
			break
		}
	}

	if !hasFlow {
		return fmt.Errorf("resource `waf_log_post_cls_flow` not found flowId, please check if it has been deleted.")
	}

	d.SetId(strings.Join([]string{helper.Int64ToStr(flowId), helper.Int64ToStr(logType)}, tccommon.FILED_SP))

	return resourceTencentCloudWafLogPostClsFlowRead(d, meta)
}

func resourceTencentCloudWafLogPostClsFlowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_cls_flow.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		hasFlow bool
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	flowId := idSplit[0]
	logType := idSplit[1]

	respData, err := service.DescribeWafLogPostClsFlowById(ctx, helper.StrToInt64(logType))
	if err != nil {
		return err
	}

	if respData == nil || len(respData.PostCLSFlows) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_log_post_cls_flow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, item := range respData.PostCLSFlows {
		if *item.FlowId == helper.StrToInt64(flowId) && *item.LogType == helper.StrToInt64(logType) {
			if item.CLSRegion != nil {
				_ = d.Set("cls_region", item.CLSRegion)
			}

			if item.LogsetName != nil {
				_ = d.Set("logset_name", item.LogsetName)
			}

			if item.LogType != nil {
				_ = d.Set("log_type", item.LogType)
			}

			if item.LogTopicName != nil {
				_ = d.Set("log_topic_name", item.LogTopicName)
			}

			if item.FlowId != nil {
				_ = d.Set("flow_id", item.FlowId)
			}

			if item.LogsetID != nil {
				_ = d.Set("logset_id", item.LogsetID)
			}

			if item.LogTopicID != nil {
				_ = d.Set("log_topic_id", item.LogTopicID)
			}

			if item.Status != nil {
				_ = d.Set("status", item.Status)
			}

			hasFlow = true
			break
		}
	}

	if !hasFlow {
		return fmt.Errorf("resource `waf_log_post_cls_flow` not found flowId, please check if it has been deleted.")
	}

	return nil
}

func resourceTencentCloudWafLogPostClsFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_cls_flow.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = wafv20180125.NewCreatePostCLSFlowRequest()
	)

	immutableArgs := []string{"log_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	logType := idSplit[1]

	if v, ok := d.GetOk("cls_region"); ok {
		request.CLSRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_topic_name"); ok {
		request.LogTopicName = helper.String(v.(string))
	}

	request.LogType = helper.StrToInt64Point(logType)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreatePostCLSFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf log post cls flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafLogPostClsFlowRead(d, meta)
}

func resourceTencentCloudWafLogPostClsFlowDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_cls_flow.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewDestroyPostCLSFlowRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	flowId := idSplit[0]
	logType := idSplit[1]

	request.FlowId = helper.StrToInt64Point(flowId)
	request.LogType = helper.StrToInt64Point(logType)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DestroyPostCLSFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf log post cls flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
