package tke

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func TKEGpuArgsSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mig_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether to enable MIG.",
		},
		"driver": {
			Type:         schema.TypeMap,
			Optional:     true,
			ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
			Description:  "GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
		},
		"cuda": {
			Type:         schema.TypeMap,
			Optional:     true,
			ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
			Description:  "CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
		},
		"cudnn": {
			Type:         schema.TypeMap,
			Optional:     true,
			ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
			Description: "cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`." +
				" `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.",
		},
		"custom_driver": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.",
		},
	}
}

func TkeInstanceAdvancedSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mount_target": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Mount target. Default is not mounting.",
		},
		"docker_graph_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "/var/lib/docker",
			Description: "Docker graph path. Default is `/var/lib/docker`.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			MaxItems:    11,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:         schema.TypeString,
						ForceNew:     true,
						Optional:     true,
						Default:      svcas.SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"file_system": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Default:     "",
						Description: "File system, e.g. `ext3/ext4/xfs`.",
					},
					"auto_format_and_mount": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
						Description: "Indicate whether to auto format and mount or not. Default is `false`.",
					},
					"mount_target": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Default:     "",
						Description: "Mount target.",
					},
					"disk_partition": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.",
					},
				},
			},
		},
		"extra_args": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Custom parameter information related to the node. This is a white-list parameter.",
		},
		"user_data": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Base64-encoded User Data text, the length limit is 16KB.",
		},
		"is_schedule": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "Indicate to schedule the adding node or not. Default is true.",
		},
		"desired_pod_num": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Indicate to set desired pod number in node. valid when the cluster is podCIDR.",
		},
		"gpu_args": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: TKEGpuArgsSetting(),
			},
			Description: "GPU driver parameters.",
		},
	}
}

func ResourceTencentCloudTkeClusterAttachment() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "ID of the cluster.",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the CVM instance, this cvm will reinstall the system.",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: tccommon.ValidateAsConfigPassword,
			Description:  "Password to access, should be set if `key_ids` not set.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.",
		},
		"hostname": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Description: "The host name of the attached instance. " +
				"Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. " +
				"Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. " +
				"Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
		},
		"worker_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeInstanceAdvancedSetting(),
			},
			Description: "Deploy the machine configuration information of the 'WORKER', commonly used to attach existing instances.",
		},
		"worker_config_overrides": {
			Type:     schema.TypeList,
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeInstanceAdvancedSetting(),
			},
			Description: "Override variable worker_config, commonly used to attach existing instances.",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Labels of tke attachment exits CVM.",
		},
		"unschedulable": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
		},
		//compute
		"security_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "A list of security group IDs after attach to cluster.",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the node.",
		},
	}

	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterAttachmentCreate,
		Read:   resourceTencentCloudTkeClusterAttachmentRead,
		Delete: resourceTencentCloudTkeClusterAttachmentDelete,
		Schema: schemaBody,
	}
}

