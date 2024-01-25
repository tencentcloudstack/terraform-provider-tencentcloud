// Code generated by iacg; DO NOT EDIT.
package cls

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConfigCreate,
		Read:   resourceTencentCloudClsConfigRead,
		Update: resourceTencentCloudClsConfigUpdate,
		Delete: resourceTencentCloudClsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection configuration name.",
			},

			"output": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log topic ID (TopicId) of collection configuration.",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log collection path containing the filename.",
			},

			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the log to be collected. Valid values: json_log: log in JSON format; delimiter_log: log in delimited format; minimalist_log: minimalist log; multiline_log: log in multi-line format; fullregex_log: log in full regex format. Default value: minimalist_log.",
			},

			"extract_rule": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Extraction rule. If ExtractRule is set, LogType must be set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time field key name. time_key and time_format must appear in pair.",
						},
						"time_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time field format. For more information, please see the output parameters of the time format description of the strftime function in C language.",
						},
						"delimiter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Delimiter for delimited log, which is valid only if log_type is delimiter_log.",
						},
						"log_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Full log matching rule, which is valid only if log_type is fullregex_log.",
						},
						"begin_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "First-Line matching rule, which is valid only if log_type is multiline_log or fullregex_log.",
						},
						"keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"filter_key_regex": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Log keys to be filtered and the corresponding regex.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Log key to be filtered.",
									},
									"regex": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filter rule regex corresponding to key.",
									},
								},
							},
						},
						"un_match_up_load_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to upload the logs that failed to be parsed. Valid values: true: yes; false: no.",
						},
						"un_match_log_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unmatched log key.",
						},
						"backtracking": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).",
						},
						"is_gbk": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "GBK encoding. Default 0.",
						},
						"json_standard": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "standard json. Default 0.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog protocol, tcp or udp.",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog system log collection specifies the address and port that the collector listens to.",
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
							Description: "metadata tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "tag value.",
									},
								},
							},
						},
					},
				},
			},

			"exclude_paths": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Collection path blocklist.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type. Valid values: File, Path.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specific content corresponding to Type.",
						},
					},
				},
			},

			"user_define_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom collection rule, which is a serialized JSON string.",
			},
		},
	}
}

func resourceTencentCloudClsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		configId string
	)
	var (
		request  = cls.NewCreateConfigRequest()
		response = cls.NewCreateConfigResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output"); ok {
		request.Output = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if extractRuleMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
		extractRuleInfo := cls.ExtractRuleInfo{}
		if v, ok := extractRuleMap["time_key"]; ok {
			extractRuleInfo.TimeKey = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["time_format"]; ok {
			extractRuleInfo.TimeFormat = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["delimiter"]; ok {
			extractRuleInfo.Delimiter = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["log_regex"]; ok {
			extractRuleInfo.LogRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["begin_regex"]; ok {
			extractRuleInfo.BeginRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["keys"]; ok {
			keysSet := v.(*schema.Set).List()
			for i := range keysSet {
				keys := keysSet[i].(string)
				extractRuleInfo.Keys = append(extractRuleInfo.Keys, helper.String(keys))
			}
		}
		if v, ok := extractRuleMap["filter_key_regex"]; ok {
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
		if v, ok := extractRuleMap["un_match_up_load_switch"]; ok {
			extractRuleInfo.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
		}
		if v, ok := extractRuleMap["un_match_log_key"]; ok {
			extractRuleInfo.UnMatchLogKey = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["backtracking"]; ok {
			extractRuleInfo.Backtracking = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleMap["is_gbk"]; ok {
			extractRuleInfo.IsGBK = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleMap["json_standard"]; ok {
			extractRuleInfo.JsonStandard = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleMap["protocol"]; ok {
			extractRuleInfo.Protocol = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["address"]; ok {
			extractRuleInfo.Address = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["parse_protocol"]; ok {
			extractRuleInfo.ParseProtocol = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["metadata_type"]; ok {
			extractRuleInfo.MetadataType = helper.IntInt64(v.(int))
		}
		if v, ok := extractRuleMap["path_regex"]; ok {
			extractRuleInfo.PathRegex = helper.String(v.(string))
		}
		if v, ok := extractRuleMap["meta_tags"]; ok {
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
		request.ExtractRule = &extractRuleInfo
	}

	if v, ok := d.GetOk("exclude_paths"); ok {
		for _, item := range v.([]interface{}) {
			excludePathsMap := item.(map[string]interface{})
			excludePathInfo := cls.ExcludePathInfo{}
			if v, ok := excludePathsMap["type"]; ok {
				excludePathInfo.Type = helper.String(v.(string))
			}
			if v, ok := excludePathsMap["value"]; ok {
				excludePathInfo.Value = helper.String(v.(string))
			}
			request.ExcludePaths = append(request.ExcludePaths, &excludePathInfo)
		}
	}

	if v, ok := d.GetOk("user_define_rule"); ok {
		request.UserDefineRule = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls config failed, reason:%+v", logId, err)
		return err
	}

	configId = *response.Response.ConfigId
	d.SetId(configId)

	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	configId := d.Id()

	respData, err := service.DescribeClsConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cls_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Output != nil {
		_ = d.Set("output", respData.Output)
	}

	if respData.Path != nil {
		_ = d.Set("path", respData.Path)
	}

	if respData.LogType != nil {
		_ = d.Set("log_type", respData.LogType)
	}

	extractRuleMap := map[string]interface{}{}

	if respData.ExtractRule != nil {
		if respData.ExtractRule.TimeKey != nil {
			extractRuleMap["time_key"] = respData.ExtractRule.TimeKey
		}

		if respData.ExtractRule.TimeFormat != nil {
			extractRuleMap["time_format"] = respData.ExtractRule.TimeFormat
		}

		if respData.ExtractRule.Delimiter != nil {
			extractRuleMap["delimiter"] = respData.ExtractRule.Delimiter
		}

		if respData.ExtractRule.LogRegex != nil {
			extractRuleMap["log_regex"] = respData.ExtractRule.LogRegex
		}

		if respData.ExtractRule.BeginRegex != nil {
			extractRuleMap["begin_regex"] = respData.ExtractRule.BeginRegex
		}

		if respData.ExtractRule.Keys != nil {
			extractRuleMap["keys"] = respData.ExtractRule.Keys
		}

		filterKeyRegexList := make([]map[string]interface{}, 0, len(respData.ExtractRule.FilterKeyRegex))
		if respData.ExtractRule.FilterKeyRegex != nil {
			for _, filterKeyRegex := range respData.ExtractRule.FilterKeyRegex {
				filterKeyRegexMap := map[string]interface{}{}

				if filterKeyRegex.Key != nil {
					filterKeyRegexMap["key"] = filterKeyRegex.Key
				}

				if filterKeyRegex.Regex != nil {
					filterKeyRegexMap["regex"] = filterKeyRegex.Regex
				}

				filterKeyRegexList = append(filterKeyRegexList, filterKeyRegexMap)
			}

			extractRuleMap["filter_key_regex"] = filterKeyRegexList
		}

		if respData.ExtractRule.UnMatchUpLoadSwitch != nil {
			extractRuleMap["un_match_up_load_switch"] = respData.ExtractRule.UnMatchUpLoadSwitch
		}

		if respData.ExtractRule.UnMatchLogKey != nil {
			extractRuleMap["un_match_log_key"] = respData.ExtractRule.UnMatchLogKey
		}

		if respData.ExtractRule.Backtracking != nil {
			extractRuleMap["backtracking"] = respData.ExtractRule.Backtracking
		}

		if respData.ExtractRule.IsGBK != nil {
			extractRuleMap["is_gbk"] = respData.ExtractRule.IsGBK
		}

		if respData.ExtractRule.JsonStandard != nil {
			extractRuleMap["json_standard"] = respData.ExtractRule.JsonStandard
		}

		if respData.ExtractRule.Protocol != nil {
			extractRuleMap["protocol"] = respData.ExtractRule.Protocol
		}

		if respData.ExtractRule.Address != nil {
			extractRuleMap["address"] = respData.ExtractRule.Address
		}

		if respData.ExtractRule.ParseProtocol != nil {
			extractRuleMap["parse_protocol"] = respData.ExtractRule.ParseProtocol
		}

		if respData.ExtractRule.MetadataType != nil {
			extractRuleMap["metadata_type"] = respData.ExtractRule.MetadataType
		}

		if respData.ExtractRule.PathRegex != nil {
			extractRuleMap["path_regex"] = respData.ExtractRule.PathRegex
		}

		metaTagsList := make([]map[string]interface{}, 0, len(respData.ExtractRule.MetaTags))
		if respData.ExtractRule.MetaTags != nil {
			for _, metaTags := range respData.ExtractRule.MetaTags {
				metaTagsMap := map[string]interface{}{}

				if metaTags.Key != nil {
					metaTagsMap["key"] = metaTags.Key
				}

				if metaTags.Value != nil {
					metaTagsMap["value"] = metaTags.Value
				}

				metaTagsList = append(metaTagsList, metaTagsMap)
			}

			extractRuleMap["meta_tags"] = metaTagsList
		}

		_ = d.Set("extract_rule", []interface{}{extractRuleMap})
	}

	excludePathsList := make([]map[string]interface{}, 0, len(respData.ExcludePaths))
	if respData.ExcludePaths != nil {
		for _, excludePaths := range respData.ExcludePaths {
			excludePathsMap := map[string]interface{}{}

			if excludePaths.Type != nil {
				excludePathsMap["type"] = excludePaths.Type
			}

			if excludePaths.Value != nil {
				excludePathsMap["value"] = excludePaths.Value
			}

			excludePathsList = append(excludePathsList, excludePathsMap)
		}

		_ = d.Set("exclude_paths", excludePathsList)
	}

	if respData.UserDefineRule != nil {
		_ = d.Set("user_define_rule", respData.UserDefineRule)
	}

	return nil
}

func resourceTencentCloudClsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	configId := d.Id()

	needChange := false
	mutableArgs := []string{"name", "output", "path", "log_type", "extract_rule", "exclude_paths", "user_define_rule"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyConfigRequest()

		request.ConfigId = &configId

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("output"); ok {
			request.Output = helper.String(v.(string))
		}

		if v, ok := d.GetOk("path"); ok {
			request.Path = helper.String(v.(string))
		}

		if v, ok := d.GetOk("log_type"); ok {
			request.LogType = helper.String(v.(string))
		}

		if extractRuleMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
			extractRuleInfo := cls.ExtractRuleInfo{}
			if v, ok := extractRuleMap["time_key"]; ok {
				extractRuleInfo.TimeKey = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["time_format"]; ok {
				extractRuleInfo.TimeFormat = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["delimiter"]; ok {
				extractRuleInfo.Delimiter = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["log_regex"]; ok {
				extractRuleInfo.LogRegex = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["begin_regex"]; ok {
				extractRuleInfo.BeginRegex = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["keys"]; ok {
				keysSet := v.(*schema.Set).List()
				for i := range keysSet {
					keys := keysSet[i].(string)
					extractRuleInfo.Keys = append(extractRuleInfo.Keys, helper.String(keys))
				}
			}
			if v, ok := extractRuleMap["filter_key_regex"]; ok {
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
			if v, ok := extractRuleMap["un_match_up_load_switch"]; ok {
				extractRuleInfo.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
			}
			if v, ok := extractRuleMap["un_match_log_key"]; ok {
				extractRuleInfo.UnMatchLogKey = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["backtracking"]; ok {
				extractRuleInfo.Backtracking = helper.IntInt64(v.(int))
			}
			if v, ok := extractRuleMap["is_gbk"]; ok {
				extractRuleInfo.IsGBK = helper.IntInt64(v.(int))
			}
			if v, ok := extractRuleMap["json_standard"]; ok {
				extractRuleInfo.JsonStandard = helper.IntInt64(v.(int))
			}
			if v, ok := extractRuleMap["protocol"]; ok {
				extractRuleInfo.Protocol = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["address"]; ok {
				extractRuleInfo.Address = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["parse_protocol"]; ok {
				extractRuleInfo.ParseProtocol = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["metadata_type"]; ok {
				extractRuleInfo.MetadataType = helper.IntInt64(v.(int))
			}
			if v, ok := extractRuleMap["path_regex"]; ok {
				extractRuleInfo.PathRegex = helper.String(v.(string))
			}
			if v, ok := extractRuleMap["meta_tags"]; ok {
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
			request.ExtractRule = &extractRuleInfo
		}

		if v, ok := d.GetOk("exclude_paths"); ok {
			for _, item := range v.([]interface{}) {
				excludePathsMap := item.(map[string]interface{})
				excludePathInfo := cls.ExcludePathInfo{}
				if v, ok := excludePathsMap["type"]; ok {
					excludePathInfo.Type = helper.String(v.(string))
				}
				if v, ok := excludePathsMap["value"]; ok {
					excludePathInfo.Value = helper.String(v.(string))
				}
				request.ExcludePaths = append(request.ExcludePaths, &excludePathInfo)
			}
		}

		if v, ok := d.GetOk("user_define_rule"); ok {
			request.UserDefineRule = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyConfig(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cls config failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	configId := d.Id()

	if err := service.DeleteClsConfigById(ctx, configId); err != nil {
		return err
	}

	return nil
}
