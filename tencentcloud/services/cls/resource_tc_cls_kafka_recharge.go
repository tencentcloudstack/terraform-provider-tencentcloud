package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsKafkaRecharge() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsKafkaRechargeCreate,
		Read:   resourceTencentCloudClsKafkaRechargeRead,
		Update: resourceTencentCloudClsKafkaRechargeUpdate,
		Delete: resourceTencentCloudClsKafkaRechargeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "recharge for cls TopicId.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "kafka recharge name.",
			},

			"kafka_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "kafka recharge type, 0 for CKafka, 1 fro user define Kafka.",
			},

			"user_kafka_topics": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "user need recharge kafka topic list.",
			},

			"offset": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The translation is: -2: Earliest (default) -1: Latest.",
			},

			"kafka_instance": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CKafka Instance id.",
			},

			"server_addr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Server addr.",
			},

			"is_encryption_addr": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "ServerAddr is encryption addr.",
			},

			"protocol": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "encryption protocol.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "protocol type.",
						},
						"mechanism": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "encryption type.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "user password.",
						},
					},
				},
			},

			"consumer_group_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "user consumer group name.",
			},

			"log_recharge_rule": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "log recharge rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recharge_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "recharge type.",
						},
						"encoding_format": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "encoding format.",
						},
						"default_time_switch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "user default time.",
						},
						"log_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "log regex.",
						},
						"un_match_log_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "is push parse failed log.",
						},
						"un_match_log_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "parse failed log key.",
						},
						"un_match_log_time_src": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "parse failed log time from.",
						},
						"default_time_src": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "default time from.",
						},
						"time_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time key.",
						},
						"time_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time regex.",
						},
						"time_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time format.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time zone.",
						},
						"metadata": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "metadata.",
						},
						"keys": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "log key list.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsKafkaRechargeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_recharge.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateKafkaRechargeRequest()
		response = cls.NewCreateKafkaRechargeResponse()
		id       string
		topicId  string
	)
	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("kafka_type"); ok {
		request.KafkaType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user_kafka_topics"); ok {
		request.UserKafkaTopics = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("offset"); ok {
		request.Offset = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("kafka_instance"); ok {
		request.KafkaInstance = helper.String(v.(string))
	}

	if v, ok := d.GetOk("server_addr"); ok {
		request.ServerAddr = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_encryption_addr"); ok {
		request.IsEncryptionAddr = helper.Bool(v.(bool))
	}
	if dMap, ok := helper.InterfacesHeadMap(d, "protocol"); ok {
		kafkaProtocolInfo := cls.KafkaProtocolInfo{}
		if v, ok := dMap["protocol"]; ok {
			kafkaProtocolInfo.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["mechanism"]; ok {
			kafkaProtocolInfo.Mechanism = helper.String(v.(string))
		}
		if v, ok := dMap["user_name"]; ok {
			kafkaProtocolInfo.UserName = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			kafkaProtocolInfo.Password = helper.String(v.(string))
		}
		request.Protocol = &kafkaProtocolInfo
	}

	if v, ok := d.GetOk("consumer_group_name"); ok {
		request.ConsumerGroupName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "log_recharge_rule"); ok {
		logRechargeRuleInfo := cls.LogRechargeRuleInfo{}
		if v, ok := dMap["recharge_type"]; ok {
			logRechargeRuleInfo.RechargeType = helper.String(v.(string))
		}
		if v, ok := dMap["encoding_format"]; ok {
			logRechargeRuleInfo.EncodingFormat = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["default_time_switch"]; ok {
			logRechargeRuleInfo.DefaultTimeSwitch = helper.Bool(v.(bool))
		}
		if v, ok := dMap["log_regex"]; ok {
			logRechargeRuleInfo.LogRegex = helper.String(v.(string))
		}
		if v, ok := dMap["un_match_log_switch"]; ok {
			logRechargeRuleInfo.UnMatchLogSwitch = helper.Bool(v.(bool))
		}
		if v, ok := dMap["un_match_log_key"]; ok {
			logRechargeRuleInfo.UnMatchLogKey = helper.String(v.(string))
		}
		if v, ok := dMap["un_match_log_time_src"]; ok {
			logRechargeRuleInfo.UnMatchLogTimeSrc = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["default_time_src"]; ok {
			logRechargeRuleInfo.DefaultTimeSrc = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["time_key"]; ok {
			logRechargeRuleInfo.TimeKey = helper.String(v.(string))
		}
		if v, ok := dMap["time_regex"]; ok {
			logRechargeRuleInfo.TimeRegex = helper.String(v.(string))
		}
		if v, ok := dMap["time_format"]; ok {
			logRechargeRuleInfo.TimeFormat = helper.String(v.(string))
		}
		if v, ok := dMap["time_zone"]; ok {
			logRechargeRuleInfo.TimeZone = helper.String(v.(string))
		}
		if v, ok := dMap["metadata"]; ok {
			metadataSet := v.(*schema.Set).List()
			for i := range metadataSet {
				metadata := metadataSet[i].(string)
				logRechargeRuleInfo.Metadata = append(logRechargeRuleInfo.Metadata, &metadata)
			}
		}
		if v, ok := dMap["keys"]; ok {
			keysSet := v.(*schema.Set).List()
			for i := range keysSet {
				keys := keysSet[i].(string)
				logRechargeRuleInfo.Keys = append(logRechargeRuleInfo.Keys, &keys)
			}
		}
		request.LogRechargeRule = &logRechargeRuleInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateKafkaRecharge(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls kafkaRecharge failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id + tccommon.FILED_SP + topicId)

	return resourceTencentCloudClsKafkaRechargeRead(d, meta)
}

func resourceTencentCloudClsKafkaRechargeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_recharge.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	kafkaRechargeId := idSplit[0]
	kafkaTopic := idSplit[1]

	kafkaRecharge, err := service.DescribeClsKafkaRechargeById(ctx, kafkaRechargeId, kafkaTopic)
	if err != nil {
		return err
	}

	if kafkaRecharge == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsKafkaRecharge` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if kafkaRecharge.TopicId != nil {
		_ = d.Set("topic_id", kafkaRecharge.TopicId)
	}

	if kafkaRecharge.Name != nil {
		_ = d.Set("name", kafkaRecharge.Name)
	}

	if kafkaRecharge.KafkaType != nil {
		_ = d.Set("kafka_type", kafkaRecharge.KafkaType)
	}

	if kafkaRecharge.UserKafkaTopics != nil {
		_ = d.Set("user_kafka_topics", kafkaRecharge.UserKafkaTopics)
	}

	if kafkaRecharge.Offset != nil {
		_ = d.Set("offset", kafkaRecharge.Offset)
	}

	if kafkaRecharge.KafkaInstance != nil {
		_ = d.Set("kafka_instance", kafkaRecharge.KafkaInstance)
	}

	if kafkaRecharge.ServerAddr != nil {
		_ = d.Set("server_addr", kafkaRecharge.ServerAddr)
	}

	if kafkaRecharge.IsEncryptionAddr != nil {
		_ = d.Set("is_encryption_addr", kafkaRecharge.IsEncryptionAddr)
	}

	if kafkaRecharge.Protocol != nil {
		protocolMap := map[string]interface{}{}

		if kafkaRecharge.Protocol.Protocol != nil {
			protocolMap["protocol"] = kafkaRecharge.Protocol.Protocol
		}

		if kafkaRecharge.Protocol.Mechanism != nil {
			protocolMap["mechanism"] = kafkaRecharge.Protocol.Mechanism
		}

		if kafkaRecharge.Protocol.UserName != nil {
			protocolMap["user_name"] = kafkaRecharge.Protocol.UserName
		}

		if kafkaRecharge.Protocol.Password != nil {
			protocolMap["password"] = kafkaRecharge.Protocol.Password
		}

		_ = d.Set("protocol", []interface{}{protocolMap})
	}

	if kafkaRecharge.ConsumerGroupName != nil {
		_ = d.Set("consumer_group_name", kafkaRecharge.ConsumerGroupName)
	}

	if kafkaRecharge.LogRechargeRule != nil {
		logRechargeRuleMap := map[string]interface{}{}

		if kafkaRecharge.LogRechargeRule.RechargeType != nil {
			logRechargeRuleMap["recharge_type"] = kafkaRecharge.LogRechargeRule.RechargeType
		}

		if kafkaRecharge.LogRechargeRule.EncodingFormat != nil {
			logRechargeRuleMap["encoding_format"] = kafkaRecharge.LogRechargeRule.EncodingFormat
		}

		if kafkaRecharge.LogRechargeRule.DefaultTimeSwitch != nil {
			logRechargeRuleMap["default_time_switch"] = kafkaRecharge.LogRechargeRule.DefaultTimeSwitch
		}

		if kafkaRecharge.LogRechargeRule.LogRegex != nil {
			logRechargeRuleMap["log_regex"] = kafkaRecharge.LogRechargeRule.LogRegex
		}

		if kafkaRecharge.LogRechargeRule.UnMatchLogSwitch != nil {
			logRechargeRuleMap["un_match_log_switch"] = kafkaRecharge.LogRechargeRule.UnMatchLogSwitch
		}

		if kafkaRecharge.LogRechargeRule.UnMatchLogKey != nil {
			logRechargeRuleMap["un_match_log_key"] = kafkaRecharge.LogRechargeRule.UnMatchLogKey
		}

		if kafkaRecharge.LogRechargeRule.UnMatchLogTimeSrc != nil {
			logRechargeRuleMap["un_match_log_time_src"] = kafkaRecharge.LogRechargeRule.UnMatchLogTimeSrc
		}

		if kafkaRecharge.LogRechargeRule.DefaultTimeSrc != nil {
			logRechargeRuleMap["default_time_src"] = kafkaRecharge.LogRechargeRule.DefaultTimeSrc
		}

		if kafkaRecharge.LogRechargeRule.TimeKey != nil {
			logRechargeRuleMap["time_key"] = kafkaRecharge.LogRechargeRule.TimeKey
		}

		if kafkaRecharge.LogRechargeRule.TimeRegex != nil {
			logRechargeRuleMap["time_regex"] = kafkaRecharge.LogRechargeRule.TimeRegex
		}

		if kafkaRecharge.LogRechargeRule.TimeFormat != nil {
			logRechargeRuleMap["time_format"] = kafkaRecharge.LogRechargeRule.TimeFormat
		}

		if kafkaRecharge.LogRechargeRule.TimeZone != nil {
			logRechargeRuleMap["time_zone"] = kafkaRecharge.LogRechargeRule.TimeZone
		}

		if kafkaRecharge.LogRechargeRule.Metadata != nil {
			logRechargeRuleMap["metadata"] = kafkaRecharge.LogRechargeRule.Metadata
		}

		if kafkaRecharge.LogRechargeRule.Keys != nil {
			logRechargeRuleMap["keys"] = kafkaRecharge.LogRechargeRule.Keys
		}

		_ = d.Set("log_recharge_rule", []interface{}{logRechargeRuleMap})
	}

	return nil
}

func resourceTencentCloudClsKafkaRechargeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_recharge.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cls.NewModifyKafkaRechargeRequest()

	kafkaRechargeId := d.Id()

	request.Id = &kafkaRechargeId

	immutableArgs := []string{"topic_id", "name", "kafka_type", "user_kafka_topics", "offset", "kafka_instance", "server_addr", "is_encryption_addr", "protocol", "consumer_group_name", "log_recharge_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("topic_id") {
		if v, ok := d.GetOk("topic_id"); ok {
			request.TopicId = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("kafka_type") {
		if v, ok := d.GetOkExists("kafka_type"); ok {
			request.KafkaType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("user_kafka_topics") {
		if v, ok := d.GetOk("user_kafka_topics"); ok {
			request.UserKafkaTopics = helper.String(v.(string))
		}
	}

	if d.HasChange("kafka_instance") {
		if v, ok := d.GetOk("kafka_instance"); ok {
			request.KafkaInstance = helper.String(v.(string))
		}
	}

	if d.HasChange("server_addr") {
		if v, ok := d.GetOk("server_addr"); ok {
			request.ServerAddr = helper.String(v.(string))
		}
	}

	if d.HasChange("is_encryption_addr") {
		if v, ok := d.GetOkExists("is_encryption_addr"); ok {
			request.IsEncryptionAddr = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("protocol") {
		if dMap, ok := helper.InterfacesHeadMap(d, "protocol"); ok {
			kafkaProtocolInfo := cls.KafkaProtocolInfo{}
			if v, ok := dMap["protocol"]; ok {
				kafkaProtocolInfo.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["mechanism"]; ok {
				kafkaProtocolInfo.Mechanism = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				kafkaProtocolInfo.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				kafkaProtocolInfo.Password = helper.String(v.(string))
			}
			request.Protocol = &kafkaProtocolInfo
		}
	}

	if d.HasChange("consumer_group_name") {
		if v, ok := d.GetOk("consumer_group_name"); ok {
			request.ConsumerGroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("log_recharge_rule") {
		if dMap, ok := helper.InterfacesHeadMap(d, "log_recharge_rule"); ok {
			logRechargeRuleInfo := cls.LogRechargeRuleInfo{}
			if v, ok := dMap["recharge_type"]; ok {
				logRechargeRuleInfo.RechargeType = helper.String(v.(string))
			}
			if v, ok := dMap["encoding_format"]; ok {
				logRechargeRuleInfo.EncodingFormat = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["default_time_switch"]; ok {
				logRechargeRuleInfo.DefaultTimeSwitch = helper.Bool(v.(bool))
			}
			if v, ok := dMap["log_regex"]; ok {
				logRechargeRuleInfo.LogRegex = helper.String(v.(string))
			}
			if v, ok := dMap["un_match_log_switch"]; ok {
				logRechargeRuleInfo.UnMatchLogSwitch = helper.Bool(v.(bool))
			}
			if v, ok := dMap["un_match_log_key"]; ok {
				logRechargeRuleInfo.UnMatchLogKey = helper.String(v.(string))
			}
			if v, ok := dMap["un_match_log_time_src"]; ok {
				logRechargeRuleInfo.UnMatchLogTimeSrc = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["default_time_src"]; ok {
				logRechargeRuleInfo.DefaultTimeSrc = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["time_key"]; ok {
				logRechargeRuleInfo.TimeKey = helper.String(v.(string))
			}
			if v, ok := dMap["time_regex"]; ok {
				logRechargeRuleInfo.TimeRegex = helper.String(v.(string))
			}
			if v, ok := dMap["time_format"]; ok {
				logRechargeRuleInfo.TimeFormat = helper.String(v.(string))
			}
			if v, ok := dMap["time_zone"]; ok {
				logRechargeRuleInfo.TimeZone = helper.String(v.(string))
			}
			if v, ok := dMap["metadata"]; ok {
				metadataSet := v.(*schema.Set).List()
				for i := range metadataSet {
					metadata := metadataSet[i].(string)
					logRechargeRuleInfo.Metadata = append(logRechargeRuleInfo.Metadata, &metadata)
				}
			}
			if v, ok := dMap["keys"]; ok {
				keysSet := v.(*schema.Set).List()
				for i := range keysSet {
					keys := keysSet[i].(string)
					logRechargeRuleInfo.Keys = append(logRechargeRuleInfo.Keys, &keys)
				}
			}
			request.LogRechargeRule = &logRechargeRuleInfo
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyKafkaRecharge(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls kafkaRecharge failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsKafkaRechargeRead(d, meta)
}

func resourceTencentCloudClsKafkaRechargeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_kafka_recharge.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	kafkaRechargeId := idSplit[0]
	kafkaTopic := idSplit[1]

	if err := service.DeleteClsKafkaRechargeById(ctx, kafkaRechargeId, kafkaTopic); err != nil {
		return err
	}

	return nil
}
