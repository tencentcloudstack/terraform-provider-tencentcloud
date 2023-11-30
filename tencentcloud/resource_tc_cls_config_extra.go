package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsConfigExtra() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConfigExtraCreate,
		Read:   resourceTencentCloudClsConfigExtraRead,
		Delete: resourceTencentCloudClsConfigExtraDelete,
		Update: resourceTencentCloudClsConfigExtraUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection configuration name.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic ID (TopicId) of collection configuration.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type. Valid values: container_stdout; container_file; host_file.",
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Type of the log to be collected. Valid values: json_log: log in JSON format; delimiter_log: log in delimited format; minimalist_log: minimalist log; multiline_log: log in multi-line format; " +
					"fullregex_log: log in full regex format. Default value: minimalist_log.",
			},
			"config_flag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection configuration flag.",
			},
			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset Id.",
			},
			"logset_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset Name.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Topic Name.",
			},
			"host_file": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Node file config info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log file dir.",
						},
						"file_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log file name.",
						},
						"custom_labels": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Metadata info.",
						},
					},
				},
			},
			"container_file": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Container file path info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Namespace. There can be multiple namespaces, separated by separators, such as A, B.",
						},
						"container": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Container name.",
						},
						"log_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log Path.",
						},
						"file_pattern": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "log name.",
						},
						"include_labels": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Pod label info.",
						},
						"workload": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Workload info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "workload type.",
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "workload name.",
									},
									"container": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "container name.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "namespace.",
									},
								},
							},
						},
						"exclude_namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespaces to be excluded, separated by separators, such as A, B.",
						},
						"exclude_labels": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Pod label to be excluded.",
						},
					},
				},
			},
			"container_stdout": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Container stdout info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_containers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Is all containers.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespace. There can be multiple namespaces, separated by separators, such as A, B.",
						},
						"include_labels": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Pod label info.",
						},
						"workloads": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Workload info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "workload type.",
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "workload name.",
									},
									"container": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "container name.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "namespace.",
									},
								},
							},
						},
						"exclude_namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Namespaces to be excluded, separated by separators, such as A, B.",
						},
						"exclude_labels": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Pod label to be excluded.",
						},
					},
				},
			},
			"log_format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log format.",
			},
			"extract_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
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
							Computed:    true,
							Description: "Size of the data to be rewound in incremental collection mode. Default value: -1 (full collection).",
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
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Binding group id.",
			},
			"group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Binding group ids.",
			},
		},
	}
}

func resourceTencentCloudClsConfigExtraCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_config_extra.create")()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateConfigExtraRequest()
		response *cls.CreateConfigExtraResponse
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("config_flag"); ok {
		request.ConfigFlag = helper.String(v.(string))
	}
	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("host_file"); ok {
		hostFiles := make([]*cls.HostFileInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one host file.")
		}
		hostFile := cls.HostFileInfo{}
		dMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := dMap["log_path"]; ok {
			hostFile.LogPath = helper.String(v.(string))
		}
		if v, ok := dMap["file_pattern"]; ok {
			hostFile.FilePattern = helper.String(v.(string))
		}
		if v, ok := dMap["custom_labels"]; ok {
			customLabels := v.(*schema.Set).List()
			for _, customLabel := range customLabels {
				hostFile.CustomLabels = append(hostFile.CustomLabels, helper.String(customLabel.(string)))
			}
		}
		hostFiles = append(hostFiles, &hostFile)
		request.HostFile = hostFiles[0]
	}
	if v, ok := d.GetOk("container_file"); ok {
		containerFiles := make([]*cls.ContainerFileInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one container file.")
		}
		containerFile := cls.ContainerFileInfo{}
		dMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := dMap["namespace"]; ok {
			containerFile.Namespace = helper.String(v.(string))
		}
		if v, ok := dMap["container"]; ok {
			containerFile.Container = helper.String(v.(string))
		}
		if v, ok := dMap["log_path"]; ok {
			containerFile.LogPath = helper.String(v.(string))
		}
		if v, ok := dMap["file_pattern"]; ok {
			containerFile.FilePattern = helper.String(v.(string))
		}
		if v, ok := dMap["include_labels"]; ok {
			includeLabels := v.(*schema.Set).List()
			for _, includeLabel := range includeLabels {
				containerFile.IncludeLabels = append(containerFile.IncludeLabels, helper.String(includeLabel.(string)))
			}
		}
		if v, ok := dMap["workload"]; ok {
			workloads := make([]*cls.ContainerWorkLoadInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one workload.")
			}
			workload := cls.ContainerWorkLoadInfo{}
			dMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := dMap["kind"]; ok {
				workload.Kind = helper.String(v.(string))
			}
			if v, ok := dMap["name"]; ok {
				workload.Name = helper.String(v.(string))
			}
			if v, ok := dMap["container"]; ok {
				workload.Container = helper.String(v.(string))
			}
			if v, ok := dMap["namespace"]; ok {
				workload.Namespace = helper.String(v.(string))
			}
			workloads = append(workloads, &workload)
			containerFile.WorkLoad = workloads[0]
		}
		if v, ok := dMap["exclude_namespace"]; ok {
			containerFile.ExcludeNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["exclude_labels"]; ok {
			excludeLabels := v.(*schema.Set).List()
			for _, excludeLabel := range excludeLabels {
				containerFile.ExcludeLabels = append(containerFile.ExcludeLabels, helper.String(excludeLabel.(string)))
			}
		}
		containerFiles = append(containerFiles, &containerFile)
		request.ContainerFile = containerFiles[0]
	}

	if v, ok := d.GetOk("container_stdout"); ok {
		containerStdouts := make([]*cls.ContainerStdoutInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one container file.")
		}
		containerStdout := cls.ContainerStdoutInfo{}
		dMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := dMap["all_containers"]; ok {
			containerStdout.AllContainers = helper.Bool(v.(bool))
		}
		if v, ok := dMap["namespace"]; ok {
			containerStdout.Namespace = helper.String(v.(string))
		}
		if v, ok := dMap["include_labels"]; ok {
			includeLabels := v.(*schema.Set).List()
			for _, includeLabel := range includeLabels {
				containerStdout.IncludeLabels = append(containerStdout.IncludeLabels, helper.String(includeLabel.(string)))
			}
		}
		if v, ok := dMap["workloads"]; ok {
			workloads := make([]*cls.ContainerWorkLoadInfo, 0, 10)
			workload := cls.ContainerWorkLoadInfo{}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["kind"]; ok {
					workload.Kind = helper.String(v.(string))
				}
				if v, ok := dMap["name"]; ok {
					workload.Name = helper.String(v.(string))
				}
				if v, ok := dMap["container"]; ok {
					workload.Container = helper.String(v.(string))
				}
				if v, ok := dMap["namespace"]; ok {
					workload.Namespace = helper.String(v.(string))
				}
				workloads = append(workloads, &workload)
			}
			containerStdout.WorkLoads = workloads
		}
		if v, ok := dMap["exclude_namespace"]; ok {
			containerStdout.ExcludeNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["exclude_labels"]; ok {
			excludeLabels := v.(*schema.Set).List()
			for _, excludeLabel := range excludeLabels {
				containerStdout.ExcludeLabels = append(containerStdout.ExcludeLabels, helper.String(excludeLabel.(string)))
			}
		}
		containerStdouts = append(containerStdouts, &containerStdout)
		request.ContainerStdout = containerStdouts[0]
	}
	if v, ok := d.GetOk("log_format"); ok {
		request.LogFormat = helper.String(v.(string))
	}
	if v, ok := d.GetOk("extract_rule"); ok {
		extractRules := make([]*cls.ExtractRuleInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one extract rule.")
		}
		extractRule := cls.ExtractRuleInfo{}
		dMap := v.([]interface{})[0].(map[string]interface{})
		if v, ok := dMap["time_key"]; ok {
			extractRule.TimeKey = helper.String(v.(string))
		}
		if v, ok := dMap["time_format"]; ok {
			extractRule.TimeFormat = helper.String(v.(string))
		}
		if v, ok := dMap["delimiter"]; ok {
			extractRule.Delimiter = helper.String(v.(string))
		}
		if v, ok := dMap["log_regex"]; ok {
			extractRule.LogRegex = helper.String(v.(string))
		}
		if v, ok := dMap["begin_regex"]; ok {
			extractRule.BeginRegex = helper.String(v.(string))
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
		if v, ok := dMap["un_match_log_key"]; ok {
			extractRule.UnMatchLogKey = helper.String(v.(string))
		}
		if v, ok := dMap["backtracking"]; ok {
			extractRule.Backtracking = helper.IntInt64(v.(int))
		}
		extractRules = append(extractRules, &extractRule)
		request.ExtractRule = extractRules[0]
	}
	if v, ok := d.GetOk("exclude_paths"); ok {
		excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			excludePath := cls.ExcludePathInfo{}
			if v, ok := dMap["type"]; ok {
				excludePath.Type = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				excludePath.Value = helper.String(v.(string))
			}
			excludePaths = append(excludePaths, &excludePath)
		}
		request.ExcludePaths = excludePaths
	}
	if v, ok := d.GetOk("user_define_rule"); ok {
		request.UserDefineRule = helper.String(v.(string))
	}
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("group_ids"); ok {
		groupIds := v.(*schema.Set).List()
		for _, groupId := range groupIds {
			request.GroupIds = append(request.GroupIds, helper.String(groupId.(string)))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateConfigExtra(request)
		if e != nil {
			return retryError(e)
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

	id := *response.Response.ConfigExtraId
	d.SetId(id)
	return resourceTencentCloudClsConfigExtraRead(d, meta)
}

func resourceTencentCloudClsConfigExtraRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_config_extra.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	configExtraId := d.Id()

	configExtra, err := service.DescribeClsConfigExtraById(ctx, configExtraId)
	if err != nil {
		return err
	}

	if configExtra == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsConfigExtra` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configExtra.Name != nil {
		_ = d.Set("name", configExtra.Name)
	}

	if configExtra.TopicId != nil {
		_ = d.Set("topic_id", configExtra.TopicId)
	}

	if configExtra.Type != nil {
		_ = d.Set("type", configExtra.Type)
	}

	if configExtra.LogType != nil {
		_ = d.Set("log_type", configExtra.LogType)
	}

	if configExtra.ConfigFlag != nil {
		_ = d.Set("config_flag", configExtra.ConfigFlag)
	}

	if configExtra.LogsetId != nil {
		_ = d.Set("logset_id", configExtra.LogsetId)
	}

	if configExtra.LogsetName != nil {
		_ = d.Set("logset_name", configExtra.LogsetName)
	}

	if configExtra.TopicName != nil {
		_ = d.Set("topic_name", configExtra.TopicName)
	}

	if configExtra.HostFile != nil {
		hostFileMap := map[string]interface{}{}

		if configExtra.HostFile.LogPath != nil {
			hostFileMap["log_path"] = configExtra.HostFile.LogPath
		}

		if configExtra.HostFile.FilePattern != nil {
			hostFileMap["file_pattern"] = configExtra.HostFile.FilePattern
		}

		if len(configExtra.HostFile.CustomLabels) > 0 {
			hostFileMap["custom_labels"] = configExtra.HostFile.CustomLabels
		}

		_ = d.Set("host_file", []interface{}{hostFileMap})
	}

	if configExtra.ContainerFile != nil {
		containerFileMap := map[string]interface{}{}

		if configExtra.ContainerFile.Namespace != nil {
			containerFileMap["namespace"] = configExtra.ContainerFile.Namespace
		}

		if configExtra.ContainerFile.Container != nil {
			containerFileMap["container"] = configExtra.ContainerFile.Container
		}

		if configExtra.ContainerFile.LogPath != nil {
			containerFileMap["log_path"] = configExtra.ContainerFile.LogPath
		}

		if configExtra.ContainerFile.FilePattern != nil {
			containerFileMap["file_pattern"] = configExtra.ContainerFile.FilePattern
		}

		if len(configExtra.ContainerFile.IncludeLabels) > 0 {
			containerFileMap["include_labels"] = configExtra.ContainerFile.IncludeLabels
		}

		if configExtra.ContainerFile.WorkLoad != nil {
			workLoadMap := map[string]interface{}{}

			if configExtra.ContainerFile.WorkLoad.Kind != nil {
				workLoadMap["kind"] = configExtra.ContainerFile.WorkLoad.Kind
			}

			if configExtra.ContainerFile.WorkLoad.Name != nil {
				workLoadMap["name"] = configExtra.ContainerFile.WorkLoad.Name
			}

			if configExtra.ContainerFile.WorkLoad.Container != nil {
				workLoadMap["container"] = configExtra.ContainerFile.WorkLoad.Container
			}

			if configExtra.ContainerFile.WorkLoad.Namespace != nil {
				workLoadMap["namespace"] = configExtra.ContainerFile.WorkLoad.Namespace
			}

			containerFileMap["workload"] = []interface{}{workLoadMap}
		}

		if configExtra.ContainerFile.ExcludeNamespace != nil {
			containerFileMap["exclude_namespace"] = configExtra.ContainerFile.ExcludeNamespace
		}

		if len(configExtra.ContainerFile.ExcludeLabels) > 0 {
			containerFileMap["exclude_labels"] = configExtra.ContainerFile.ExcludeLabels
		}

		_ = d.Set("container_file", []interface{}{containerFileMap})
	}

	if configExtra.ContainerStdout != nil {
		containerStdoutMap := map[string]interface{}{}

		if configExtra.ContainerStdout.AllContainers != nil {
			containerStdoutMap["all_containers"] = configExtra.ContainerStdout.AllContainers
		}

		if configExtra.ContainerStdout.Namespace != nil {
			containerStdoutMap["namespace"] = configExtra.ContainerStdout.Namespace
		}

		if configExtra.ContainerStdout.IncludeLabels != nil {
			containerStdoutMap["include_labels"] = configExtra.ContainerStdout.IncludeLabels
		}

		if configExtra.ContainerStdout.WorkLoads != nil {
			workLoadsList := []interface{}{}
			for _, workLoads := range configExtra.ContainerStdout.WorkLoads {
				workLoadsMap := map[string]interface{}{}

				if workLoads.Kind != nil {
					workLoadsMap["kind"] = workLoads.Kind
				}

				if workLoads.Name != nil {
					workLoadsMap["name"] = workLoads.Name
				}

				if workLoads.Container != nil {
					workLoadsMap["container"] = workLoads.Container
				}

				if workLoads.Namespace != nil {
					workLoadsMap["namespace"] = workLoads.Namespace
				}

				workLoadsList = append(workLoadsList, workLoadsMap)
			}

			containerStdoutMap["workloads"] = workLoadsList
		}

		if configExtra.ContainerStdout.ExcludeNamespace != nil {
			containerStdoutMap["exclude_namespace"] = configExtra.ContainerStdout.ExcludeNamespace
		}

		if configExtra.ContainerStdout.ExcludeLabels != nil {
			containerStdoutMap["exclude_labels"] = configExtra.ContainerStdout.ExcludeLabels
		}

		_ = d.Set("container_stdout", []interface{}{containerStdoutMap})
	}

	if configExtra.LogFormat != nil {
		_ = d.Set("log_format", configExtra.LogFormat)
	}

	if configExtra.ExtractRule != nil {
		extractRuleMap := map[string]interface{}{}

		if configExtra.ExtractRule.TimeKey != nil {
			extractRuleMap["time_key"] = configExtra.ExtractRule.TimeKey
		}

		if configExtra.ExtractRule.TimeFormat != nil {
			extractRuleMap["time_format"] = configExtra.ExtractRule.TimeFormat
		}

		if configExtra.ExtractRule.Delimiter != nil {
			extractRuleMap["delimiter"] = configExtra.ExtractRule.Delimiter
		}

		if configExtra.ExtractRule.LogRegex != nil {
			extractRuleMap["log_regex"] = configExtra.ExtractRule.LogRegex
		}

		if configExtra.ExtractRule.BeginRegex != nil {
			extractRuleMap["begin_regex"] = configExtra.ExtractRule.BeginRegex
		}

		if len(configExtra.ExtractRule.Keys) > 0 {
			extractRuleMap["keys"] = configExtra.ExtractRule.Keys
		}

		if configExtra.ExtractRule.FilterKeyRegex != nil {
			filterKeyRegexList := []interface{}{}
			for _, filterKeyRegex := range configExtra.ExtractRule.FilterKeyRegex {
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

		if configExtra.ExtractRule.UnMatchUpLoadSwitch != nil {
			extractRuleMap["un_match_up_load_switch"] = configExtra.ExtractRule.UnMatchUpLoadSwitch
		}

		if configExtra.ExtractRule.UnMatchLogKey != nil {
			extractRuleMap["un_match_log_key"] = configExtra.ExtractRule.UnMatchLogKey
		}

		if configExtra.ExtractRule.Backtracking != nil {
			extractRuleMap["backtracking"] = configExtra.ExtractRule.Backtracking
		}

		_ = d.Set("extract_rule", []interface{}{extractRuleMap})
	}

	if configExtra.ExcludePaths != nil {
		excludePathsList := []interface{}{}
		for _, excludePath := range configExtra.ExcludePaths {
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

	if configExtra.UserDefineRule != nil {
		_ = d.Set("user_define_rule", configExtra.UserDefineRule)
	}

	if configExtra.GroupId != nil {
		_ = d.Set("group_id", configExtra.GroupId)
	}

	return nil
}

func resourceTencentCloudClsConfigExtraUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_config_extra.update")()
	logId := getLogId(contextNil)
	request := cls.NewModifyConfigExtraRequest()

	request.ConfigExtraId = helper.String(d.Id())

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}
	if d.HasChange("topic_id") {
		if v, ok := d.GetOk("topic_id"); ok {
			request.TopicId = helper.String(v.(string))
		}
	}
	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}
	if d.HasChange("config_flag") {
		if v, ok := d.GetOk("config_flag"); ok {
			request.ConfigFlag = helper.String(v.(string))
		}
	}
	if d.HasChange("logset_id") {
		if v, ok := d.GetOk("logset_id"); ok {
			request.LogsetId = helper.String(v.(string))
		}
	}
	if d.HasChange("logset_name") {
		if v, ok := d.GetOk("logset_name"); ok {
			request.LogsetName = helper.String(v.(string))
		}
	}
	if d.HasChange("topic_name") {
		if v, ok := d.GetOk("topic_name"); ok {
			request.TopicName = helper.String(v.(string))
		}
	}
	if d.HasChange("host_file") {
		if v, ok := d.GetOk("host_file"); ok {
			hostFiles := make([]*cls.HostFileInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one host file.")
			}
			hostFile := cls.HostFileInfo{}
			dMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := dMap["log_path"]; ok {
				hostFile.LogPath = helper.String(v.(string))
			}
			if v, ok := dMap["file_pattern"]; ok {
				hostFile.FilePattern = helper.String(v.(string))
			}
			if v, ok := dMap["custom_labels"]; ok {
				customLabels := v.(*schema.Set).List()
				for _, customLabel := range customLabels {
					hostFile.CustomLabels = append(hostFile.CustomLabels, helper.String(customLabel.(string)))
				}
			}
			hostFiles = append(hostFiles, &hostFile)
			request.HostFile = hostFiles[0]
		}
	}
	if d.HasChange("container_file") {
		if v, ok := d.GetOk("container_file"); ok {
			containerFiles := make([]*cls.ContainerFileInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one container file.")
			}
			containerFile := cls.ContainerFileInfo{}
			dMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := dMap["namespace"]; ok {
				containerFile.Namespace = helper.String(v.(string))
			}
			if v, ok := dMap["container"]; ok {
				containerFile.Container = helper.String(v.(string))
			}
			if v, ok := dMap["log_path"]; ok {
				containerFile.LogPath = helper.String(v.(string))
			}
			if v, ok := dMap["file_pattern"]; ok {
				containerFile.FilePattern = helper.String(v.(string))
			}
			if v, ok := dMap["include_labels"]; ok {
				includeLabels := v.(*schema.Set).List()
				for _, includeLabel := range includeLabels {
					containerFile.IncludeLabels = append(containerFile.IncludeLabels, helper.String(includeLabel.(string)))
				}
			}
			if v, ok := dMap["workload"]; ok {
				workloads := make([]*cls.ContainerWorkLoadInfo, 0, 10)
				if len(v.([]interface{})) != 1 {
					return fmt.Errorf("need only one workload.")
				}
				workload := cls.ContainerWorkLoadInfo{}
				dMap := v.([]interface{})[0].(map[string]interface{})
				if v, ok := dMap["kind"]; ok {
					workload.Kind = helper.String(v.(string))
				}
				if v, ok := dMap["name"]; ok {
					workload.Name = helper.String(v.(string))
				}
				if v, ok := dMap["container"]; ok {
					workload.Container = helper.String(v.(string))
				}
				if v, ok := dMap["namespace"]; ok {
					workload.Namespace = helper.String(v.(string))
				}
				workloads = append(workloads, &workload)
				containerFile.WorkLoad = workloads[0]
			}
			if v, ok := dMap["exclude_namespace"]; ok {
				containerFile.ExcludeNamespace = helper.String(v.(string))
			}
			if v, ok := dMap["exclude_labels"]; ok {
				excludeLabels := v.(*schema.Set).List()
				for _, excludeLabel := range excludeLabels {
					containerFile.ExcludeLabels = append(containerFile.ExcludeLabels, helper.String(excludeLabel.(string)))
				}
			}
			containerFiles = append(containerFiles, &containerFile)
			request.ContainerFile = containerFiles[0]
		}
	}
	if d.HasChange("container_stdout") {
		if v, ok := d.GetOk("container_stdout"); ok {
			containerStdouts := make([]*cls.ContainerStdoutInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one container file.")
			}
			containerStdout := cls.ContainerStdoutInfo{}
			dMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := dMap["all_containers"]; ok {
				containerStdout.AllContainers = helper.Bool(v.(bool))
			}
			if v, ok := dMap["namespace"]; ok {
				containerStdout.Namespace = helper.String(v.(string))
			}
			if v, ok := dMap["include_labels"]; ok {
				includeLabels := v.(*schema.Set).List()
				for _, includeLabel := range includeLabels {
					containerStdout.IncludeLabels = append(containerStdout.IncludeLabels, helper.String(includeLabel.(string)))
				}
			}
			if v, ok := dMap["workloads"]; ok {
				workloads := make([]*cls.ContainerWorkLoadInfo, 0, 10)
				workload := cls.ContainerWorkLoadInfo{}
				for _, item := range v.([]interface{}) {
					dMap := item.(map[string]interface{})
					if v, ok := dMap["kind"]; ok {
						workload.Kind = helper.String(v.(string))
					}
					if v, ok := dMap["name"]; ok {
						workload.Name = helper.String(v.(string))
					}
					if v, ok := dMap["container"]; ok {
						workload.Container = helper.String(v.(string))
					}
					if v, ok := dMap["namespace"]; ok {
						workload.Namespace = helper.String(v.(string))
					}
					workloads = append(workloads, &workload)
				}
				containerStdout.WorkLoads = workloads
			}
			if v, ok := dMap["exclude_namespace"]; ok {
				containerStdout.ExcludeNamespace = helper.String(v.(string))
			}
			if v, ok := dMap["exclude_labels"]; ok {
				excludeLabels := v.(*schema.Set).List()
				for _, excludeLabel := range excludeLabels {
					containerStdout.ExcludeLabels = append(containerStdout.ExcludeLabels, helper.String(excludeLabel.(string)))
				}
			}
			containerStdouts = append(containerStdouts, &containerStdout)
			request.ContainerStdout = containerStdouts[0]
		}
	}
	if d.HasChange("log_format") {
		if v, ok := d.GetOk("log_format"); ok {
			request.LogFormat = helper.String(v.(string))
		}
	}

	if d.HasChange("log_type") || d.HasChange("extract_rule") {
		if v, ok := d.GetOk("log_type"); ok {
			request.LogType = helper.String(v.(string))
		}
		if v, ok := d.GetOk("extract_rule"); ok {
			extractRules := make([]*cls.ExtractRuleInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one extract rule.")
			}
			extractRule := cls.ExtractRuleInfo{}
			dMap := v.([]interface{})[0].(map[string]interface{})
			if v, ok := dMap["time_key"]; ok {
				extractRule.TimeKey = helper.String(v.(string))
			}
			if v, ok := dMap["time_format"]; ok {
				extractRule.TimeFormat = helper.String(v.(string))
			}
			if v, ok := dMap["delimiter"]; ok {
				extractRule.Delimiter = helper.String(v.(string))
			}
			if v, ok := dMap["log_regex"]; ok {
				extractRule.LogRegex = helper.String(v.(string))
			}
			if v, ok := dMap["begin_regex"]; ok {
				extractRule.BeginRegex = helper.String(v.(string))
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
			if v, ok := dMap["un_match_log_key"]; ok {
				extractRule.UnMatchLogKey = helper.String(v.(string))
			}
			if v, ok := dMap["backtracking"]; ok {
				extractRule.Backtracking = helper.IntInt64(v.(int))
			}
			extractRules = append(extractRules, &extractRule)
			request.ExtractRule = extractRules[0]
		}
	}
	if d.HasChange("exclude_paths") {
		if v, ok := d.GetOk("exclude_paths"); ok {
			excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				excludePath := cls.ExcludePathInfo{}
				if v, ok := dMap["type"]; ok {
					excludePath.Type = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					excludePath.Value = helper.String(v.(string))
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
	if d.HasChange("group_id") {
		if v, ok := d.GetOk("group_id"); ok {
			request.GroupId = helper.String(v.(string))
		}
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyConfigExtra(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsConfigExtraRead(d, meta)
}

func resourceTencentCloudClsConfigExtraDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_config.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsConfigExtra(ctx, id); err != nil {
		return err
	}

	return nil
}
