/*
Provides a resource to create a tem workload

Example Usage

```hcl
resource "tencentcloud_tem_workload" "workload" {
  application_id = "app-xxx"
  environment_id = "en-xxx"
  deploy_version = "hello-world"
  deploy_mode = "IMAGE"
  img_repo = "tem_demo/tem_demo"
  init_pod_num = 1
  cpu_spec =
  memory_spec =
  post_start = &lt;nil&gt;
  pre_stop = &lt;nil&gt;
  security_group_ids =
  repo_type = 3
  repo_server = &lt;nil&gt;
  tcr_instance_id = &lt;nil&gt;
  env_conf {
		key = "key"
		value = "value"
		type = "default"
		config = "config-name"
		secret = "secret-name"

  }
  storage_confs {
		storage_vol_name = "xxx"
		storage_vol_ip = "0.0.0.0"
		storage_vol_path = "/"

  }
  storage_mount_confs {
		volume_name = "xxx"
		mount_path = "/"

  }
  liveness {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  readiness {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  startup_probe {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  deploy_strategy_conf {
		deploy_strategy_type = 0
		beta_batch_num = 0
		total_batch_count = 1
		batch_interval = 200
		min_available = -1
		force = true

  }
}
```

Import

tem workload can be imported using the id, e.g.

```
terraform import tencentcloud_tem_workload.workload workload_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTemWorkload() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemWorkloadCreate,
		Read:   resourceTencentCloudTemWorkloadRead,
		Update: resourceTencentCloudTemWorkloadUpdate,
		Delete: resourceTencentCloudTemWorkloadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment ID.",
			},

			"deploy_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy version.",
			},

			"deploy_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy mode, support IMAGE.",
			},

			"img_repo": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Repository name.",
			},

			"init_pod_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Initial pod number.",
			},

			"cpu_spec": {
				Required:    true,
				Type:        schema.TypeFloat,
				Description: "Cpu.",
			},

			"memory_spec": {
				Required:    true,
				Type:        schema.TypeFloat,
				Description: "Mem.",
			},

			"post_start": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Mem.",
			},

			"pre_stop": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Mem.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security groups.",
			},

			"repo_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Repo type when deploy: 0: tcr personal; 1: tcr enterprise; 2: public repository; 3: tem host tcr; 4: demo repo.",
			},

			"repo_server": {
				Optional:    true,
				Description: "Repo server addr when deploy by image.",
			},

			"tcr_instance_id": {
				Optional:    true,
				Description: "Tcr instance id when deploy by image.",
			},

			"env_conf": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Env key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Env value.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Env type.",
						},
						"config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Referenced config name.",
						},
						"secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Referenced secret name.",
						},
					},
				},
			},

			"storage_confs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Storage configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_vol_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Volume name.",
						},
						"storage_vol_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Volume ip.",
						},
						"storage_vol_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Volume path.",
						},
					},
				},
			},

			"storage_mount_confs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Storage mount configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Volume name.",
						},
						"mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mount path.",
						},
					},
				},
			},

			"liveness": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Liveness config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Liveness check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Initial delay seconds for liveness check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout seconds for liveness check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Period seconds for liveness check.",
						},
					},
				},
			},

			"readiness": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Readiness check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Initial delay seconds for readiness check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout seconds for readiness check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Period seconds for readiness check.",
						},
					},
				},
			},

			"startup_probe": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Check type, support HttpGet, TcpSocket and Exec.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path.",
						},
						"exec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Startup check port.",
						},
						"initial_delay_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Initial delay seconds for startup check.",
						},
						"timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout seconds for startup check.",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Period seconds for startup check.",
						},
					},
				},
			},

			"deploy_strategy_conf": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Deploy strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deploy_strategy_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Strategy type, 0 means auto, 1 means manual, 2 means manual with beta batch.",
						},
						"beta_batch_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Beta batch number.",
						},
						"total_batch_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Total batch number.",
						},
						"batch_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Interval between batches.",
						},
						"min_available": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimal availabe instances duration deployment.",
						},
						"force": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Force update.",
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
		response      = tem.NewDeployApplicationResponse()
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

	if v, ok := d.GetOkExists("init_pod_num"); ok {
		request.InitPodNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("cpu_spec"); ok {
		request.CpuSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("memory_spec"); ok {
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

	if v, ok := d.GetOkExists("repo_type"); ok {
		request.RepoType = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("repo_server"); v != nil {
	}

	if v, _ := d.GetOk("tcr_instance_id"); v != nil {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem workload failed, reason:%+v", logId, err)
		return err
	}

	applicationId = *response.Response.ApplicationId
	d.SetId(strings.Join([]string{applicationId, environmentId}, FILED_SP))

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
	applicationId := idSplit[0]
	environmentId := idSplit[1]

	workload, err := service.DescribeTemWorkloadById(ctx, applicationId, environmentId)
	if err != nil {
		return err
	}

	if workload == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemWorkload` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
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
		_ = d.Set("img_repo", workload.ImgRepo)
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

			if workload.EnvConf.Key != nil {
				envConfMap["key"] = workload.EnvConf.Key
			}

			if workload.EnvConf.Value != nil {
				envConfMap["value"] = workload.EnvConf.Value
			}

			if workload.EnvConf.Type != nil {
				envConfMap["type"] = workload.EnvConf.Type
			}

			if workload.EnvConf.Config != nil {
				envConfMap["config"] = workload.EnvConf.Config
			}

			if workload.EnvConf.Secret != nil {
				envConfMap["secret"] = workload.EnvConf.Secret
			}

			envConfList = append(envConfList, envConfMap)
		}

		_ = d.Set("env_conf", envConfList)

	}

	if workload.StorageConfs != nil {
		storageConfsList := []interface{}{}
		for _, storageConfs := range workload.StorageConfs {
			storageConfsMap := map[string]interface{}{}

			if workload.StorageConfs.StorageVolName != nil {
				storageConfsMap["storage_vol_name"] = workload.StorageConfs.StorageVolName
			}

			if workload.StorageConfs.StorageVolIp != nil {
				storageConfsMap["storage_vol_ip"] = workload.StorageConfs.StorageVolIp
			}

			if workload.StorageConfs.StorageVolPath != nil {
				storageConfsMap["storage_vol_path"] = workload.StorageConfs.StorageVolPath
			}

			storageConfsList = append(storageConfsList, storageConfsMap)
		}

		_ = d.Set("storage_confs", storageConfsList)

	}

	if workload.StorageMountConfs != nil {
		storageMountConfsList := []interface{}{}
		for _, storageMountConfs := range workload.StorageMountConfs {
			storageMountConfsMap := map[string]interface{}{}

			if workload.StorageMountConfs.VolumeName != nil {
				storageMountConfsMap["volume_name"] = workload.StorageMountConfs.VolumeName
			}

			if workload.StorageMountConfs.MountPath != nil {
				storageMountConfsMap["mount_path"] = workload.StorageMountConfs.MountPath
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

	if workload.DeployStrategyConf != nil {
		deployStrategyConfMap := map[string]interface{}{}

		if workload.DeployStrategyConf.DeployStrategyType != nil {
			deployStrategyConfMap["deploy_strategy_type"] = workload.DeployStrategyConf.DeployStrategyType
		}

		if workload.DeployStrategyConf.BetaBatchNum != nil {
			deployStrategyConfMap["beta_batch_num"] = workload.DeployStrategyConf.BetaBatchNum
		}

		if workload.DeployStrategyConf.TotalBatchCount != nil {
			deployStrategyConfMap["total_batch_count"] = workload.DeployStrategyConf.TotalBatchCount
		}

		if workload.DeployStrategyConf.BatchInterval != nil {
			deployStrategyConfMap["batch_interval"] = workload.DeployStrategyConf.BatchInterval
		}

		if workload.DeployStrategyConf.MinAvailable != nil {
			deployStrategyConfMap["min_available"] = workload.DeployStrategyConf.MinAvailable
		}

		if workload.DeployStrategyConf.Force != nil {
			deployStrategyConfMap["force"] = workload.DeployStrategyConf.Force
		}

		_ = d.Set("deploy_strategy_conf", []interface{}{deployStrategyConfMap})
	}

	return nil
}

func resourceTencentCloudTemWorkloadUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_workload.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewDeployApplicationRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	applicationId := idSplit[0]
	environmentId := idSplit[1]

	request.ApplicationId = &applicationId
	request.EnvironmentId = &environmentId

	immutableArgs := []string{"application_id", "environment_id", "deploy_version", "deploy_mode", "img_repo", "init_pod_num", "cpu_spec", "memory_spec", "post_start", "pre_stop", "security_group_ids", "repo_type", "repo_server", "tcr_instance_id", "env_conf", "storage_confs", "storage_mount_confs", "liveness", "readiness", "startup_probe", "deploy_strategy_conf"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("deploy_version") {
		if v, ok := d.GetOk("deploy_version"); ok {
			request.DeployVersion = helper.String(v.(string))
		}
	}

	if d.HasChange("deploy_mode") {
		if v, ok := d.GetOk("deploy_mode"); ok {
			request.DeployMode = helper.String(v.(string))
		}
	}

	if d.HasChange("img_repo") {
		if v, ok := d.GetOk("img_repo"); ok {
			request.ImgRepo = helper.String(v.(string))
		}
	}

	if d.HasChange("init_pod_num") {
		if v, ok := d.GetOkExists("init_pod_num"); ok {
			request.InitPodNum = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("cpu_spec") {
		if v, ok := d.GetOkExists("cpu_spec"); ok {
			request.CpuSpec = helper.Float64(v.(float64))
		}
	}

	if d.HasChange("memory_spec") {
		if v, ok := d.GetOkExists("memory_spec"); ok {
			request.MemorySpec = helper.Float64(v.(float64))
		}
	}

	if d.HasChange("post_start") {
		if v, ok := d.GetOk("post_start"); ok {
			request.PostStart = helper.String(v.(string))
		}
	}

	if d.HasChange("pre_stop") {
		if v, ok := d.GetOk("pre_stop"); ok {
			request.PreStop = helper.String(v.(string))
		}
	}

	if d.HasChange("security_group_ids") {
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroupIdsSet := v.(*schema.Set).List()
			for i := range securityGroupIdsSet {
				securityGroupIds := securityGroupIdsSet[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
			}
		}
	}

	if d.HasChange("repo_type") {
		if v, ok := d.GetOkExists("repo_type"); ok {
			request.RepoType = helper.IntInt64(v.(int))
		}
	}

	if v, _ := d.GetOk("repo_server"); v != nil {
	}

	if v, _ := d.GetOk("tcr_instance_id"); v != nil {
	}

	if d.HasChange("env_conf") {
		if v, ok := d.GetOk("env_conf"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("storage_confs") {
		if v, ok := d.GetOk("storage_confs"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("storage_mount_confs") {
		if v, ok := d.GetOk("storage_mount_confs"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("liveness") {
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
	}

	if d.HasChange("readiness") {
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
	}

	if d.HasChange("startup_probe") {
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
	}

	if d.HasChange("deploy_strategy_conf") {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem workload failed, reason:%+v", logId, err)
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
	applicationId := idSplit[0]
	environmentId := idSplit[1]

	if err := service.DeleteTemWorkloadById(ctx, applicationId, environmentId); err != nil {
		return err
	}

	return nil
}
