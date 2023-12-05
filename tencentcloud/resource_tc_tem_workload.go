package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemWorkload() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemWorkloadRead,
		Create: resourceTencentCloudTemWorkloadCreate,
		Update: resourceTencentCloudTemWorkloadUpdate,
		Delete: resourceTencentCloudTemWorkloadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "application ID.",
			},

			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "environment ID.",
			},

			"deploy_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "deploy version.",
			},

			"deploy_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "deploy mode, support IMAGE.",
			},

			"img_repo": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "repository name.",
			},

			"init_pod_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "initial pod number.",
			},

			"cpu_spec": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "cpu.",
			},

			"memory_spec": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "mem.",
			},

			"post_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "mem.",
			},

			"pre_stop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "mem.",
			},

			"security_group_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "security groups.",
			},

			"repo_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "repo type when deploy: 0: tcr personal; 1: tcr enterprise; 2: public repository; 3: tem host tcr; 4: demo repo.",
			},

			"repo_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "repo server addr when deploy by image.",
			},

			"tcr_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "tcr instance id when deploy by image.",
			},

			"env_conf": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "env key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "env value.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "env type, support default, referenced.",
						},
						"config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "referenced config name when type=referenced.",
						},
						"secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "referenced secret name when type=referenced.",
						},
					},
				},
			},

			"storage_confs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "storage configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_vol_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume name.",
						},
						"storage_vol_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume ip.",
						},
						"storage_vol_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume path.",
						},
					},
				},
			},

			"storage_mount_confs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "storage mount configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "volume name.",
						},
						"mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mount path.",
						},
					},
				},
			},

			"liveness": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "liveness config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "liveness check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "initial delay seconds for liveness check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "timeout seconds for liveness check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "period seconds for liveness check.",
						},
					},
				},
			},

			"readiness": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "readiness check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "initial delay seconds for readiness check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "timeout seconds for readiness check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "period seconds for readiness check.",
						},
					},
				},
			},

			"startup_probe": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "startup check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "initial delay seconds for startup check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "timeout seconds for startup check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "period seconds for startup check.",
						},
					},
				},
			},

			"deploy_strategy_conf": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "deploy strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deploy_strategy_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "strategy type, 0 means auto, 1 means manual, 2 means manual with beta batch.",
						},
						"beta_batch_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "beta batch number.",
						},
						"total_batch_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "total batch number.",
						},
						"batch_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "interval between batches.",
						},
						"min_available": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "minimal available instances duration deployment.",
						},
						"force": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "force update.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemWorkloadCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_workload.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewDeployApplicationRequest()
		applicationId string
		environmentId string
	)

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deploy_version"); ok {
		request.DeployVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deploy_mode"); ok {
		request.DeployMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("img_repo"); ok {
		request.ImgRepo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("init_pod_num"); ok {
		request.InitPodNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cpu_spec"); ok {
		request.CpuSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOk("memory_spec"); ok {
		request.MemorySpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOk("post_start"); ok {
		request.PostStart = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pre_stop"); ok {
		request.PreStop = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("repo_type"); ok {
		request.RepoType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("repo_server"); ok {
		request.RepoServer = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tcr_instance_id"); ok {
		request.TcrInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("env_conf"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			pair := tem.Pair{}
			if v, ok := dMap["key"]; ok {
				pair.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				pair.Value = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				pair.Type = helper.String(v.(string))
			}
			if v, ok := dMap["config"]; ok {
				pair.Config = helper.String(v.(string))
			}
			if v, ok := dMap["secret"]; ok {
				pair.Secret = helper.String(v.(string))
			}
			request.EnvConf = append(request.EnvConf, &pair)
		}
	}

	if v, ok := d.GetOk("storage_confs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			storageConf := tem.StorageConf{}
			if v, ok := dMap["storage_vol_name"]; ok {
				storageConf.StorageVolName = helper.String(v.(string))
			}
			if v, ok := dMap["storage_vol_ip"]; ok {
				storageConf.StorageVolIp = helper.String(v.(string))
			}
			if v, ok := dMap["storage_vol_path"]; ok {
				storageConf.StorageVolPath = helper.String(v.(string))
			}
			request.StorageConfs = append(request.StorageConfs, &storageConf)

		}
	}

	if v, ok := d.GetOk("storage_mount_confs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			storageMountConf := tem.StorageMountConf{}
			if v, ok := dMap["volume_name"]; ok {
				storageMountConf.VolumeName = helper.String(v.(string))
			}
			if v, ok := dMap["mount_path"]; ok {
				storageMountConf.MountPath = helper.String(v.(string))
			}
			request.StorageMountConfs = append(request.StorageMountConfs, &storageMountConf)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "liveness"); ok {
		healthCheckConfig := tem.HealthCheckConfig{}
		if v, ok := dMap["type"]; ok {
			healthCheckConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["protocol"]; ok {
			healthCheckConfig.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["path"]; ok {
			healthCheckConfig.Path = helper.String(v.(string))
		}
		if v, ok := dMap["exec"]; ok {
			healthCheckConfig.Exec = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			healthCheckConfig.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["initial_delay_seconds"]; ok {
			healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["timeout_seconds"]; ok {
			healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["period_seconds"]; ok {
			healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
		}
		request.Liveness = &healthCheckConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "readiness"); ok {
		healthCheckConfig := tem.HealthCheckConfig{}
		if v, ok := dMap["type"]; ok {
			healthCheckConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["protocol"]; ok {
			healthCheckConfig.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["path"]; ok {
			healthCheckConfig.Path = helper.String(v.(string))
		}
		if v, ok := dMap["exec"]; ok {
			healthCheckConfig.Exec = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			healthCheckConfig.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["initial_delay_seconds"]; ok {
			healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["timeout_seconds"]; ok {
			healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["period_seconds"]; ok {
			healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
		}
		request.Readiness = &healthCheckConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "startup_probe"); ok {
		healthCheckConfig := tem.HealthCheckConfig{}
		if v, ok := dMap["type"]; ok {
			healthCheckConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["protocol"]; ok {
			healthCheckConfig.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["path"]; ok {
			healthCheckConfig.Path = helper.String(v.(string))
		}
		if v, ok := dMap["exec"]; ok {
			healthCheckConfig.Exec = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			healthCheckConfig.Port = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["initial_delay_seconds"]; ok {
			healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["timeout_seconds"]; ok {
			healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["period_seconds"]; ok {
			healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
		}
		request.StartupProbe = &healthCheckConfig

	}

	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_strategy_conf"); ok {
		deployStrategyConf := tem.DeployStrategyConf{}
		if v, ok := dMap["deploy_strategy_type"]; ok {
			deployStrategyConf.DeployStrategyType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["beta_batch_num"]; ok {
			deployStrategyConf.BetaBatchNum = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["total_batch_count"]; ok {
			deployStrategyConf.TotalBatchCount = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["batch_interval"]; ok {
			deployStrategyConf.BatchInterval = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["min_available"]; ok {
			deployStrategyConf.MinAvailable = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["force"]; ok {
			deployStrategyConf.Force = helper.Bool(v.(bool))
		}
		request.DeployStrategyConf = &deployStrategyConf

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().DeployApplication(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem workload failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(environmentId + FILED_SP + applicationId)
	return resourceTencentCloudTemWorkloadRead(d, meta)
}

func resourceTencentCloudTemWorkloadRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_workload.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	workloads, err := service.DescribeTemWorkload(ctx, environmentId, applicationId)

	workload := workloads.Result

	if err != nil {
		return err
	}

	if workload == nil {
		d.SetId("")
		return fmt.Errorf("resource `workload` does not exist")
	}

	if workload.ApplicationId != nil {
		_ = d.Set("application_id", workload.ApplicationId)
	}

	if workload.EnvironmentId != nil {
		_ = d.Set("environment_id", workload.EnvironmentId)
	}

	if workload.DeployVersion != nil {
		_ = d.Set("deploy_version", workload.DeployVersion)
	}

	if workload.DeployMode != nil {
		_ = d.Set("deploy_mode", workload.DeployMode)
	}

	if workload.ImgRepo != nil {
		_ = d.Set("img_repo", workload.ImgName)
	}

	if workload.InitPodNum != nil {
		_ = d.Set("init_pod_num", workload.InitPodNum)
	}

	if workload.CpuSpec != nil {
		_ = d.Set("cpu_spec", workload.CpuSpec)
	}

	if workload.MemorySpec != nil {
		_ = d.Set("memory_spec", workload.MemorySpec)
	}

	if workload.PostStart != nil {
		_ = d.Set("post_start", workload.PostStart)
	}

	if workload.PreStop != nil {
		_ = d.Set("pre_stop", workload.PreStop)
	}

	if workload.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", workload.SecurityGroupIds)
	}

	if workload.RepoType != nil {
		_ = d.Set("repo_type", workload.RepoType)
	}

	if workload.RepoServer != nil {
		_ = d.Set("repo_server", workload.RepoServer)
	}

	if workload.TcrInstanceId != nil {
		_ = d.Set("tcr_instance_id", workload.TcrInstanceId)
	}

	if workload.EnvConf != nil {
		envConfList := []interface{}{}
		for _, envConf := range workload.EnvConf {
			envConfMap := map[string]interface{}{}
			if envConf.Type != nil {
				if *envConf.Type != "reserved" {
					envConfMap["type"] = envConf.Type
					if envConf.Key != nil {
						envConfMap["key"] = envConf.Key
					}
					if envConf.Value != nil {
						envConfMap["value"] = envConf.Value
					}
					if envConf.Config != nil {
						envConfMap["config"] = envConf.Config
					}
					if envConf.Secret != nil {
						envConfMap["secret"] = envConf.Secret
					}
					envConfList = append(envConfList, envConfMap)
				}
			}

		}
		_ = d.Set("env_conf", envConfList)
	}

	if workload.StorageConfs != nil {
		storageConfsList := []interface{}{}
		for _, storageConfs := range workload.StorageConfs {
			storageConfsMap := map[string]interface{}{}
			if storageConfs.StorageVolName != nil {
				storageConfsMap["storage_vol_name"] = storageConfs.StorageVolName
			}
			if storageConfs.StorageVolIp != nil {
				storageConfsMap["storage_vol_ip"] = storageConfs.StorageVolIp
			}
			if storageConfs.StorageVolPath != nil {
				storageConfsMap["storage_vol_path"] = storageConfs.StorageVolPath
			}

			storageConfsList = append(storageConfsList, storageConfsMap)
		}
		_ = d.Set("storage_confs", storageConfsList)
	}

	if workload.StorageMountConfs != nil {
		storageMountConfsList := []interface{}{}
		for _, storageMountConfs := range workload.StorageMountConfs {
			storageMountConfsMap := map[string]interface{}{}
			if storageMountConfs.VolumeName != nil {
				storageMountConfsMap["volume_name"] = storageMountConfs.VolumeName
			}
			if storageMountConfs.MountPath != nil {
				storageMountConfsMap["mount_path"] = storageMountConfs.MountPath
			}

			storageMountConfsList = append(storageMountConfsList, storageMountConfsMap)
		}
		_ = d.Set("storage_mount_confs", storageMountConfsList)
	}

	if workload.Liveness != nil {
		livenessMap := map[string]interface{}{}
		if workload.Liveness.Type != nil {
			livenessMap["type"] = workload.Liveness.Type
		}
		if workload.Liveness.Protocol != nil {
			livenessMap["protocol"] = workload.Liveness.Protocol
		}
		if workload.Liveness.Path != nil {
			livenessMap["path"] = workload.Liveness.Path
		}
		if workload.Liveness.Exec != nil {
			livenessMap["exec"] = workload.Liveness.Exec
		}
		if workload.Liveness.Port != nil {
			livenessMap["port"] = workload.Liveness.Port
		}
		if workload.Liveness.InitialDelaySeconds != nil {
			livenessMap["initial_delay_seconds"] = workload.Liveness.InitialDelaySeconds
		}
		if workload.Liveness.TimeoutSeconds != nil {
			livenessMap["timeout_seconds"] = workload.Liveness.TimeoutSeconds
		}
		if workload.Liveness.PeriodSeconds != nil {
			livenessMap["period_seconds"] = workload.Liveness.PeriodSeconds
		}

		_ = d.Set("liveness", []interface{}{livenessMap})
	}

	if workload.Readiness != nil {
		readinessMap := map[string]interface{}{}
		if workload.Readiness.Type != nil {
			readinessMap["type"] = workload.Readiness.Type
		}
		if workload.Readiness.Protocol != nil {
			readinessMap["protocol"] = workload.Readiness.Protocol
		}
		if workload.Readiness.Path != nil {
			readinessMap["path"] = workload.Readiness.Path
		}
		if workload.Readiness.Exec != nil {
			readinessMap["exec"] = workload.Readiness.Exec
		}
		if workload.Readiness.Port != nil {
			readinessMap["port"] = workload.Readiness.Port
		}
		if workload.Readiness.InitialDelaySeconds != nil {
			readinessMap["initial_delay_seconds"] = workload.Readiness.InitialDelaySeconds
		}
		if workload.Readiness.TimeoutSeconds != nil {
			readinessMap["timeout_seconds"] = workload.Readiness.TimeoutSeconds
		}
		if workload.Readiness.PeriodSeconds != nil {
			readinessMap["period_seconds"] = workload.Readiness.PeriodSeconds
		}

		_ = d.Set("readiness", []interface{}{readinessMap})
	}

	if workload.StartupProbe != nil {
		startupProbeMap := map[string]interface{}{}
		if workload.StartupProbe.Type != nil {
			startupProbeMap["type"] = workload.StartupProbe.Type
		}
		if workload.StartupProbe.Protocol != nil {
			startupProbeMap["protocol"] = workload.StartupProbe.Protocol
		}
		if workload.StartupProbe.Path != nil {
			startupProbeMap["path"] = workload.StartupProbe.Path
		}
		if workload.StartupProbe.Exec != nil {
			startupProbeMap["exec"] = workload.StartupProbe.Exec
		}
		if workload.StartupProbe.Port != nil {
			startupProbeMap["port"] = workload.StartupProbe.Port
		}
		if workload.StartupProbe.InitialDelaySeconds != nil {
			startupProbeMap["initial_delay_seconds"] = workload.StartupProbe.InitialDelaySeconds
		}
		if workload.StartupProbe.TimeoutSeconds != nil {
			startupProbeMap["timeout_seconds"] = workload.StartupProbe.TimeoutSeconds
		}
		if workload.StartupProbe.PeriodSeconds != nil {
			startupProbeMap["period_seconds"] = workload.StartupProbe.PeriodSeconds
		}

		_ = d.Set("startup_probe", []interface{}{startupProbeMap})
	}

	return nil
}

func resourceTencentCloudTemWorkloadUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_workload.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tem.NewDeployApplicationRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	if d.HasChange("deploy_version") || d.HasChange("deploy_mode") || d.HasChange("img_repo") || d.HasChange("init_pod_num") ||
		d.HasChange("cpu_spec") || d.HasChange("memory_spec") || d.HasChange("post_start") || d.HasChange("pre_stop") || d.HasChange("security_group_ids") ||
		d.HasChange("repo_type") || d.HasChange("repo_server") || d.HasChange("tcr_instance_id") || d.HasChange("env_conf") || d.HasChange("storage_confs") ||
		d.HasChange("storage_mount_confs") || d.HasChange("liveness") || d.HasChange("readiness") || d.HasChange("startup_probe") || d.HasChange("deploy_strategy_conf") {

		if v, ok := d.GetOk("deploy_version"); ok {
			request.DeployVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("deploy_mode"); ok {
			request.DeployMode = helper.String(v.(string))
		}

		if v, ok := d.GetOk("img_repo"); ok {
			request.ImgRepo = helper.String(v.(string))
		}

		if v, ok := d.GetOk("init_pod_num"); ok {
			request.InitPodNum = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("cpu_spec"); ok {
			request.CpuSpec = helper.Float64(v.(float64))
		}

		if v, ok := d.GetOk("memory_spec"); ok {
			request.MemorySpec = helper.Float64(v.(float64))
		}

		if v, ok := d.GetOk("post_start"); ok {
			request.PostStart = helper.String(v.(string))
		}

		if v, ok := d.GetOk("pre_stop"); ok {
			request.PreStop = helper.String(v.(string))
		}

		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroupIdsSet := v.(*schema.Set).List()
			for i := range securityGroupIdsSet {
				securityGroupIds := securityGroupIdsSet[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
			}
		}

		if v, ok := d.GetOk("repo_type"); ok {
			request.RepoType = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("repo_server"); ok {
			request.RepoServer = helper.String(v.(string))
		}

		if v, ok := d.GetOk("tcr_instance_id"); ok {
			request.TcrInstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("env_conf"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				pair := tem.Pair{}
				if v, ok := dMap["key"]; ok {
					pair.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					pair.Value = helper.String(v.(string))
				}
				if v, ok := dMap["type"]; ok {
					pair.Type = helper.String(v.(string))
				}
				if v, ok := dMap["config"]; ok {
					pair.Config = helper.String(v.(string))
				}
				if v, ok := dMap["secret"]; ok {
					pair.Secret = helper.String(v.(string))
				}
				request.EnvConf = append(request.EnvConf, &pair)
			}
		}

		if v, ok := d.GetOk("storage_confs"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				storageConf := tem.StorageConf{}
				if v, ok := dMap["storage_vol_name"]; ok {
					storageConf.StorageVolName = helper.String(v.(string))
				}
				if v, ok := dMap["storage_vol_ip"]; ok {
					storageConf.StorageVolIp = helper.String(v.(string))
				}
				if v, ok := dMap["storage_vol_path"]; ok {
					storageConf.StorageVolPath = helper.String(v.(string))
				}
				request.StorageConfs = append(request.StorageConfs, &storageConf)
			}
		}

		if v, ok := d.GetOk("storage_mount_confs"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				storageMountConf := tem.StorageMountConf{}
				if v, ok := dMap["volume_name"]; ok {
					storageMountConf.VolumeName = helper.String(v.(string))
				}
				if v, ok := dMap["mount_path"]; ok {
					storageMountConf.MountPath = helper.String(v.(string))
				}
				request.StorageMountConfs = append(request.StorageMountConfs, &storageMountConf)
			}
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "liveness"); ok {
			healthCheckConfig := tem.HealthCheckConfig{}
			if v, ok := dMap["type"]; ok {
				healthCheckConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				healthCheckConfig.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["path"]; ok {
				healthCheckConfig.Path = helper.String(v.(string))
			}
			if v, ok := dMap["exec"]; ok {
				healthCheckConfig.Exec = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				healthCheckConfig.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["initial_delay_seconds"]; ok {
				healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["timeout_seconds"]; ok {
				healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["period_seconds"]; ok {
				healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
			}
			request.Liveness = &healthCheckConfig
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "readiness"); ok {
			healthCheckConfig := tem.HealthCheckConfig{}
			if v, ok := dMap["type"]; ok {
				healthCheckConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				healthCheckConfig.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["path"]; ok {
				healthCheckConfig.Path = helper.String(v.(string))
			}
			if v, ok := dMap["exec"]; ok {
				healthCheckConfig.Exec = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				healthCheckConfig.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["initial_delay_seconds"]; ok {
				healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["timeout_seconds"]; ok {
				healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["period_seconds"]; ok {
				healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
			}
			request.Readiness = &healthCheckConfig
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "startup_probe"); ok {
			healthCheckConfig := tem.HealthCheckConfig{}
			if v, ok := dMap["type"]; ok {
				healthCheckConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				healthCheckConfig.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["path"]; ok {
				healthCheckConfig.Path = helper.String(v.(string))
			}
			if v, ok := dMap["exec"]; ok {
				healthCheckConfig.Exec = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				healthCheckConfig.Port = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["initial_delay_seconds"]; ok {
				healthCheckConfig.InitialDelaySeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["timeout_seconds"]; ok {
				healthCheckConfig.TimeoutSeconds = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["period_seconds"]; ok {
				healthCheckConfig.PeriodSeconds = helper.IntInt64(v.(int))
			}
			request.StartupProbe = &healthCheckConfig
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "deploy_strategy_conf"); ok {
			deployStrategyConf := tem.DeployStrategyConf{}
			if v, ok := dMap["deploy_strategy_type"]; ok {
				deployStrategyConf.DeployStrategyType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["beta_batch_num"]; ok {
				deployStrategyConf.BetaBatchNum = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["total_batch_count"]; ok {
				deployStrategyConf.TotalBatchCount = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["batch_interval"]; ok {
				deployStrategyConf.BatchInterval = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["min_available"]; ok {
				deployStrategyConf.MinAvailable = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["force"]; ok {
				deployStrategyConf.Force = helper.Bool(v.(bool))
			}
			request.DeployStrategyConf = &deployStrategyConf
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().DeployApplication(request)
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

	return resourceTencentCloudTemWorkloadRead(d, meta)
}

func resourceTencentCloudTemWorkloadDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_workload.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	if err := service.DeleteTemWorkloadById(ctx, environmentId, applicationId); err != nil {
		return err
	}

	return nil
}
