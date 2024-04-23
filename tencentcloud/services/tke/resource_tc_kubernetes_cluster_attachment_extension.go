package tke

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
