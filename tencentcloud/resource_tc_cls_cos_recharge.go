/*
Provides a resource to create a cls cos_recharge

Example Usage

```hcl
resource "tencentcloud_cls_cos_recharge" "cos_recharge" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  logset_id = "5cd3a17e-fb0b-418c-afd7-77b365397427"
  name = "test"
  bucket = "test-12345677"
  bucket_region = "ap-guangzhou"
  prefix = "/path"
  log_type = "json_log"
  compress = "gzip"
  extract_rule_info {
		time_key = "time"
		time_format = "YYYY-MM-DD HH:MM:SS"
		delimiter = ","
		log_regex = "*"
		begin_regex = "^*"
		keys =
		filter_key_regex {
			key = "testKey"
			regex = "testValue"
		}
		un_match_up_load_switch = false
		un_match_log_key = "test"
		backtracking = -1
		is_g_b_k = 0
		json_standard = 1
		protocol = "tcp"
		address = "127.0.0.1:9000"
		parse_protocol = "rfc3164"
		metadata_type = 0
		path_regex = "null"
		meta_tags {
			key = "testKey"
			value = "testValue"
		}

  }
}
```

Import

cls cos_recharge can be imported using the id, e.g.

```
terraform import tencentcloud_cls_cos_recharge.cos_recharge cos_recharge_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudClsCosRecharge() *schema.Resource {
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
				Description: "Topic id.",
			},

			"logset_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Logset id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Recharge name.",
			},

			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos bucket.",
			},

			"bucket_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos bucket region.",
			},

			"prefix": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos file prefix.",
			},

			"log_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log type.",
			},

			"compress": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Supported gzip, lzop, snappy.",
			},

			"extract_rule_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Extract rule info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time key.",
						},
						"time_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time format.",
						},
						"delimiter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log delimiter.",
						},
						"log_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log regex.",
						},
						"begin_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Begin line regex.",
						},
						"keys": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Key list.",
						},
						"filter_key_regex": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Rules that need to filter logs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Need filter log key.",
									},
									"regex": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Need filter log regex.",
									},
								},
							},
						},
						"un_match_up_load_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to upload the parsing failure log.",
						},
						"un_match_log_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parsing failure log key.",
						},
						"backtracking": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Backtracking data volume in incremental acquisition mode.",
						},
						"is_g_b_k": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Gbk encoding.",
						},
						"json_standard": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Is standard json.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Syslog protocol.",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Syslog address.",
						},
						"parse_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parse protocol.",
						},
						"metadata_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Metadata type.",
						},
						"path_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metadata path regex.",
						},
						"meta_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Metadata tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Metadata key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Metadata value.",
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
	defer logElapsed("resource.tencentcloud_cls_cos_recharge.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateCosRechargeRequest()
		response = cls.NewCreateCosRechargeResponse()
		topicId  string
		id       string
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
		if v, ok := dMap["is_g_b_k"]; ok {
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateCosRecharge(request)
		if e != nil {
			return retryError(e)
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

	topicId = *response.Response.TopicId
	d.SetId(strings.Join([]string{topicId, id}, FILED_SP))

	return resourceTencentCloudClsCosRechargeRead(d, meta)
}

func resourceTencentCloudClsCosRechargeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_recharge.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	id := idSplit[1]

	cosRecharge, err := service.DescribeClsCosRechargeById(ctx, topicId, id)
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

			extractRuleInfoMap["filter_key_regex"] = []interface{}{filterKeyRegexList}
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
			extractRuleInfoMap["is_g_b_k"] = cosRecharge.ExtractRuleInfo.IsGBK
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

			extractRuleInfoMap["meta_tags"] = []interface{}{metaTagsList}
		}

		_ = d.Set("extract_rule_info", []interface{}{extractRuleInfoMap})
	}

	return nil
}

func resourceTencentCloudClsCosRechargeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_recharge.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyCosRechargeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	id := idSplit[1]

	request.TopicId = &topicId
	request.Id = &id

	immutableArgs := []string{"topic_id", "logset_id", "name", "bucket", "bucket_region", "prefix", "log_type", "compress", "extract_rule_info"}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyCosRecharge(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_cls_cos_recharge.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicId := idSplit[0]
	id := idSplit[1]

	if err := service.DeleteClsCosRechargeById(ctx, topicId, id); err != nil {
		return err
	}

	return nil
}
