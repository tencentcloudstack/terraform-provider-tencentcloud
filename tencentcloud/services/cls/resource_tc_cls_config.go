package cls

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConfigCreate,
		Read:   resourceTencentCloudClsConfigRead,
		Delete: resourceTencentCloudClsConfigDelete,
		Update: resourceTencentCloudClsConfigUpdate,
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
				Description: "Log collection path containing the filename. Required for document collection.",
			},
			"log_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Type of the log to be collected. Valid values: json_log: log in JSON format; delimiter_log: log in delimited format; minimalist_log: minimalist log; multiline_log: log in multi-line format; " +
					"fullregex_log: log in full regex format. Default value: minimalist_log.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.",
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
							Description: "Unmatched log key. Required when UnMatchUpLoadSwitch is true.",
						},
						"backtracking": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).",
						},
						"is_gbk": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "GBK encoding. Default 0. Note: - Currently, when the value is 0, it means UTF-8 encoding.",
						},
						"json_standard": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "standard json. Default 0.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog protocol, tcp or udp. The value can be tcp or udp. It is effective only when LogType is service_syslog. Other types do not need to be filled in.",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog system log collection specifies the address and port that the collector listens to. This parameter is only valid when LogType is service_syslog. It does not need to be filled in for other types.",
						},
						"parse_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "parse protocol. This parameter is only valid when LogType is service_syslog. It does not need to be filled in for other types.",
						},
						"metadata_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "metadata type. 0: Do not use metadata information; 1: Use machine group metadata; 2: Use user-defined metadata; 3: Use collection configuration path. Note: COS import does not support this field.",
						},
						"path_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "metadata path regex.",
						},
						"meta_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "metadata tags. Note: - Required when MetadataType is 2. - COS import does not support this field.",
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
				Description: "Custom collection rule, which is a serialized JSON string. Required when LogType is user_define_log.",
			},
		},
	}
}

func resourceTencentCloudClsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateConfigRequest()
		response *cls.CreateConfigResponse
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

	if dMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
		extractRule := cls.ExtractRuleInfo{}
		if v, ok := dMap["time_key"].(string); ok && v != "" {
			extractRule.TimeKey = helper.String(v)
		}
		if v, ok := dMap["time_format"].(string); ok && v != "" {
			extractRule.TimeFormat = helper.String(v)
		}
		if v, ok := dMap["delimiter"].(string); ok && v != "" {
			extractRule.Delimiter = helper.String(v)
		}
		if v, ok := dMap["log_regex"].(string); ok && v != "" {
			extractRule.LogRegex = helper.String(v)
		}
		if v, ok := dMap["begin_regex"].(string); ok && v != "" {
			extractRule.BeginRegex = helper.String(v)
		}
		if v, ok := dMap["keys"]; ok {
			keys := v.(*schema.Set).List()
			for _, key := range keys {
				extractRule.Keys = append(extractRule.Keys, helper.String(key.(string)))
			}
		}
		if v, ok := dMap["filter_key_regex"]; ok {
			keyRegexs := make([]*cls.KeyRegexInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				keyRegex := cls.KeyRegexInfo{}
				if v, ok := dMap["key"]; ok {
					keyRegex.Key = helper.String(v.(string))
				}
				if v, ok := dMap["regex"]; ok {
					keyRegex.Regex = helper.String(v.(string))
				}
				keyRegexs = append(keyRegexs, &keyRegex)
			}
			extractRule.FilterKeyRegex = keyRegexs
		}
		if v, ok := dMap["un_match_up_load_switch"]; ok {
			extractRule.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
		}
		if v, ok := dMap["un_match_log_key"].(string); ok && v != "" {
			extractRule.UnMatchLogKey = helper.String(v)
		}
		if v, ok := dMap["backtracking"]; ok {
			extractRule.Backtracking = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["is_gbk"]; ok {
			extractRule.IsGBK = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["json_standard"]; ok {
			extractRule.JsonStandard = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["protocol"].(string); ok && v != "" {
			extractRule.Protocol = helper.String(v)
		}
		if v, ok := dMap["address"].(string); ok && v != "" {
			extractRule.Address = helper.String(v)
		}
		if v, ok := dMap["parse_protocol"].(string); ok && v != "" {
			extractRule.ParseProtocol = helper.String(v)
		}
		if v, ok := dMap["metadata_type"]; ok {
			extractRule.MetadataType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["path_regex"].(string); ok && v != "" {
			extractRule.PathRegex = helper.String(v)
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
				extractRule.MetaTags = append(extractRule.MetaTags, &metaTagInfo)
			}
		}
		request.ExtractRule = &extractRule
	}
	if v, ok := d.GetOk("exclude_paths"); ok {
		excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			excludePath := cls.ExcludePathInfo{}
			if v, ok := dMap["type"].(string); ok && v != "" {
				excludePath.Type = helper.String(v)
			}
			if v, ok := dMap["value"].(string); ok && v != "" {
				excludePath.Value = helper.String(v)
			}
			excludePaths = append(excludePaths, &excludePath)
		}
		request.ExcludePaths = excludePaths
	}
	if v, ok := d.GetOk("user_define_rule"); ok {
		request.UserDefineRule = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls config extra failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.ConfigId
	d.SetId(id)
	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	configId := d.Id()

	config, err := service.DescribeClsConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if config == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if config.Name != nil {
		_ = d.Set("name", config.Name)
	}

	if config.Output != nil {
		_ = d.Set("output", config.Output)
	}

	if config.Path != nil {
		_ = d.Set("path", config.Path)
	}

	if config.LogType != nil {
		_ = d.Set("log_type", config.LogType)
	}

	if config.ExtractRule != nil {
		extractRuleMap := map[string]interface{}{}

		if config.ExtractRule.TimeKey != nil {
			extractRuleMap["time_key"] = config.ExtractRule.TimeKey
		}

		if config.ExtractRule.TimeFormat != nil {
			extractRuleMap["time_format"] = config.ExtractRule.TimeFormat
		}

		if config.ExtractRule.Delimiter != nil {
			extractRuleMap["delimiter"] = config.ExtractRule.Delimiter
		}

		if config.ExtractRule.LogRegex != nil {
			extractRuleMap["log_regex"] = config.ExtractRule.LogRegex
		}

		if config.ExtractRule.BeginRegex != nil {
			extractRuleMap["begin_regex"] = config.ExtractRule.BeginRegex
		}

		if config.ExtractRule.Keys != nil {
			extractRuleMap["keys"] = config.ExtractRule.Keys
		}

		if config.ExtractRule.FilterKeyRegex != nil {
			filterKeyRegexList := []interface{}{}
			for _, filterKeyRegex := range config.ExtractRule.FilterKeyRegex {
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

		if config.ExtractRule.UnMatchUpLoadSwitch != nil {
			extractRuleMap["un_match_up_load_switch"] = config.ExtractRule.UnMatchUpLoadSwitch
		}

		if config.ExtractRule.UnMatchLogKey != nil {
			extractRuleMap["un_match_log_key"] = config.ExtractRule.UnMatchLogKey
		}

		if config.ExtractRule.Backtracking != nil {
			extractRuleMap["backtracking"] = config.ExtractRule.Backtracking
		}

		if config.ExtractRule.IsGBK != nil {
			extractRuleMap["is_gbk"] = config.ExtractRule.IsGBK
		}

		if config.ExtractRule.JsonStandard != nil {
			extractRuleMap["json_standard"] = config.ExtractRule.JsonStandard
		}

		if config.ExtractRule.Protocol != nil {
			extractRuleMap["protocol"] = config.ExtractRule.Protocol
		}

		if config.ExtractRule.Address != nil {
			extractRuleMap["address"] = config.ExtractRule.Address
		}

		if config.ExtractRule.ParseProtocol != nil {
			extractRuleMap["parse_protocol"] = config.ExtractRule.ParseProtocol
		}

		if config.ExtractRule.MetadataType != nil {
			extractRuleMap["metadata_type"] = config.ExtractRule.MetadataType
		}

		if config.ExtractRule.PathRegex != nil {
			extractRuleMap["path_regex"] = config.ExtractRule.PathRegex
		}

		if config.ExtractRule.MetaTags != nil {
			metaTagsList := []interface{}{}
			for _, metaTags := range config.ExtractRule.MetaTags {
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

	if config.ExcludePaths != nil {
		excludePathsList := []interface{}{}
		for _, excludePath := range config.ExcludePaths {
			excludePathsMap := map[string]interface{}{}

			if excludePath.Type != nil {
				excludePathsMap["type"] = excludePath.Type
			}

			if excludePath.Value != nil {
				excludePathsMap["value"] = excludePath.Value
			}

			excludePathsList = append(excludePathsList, excludePathsMap)
		}

		_ = d.Set("exclude_paths", excludePathsList)
	}

	if config.UserDefineRule != nil {
		_ = d.Set("user_define_rule", config.UserDefineRule)
	}

	return nil
}

func resourceTencentCloudClsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cls.NewModifyConfigRequest()

	request.ConfigId = helper.String(d.Id())

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}
	if d.HasChange("output") {
		if v, ok := d.GetOk("output"); ok {
			request.Output = helper.String(v.(string))
		}
	}
	if d.HasChange("path") {
		if v, ok := d.GetOk("path"); ok {
			request.Path = helper.String(v.(string))
		}
	}
	if d.HasChange("log_type") || d.HasChange("extract_rule") {
		if v, ok := d.GetOk("log_type"); ok {
			request.LogType = helper.String(v.(string))
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
			extractRule := cls.ExtractRuleInfo{}
			if v, ok := dMap["time_key"].(string); ok && v != "" {
				extractRule.TimeKey = helper.String(v)
			}
			if v, ok := dMap["time_format"].(string); ok && v != "" {
				extractRule.TimeFormat = helper.String(v)
			}
			if v, ok := dMap["delimiter"].(string); ok && v != "" {
				extractRule.Delimiter = helper.String(v)
			}
			if v, ok := dMap["log_regex"].(string); ok && v != "" {
				extractRule.LogRegex = helper.String(v)
			}
			if v, ok := dMap["begin_regex"].(string); ok && v != "" {
				extractRule.BeginRegex = helper.String(v)
			}
			if v, ok := dMap["keys"]; ok {
				keys := v.(*schema.Set).List()
				for _, key := range keys {
					extractRule.Keys = append(extractRule.Keys, helper.String(key.(string)))
				}
			}
			if v, ok := dMap["filter_key_regex"]; ok {
				keyRegexs := make([]*cls.KeyRegexInfo, 0, 10)
				for _, item := range v.([]interface{}) {
					dMap := item.(map[string]interface{})
					keyRegex := cls.KeyRegexInfo{}
					if v, ok := dMap["key"]; ok {
						keyRegex.Key = helper.String(v.(string))
					}
					if v, ok := dMap["regex"]; ok {
						keyRegex.Regex = helper.String(v.(string))
					}
					keyRegexs = append(keyRegexs, &keyRegex)
				}
				extractRule.FilterKeyRegex = keyRegexs
			}
			if v, ok := dMap["un_match_up_load_switch"]; ok {
				extractRule.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
			}
			if v, ok := dMap["un_match_log_key"].(string); ok && v != "" {
				extractRule.UnMatchLogKey = helper.String(v)
			}
			if v, ok := dMap["backtracking"]; ok {
				extractRule.Backtracking = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["is_gbk"]; ok {
				extractRule.IsGBK = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["json_standard"]; ok {
				extractRule.JsonStandard = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["protocol"].(string); ok && v != "" {
				extractRule.Protocol = helper.String(v)
			}
			if v, ok := dMap["address"].(string); ok && v != "" {
				extractRule.Address = helper.String(v)
			}
			if v, ok := dMap["parse_protocol"].(string); ok && v != "" {
				extractRule.ParseProtocol = helper.String(v)
			}
			if v, ok := dMap["metadata_type"]; ok {
				extractRule.MetadataType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["path_regex"].(string); ok && v != "" {
				extractRule.PathRegex = helper.String(v)
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
					extractRule.MetaTags = append(extractRule.MetaTags, &metaTagInfo)
				}
			}
			request.ExtractRule = &extractRule
		}
	}
	if d.HasChange("exclude_paths") {
		if v, ok := d.GetOk("exclude_paths"); ok {
			excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				excludePath := cls.ExcludePathInfo{}
				if v, ok := dMap["type"].(string); ok && v != "" {
					excludePath.Type = helper.String(v)
				}
				if v, ok := dMap["value"].(string); ok && v != "" {
					excludePath.Value = helper.String(v)
				}
				excludePaths = append(excludePaths, &excludePath)
			}
			request.ExcludePaths = excludePaths
		}
	}

	if d.HasChange("user_define_rule") {
		if v, ok := d.GetOk("user_define_rule"); ok {
			request.UserDefineRule = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if err := service.DeleteClsConfig(ctx, id); err != nil {
		return err
	}

	return nil
}
