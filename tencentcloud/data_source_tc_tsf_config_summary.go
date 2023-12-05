package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfConfigSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfConfigSummaryRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Application ID. If not passed, the query will be for all.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query keyword, fuzzy query: application name, configuration item name. If not passed, the query will be for all.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Order term. support Sort by time: creation_time; or Sort by name: config_name.",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Pass 0 for ascending order and 1 for descending order.",
			},

			"config_tag_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "config tag list.",
			},

			"disable_program_auth_check": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to disable dataset authentication.",
			},

			"config_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Config Id List.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "config Page Item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "config list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item ID.Note: This field may return null, indicating that no valid value was found.",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration name.Note: This field may return null, indicating that no valid value was found.",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration version. Note: This field may return null, indicating that no valid value was found.",
									},
									"config_version_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration version description.Note: This field may return null, indicating that no valid value was found.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration value.Note: This field may return null, indicating that no valid value was found.",
									},
									"config_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Config type. Note: This field may return null, indicating that no valid value was found.",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time.Note: This field may return null, indicating that no valid value was found.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID.Note: This field may return null, indicating that no valid value was found.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application Name. Note: This field may return null, indicating that no valid value was found.",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deletion flag, true: deletable; false: not deletable.Note: This field may return null, indicating that no valid value was found.",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last update time.Note: This field may return null, indicating that no valid value was found.",
									},
									"config_version_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Configure version count.Note: This field may return null, indicating that no valid value was found.",
									},
								},
							},
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

func dataSourceTencentCloudTsfConfigSummaryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_config_summary.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("config_tag_list"); ok {
		configTagListSet := v.(*schema.Set).List()
		paramMap["ConfigTagList"] = helper.InterfacesStringsPoint(configTagListSet)
	}

	if v, _ := d.GetOk("disable_program_auth_check"); v != nil {
		paramMap["DisableProgramAuthCheck"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("config_id_list"); ok {
		configIdListSet := v.(*schema.Set).List()
		paramMap["ConfigIdList"] = helper.InterfacesStringsPoint(configIdListSet)
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var config *tsf.TsfPageConfig
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfConfigSummaryByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		config = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(config.Content))
	tsfPageConfigMap := map[string]interface{}{}
	if config != nil {
		if config.TotalCount != nil {
			tsfPageConfigMap["total_count"] = config.TotalCount
		}

		if config.Content != nil {
			contentList := []interface{}{}
			for _, content := range config.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.ConfigVersion != nil {
					contentMap["config_version"] = content.ConfigVersion
				}

				if content.ConfigVersionDesc != nil {
					contentMap["config_version_desc"] = content.ConfigVersionDesc
				}

				if content.ConfigValue != nil {
					contentMap["config_value"] = content.ConfigValue
				}

				if content.ConfigType != nil {
					contentMap["config_type"] = content.ConfigType
				}

				if content.CreationTime != nil {
					contentMap["creation_time"] = content.CreationTime
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.DeleteFlag != nil {
					contentMap["delete_flag"] = content.DeleteFlag
				}

				if content.LastUpdateTime != nil {
					contentMap["last_update_time"] = content.LastUpdateTime
				}

				if content.ConfigVersionCount != nil {
					contentMap["config_version_count"] = content.ConfigVersionCount
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ConfigId)
			}

			tsfPageConfigMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageConfigMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageConfigMap); e != nil {
			return e
		}
	}
	return nil
}
