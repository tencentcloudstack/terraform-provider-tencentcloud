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

func ResourceTencentCloudWafLogPostCkafkaFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafLogPostCkafkaFlowCreate,
		Read:   resourceTencentCloudWafLogPostCkafkaFlowRead,
		Update: resourceTencentCloudWafLogPostCkafkaFlowUpdate,
		Delete: resourceTencentCloudWafLogPostCkafkaFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ckafka_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where CKafka is located for delivery.",
			},

			"ckafka_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CKafka ID.",
			},

			"brokers": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The supporting environment is IP:PORT, The external network environment is domain:PORT.",
			},

			"compression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Default to none, supports snappy, gzip, and lz4 compression, recommended snappy.",
			},

			"vip_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "1. External network TGW, 2. Supporting environment, default is supporting environment.",
			},

			"log_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "1- Access log, 2- Attack log, the default is access log.",
			},

			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Theme name, default not to pass or pass empty string, default value is waf_post_access_log.",
			},

			"kafka_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version number of Kafka cluster.",
			},

			"sasl_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable SASL verification, default not enabled, 0-off, 1-on.",
			},

			"sasl_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SASL username.",
			},

			"sasl_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "SASL password.",
			},

			"write_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Enable access to certain fields of the log and check if they have been delivered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_headers": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "1: Enable 0: Do not enable.",
						},

						"enable_body": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "1: Enable 0: Do not enable.",
						},

						"enable_bot": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "1: Enable 0: Do not enable.",
						},
					},
				},
			},

			"flow_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID for post cls flow.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status 0- Off 1- On.",
			},
		},
	}
}

