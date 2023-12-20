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

func ResourceTencentCloudClsCosRecharge() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsCosRechargeCreate,
		Read:   resourceTencentCloudClsCosRechargeRead,
		Update: resourceTencentCloudClsCosRechargeUpdate,
		Delete: resourceTencentCloudClsCosRechargeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "topic id.",
			},

			"logset_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "logset id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "recharge name.",
			},

			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "cos bucket.",
			},

			"bucket_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "cos bucket region.",
			},

			"prefix": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "cos file prefix.",
			},

			"log_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "log type.",
			},

			"compress": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "supported gzip, lzop, snappy.",
			},

			"extract_rule_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "extract rule info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time key.",
						},
						"time_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "time format.",
						},
						"delimiter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "log delimiter.",
						},
						"log_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "log regex.",
						},
						"begin_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "begin line regex.",
						},
						"keys": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "key list.",
						},
						"filter_key_regex": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "rules that need to filter logs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "need filter log key.",
									},
									"regex": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "need filter log regex.",
									},
								},
							},
						},
						"un_match_up_load_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether to upload the parsing failure log.",
						},
						"un_match_log_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "parsing failure log key.",
						},
						"backtracking": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "backtracking data volume in incremental acquisition mode.",
						},
						"is_gbk": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "gbk encoding.",
						},
						"json_standard": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "is standard json.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog protocol.",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog address.",
						},
						"parse_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "parse protocol.",
						},
						"metadata_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "metadata type.",
						},
						"path_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "metadata path regex.",
						},
						"meta_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "metadata tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "metadata key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "metadata value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsCosRechargeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_recharge.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cls.NewCreateCosRechargeRequest()
		response   = cls.NewCreateCosRechargeResponse()
		topicId    string
		reChargeId string
	)
	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket"); ok {
		request.Bucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket_region"); ok {
		request.BucketRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("prefix"); ok {
		request.Prefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compress"); ok {
		request.Compress = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "extract_rule_info"); ok {
		extractRuleInfo := cls.ExtractRuleInfo{}
		if v, ok := dMap["time_key"]; ok {
			extractRuleInfo.TimeKey = helper.String(v.(string))
		}
		if v, ok := dMap["time_format"]; ok {
			extractRuleInfo.TimeFormat = helper.String(v.(string))
		}
		if v, ok := dMap["delimiter"]; ok {
			extractRuleInfo.Delimiter = helper.String(v.(string))
		}
		if v, ok := dMap["log_regex"]; ok {
			extractRuleInfo.LogRegex = helper.String(v.(string))
		}
		if v, ok := dMap["begin_regex"]; ok {
			extractRuleInfo.BeginRegex = helper.String(v.(string))
		}
		if v, ok := dMap["keys"]; ok {
			keysSet := v.(*schema.Set).List()
			for i := range keysSet {
				keys := keysSet[i].(string)
				extractRuleInfo.Keys = append(extractRuleInfo.Keys, &keys)
			}
		}
		if v, ok := dMap["filter_key_regex"]; ok {
			for _, item := range v.([]interface{}) {
				filterKeyRegexMap := item.(map[string]interface{})
				keyRegexInfo := cls.KeyRegexInfo{}
				if v, ok := filterKeyRegexMap["key"]; ok {
					keyRegexInfo.Key = helper.String(v.(string))
				}
				if v, ok := filterKeyRegexMap["regex"]; ok {
					keyRegexInfo.Regex = helper.String(v.(string))
				}
				extractRuleInfo.FilterKeyRegex = append(extractRuleInfo.FilterKeyRegex, &keyRegexInfo)
			}
		}
		if v, ok := dMap["un_match_up_load_switch"]; ok {
			extractRuleInfo.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
		}
		if v, ok := dMap["un_match_log_key"]; ok {
			extractRuleInfo.UnMatchLogKey = helper.String(v.(string))
		}
		if v, ok := dMap["backtracking"]; ok {
			extractRuleInfo.Backtracking = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["is_gbk"]; ok {
			extractRuleInfo.IsGBK = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["json_standard"]; ok {
			extractRuleInfo.JsonStandard = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["protocol"]; ok {
			extractRuleInfo.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["address"]; ok {
			extractRuleInfo.Address = helper.String(v.(string))
		}
		if v, ok := dMap["parse_protocol"]; ok {
			extractRuleInfo.ParseProtocol = helper.String(v.(string))
		}
		if v, ok := dMap["metadata_type"]; ok {
			extractRuleInfo.MetadataType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["path_regex"]; ok {
			extractRuleInfo.PathRegex = helper.String(v.(string))
		}
		if v, ok := dMap["meta_tags"]; ok {
			for _, item := range v.([]interface{}) {
				metaTagsMap := item.(map[string]interface{})
				metaTagInfo := cls.MetaTagInfo{}
				if v, ok := metaTagsMap["key"]; ok {
					metaTagInfo.Key = helper.String(v.(string))
				}
				if v, ok := metaTagsMap["value"]; ok {
					metaTagInfo.Value = helper.String(v.(string))
				}
				extractRuleInfo.MetaTags = append(extractRuleInfo.MetaTags, &metaTagInfo)
			}
		}
		request.ExtractRuleInfo = &extractRuleInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateCosRecharge(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls cosRecharge failed, reason:%+v", logId, err)
		return err
	}

	reChargeId = *response.Response.Id

	d.SetId(topicId + tccommon.FILED_SP + reChargeId)

	return resourceTencentCloudClsCosRechargeRead(d, meta)
}

func resourceTencentCloudClsCosRechargeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_recharge.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	rechargeId := idSplit[1]

	cosRecharge, err := service.DescribeClsCosRechargeById(ctx, topicId, rechargeId)
	if err != nil {
		return err
	}

	if cosRecharge == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsCosRecharge` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cosRecharge.TopicId != nil {
		_ = d.Set("topic_id", cosRecharge.TopicId)
	}

	if cosRecharge.LogsetId != nil {
		_ = d.Set("logset_id", cosRecharge.LogsetId)
	}

	if cosRecharge.Name != nil {
		_ = d.Set("name", cosRecharge.Name)
	}

	if cosRecharge.Bucket != nil {
		_ = d.Set("bucket", cosRecharge.Bucket)
	}

	if cosRecharge.BucketRegion != nil {
		_ = d.Set("bucket_region", cosRecharge.BucketRegion)
	}

	if cosRecharge.Prefix != nil {
		_ = d.Set("prefix", cosRecharge.Prefix)
	}

	if cosRecharge.LogType != nil {
		_ = d.Set("log_type", cosRecharge.LogType)
	}

	if cosRecharge.Compress != nil {
		_ = d.Set("compress", cosRecharge.Compress)
	}

	if cosRecharge.ExtractRuleInfo != nil {
		extractRuleInfoMap := map[string]interface{}{}

		if cosRecharge.ExtractRuleInfo.TimeKey != nil {
			extractRuleInfoMap["time_key"] = cosRecharge.ExtractRuleInfo.TimeKey
		}

		if cosRecharge.ExtractRuleInfo.TimeFormat != nil {
			extractRuleInfoMap["time_format"] = cosRecharge.ExtractRuleInfo.TimeFormat
		}

		if cosRecharge.ExtractRuleInfo.Delimiter != nil {
			extractRuleInfoMap["delimiter"] = cosRecharge.ExtractRuleInfo.Delimiter
		}

		if cosRecharge.ExtractRuleInfo.LogRegex != nil {
			extractRuleInfoMap["log_regex"] = cosRecharge.ExtractRuleInfo.LogRegex
		}

		if cosRecharge.ExtractRuleInfo.BeginRegex != nil {
			extractRuleInfoMap["begin_regex"] = cosRecharge.ExtractRuleInfo.BeginRegex
		}

		if cosRecharge.ExtractRuleInfo.Keys != nil {
			extractRuleInfoMap["keys"] = cosRecharge.ExtractRuleInfo.Keys
		}

		if cosRecharge.ExtractRuleInfo.FilterKeyRegex != nil {
			filterKeyRegexList := []interface{}{}
			for _, filterKeyRegex := range cosRecharge.ExtractRuleInfo.FilterKeyRegex {
				filterKeyRegexMap := map[string]interface{}{}

				if filterKeyRegex.Key != nil {
					filterKeyRegexMap["key"] = filterKeyRegex.Key
				}

				if filterKeyRegex.Regex != nil {
					filterKeyRegexMap["regex"] = filterKeyRegex.Regex
				}

				filterKeyRegexList = append(filterKeyRegexList, filterKeyRegexMap)
			}

			extractRuleInfoMap["filter_key_regex"] = filterKeyRegexList
		}

		if cosRecharge.ExtractRuleInfo.UnMatchUpLoadSwitch != nil {
			extractRuleInfoMap["un_match_up_load_switch"] = cosRecharge.ExtractRuleInfo.UnMatchUpLoadSwitch
		}

		if cosRecharge.ExtractRuleInfo.UnMatchLogKey != nil {
			extractRuleInfoMap["un_match_log_key"] = cosRecharge.ExtractRuleInfo.UnMatchLogKey
		}

		if cosRecharge.ExtractRuleInfo.Backtracking != nil {
			extractRuleInfoMap["backtracking"] = cosRecharge.ExtractRuleInfo.Backtracking
		}

		if cosRecharge.ExtractRuleInfo.IsGBK != nil {
			extractRuleInfoMap["is_gbk"] = cosRecharge.ExtractRuleInfo.IsGBK
		}

		if cosRecharge.ExtractRuleInfo.JsonStandard != nil {
			extractRuleInfoMap["json_standard"] = cosRecharge.ExtractRuleInfo.JsonStandard
		}

		if cosRecharge.ExtractRuleInfo.Protocol != nil {
			extractRuleInfoMap["protocol"] = cosRecharge.ExtractRuleInfo.Protocol
		}

		if cosRecharge.ExtractRuleInfo.Address != nil {
			extractRuleInfoMap["address"] = cosRecharge.ExtractRuleInfo.Address
		}

		if cosRecharge.ExtractRuleInfo.ParseProtocol != nil {
			extractRuleInfoMap["parse_protocol"] = cosRecharge.ExtractRuleInfo.ParseProtocol
		}

		if cosRecharge.ExtractRuleInfo.MetadataType != nil {
			extractRuleInfoMap["metadata_type"] = cosRecharge.ExtractRuleInfo.MetadataType
		}

		if cosRecharge.ExtractRuleInfo.PathRegex != nil {
			extractRuleInfoMap["path_regex"] = cosRecharge.ExtractRuleInfo.PathRegex
		}

		if cosRecharge.ExtractRuleInfo.MetaTags != nil {
			metaTagsList := []interface{}{}
			for _, metaTags := range cosRecharge.ExtractRuleInfo.MetaTags {
				metaTagsMap := map[string]interface{}{}

				if metaTags.Key != nil {
					metaTagsMap["key"] = metaTags.Key
				}

				if metaTags.Value != nil {
					metaTagsMap["value"] = metaTags.Value
				}

				metaTagsList = append(metaTagsList, metaTagsMap)
			}

			extractRuleInfoMap["meta_tags"] = metaTagsList
		}

		_ = d.Set("extract_rule_info", []interface{}{extractRuleInfoMap})
	}

	return nil
}

func resourceTencentCloudClsCosRechargeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_recharge.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cls.NewModifyCosRechargeRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	id := idSplit[1]

	immutableArgs := []string{
		"logset_id", "bucket", "bucket_region", "prefix",
		"log_type", "compress", "extract_rule_info",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TopicId = &topicId
	request.Id = &id

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyCosRecharge(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls cosRecharge failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsCosRechargeRead(d, meta)
}

func resourceTencentCloudClsCosRechargeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_recharge.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
