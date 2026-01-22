package dlc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcDataEngine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcDataEngineCreate,
		Read:   resourceTencentCloudDlcDataEngineRead,
		Update: resourceTencentCloudDlcDataEngineUpdate,
		Delete: resourceTencentCloudDlcDataEngineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"engine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The engine type. Valid values: `spark` and `presto`.",
			},

			"data_engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of the virtual cluster.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The cluster type. Valid values: `spark_private`, `presto_private`, `presto_cu`, and `spark_cu`.",
			},

			"mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The billing mode. Valid values: `0` (shared engine), `1` (pay-as-you-go), and `2` (monthly subscription).",
			},

			"auto_resume": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically start the clusters.",
			},

			"size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Cluster size. Required when updating.",
			},

			"min_clusters": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of clusters.",
			},

			"max_clusters": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of clusters.",
			},

			"default_data_engine": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is the default virtual cluster.",
			},

			"cidr_block": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The VPC CIDR block.",
			},

			"message": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The description.",
			},

			"pay_mode": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The pay mode. Valid value: `0` (postpaid, default) and `1` (prepaid) (currently not available).",
			},

			"time_span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The usage duration of the resource. Postpaid: Fill in 3,600 as a fixed figure; prepaid: fill in a figure equal to or bigger than 1 which means purchasing resources for one month. The maximum figure is not bigger than 120. The default value is 1.",
			},

			"time_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The unit of the resource period. Valid values: `s` (default) for the postpaid mode and `m` for the prepaid mode.",
			},

			"auto_renew": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The auto-renewal status of the resource. For the postpaid mode, no renewal is required, and the value is fixed to `0`. For the prepaid mode, valid values are `0` (manual), `1` (auto), and `2` (no renewal). If this parameter is set to `0` for a key account in the prepaid mode, auto-renewal applies. It defaults to `0`.",
			},

			"auto_suspend": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically suspend clusters. Valid values: `false` (default, no) and `true` (yes).",
			},

			"crontab_resume_suspend": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable scheduled start and suspension of clusters. Valid values: `0` (disable) and `1` (enable). Note: This policy and the auto-suspension policy are mutually exclusive.",
			},

			"crontab_resume_suspend_strategy": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The complex policy for scheduled start and suspension, including the start/suspension time and suspension policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resume_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scheduled starting time, such as 8: 00 a.m. on Monday and Wednesday.",
						},
						"suspend_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scheduled suspension time, such as 8: 00 p.m. on Monday and Wednesday.",
						},
						"suspend_strategy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The suspension setting. Valid values: `0` (suspension after task end, default) and `1` (force suspension).",
						},
					},
				},
			},

			"engine_exec_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The type of tasks to be executed by the engine, which defaults to SQL. Valid values: `SQL` and `BATCH`.",
			},

			"max_concurrency": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The max task concurrency of a cluster, which defaults to 5.",
			},

			"tolerable_queue_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The task queue time limit, which defaults to 0. When the actual queue time exceeds the value set here, scale-out may be triggered. Setting this parameter to 0 represents that scale-out may be triggered immediately after a task queues up.",
			},

			"auto_suspend_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The cluster auto-suspension time, which defaults to 10 min.",
			},

			"resource_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The resource type. Valid values: `Standard_CU` (standard) and `Memory_CU` (memory).",
			},

			"data_engine_config_pairs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The advanced configurations of clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration items.",
						},
						"config_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration value.",
						},
					},
				},
			},

			"image_version_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The version name of cluster image, such as SuperSQL-P 1.1 and SuperSQL-S 3.2. If no value is passed in, a cluster is created using the latest image version.",
			},

			"main_cluster_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The primary cluster, which is specified when a failover cluster is created.",
			},

			"elastic_switch": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the scaling feature for a monthly subscribed Spark job cluster.",
			},

			"elastic_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The upper limit (in CUs) for scaling of the monthly subscribed Spark job cluster.",
			},

			"session_resource_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The session resource configuration template for a Spark job cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"driver_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The driver size. Valid values for the standard resource type: `small`, `medium`, `large`, and `xlarge`. Valid values for the memory resource type: `m.small`, `m.medium`, `m.large`, and `m.xlarge`.",
						},
						"executor_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The executor size. Valid values for the standard resource type: `small`, `medium`, `large`, and `xlarge`. Valid values for the memory resource type: `m.small`, `m.medium`, `m.large`, and `m.xlarge`.",
						},
						"executor_nums": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The executor count. The minimum value is 1 and the maximum value is less than the cluster specification.",
						},
						"executor_max_numbers": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum executor count (in dynamic mode). The minimum value is 1 and the maximum value is less than the cluster specification. If you set `ExecutorMaxNumbers` to a value smaller than that of `ExecutorNums`, the value of `ExecutorMaxNumbers` is automatically changed to that of `ExecutorNums`.",
						},
						"running_time_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The running time parameters of the session resource configuration template for a Spark job cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration items.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration value.",
									},
								},
							},
						},
					},
				},
			},

			"auto_authorization": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Automatic authorization.",
			},

			"engine_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Engine network ID.",
			},

			"engine_generation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Generation of the engine. SuperSQL means the supersql engine while Native means the standard engine. It is SuperSQL by default.",
			},

			// computed
			"data_engine_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data engine ID.",
			},
		},
	}
}

func resourceTencentCloudDlcDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = dlc.NewCreateDataEngineRequest()
		dataEngineId   string
		dataEngineName string
	)

	if v, ok := d.GetOk("engine_type"); ok {
		request.EngineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_engine_name"); ok {
		request.DataEngineName = helper.String(v.(string))
		dataEngineName = v.(string)
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("mode"); ok {
		request.Mode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_resume"); ok {
		request.AutoResume = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("size"); ok {
		request.Size = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("min_clusters"); ok {
		request.MinClusters = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_clusters"); ok {
		request.MaxClusters = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("default_data_engine"); ok {
		request.DefaultDataEngine = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cidr_block"); ok {
		request.CidrBlock = helper.String(v.(string))
	}

	if v, ok := d.GetOk("message"); ok {
		request.Message = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("pay_mode"); ok {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request.AutoRenew = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_suspend"); ok {
		request.AutoSuspend = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("crontab_resume_suspend"); ok {
		request.CrontabResumeSuspend = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "crontab_resume_suspend_strategy"); ok {
		crontabResumeSuspendStrategy := dlc.CrontabResumeSuspendStrategy{}
		if v, ok := dMap["resume_time"]; ok {
			crontabResumeSuspendStrategy.ResumeTime = helper.String(v.(string))
		}

		if v, ok := dMap["suspend_time"]; ok {
			crontabResumeSuspendStrategy.SuspendTime = helper.String(v.(string))
		}

		if v, ok := dMap["suspend_strategy"]; ok {
			crontabResumeSuspendStrategy.SuspendStrategy = helper.IntInt64(v.(int))
		}

		request.CrontabResumeSuspendStrategy = &crontabResumeSuspendStrategy
	}

	if v, ok := d.GetOk("engine_exec_type"); ok {
		request.EngineExecType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("tolerable_queue_time"); ok {
		request.TolerableQueueTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_suspend_time"); ok {
		request.AutoSuspendTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_engine_config_pairs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataEngineConfigPair := dlc.DataEngineConfigPair{}
			if v, ok := dMap["config_item"]; ok {
				dataEngineConfigPair.ConfigItem = helper.String(v.(string))
			}

			if v, ok := dMap["config_value"]; ok {
				dataEngineConfigPair.ConfigItem = helper.String(v.(string))
			}

			request.DataEngineConfigPairs = append(request.DataEngineConfigPairs, &dataEngineConfigPair)
		}
	}

	if v, ok := d.GetOk("image_version_name"); ok {
		request.ImageVersionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("main_cluster_name"); ok {
		request.MainClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("elastic_switch"); ok {
		request.ElasticSwitch = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("elastic_limit"); ok {
		request.ElasticLimit = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "session_resource_template"); ok {
		sessionResourceTemplate := dlc.SessionResourceTemplate{}
		if v, ok := dMap["driver_size"]; ok {
			sessionResourceTemplate.DriverSize = helper.String(v.(string))
		}

		if v, ok := dMap["executor_size"]; ok {
			sessionResourceTemplate.ExecutorSize = helper.String(v.(string))
		}

		if v, ok := dMap["executor_nums"]; ok {
			sessionResourceTemplate.ExecutorNums = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["executor_max_numbers"]; ok {
			sessionResourceTemplate.ExecutorMaxNumbers = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["running_time_parameters"]; ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				dataEngineConfigPair := dlc.DataEngineConfigPair{}
				if v, ok := dMap["config_item"]; ok {
					dataEngineConfigPair.ConfigItem = helper.String(v.(string))
				}

				if v, ok := dMap["config_value"]; ok {
					dataEngineConfigPair.ConfigItem = helper.String(v.(string))
				}

				sessionResourceTemplate.RunningTimeParameters = append(sessionResourceTemplate.RunningTimeParameters, &dataEngineConfigPair)
			}
		}

		request.SessionResourceTemplate = &sessionResourceTemplate
	}

	if v, ok := d.GetOkExists("auto_authorization"); ok {
		request.AutoAuthorization = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("engine_network_id"); ok {
		request.EngineNetworkId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_generation"); ok {
		request.EngineGeneration = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc data engine failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc data engine failed, reason:%+v", logId, err)
		return err
	}

	// get dataEngineId
	describeRequest := dlc.NewDescribeDataEngineRequest()
	describeRequest.DataEngineName = helper.String(dataEngineName)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDataEngine(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DataEngine == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] resopse is null, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}

		if result.Response.DataEngine.DataEngineId == nil {
			return resource.NonRetryableError(fmt.Errorf("DataEngineId is nil."))
		}

		dataEngineId = *result.Response.DataEngine.DataEngineId
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc dataEngine failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineName + tccommon.FILED_SP + dataEngineId)

	// wait
	err = resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDataEngine(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DataEngine == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] resopse is null, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}

		if *result.Response.DataEngine.State != int64(2) && *result.Response.DataEngine.State != int64(1) {
			e = fmt.Errorf("[DEBUG]%s api[%s] status [%v] not ready , request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), *result.Response.DataEngine.State, describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc dataEngine failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDlcDataEngineRead(d, meta)
}

func resourceTencentCloudDlcDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	dataEngineName := idSplit[0]
	dataEngine, err := service.DescribeDlcDataEngineByName(ctx, dataEngineName)
	if err != nil {
		return err
	}

	if dataEngine == nil {
		log.Printf("[WARN]%s resource `DlcDataEngine` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if dataEngine.EngineType != nil {
		_ = d.Set("engine_type", dataEngine.EngineType)
	}

	if dataEngine.DataEngineName != nil {
		_ = d.Set("data_engine_name", dataEngine.DataEngineName)
	}

	if dataEngine.ClusterType != nil {
		_ = d.Set("cluster_type", dataEngine.ClusterType)
	}

	if dataEngine.Mode != nil {
		_ = d.Set("mode", dataEngine.Mode)
	}

	if dataEngine.AutoResume != nil {
		_ = d.Set("auto_resume", dataEngine.AutoResume)
	}

	if dataEngine.Size != nil {
		_ = d.Set("size", dataEngine.Size)
	}

	if dataEngine.MinClusters != nil {
		_ = d.Set("min_clusters", dataEngine.MinClusters)
	}

	if dataEngine.MaxClusters != nil {
		_ = d.Set("max_clusters", dataEngine.MaxClusters)
	}

	if dataEngine.DefaultDataEngine != nil {
		_ = d.Set("default_data_engine", dataEngine.DefaultDataEngine)
	}

	if dataEngine.CidrBlock != nil {
		_ = d.Set("cidr_block", dataEngine.CidrBlock)
	}

	if dataEngine.Message != nil {
		_ = d.Set("message", dataEngine.Message)
	}

	if dataEngine.Mode != nil {
		if *dataEngine.Mode == 1 {
			_ = d.Set("pay_mode", 0)
		} else if *dataEngine.Mode == 2 {
			_ = d.Set("pay_mode", 1)
		}
	}

	if dataEngine.RenewFlag != nil {
		_ = d.Set("auto_renew", dataEngine.RenewFlag)
	}

	if dataEngine.AutoSuspend != nil {
		_ = d.Set("auto_suspend", dataEngine.AutoSuspend)
	}

	if dataEngine.CrontabResumeSuspend != nil {
		_ = d.Set("crontab_resume_suspend", dataEngine.CrontabResumeSuspend)
	}

	if dataEngine.CrontabResumeSuspendStrategy != nil {
		crontabResumeSuspendStrategyMap := map[string]interface{}{}
		if dataEngine.CrontabResumeSuspendStrategy.ResumeTime != nil {
			crontabResumeSuspendStrategyMap["resume_time"] = dataEngine.CrontabResumeSuspendStrategy.ResumeTime
		}

		if dataEngine.CrontabResumeSuspendStrategy.SuspendTime != nil {
			crontabResumeSuspendStrategyMap["suspend_time"] = dataEngine.CrontabResumeSuspendStrategy.SuspendTime
		}

		if dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy != nil {
			crontabResumeSuspendStrategyMap["suspend_strategy"] = dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy
		}

		_ = d.Set("crontab_resume_suspend_strategy", []interface{}{crontabResumeSuspendStrategyMap})
	}

	if dataEngine.EngineExecType != nil {
		_ = d.Set("engine_exec_type", dataEngine.EngineExecType)
	}

	if dataEngine.MaxConcurrency != nil {
		_ = d.Set("max_concurrency", dataEngine.MaxConcurrency)
	}

	if dataEngine.TolerableQueueTime != nil {
		_ = d.Set("tolerable_queue_time", dataEngine.TolerableQueueTime)
	}

	if dataEngine.AutoSuspendTime != nil {
		_ = d.Set("auto_suspend_time", dataEngine.AutoSuspendTime)
	}

	if dataEngine.ResourceType != nil {
		_ = d.Set("resource_type", dataEngine.ResourceType)
	}

	if dataEngine.ImageVersionName != nil {
		_ = d.Set("image_version_name", dataEngine.ImageVersionName)
	}

	if dataEngine.ElasticSwitch != nil {
		_ = d.Set("elastic_switch", dataEngine.ElasticSwitch)
	}

	if dataEngine.ElasticLimit != nil {
		_ = d.Set("elastic_limit", dataEngine.ElasticLimit)
	}

	if dataEngine.SessionResourceTemplate != nil {
		sessionResourceTemplateMap := map[string]interface{}{}
		if dataEngine.SessionResourceTemplate.DriverSize != nil {
			sessionResourceTemplateMap["driver_size"] = dataEngine.SessionResourceTemplate.DriverSize
		}

		if dataEngine.SessionResourceTemplate.ExecutorSize != nil {
			sessionResourceTemplateMap["executor_size"] = dataEngine.SessionResourceTemplate.ExecutorSize
		}

		if dataEngine.SessionResourceTemplate.ExecutorNums != nil {
			sessionResourceTemplateMap["executor_nums"] = dataEngine.SessionResourceTemplate.ExecutorNums
		}

		if dataEngine.SessionResourceTemplate.ExecutorMaxNumbers != nil {
			sessionResourceTemplateMap["executor_max_numbers"] = dataEngine.SessionResourceTemplate.ExecutorMaxNumbers
		}

		if dataEngine.SessionResourceTemplate.RunningTimeParameters != nil {
			temList := make([]map[string]interface{}, 0, len(dataEngine.SessionResourceTemplate.RunningTimeParameters))
			for _, item := range dataEngine.SessionResourceTemplate.RunningTimeParameters {
				dMap := make(map[string]interface{})
				if item.ConfigItem != nil {
					dMap["config_item"] = *item.ConfigItem
				}

				if item.ConfigValue != nil {
					dMap["config_value"] = *item.ConfigValue
				}

				temList = append(temList, dMap)
			}

			sessionResourceTemplateMap["running_time_parameters"] = temList
		}

		_ = d.Set("session_resource_template", []interface{}{sessionResourceTemplateMap})
	}

	if dataEngine.AutoAuthorization != nil {
		_ = d.Set("auto_authorization", dataEngine.AutoAuthorization)
	}

	if dataEngine.EngineNetworkId != nil {
		_ = d.Set("engine_network_id", dataEngine.EngineNetworkId)
	}

	if dataEngine.EngineGeneration != nil {
		_ = d.Set("engine_generation", dataEngine.EngineGeneration)
	}

	if dataEngine.DataEngineId != nil {
		_ = d.Set("data_engine_id", dataEngine.DataEngineId)
	}

	return nil
}

func resourceTencentCloudDlcDataEngineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dlc.NewUpdateDataEngineRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dataEngineName := idSplit[0]

	immutableArgs := []string{"engine_type", "data_engine_name", "cluster_type", "mode", "default_data_engine", "cidr_block",
		"pay_mode", "time_span", "time_unit", "auto_renew", "engine_exec_type", "tolerable_queue_time",
		"resource_type", "data_engine_config_pairs", "image_version_name", "main_cluster_name", "auto_authorization",
		"engine_network_id", "engine_generation"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOkExists("size"); ok {
		request.Size = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("min_clusters"); ok {
		request.MinClusters = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_clusters"); ok {
		request.MaxClusters = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_resume"); ok {
		request.AutoResume = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("message"); ok {
		request.Message = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_suspend"); ok {
		request.AutoSuspend = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("crontab_resume_suspend"); ok {
		request.CrontabResumeSuspend = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "crontab_resume_suspend_strategy"); ok {
		crontabResumeSuspendStrategy := dlc.CrontabResumeSuspendStrategy{}
		if v, ok := dMap["resume_time"]; ok {
			crontabResumeSuspendStrategy.ResumeTime = helper.String(v.(string))
		}

		if v, ok := dMap["suspend_time"]; ok {
			crontabResumeSuspendStrategy.SuspendTime = helper.String(v.(string))
		}

		if v, ok := dMap["suspend_strategy"]; ok {
			crontabResumeSuspendStrategy.SuspendStrategy = helper.IntInt64(v.(int))
		}

		request.CrontabResumeSuspendStrategy = &crontabResumeSuspendStrategy
	}

	if v, ok := d.GetOkExists("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("tolerable_queue_time"); ok {
		request.TolerableQueueTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_suspend_time"); ok {
		request.AutoSuspendTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("elastic_switch"); ok {
		request.ElasticSwitch = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("elastic_limit"); ok {
		request.ElasticLimit = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "session_resource_template"); ok {
		sessionResourceTemplate := dlc.SessionResourceTemplate{}
		if v, ok := dMap["driver_size"]; ok {
			sessionResourceTemplate.DriverSize = helper.String(v.(string))
		}

		if v, ok := dMap["executor_size"]; ok {
			sessionResourceTemplate.ExecutorSize = helper.String(v.(string))
		}

		if v, ok := dMap["executor_nums"]; ok {
			sessionResourceTemplate.ExecutorNums = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["executor_max_numbers"]; ok {
			sessionResourceTemplate.ExecutorMaxNumbers = helper.IntUint64(v.(int))
		}

		if v, ok := dMap["running_time_parameters"]; ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				dataEngineConfigPair := dlc.DataEngineConfigPair{}
				if v, ok := dMap["config_item"]; ok {
					dataEngineConfigPair.ConfigItem = helper.String(v.(string))
				}

				if v, ok := dMap["config_value"]; ok {
					dataEngineConfigPair.ConfigItem = helper.String(v.(string))
				}

				sessionResourceTemplate.RunningTimeParameters = append(sessionResourceTemplate.RunningTimeParameters, &dataEngineConfigPair)
			}
		}

		request.SessionResourceTemplate = &sessionResourceTemplate
	}

	request.DataEngineName = &dataEngineName
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dlc data engine failed, reason:%+v", logId, err)
		return err
	}

	// wait
	describeRequest := dlc.NewDescribeDataEngineRequest()
	describeRequest.DataEngineName = helper.String(dataEngineName)
	err = resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDataEngine(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.DataEngine == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] resopse is null, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}

		if *result.Response.DataEngine.State != int64(2) && *result.Response.DataEngine.State != int64(1) {
			e = fmt.Errorf("[DEBUG]%s api[%s] status [%v] not ready , request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), *result.Response.DataEngine.State, describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dlc dataEngine failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDlcDataEngineRead(d, meta)
}

func resourceTencentCloudDlcDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dataEngineName := idSplit[0]

	if err := service.DeleteDlcDataEngineByName(ctx, dataEngineName); err != nil {
		return err
	}

	return nil
}