func resourceTencentCloudWafLogPostCkafkaFlowCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_ckafka_flow.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = wafv20180125.NewCreatePostCKafkaFlowRequest()
		flowId  int64
		logType int64
		hasFlow bool
	)

	if v, ok := d.GetOk("ckafka_region"); ok {
		request.CKafkaRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ckafka_id"); ok {
		request.CKafkaID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("brokers"); ok {
		request.Brokers = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compression"); ok {
		request.Compression = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("vip_type"); ok {
		request.VipType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("log_type"); ok {
		request.LogType = helper.IntInt64(v.(int))
		logType = int64(v.(int))
	}

	if v, ok := d.GetOk("topic"); ok {
		request.Topic = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		request.KafkaVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sasl_enable"); ok {
		request.SASLEnable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sasl_user"); ok {
		request.SASLUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sasl_password"); ok {
		request.SASLPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("write_config"); ok {
		for _, item := range v.([]interface{}) {
			if dMap, ok := item.(map[string]interface{}); ok && dMap != nil {
				config := wafv20180125.FieldWriteConfig{}
				if v, ok := dMap["enable_headers"].(int); ok {
					config.EnableHeaders = helper.IntInt64(v)
				}

				if v, ok := dMap["enable_body"].(int); ok {
					config.EnableBody = helper.IntInt64(v)
				}

				if v, ok := dMap["enable_bot"].(int); ok {
					config.EnableBot = helper.IntInt64(v)
				}

				request.WriteConfig = &config
			}
		}
	}

	// check unique first
	respData, err := service.DescribeWafLogPostCkafkaFlowById(ctx, logType)
	if err != nil {
		return err
	}

	if respData != nil && len(respData.PostCKafkaFlows) != 0 {
		return fmt.Errorf("In the case of `log_type` is %d, only one resource can be created and cannot be created multiple times.", logType)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreatePostCKafkaFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create waf log post ckafka flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// get flowId
	respData, err = service.DescribeWafLogPostCkafkaFlowById(ctx, logType)
	if err != nil {
		return err
	}

	if respData == nil || len(respData.PostCKafkaFlows) == 0 {
		return fmt.Errorf("resource `tencentcloud_waf_log_post_ckafka_flow` not found, please check if it has been deleted.")
	}

	for _, item := range respData.PostCKafkaFlows {
		if *item.LogType == logType {
			flowId = *item.FlowId
			hasFlow = true
			break
		}
	}

	if !hasFlow {
		return fmt.Errorf("resource `tencentcloud_waf_log_post_ckafka_flow` not found flowId, please check if it has been deleted.")
	}

	d.SetId(strings.Join([]string{helper.Int64ToStr(flowId), helper.Int64ToStr(logType)}, tccommon.FILED_SP))

	return resourceTencentCloudWafLogPostCkafkaFlowRead(d, meta)
}

func resourceTencentCloudWafLogPostCkafkaFlowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_ckafka_flow.read")()
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

	respData, err := service.DescribeWafLogPostCkafkaFlowById(ctx, helper.StrToInt64(logType))
	if err != nil {
		return err
	}

	if respData == nil || len(respData.PostCKafkaFlows) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `waf_log_post_ckafka_flow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	for _, item := range respData.PostCKafkaFlows {
		if *item.FlowId == helper.StrToInt64(flowId) && *item.LogType == helper.StrToInt64(logType) {
			if item.CKafkaRegion != nil {
				_ = d.Set("ckafka_region", item.CKafkaRegion)
			}

			if item.CKafkaID != nil {
				_ = d.Set("ckafka_id", item.CKafkaID)
			}

			if item.Brokers != nil {
				_ = d.Set("brokers", item.Brokers)
			}

			if item.Compression != nil {
				_ = d.Set("compression", item.Compression)
			}

			if item.VipType != nil {
				_ = d.Set("vip_type", item.VipType)
			}

			if item.LogType != nil {
				_ = d.Set("log_type", item.LogType)
			}

			if item.Topic != nil {
				_ = d.Set("topic", item.Topic)
			}

			if item.Version != nil {
				_ = d.Set("kafka_version", item.Version)
			}

			if item.SASLEnable != nil {
				_ = d.Set("sasl_enable", item.SASLEnable)
			}

			if item.SASLUser != nil {
				_ = d.Set("sasl_user", item.SASLUser)
			}

			if item.SASLPassword != nil {
				_ = d.Set("sasl_password", item.SASLPassword)
			}

			if item.WriteConfig != nil {
				tmpList := make([]map[string]interface{}, 0, 1)
				dMap := make(map[string]interface{})
				if item.WriteConfig.EnableHeaders != nil {
					dMap["enable_headers"] = item.WriteConfig.EnableHeaders
				}

				if item.WriteConfig.EnableBody != nil {
					dMap["enable_body"] = item.WriteConfig.EnableBody
				}

				if item.WriteConfig.EnableBot != nil {
					dMap["enable_bot"] = item.WriteConfig.EnableBot
				}

				tmpList = append(tmpList, dMap)
				_ = d.Set("write_config", tmpList)
			}

			if item.FlowId != nil {
				_ = d.Set("flow_id", item.FlowId)
			}

			if item.Status != nil {
				_ = d.Set("status", item.Status)
			}

			hasFlow = true
			break
		}
	}

	if !hasFlow {
		return fmt.Errorf("resource `waf_log_post_ckafka_flow` not found flowId, please check if it has been deleted.")
	}

	return nil
}

func resourceTencentCloudWafLogPostCkafkaFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_ckafka_flow.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = wafv20180125.NewCreatePostCKafkaFlowRequest()
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

	if v, ok := d.GetOk("ckafka_region"); ok {
		request.CKafkaRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ckafka_id"); ok {
		request.CKafkaID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("brokers"); ok {
		request.Brokers = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compression"); ok {
		request.Compression = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("vip_type"); ok {
		request.VipType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topic"); ok {
		request.Topic = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		request.KafkaVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sasl_enable"); ok {
		request.SASLEnable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sasl_user"); ok {
		request.SASLUser = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sasl_password"); ok {
		request.SASLPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("write_config"); ok {
		for _, item := range v.([]interface{}) {
			if dMap, ok := item.(map[string]interface{}); ok && dMap != nil {
				config := wafv20180125.FieldWriteConfig{}
				if v, ok := dMap["enable_headers"].(int); ok {
					config.EnableHeaders = helper.IntInt64(v)
				}

				if v, ok := dMap["enable_body"].(int); ok {
					config.EnableBody = helper.IntInt64(v)
				}

				if v, ok := dMap["enable_bot"].(int); ok {
					config.EnableBot = helper.IntInt64(v)
				}

				request.WriteConfig = &config
			}
		}
	}

	request.LogType = helper.StrToInt64Point(logType)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().CreatePostCKafkaFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf log post ckafka flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudWafLogPostCkafkaFlowRead(d, meta)
}

func resourceTencentCloudWafLogPostCkafkaFlowDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_log_post_ckafka_flow.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wafv20180125.NewDestroyPostCKafkaFlowRequest()
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().DestroyPostCKafkaFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf log post ckafka flow failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