func tkeGetInstanceAdvancedPara(dMap map[string]interface{}, meta interface{}) (setting tke.InstanceAdvancedSettings) {
	setting = tke.InstanceAdvancedSettings{}
	if v, ok := dMap["mount_target"]; ok {
		setting.MountTarget = helper.String(v.(string))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		setting.DataDisks = make([]*tke.DataDisk, len(dataDisks))
		for i, d := range dataDisks {
			value := d.(map[string]interface{})
			var diskType, fileSystem, mountTarget, diskPartition string
			if v, ok := value["disk_type"].(string); ok {
				diskType = v
			}
			if v, ok := value["file_system"].(string); ok {
				fileSystem = v
			}
			if v, ok := value["mount_target"].(string); ok {
				mountTarget = v
			}
			if v, ok := value["disk_partition"].(string); ok {
				diskPartition = v
			}

			diskSize := int64(value["disk_size"].(int))
			autoFormatAndMount := value["auto_format_and_mount"].(bool)
			dataDisk := &tke.DataDisk{
				DiskType:           &diskType,
				FileSystem:         &fileSystem,
				AutoFormatAndMount: &autoFormatAndMount,
				MountTarget:        &mountTarget,
				DiskPartition:      &diskPartition,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			setting.DataDisks[i] = dataDisk
		}
	}
	if v, ok := dMap["is_schedule"]; ok {
		setting.Unschedulable = helper.BoolToInt64Ptr(!v.(bool))
	}

	if v, ok := dMap["user_data"]; ok {
		setting.UserScript = helper.String(v.(string))
	}

	if v, ok := dMap["docker_graph_path"]; ok {
		setting.DockerGraphPath = helper.String(v.(string))
	}

	if v, ok := dMap["desired_pod_num"]; ok {
		setting.DesiredPodNumber = helper.Int64(int64(v.(int)))
	}

	if temp, ok := dMap["extra_args"]; ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		clusterExtraArgs := tke.InstanceExtraArgs{}
		clusterExtraArgs.Kubelet = make([]*string, 0)
		for i := range extraArgs {
			clusterExtraArgs.Kubelet = append(clusterExtraArgs.Kubelet, &extraArgs[i])
		}
		setting.ExtraArgs = &clusterExtraArgs
	}

	// get gpu_args
	if v, ok := dMap["gpu_args"]; ok && len(v.([]interface{})) > 0 {
		gpuArgs := v.([]interface{})[0].(map[string]interface{})

		var (
			migEnable    = gpuArgs["mig_enable"].(bool)
			driver       = gpuArgs["driver"].(map[string]interface{})
			cuda         = gpuArgs["cuda"].(map[string]interface{})
			cudnn        = gpuArgs["cudnn"].(map[string]interface{})
			customDriver = gpuArgs["custom_driver"].(map[string]interface{})
		)
		tkeGpuArgs := tke.GPUArgs{}
		tkeGpuArgs.MIGEnable = &migEnable
		if len(driver) > 0 {
			tkeGpuArgs.Driver = &tke.DriverVersion{
				Version: helper.String(driver["version"].(string)),
				Name:    helper.String(driver["name"].(string)),
			}
		}
		if len(cuda) > 0 {
			tkeGpuArgs.CUDA = &tke.DriverVersion{
				Version: helper.String(cuda["version"].(string)),
				Name:    helper.String(cuda["name"].(string)),
			}
		}
		if len(cudnn) > 0 {
			tkeGpuArgs.CUDNN = &tke.CUDNN{
				Version: helper.String(cudnn["version"].(string)),
				Name:    helper.String(cudnn["name"].(string)),
			}
			if cudnn["doc_name"] != nil {
				tkeGpuArgs.CUDNN.DocName = helper.String(cudnn["doc_name"].(string))
			}
			if cudnn["dev_name"] != nil {
				tkeGpuArgs.CUDNN.DevName = helper.String(cudnn["dev_name"].(string))
			}
		}
		if len(customDriver) > 0 {
			tkeGpuArgs.CustomDriver = &tke.CustomDriver{
				Address: helper.String(customDriver["address"].(string)),
			}
		}
		setting.GPUArgs = &tkeGpuArgs
	}

	return setting
}
func resourceTencentCloudTkeClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	/*tke has been deleted*/
	_, has, err := tkeService.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = tkeService.DescribeCluster(ctx, clusterId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return nil
	}
	if !has {
		d.SetId("")
		return nil
	}

	/*cvm has been deleted*/
	var instance *cvm.Instance
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, err = cvmService.DescribeInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	instanceState := ""
	has = false
	/*attachment has been  deleted*/

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, workers, err := tkeService.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				has = true
				instanceState = worker.InstanceState
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						worker.InstanceId, clusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						worker.InstanceId, worker.InstanceState))
				}
				_ = d.Set("unschedulable", worker.InstanceAdvancedSettings.Unschedulable)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	}

	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	_ = d.Set("state", instanceState)
	return nil
}

func resourceTencentCloudTkeClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	request := tke.NewAddExistedInstancesRequest()

	instanceId := helper.String(d.Get("instance_id").(string))
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceIds = []*string{instanceId}
	request.LoginSettings = &tke.LoginSettings{}

	var loginSettingsNumbers = 0

	if v, ok := d.GetOk("key_ids"); ok {
		request.LoginSettings.KeyIds = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		loginSettingsNumbers++
	}

	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
		loginSettingsNumbers++
	}

	if loginSettingsNumbers != 1 {
		return fmt.Errorf("parameters `key_ids` and `password` must set and only set one")
	}

	request.InstanceAdvancedSettings = &tke.InstanceAdvancedSettings{}
	if workConfig, ok := d.GetOk("worker_config"); ok {
		workConfigList := workConfig.([]interface{})
		if len(workConfigList) == 1 {
			workConfigPara := workConfigList[0].(map[string]interface{})
			setting := tkeGetInstanceAdvancedPara(workConfigPara, meta)
			request.InstanceAdvancedSettings = &setting
		}
	}

	request.InstanceAdvancedSettings.Labels = GetTkeLabels(d, "labels")
	if hostName, ok := d.GetOk("hostname"); ok {
		hostNameStr := hostName.(string)
		request.HostName = &hostNameStr
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		request.InstanceAdvancedSettings.Unschedulable = helper.Int64(v.(int64))
	}

	if workConfigOverrides, ok := d.GetOk("worker_config_overrides"); ok {
		workConfigOverrideList := workConfigOverrides.([]interface{})
		request.InstanceAdvancedSettingsOverrides = make([]*tke.InstanceAdvancedSettings, 0, len(workConfigOverrideList))
		for _, conf := range workConfigOverrideList {
			workConfigPara := conf.(map[string]interface{})
			setting := tkeGetInstanceAdvancedPara(workConfigPara, meta)
			request.InstanceAdvancedSettingsOverrides = append(request.InstanceAdvancedSettingsOverrides, &setting)
		}
	}

	/*cvm has been  attached*/
	var err error
	_, workers, err := tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has := false
	for _, worker := range workers {
		if worker.InstanceId == *instanceId {
			has = true
		}
	}
	if has {
		return fmt.Errorf("instance %s has been attached to cluster %s,can not attach again", *instanceId, *request.ClusterId)
	}

	var response *tke.AddExistedInstancesResponse

	if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = tkeService.client.UseTkeClient().AddExistedInstances(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("add existed instance %s to cluster %s error,reason %v", *instanceId, *request.ClusterId, err)
	}
	var success = false
	for _, v := range response.Response.SuccInstanceIds {
		if *v == *instanceId {
			d.SetId(*instanceId + "_" + *request.ClusterId)
			success = true
		}
	}

	if !success {
		return fmt.Errorf("add existed instance %s to cluster %s error, instance not in success instanceIds", *instanceId, *request.ClusterId)
	}

	/*wait for cvm status*/
	if err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, *instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance != nil && *instance.InstanceState == svccvm.CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance %s status is %s, retry...", *instanceId, *instance.InstanceState))
	}); err != nil {
		return err
	}

	/*wait for tke init ok */
	err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		has := false
		for _, worker := range workers {
			if worker.InstanceId == *instanceId {
				has = true
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						*instanceId, *request.ClusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						*instanceId, worker.InstanceState))
				}

			}
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s not exist in tke instance list", *instanceId))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTkeClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	request := tke.NewDeleteClusterInstancesRequest()

	request.ClusterId = &clusterId
	request.InstanceIds = []*string{
		&instanceId,
	}
	request.InstanceDeleteMode = helper.String("retain")

	var err error

	if err = resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := tkeService.client.UseTkeClient().DeleteClusterInstances(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
			if e.GetCode() == "InternalError.Param" &&
				strings.Contains(e.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
				return nil
			}
		}

		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
