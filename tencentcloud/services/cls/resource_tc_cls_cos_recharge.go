// Code generated by iacg; DO NOT EDIT.
package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "topic id.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "logset id.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "recharge name.",
			},

			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cos bucket.",
			},

			"bucket_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cos bucket region.",
			},

			"prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cos file prefix.",
			},

			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "log type.",
			},

			"compress": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "supported gzip, lzop, snappy.",
			},

			"extract_rule_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
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
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "key list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
		topicId    string
		rechargeId string
	)
	var (
		request  = cls.NewCreateCosRechargeRequest()
		response = cls.NewCreateCosRechargeResponse()
	)

	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
	}

	if v, ok := d.GetOk("topic_id"); ok {
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

	if extractRuleInfoMap, ok := helper.InterfacesHeadMap(d, "extract_rule_info"); ok {
		extractRuleInfo := cls.ExtractRuleInfo{}
		if v, ok := extractRuleInfoMap["time_key"]; ok {
			extractRuleInfo.TimeKey = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["time_format"]; ok {
			extractRuleInfo.TimeFormat = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["delimiter"]; ok {
			extractRuleInfo.Delimiter = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["log_regex"]; ok {
			extractRuleInfo.LogRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["begin_regex"]; ok {
			extractRuleInfo.BeginRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["keys"]; ok {
			keysSet := v.(*schema.Set).List()
			for i := range keysSet {
				keys := keysSet[i].(string)
				extractRuleInfo.Keys = append(extractRuleInfo.Keys, helper.String(keys))
			}
		}
		if v, ok := extractRuleInfoMap["filter_key_regex"]; ok {
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
		if v, ok := extractRuleInfoMap["un_match_up_load_switch"]; ok {
			extractRuleInfo.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
		}
		if v, ok := extractRuleInfoMap["un_match_log_key"]; ok {
			extractRuleInfo.UnMatchLogKey = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["backtracking"]; ok {
			extractRuleInfo.Backtracking = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleInfoMap["is_gbk"]; ok {
			extractRuleInfo.IsGBK = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleInfoMap["json_standard"]; ok {
			extractRuleInfo.JsonStandard = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleInfoMap["protocol"]; ok {
			extractRuleInfo.Protocol = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["address"]; ok {
			extractRuleInfo.Address = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["parse_protocol"]; ok {
			extractRuleInfo.ParseProtocol = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["metadata_type"]; ok {
			extractRuleInfo.MetadataType = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleInfoMap["path_regex"]; ok {
			extractRuleInfo.PathRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleInfoMap["meta_tags"]; ok {
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
		log.Printf("[CRITAL]%s create cls cos recharge failed, reason:%+v", logId, err)
		return err
	}

	rechargeId = *response.Response.Id
	d.SetId(strings.Join([]string{topicId, rechargeId}, tccommon.FILED_SP))

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

	_ = d.Set("topic_id", topicId)

	respData, err := service.DescribeClsCosRechargeById(ctx, topicId, rechargeId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cls_cos_recharge` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.TopicId != nil {
		_ = d.Set("topic_id", respData.TopicId)
	}

	if respData.LogsetId != nil {
		_ = d.Set("logset_id", respData.LogsetId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Bucket != nil {
		_ = d.Set("bucket", respData.Bucket)
	}

	if respData.BucketRegion != nil {
		_ = d.Set("bucket_region", respData.BucketRegion)
	}

	if respData.Prefix != nil {
		_ = d.Set("prefix", respData.Prefix)
	}

	if respData.LogType != nil {
		_ = d.Set("log_type", respData.LogType)
	}

	if respData.Compress != nil {
		_ = d.Set("compress", respData.Compress)
	}

	extractRuleInfoMap := map[string]interface{}{}

	if respData.ExtractRuleInfo != nil {
		if respData.ExtractRuleInfo.TimeKey != nil {
			extractRuleInfoMap["time_key"] = respData.ExtractRuleInfo.TimeKey
		}

		if respData.ExtractRuleInfo.TimeFormat != nil {
			extractRuleInfoMap["time_format"] = respData.ExtractRuleInfo.TimeFormat
		}

		if respData.ExtractRuleInfo.Delimiter != nil {
			extractRuleInfoMap["delimiter"] = respData.ExtractRuleInfo.Delimiter
		}

		if respData.ExtractRuleInfo.LogRegex != nil {
			extractRuleInfoMap["log_regex"] = respData.ExtractRuleInfo.LogRegex
		}

		if respData.ExtractRuleInfo.BeginRegex != nil {
			extractRuleInfoMap["begin_regex"] = respData.ExtractRuleInfo.BeginRegex
		}

		if respData.ExtractRuleInfo.Keys != nil {
			extractRuleInfoMap["keys"] = respData.ExtractRuleInfo.Keys
		}

		filterKeyRegexList := make([]map[string]interface{}, 0, len(respData.ExtractRuleInfo.FilterKeyRegex))
		if respData.ExtractRuleInfo.FilterKeyRegex != nil {
			for _, filterKeyRegex := range respData.ExtractRuleInfo.FilterKeyRegex {
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

		if respData.ExtractRuleInfo.UnMatchUpLoadSwitch != nil {
			extractRuleInfoMap["un_match_up_load_switch"] = respData.ExtractRuleInfo.UnMatchUpLoadSwitch
		}

		if respData.ExtractRuleInfo.UnMatchLogKey != nil {
			extractRuleInfoMap["un_match_log_key"] = respData.ExtractRuleInfo.UnMatchLogKey
		}

		if respData.ExtractRuleInfo.Backtracking != nil {
			extractRuleInfoMap["backtracking"] = respData.ExtractRuleInfo.Backtracking
		}

		if respData.ExtractRuleInfo.IsGBK != nil {
			extractRuleInfoMap["is_gbk"] = respData.ExtractRuleInfo.IsGBK
		}

		if respData.ExtractRuleInfo.JsonStandard != nil {
			extractRuleInfoMap["json_standard"] = respData.ExtractRuleInfo.JsonStandard
		}

		if respData.ExtractRuleInfo.Protocol != nil {
			extractRuleInfoMap["protocol"] = respData.ExtractRuleInfo.Protocol
		}

		if respData.ExtractRuleInfo.Address != nil {
			extractRuleInfoMap["address"] = respData.ExtractRuleInfo.Address
		}

		if respData.ExtractRuleInfo.ParseProtocol != nil {
			extractRuleInfoMap["parse_protocol"] = respData.ExtractRuleInfo.ParseProtocol
		}

		if respData.ExtractRuleInfo.MetadataType != nil {
			extractRuleInfoMap["metadata_type"] = respData.ExtractRuleInfo.MetadataType
		}

		if respData.ExtractRuleInfo.PathRegex != nil {
			extractRuleInfoMap["path_regex"] = respData.ExtractRuleInfo.PathRegex
		}

		metaTagsList := make([]map[string]interface{}, 0, len(respData.ExtractRuleInfo.MetaTags))
		if respData.ExtractRuleInfo.MetaTags != nil {
			for _, metaTags := range respData.ExtractRuleInfo.MetaTags {
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

	immutableArgs := []string{"topic_id", "logset_id", "bucket", "bucket_region", "prefix", "log_type", "compress", "extract_rule_info"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	rechargeId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyCosRechargeRequest()

		request.TopicId = &topicId

		request.Id = &rechargeId

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
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
			log.Printf("[CRITAL]%s update cls cos recharge failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsCosRechargeRead(d, meta)
}

func resourceTencentCloudClsCosRechargeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_recharge.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
