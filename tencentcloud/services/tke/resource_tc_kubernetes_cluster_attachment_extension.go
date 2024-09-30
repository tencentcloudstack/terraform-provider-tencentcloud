package tke

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func resourceTencentCloudKubernetesClusterAttachmentCreatePostFillRequest0(ctx context.Context, req *tke.AddExistedInstancesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	// key_ids
	var loginSettingsNumbers = 0
	if v, ok := d.GetOk("key_ids"); ok {
		req.LoginSettings.KeyIds = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		loginSettingsNumbers++
	}
	// password
	if req.LoginSettings.Password != nil {
		loginSettingsNumbers++
	}
	if loginSettingsNumbers != 1 {
		return fmt.Errorf("parameters `key_ids` and `password` must set and only set one")
	}

	if req.InstanceAdvancedSettings == nil {
		req.InstanceAdvancedSettings = &tke.InstanceAdvancedSettings{}
	}
	// labels
	req.InstanceAdvancedSettings.Labels = GetTkeLabels(d, "labels")

	if dMap, ok := helper.InterfacesHeadMap(d, "worker_config"); ok {
		completeInstanceAdvancedSettings(dMap, req.InstanceAdvancedSettings)
	}
	if v, ok := d.GetOk("worker_config_overrides"); ok {
		for i, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			completeInstanceAdvancedSettings(dMap, req.InstanceAdvancedSettingsOverrides[i])
		}
	}

	// only this unschedulable is valid, the is_schedule of worker_config and worker_config_overrides was deprecated.
	if v, ok := d.GetOkExists("unschedulable"); ok {
		req.InstanceAdvancedSettings.Unschedulable = helper.IntInt64(v.(int))
	}

	// 检查是否已经绑定
	if hasAttached, err := nodeHasAttachedToCluster(ctx, instanceId, *req.ClusterId); err != nil {
		return err
	} else if hasAttached {
		return fmt.Errorf("instance %s has been attached to cluster %s,can not attach again", instanceId, *req.ClusterId)
	}

	return nil
}

func resourceTencentCloudKubernetesClusterAttachmentCreateRequestOnError0(ctx context.Context, req *tke.AddExistedInstancesRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0(ctx context.Context, resp *tke.AddExistedInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return fmt.Errorf("provider meta can not be nil")
	}
	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	// 检查响应
	var success = false
	for _, v := range resp.Response.SuccInstanceIds {
		if *v == instanceId {
			success = true
		}
	}
	if !success {
		return fmt.Errorf("add existed instance %s to cluster %s error, instance not in success instanceIds", instanceId, clusterId)
	}

	/*wait for cvm status*/
	if err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance != nil && *instance.InstanceState == svccvm.CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance %s status is %s, retry...", instanceId, *instance.InstanceState))
	}); err != nil {
		return err
	}

	/*wait for tke init ok */
	return resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, workers, err := tkeService.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		has := false
		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				has = true
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						instanceId, clusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						instanceId, worker.InstanceState))
				}

			}
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s not exist in tke instance list", instanceId))
		}
		return nil
	})
}

