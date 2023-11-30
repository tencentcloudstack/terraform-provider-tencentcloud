package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClsMachineGroupConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsMachineGroupConfigsRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group id.",
			},

			"configs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "scrape config list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "scrape config id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "scrape config name.",
						},
						"log_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "style of log format.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "scrape log path.",
						},
						"log_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "log type.",
						},
						"extract_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extraction rule. If ExtractRule is set, LogType must be set.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time field key name. time_key and time_format must appear in pair.",
									},
									"time_format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time field format. For more information, please see the output parameters of the time format description of the strftime function in C language.",
									},
									"delimiter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Delimiter for delimited log, which is valid only if log_type is delimiter_log.",
									},
									"log_regex": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Full log matching rule, which is valid only if log_type is fullregex_log.",
									},
									"begin_regex": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "First-Line matching rule, which is valid only if log_type is multiline_log or fullregex_log.",
									},
									"keys": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Key name of each extracted field. An empty key indicates to discard the field. This parameter is valid only if log_type is delimiter_log. json_log logs use the key of JSON itself.",
									},
									"filter_key_regex": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Log keys to be filtered and the corresponding regex.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Log key to be filtered.",
												},
												"regex": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Filter rule regex corresponding to key.",
												},
											},
										},
									},
									"un_match_up_load_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to upload the logs that failed to be parsed. Valid values: true: yes; false: no.",
									},
									"un_match_log_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unmatched log key.",
									},
									"backtracking": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).",
									},
									"is_gbk": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "GBK encoding. Default 0.",
									},
									"json_standard": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "standard json. Default 0.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "syslog protocol, tcp or udp.",
									},
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "syslog system log collection specifies the address and port that the collector listens to.",
									},
									"parse_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "parse protocol.",
									},
									"metadata_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "metadata type.",
									},
									"path_regex": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "metadata path regex.",
									},
									"meta_tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "metadata tags.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "tag key.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
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
							Computed:    true,
							Description: "Collection path blocklist.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type. Valid values: File, Path.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specific content corresponding to Type.",
									},
								},
							},
						},
						"output": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "topicid.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"user_define_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "user define rule.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudClsMachineGroupConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cls_machine_group_configs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var configs []*cls.ConfigInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsMachineGroupConfigsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		configs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(configs))
	tmpList := make([]map[string]interface{}, 0, len(configs))

	if configs != nil {
		for _, configInfo := range configs {
			configInfoMap := map[string]interface{}{}

			if configInfo.ConfigId != nil {
				configInfoMap["config_id"] = configInfo.ConfigId
			}

			if configInfo.Name != nil {
				configInfoMap["name"] = configInfo.Name
			}

			if configInfo.LogFormat != nil {
				configInfoMap["log_format"] = configInfo.LogFormat
			}

			if configInfo.Path != nil {
				configInfoMap["path"] = configInfo.Path
			}

			if configInfo.LogType != nil {
				configInfoMap["log_type"] = configInfo.LogType
			}

			if configInfo.ExtractRule != nil {
				extractRuleMap := map[string]interface{}{}

				if configInfo.ExtractRule.TimeKey != nil {
					extractRuleMap["time_key"] = configInfo.ExtractRule.TimeKey
				}

				if configInfo.ExtractRule.TimeFormat != nil {
					extractRuleMap["time_format"] = configInfo.ExtractRule.TimeFormat
				}

				if configInfo.ExtractRule.Delimiter != nil {
					extractRuleMap["delimiter"] = configInfo.ExtractRule.Delimiter
				}

				if configInfo.ExtractRule.LogRegex != nil {
					extractRuleMap["log_regex"] = configInfo.ExtractRule.LogRegex
				}

				if configInfo.ExtractRule.BeginRegex != nil {
					extractRuleMap["begin_regex"] = configInfo.ExtractRule.BeginRegex
				}

				if configInfo.ExtractRule.Keys != nil {
					extractRuleMap["keys"] = configInfo.ExtractRule.Keys
				}

				if configInfo.ExtractRule.FilterKeyRegex != nil {
					filterKeyRegexList := []interface{}{}
					for _, filterKeyRegex := range configInfo.ExtractRule.FilterKeyRegex {
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

				if configInfo.ExtractRule.UnMatchUpLoadSwitch != nil {
					extractRuleMap["un_match_up_load_switch"] = configInfo.ExtractRule.UnMatchUpLoadSwitch
				}

				if configInfo.ExtractRule.UnMatchLogKey != nil {
					extractRuleMap["un_match_log_key"] = configInfo.ExtractRule.UnMatchLogKey
				}

				if configInfo.ExtractRule.Backtracking != nil {
					extractRuleMap["backtracking"] = configInfo.ExtractRule.Backtracking
				}

				if configInfo.ExtractRule.IsGBK != nil {
					extractRuleMap["is_gbk"] = configInfo.ExtractRule.IsGBK
				}

				if configInfo.ExtractRule.JsonStandard != nil {
					extractRuleMap["json_standard"] = configInfo.ExtractRule.JsonStandard
				}

				if configInfo.ExtractRule.Protocol != nil {
					extractRuleMap["protocol"] = configInfo.ExtractRule.Protocol
				}

				if configInfo.ExtractRule.Address != nil {
					extractRuleMap["address"] = configInfo.ExtractRule.Address
				}

				if configInfo.ExtractRule.ParseProtocol != nil {
					extractRuleMap["parse_protocol"] = configInfo.ExtractRule.ParseProtocol
				}

				if configInfo.ExtractRule.MetadataType != nil {
					extractRuleMap["metadata_type"] = configInfo.ExtractRule.MetadataType
				}

				if configInfo.ExtractRule.PathRegex != nil {
					extractRuleMap["path_regex"] = configInfo.ExtractRule.PathRegex
				}

				if configInfo.ExtractRule.MetaTags != nil {
					metaTagsList := []interface{}{}
					for _, metaTags := range configInfo.ExtractRule.MetaTags {
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

				configInfoMap["extract_rule"] = []interface{}{extractRuleMap}
			}

			if configInfo.ExcludePaths != nil {
				excludePathsList := []interface{}{}
				for _, excludePaths := range configInfo.ExcludePaths {
					excludePathsMap := map[string]interface{}{}

					if excludePaths.Type != nil {
						excludePathsMap["type"] = excludePaths.Type
					}

					if excludePaths.Value != nil {
						excludePathsMap["value"] = excludePaths.Value
					}

					excludePathsList = append(excludePathsList, excludePathsMap)
				}

				configInfoMap["exclude_paths"] = excludePathsList
			}

			if configInfo.Output != nil {
				configInfoMap["output"] = configInfo.Output
			}

			if configInfo.UpdateTime != nil {
				configInfoMap["update_time"] = configInfo.UpdateTime
			}

			if configInfo.CreateTime != nil {
				configInfoMap["create_time"] = configInfo.CreateTime
			}

			if configInfo.UserDefineRule != nil {
				configInfoMap["user_define_rule"] = configInfo.UserDefineRule
			}

			ids = append(ids, *configInfo.ConfigId)
			tmpList = append(tmpList, configInfoMap)
		}

		_ = d.Set("configs", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
