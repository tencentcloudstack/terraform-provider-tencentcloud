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
				Description: "Engine type, only support: spark/presto.",
			},

			"data_engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine name.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine cluster type, only support: spark_cu/presto_cu.",
			},

			"mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Engine mode, only support 1: ByAmount, 2: YearlyAndMonthly.",
			},

			"auto_resume": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically start the cluster, prepay not support.",
			},

			"size": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Cluster size. Required when updating.",
			},

			"min_clusters": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine min size, greater than or equal to 1 and MaxClusters bigger than MinClusters.",
			},

			"max_clusters": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine max cluster size, MaxClusters less than or equal to 10 and MaxClusters bigger than MinClusters.",
			},

			"default_data_engine": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is the default virtual cluster.",
			},

			"cidr_block": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine VPC network segment, just like 192.0.2.1/24.",
			},

			"message": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine description information.",
			},

			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine pay mode type, only support 0: postPay, 1: prePay(default).",
			},

			"time_span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine TimeSpan, prePay: minimum of 1, representing one month of purchasing resources, with a maximum of 120, default 3600, postPay: fixed fee of 3600.",
			},

			"time_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine TimeUnit, prePay: use m(default), postPay: use h.",
			},

			"auto_renew": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine auto renew, only support 0: Default, 1: AutoRenewON, 2: AutoRenewOFF.",
			},

			"auto_suspend": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically suspend the cluster, prepay not support.",
			},

			"crontab_resume_suspend": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Engine crontab resume or suspend strategy, only support: 0: Wait(default), 1: Kill.",
			},

			"crontab_resume_suspend_strategy": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Engine auto suspend strategy, when AutoSuspend is true, CrontabResumeSuspend must stop.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resume_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scheduled pull-up time: For example: 8 o&amp;#39;clock on Monday is expressed as 1000000-08:00:00.",
						},
						"suspend_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scheduled suspension time: For example: 20 o&amp;#39;clock on Monday is expressed as 1000000-20:00:00.",
						},
						"suspend_strategy": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Suspend configuration: 0 (default): wait for the task to end before suspending, 1: force suspend.",
						},
					},
				},
			},

			"engine_exec_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine exec type, only support SQL(default) or BATCH.",
			},

			"max_concurrency": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Maximum number of concurrent tasks in a single cluster, default 5.",
			},

			"tolerable_queue_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Tolerable queuing time, default 0. scaling may be triggered when tasks are queued for longer than the tolerable time. if this parameter is 0, it means that capacity expansion may be triggered immediately once a task is queued.",
			},

			"auto_suspend_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Cluster automatic suspension time, default 10 minutes.",
			},

			"resource_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Engine resource type not match, only support: Standard_CU/Memory_CU(only BATCH ExecType).",
			},

			"data_engine_config_pairs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Cluster advanced configuration.",
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
				Description: "Cluster image version name. Such as SuperSQL-P 1.1; SuperSQL-S 3.2, etc., do not upload, and create a cluster with the latest mirror version by default.",
			},

			"main_cluster_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Primary cluster name, specified when creating a disaster recovery cluster.",
			},

			"elastic_switch": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "For spark Batch ExecType, yearly and monthly cluster whether to enable elasticity.",
			},

			"elastic_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "For spark Batch ExecType, yearly and monthly cluster elastic limit.",
			},

			"session_resource_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "For spark Batch ExecType, cluster session resource configuration template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"driver_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
						},
						"executor_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.",
						},
						"executor_nums": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.",
						},
						"executor_max_numbers": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDlcDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = dlc.NewCreateDataEngineRequest()
		response       = dlc.NewCreateDataEngineResponse()
		dataEngineId   string
		dataEngineName string
	)
	if v, ok := d.GetOk("engine_type"); ok {
		request.EngineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
		request.DataEngineName = helper.String(v.(string))
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
		request.SessionResourceTemplate = &sessionResourceTemplate
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dlc dataEngine failed, reason:%+v", logId, err)
		return err
	}

	dataEngineId = *response.Response.DataEngineId
	d.SetId(dataEngineName + tccommon.FILED_SP + dataEngineId)

	describeRequest := dlc.NewDescribeDataEngineRequest()
	describeRequest.DataEngineName = helper.String(dataEngineName)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDataEngine(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(),
				describeRequest.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || result.Response.DataEngine == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] resopse is null, request body [%s], response body [%s]\n", logId,
				describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}
		if *result.Response.DataEngine.State != int64(2) && *result.Response.DataEngine.State != int64(1) {
			e = fmt.Errorf("[DEBUG]%s api[%s] status [%v] not ready , request body [%s], response body [%s]\n",
				logId, describeRequest.GetAction(), *result.Response.DataEngine.State, describeRequest.ToJsonString(), result.ToJsonString())
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

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcDataEngine` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

		_ = d.Set("session_resource_template", []interface{}{sessionResourceTemplateMap})
	}

	return nil
}

func resourceTencentCloudDlcDataEngineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_data_engine.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dlc.NewUpdateDataEngineRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	dataEngineName := idSplit[0]
	request.DataEngineName = &dataEngineName

	immutableArgs := []string{"engine_type", "data_engine_name", "cluster_type", "mode", "default_data_engine", "cidr_block",
		"pay_mode", "time_span", "time_unit", "auto_renew", "engine_exec_type", "tolerable_queue_time",
		"resource_type", "data_engine_config_pairs", "image_version_name", "main_cluster_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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

	if v, ok := d.GetOk("message"); ok {
		request.Message = helper.String(v.(string))
	}

	if d.HasChange("auto_suspend") {
		if v, ok := d.GetOkExists("auto_suspend"); ok {
			request.AutoSuspend = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("crontab_resume_suspend") {
		if v, ok := d.GetOkExists("crontab_resume_suspend"); ok {
			request.CrontabResumeSuspend = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("crontab_resume_suspend_strategy") {
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
	}

	if d.HasChange("max_concurrency") {
		if v, ok := d.GetOkExists("max_concurrency"); ok {
			request.MaxConcurrency = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("tolerable_queue_time") {
		if v, ok := d.GetOkExists("tolerable_queue_time"); ok {
			request.TolerableQueueTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("auto_suspend_time") {
		if v, ok := d.GetOkExists("auto_suspend_time"); ok {
			request.AutoSuspendTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("elastic_switch") {
		if v, ok := d.GetOkExists("elastic_switch"); ok {
			request.ElasticSwitch = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("elastic_limit") {
		if v, ok := d.GetOkExists("elastic_limit"); ok {
			request.ElasticLimit = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("session_resource_template") {
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
			request.SessionResourceTemplate = &sessionResourceTemplate
		}
	}

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
		log.Printf("[CRITAL]%s update dlc dataEngine failed, reason:%+v", logId, err)
		return err
	}

	describeRequest := dlc.NewDescribeDataEngineRequest()
	describeRequest.DataEngineName = helper.String(dataEngineName)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeDataEngine(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(),
				describeRequest.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || result.Response.DataEngine == nil {
			e = fmt.Errorf("[DEBUG]%s api[%s] resopse is null, request body [%s], response body [%s]\n", logId,
				request.GetAction(), request.ToJsonString(), result.ToJsonString())
			log.Println(e)
			return resource.RetryableError(e)
		}
		if *result.Response.DataEngine.State != int64(2) && *result.Response.DataEngine.State != int64(1) {
			e = fmt.Errorf("[DEBUG]%s api[%s] status [%v] not ready , request body [%s], response body [%s]\n",
				logId, describeRequest.GetAction(), *result.Response.DataEngine.State, describeRequest.ToJsonString(), result.ToJsonString())
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