func resourceTencentCloudKubernetesClusterAttachmentDeleteRequestOnError0(ctx context.Context, e error) *resource.RetryError {
	if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
		if sdkErr.GetCode() == "InternalError.ClusterNotFound" {
			return nil
		}
		if sdkErr.GetCode() == "InternalError.Param" &&
			strings.Contains(sdkErr.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
			return nil
		}
	}
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudKubernetesClusterAttachmentReadRequestOnError2(ctx context.Context, resp *tke.Instance, e error) *resource.RetryError {
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudKubernetesClusterAttachmentReadRequestOnSuccess2(ctx context.Context, resp *tke.Instance) *resource.RetryError {
	// d := tccommon.ResourceDataFromContext(ctx)
	if resp == nil {
		return nil
	}
	var insID string
	if resp.InstanceId != nil {
		insID = *resp.InstanceId
	}
	var insState string
	if resp.InstanceState != nil {
		insState = *resp.InstanceState
	}

	if insState == "failed" {
		return resource.NonRetryableError(fmt.Errorf(
			"cvm instance %s attach to cluster fail, reason: %s",
			insID, insState,
		))
	}

	if insState != "running" {
		return resource.RetryableError(fmt.Errorf(
			"cvm instance %s in tke status is %s, retry...",
			insID, insState,
		))
	}

	// api cannot return Taints
	// if resp.InstanceAdvancedSettings.Taints != nil {
	// 	iAdvanced := resp.InstanceAdvancedSettings
	// 	taintsList := make([]map[string]interface{}, 0, len(iAdvanced.Taints))
	// 	if iAdvanced.Taints != nil {
	// 		for _, taints := range iAdvanced.Taints {
	// 			taintsMap := map[string]interface{}{}

	// 			if taints.Key != nil {
	// 				taintsMap["key"] = taints.Key
	// 			}

	// 			if taints.Value != nil {
	// 				taintsMap["value"] = taints.Value
	// 			}

	// 			if taints.Effect != nil {
	// 				taintsMap["effect"] = taints.Effect
	// 			}

	// 			taintsList = append(taintsList, taintsMap)
	// 		}

	// 		_ = d.Set("taints", taintsList)
	// 	}
	// }

	return nil
}

// nodeHasAttachedToCluster 判断节点是否已经绑定集群
func nodeHasAttachedToCluster(ctx context.Context, insID, clsID string) (bool, error) {
	meta := tccommon.ProviderMetaFromContext(ctx)
	if meta == nil {
		return false, fmt.Errorf("provider meta can not be nil")
	}
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var err error
	_, workers, err := service.DescribeClusterInstances(ctx, clsID)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clsID)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return false, err
	}

	for _, worker := range workers {
		if worker.InstanceId == insID {
			return true, nil
		}
	}
	return false, nil
}

func completeInstanceAdvancedSettings(dMap map[string]interface{}, setting *tke.InstanceAdvancedSettings) {
	// 去除不合法的磁盘大小
	for i, disk := range setting.DataDisks {
		if disk.DiskSize != nil && *disk.DiskSize <= 0 {
			setting.DataDisks[i].DiskSize = nil
		}
	}

	// is_schedule
	if v, ok := dMap["is_schedule"]; ok {
		setting.Unschedulable = helper.BoolToInt64Ptr(!v.(bool))
	}

	// extra_args
	if temp, ok := dMap["extra_args"]; ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		clusterExtraArgs := tke.InstanceExtraArgs{}
		clusterExtraArgs.Kubelet = make([]*string, 0)
		for i := range extraArgs {
			clusterExtraArgs.Kubelet = append(clusterExtraArgs.Kubelet, &extraArgs[i])
		}
		setting.ExtraArgs = &clusterExtraArgs
	}

	if v, ok := dMap["gpu_args"]; ok && len(v.([]interface{})) > 0 {
		gpuArgs := v.([]interface{})[0].(map[string]interface{})

		driver := gpuArgs["driver"].(map[string]interface{})
		if len(driver) > 0 {
			setting.GPUArgs.Driver = &tke.DriverVersion{
				Version: helper.String(driver["version"].(string)),
				Name:    helper.String(driver["name"].(string)),
			}
		}

		cuda := gpuArgs["cuda"].(map[string]interface{})
		if len(cuda) > 0 {
			setting.GPUArgs.CUDA = &tke.DriverVersion{
				Version: helper.String(cuda["version"].(string)),
				Name:    helper.String(cuda["name"].(string)),
			}
		}

		cudnn := gpuArgs["cudnn"].(map[string]interface{})
		if len(cudnn) > 0 {
			setting.GPUArgs.CUDNN = &tke.CUDNN{
				Version: helper.String(cudnn["version"].(string)),
				Name:    helper.String(cudnn["name"].(string)),
			}
			if cudnn["doc_name"] != nil {
				setting.GPUArgs.CUDNN.DocName = helper.String(cudnn["doc_name"].(string))
			}
			if cudnn["dev_name"] != nil {
				setting.GPUArgs.CUDNN.DevName = helper.String(cudnn["dev_name"].(string))
			}
		}

		customDriver := gpuArgs["custom_driver"].(map[string]interface{})
		if len(customDriver) > 0 {
			setting.GPUArgs.CustomDriver = &tke.CustomDriver{
				Address: helper.String(customDriver["address"].(string)),
			}
		}
	}
}

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
						Description:  "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
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
		"pre_start_user_script": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
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

	if v, ok := dMap["pre_start_user_script"]; ok {
		setting.PreStartUserScript = helper.String(v.(string))
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
