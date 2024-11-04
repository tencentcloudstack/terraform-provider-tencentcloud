package tke

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func resourceTencentCloudKubernetesClusterMasterAttachmentCreatePostFillRequest0(ctx context.Context, req *tkev20180525.ScaleOutClusterMasterRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	existedInstancesForNode := tkev20180525.ExistedInstancesForNode{}
	existedInstancesPara := tkev20180525.ExistedInstancesPara{}
	enhancedService := tkev20180525.EnhancedService{}
	loginSettings := tkev20180525.LoginSettings{}
	if v, ok := d.GetOk("instance_id"); ok {
		existedInstancesPara.InstanceIds = helper.Strings([]string{v.(string)})
	}

	if v, ok := d.GetOk("node_role"); ok {
		existedInstancesForNode.NodeRole = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enhanced_security_service"); ok {
		enhancedService.SecurityService = &tkev20180525.RunSecurityServiceEnabled{Enabled: helper.Bool(v.(bool))}
		existedInstancesPara.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOkExists("enhanced_monitor_service"); ok {
		enhancedService.MonitorService = &tkev20180525.RunMonitorServiceEnabled{Enabled: helper.Bool(v.(bool))}
		existedInstancesPara.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOkExists("enhanced_automation_service"); ok {
		enhancedService.AutomationService = &tkev20180525.RunAutomationServiceEnabled{Enabled: helper.Bool(v.(bool))}
		existedInstancesPara.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOk("password"); ok {
		loginSettings.Password = helper.String(v.(string))
		existedInstancesPara.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("key_ids"); ok && len(v.([]interface{})) > 0 {
		keyIds := v.([]interface{})
		loginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			loginSettings.KeyIds = append(loginSettings.KeyIds, &keyId)
		}

		existedInstancesPara.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("security_group_ids"); ok && len(v.([]interface{})) > 0 {
		sgIds := v.([]interface{})
		existedInstancesPara.SecurityGroupIds = make([]*string, 0, len(sgIds))
		for i := range sgIds {
			sgId := sgIds[i].(string)
			existedInstancesPara.SecurityGroupIds = append(existedInstancesPara.SecurityGroupIds, &sgId)
		}
	}

	if v, ok := d.GetOk("host_name"); ok {
		existedInstancesPara.HostName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("desired_pod_numbers"); ok && len(v.([]interface{})) > 0 {
		desiredPodNumbers := v.([]interface{})
		existedInstancesForNode.DesiredPodNumbers = make([]*int64, 0, len(desiredPodNumbers))
		for i := range desiredPodNumbers {
			desiredPodNumber := desiredPodNumbers[i].(int64)
			existedInstancesForNode.DesiredPodNumbers = append(existedInstancesForNode.DesiredPodNumbers, &desiredPodNumber)
		}
	}

	if v, ok := d.GetOk("master_config"); ok && len(v.([]interface{})) > 0 {
		for _, item := range v.([]interface{}) {
			instanceAdvancedSettingsOverridesMap := item.(map[string]interface{})
			instanceAdvancedSettings := tkev20180525.InstanceAdvancedSettings{}
			if v, ok := instanceAdvancedSettingsOverridesMap["mount_target"]; ok {
				instanceAdvancedSettings.MountTarget = helper.String(v.(string))
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["docker_graph_path"]; ok {
				instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["user_script"]; ok {
				instanceAdvancedSettings.UserScript = helper.String(v.(string))
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["unschedulable"]; ok {
				instanceAdvancedSettings.Unschedulable = helper.IntInt64(v.(int))
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["labels"]; ok && len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					labelsMap := item.(map[string]interface{})
					labels := tkev20180525.Label{}
					if v, ok := labelsMap["name"]; ok {
						labels.Name = helper.String(v.(string))
					}

					if v, ok := labelsMap["value"]; ok {
						labels.Value = helper.String(v.(string))
					}

					instanceAdvancedSettings.Labels = append(instanceAdvancedSettings.Labels, &labels)
				}
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["data_disk"]; ok && len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					dataDisksMap := item.(map[string]interface{})
					dataDisk := tkev20180525.DataDisk{}
					if v, ok := dataDisksMap["disk_type"]; ok {
						dataDisk.DiskType = helper.String(v.(string))
					}

					if v, ok := dataDisksMap["file_system"]; ok {
						dataDisk.FileSystem = helper.String(v.(string))
					}

					if v, ok := dataDisksMap["disk_size"]; ok {
						dataDisk.DiskSize = helper.IntInt64(v.(int))
					}

					if v, ok := dataDisksMap["auto_format_and_mount"]; ok {
						dataDisk.AutoFormatAndMount = helper.Bool(v.(bool))
					}

					if v, ok := dataDisksMap["mount_target"]; ok {
						dataDisk.MountTarget = helper.String(v.(string))
					}

					if v, ok := dataDisksMap["disk_partition"]; ok {
						dataDisk.DiskPartition = helper.String(v.(string))
					}

					instanceAdvancedSettings.DataDisks = append(instanceAdvancedSettings.DataDisks, &dataDisk)
				}
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["extra_args"]; ok && len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					extraArgsMap := item.(map[string]interface{})
					args := tkev20180525.InstanceExtraArgs{}
					if v, ok := extraArgsMap["kubelet"]; ok {
						args.Kubelet = helper.InterfacesStringsPoint(v.([]interface{}))
					}

					instanceAdvancedSettings.ExtraArgs = &args
				}
			}

			if v, ok := d.GetOkExists("desired_pod_number"); ok {
				instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["gpu_args"]; ok && len(v.([]interface{})) > 0 {
				gpuArgs := v.([]interface{})[0].(map[string]interface{})

				var (
					migEnable    = gpuArgs["mig_enable"].(bool)
					driver       = gpuArgs["driver"].(map[string]interface{})
					cuda         = gpuArgs["cuda"].(map[string]interface{})
					cudnn        = gpuArgs["cudnn"].(map[string]interface{})
					customDriver = gpuArgs["custom_driver"].(map[string]interface{})
				)

				tkeGpuArgs := tkev20180525.GPUArgs{}
				tkeGpuArgs.MIGEnable = &migEnable
				if len(driver) > 0 {
					tkeGpuArgs.Driver = &tkev20180525.DriverVersion{
						Version: helper.String(driver["version"].(string)),
						Name:    helper.String(driver["name"].(string)),
					}
				}

				if len(cuda) > 0 {
					tkeGpuArgs.CUDA = &tkev20180525.DriverVersion{
						Version: helper.String(cuda["version"].(string)),
						Name:    helper.String(cuda["name"].(string)),
					}
				}

				if len(cudnn) > 0 {
					tkeGpuArgs.CUDNN = &tkev20180525.CUDNN{
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
					tkeGpuArgs.CustomDriver = &tkev20180525.CustomDriver{
						Address: helper.String(customDriver["address"].(string)),
					}
				}

				instanceAdvancedSettings.GPUArgs = &tkeGpuArgs
			}

			if v, ok := instanceAdvancedSettingsOverridesMap["taints"]; ok && len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					taintsMap := item.(map[string]interface{})
					taint := tkev20180525.Taint{}
					if v, ok := taintsMap["key"]; ok {
						taint.Key = helper.String(v.(string))
					}

					if v, ok := taintsMap["value"]; ok {
						taint.Value = helper.String(v.(string))
					}

					if v, ok := taintsMap["effect"]; ok {
						taint.Effect = helper.String(v.(string))
					}

					instanceAdvancedSettings.Taints = append(instanceAdvancedSettings.Taints, &taint)
				}
			}

			existedInstancesForNode.InstanceAdvancedSettingsOverride = &instanceAdvancedSettings
		}

		existedInstancesForNode.ExistedInstancesPara = &existedInstancesPara
		req.ExistedInstancesForNode = []*tkev20180525.ExistedInstancesForNode{&existedInstancesForNode}
	}

	return nil
}

func resourceTencentCloudKubernetesClusterMasterAttachmentCreateRequestOnError0(ctx context.Context, req *tkev20180525.ScaleOutClusterMasterRequest, e error) *resource.RetryError {
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudKubernetesClusterMasterAttachmentCreatePostHandleResponse0(ctx context.Context, resp *tkev20180525.ScaleOutClusterMasterResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var (
		meta       = tccommon.ProviderMetaFromContext(ctx)
		tkeService = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		cvmService = svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		clusterId  string
		instanceId string
		nodeRole   string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("node_role"); ok {
		nodeRole = v.(string)
	}

	// wait for cvm status
	if err := resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

	// wait for tke init
	return resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resp, err := tkeService.DescribeKubernetesClusterMasterAttachmentById2(ctx, clusterId, instanceId, nodeRole)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}

		if len(resp.InstanceSet) != 1 {
			return resource.NonRetryableError(fmt.Errorf("tke master node cvm instance %s not exist in tke instance list", instanceId))
		}

		if *resp.InstanceSet[0].InstanceState != "running" {
			return resource.RetryableError(fmt.Errorf("tke master node cvm instance  %s in tke status is %s, retry...", instanceId, *resp.InstanceSet[0].InstanceState))
		}

		return nil
	})
}

func resourceTencentCloudKubernetesClusterMasterAttachmentReadRequestOnError2(ctx context.Context, resp *tkev20180525.DescribeClusterInstancesResponseParams, e error) *resource.RetryError {
	return tccommon.RetryError(e, tccommon.InternalError)
}

func resourceTencentCloudKubernetesClusterMasterAttachmentReadRequestOnSuccess2(ctx context.Context, resp *tkev20180525.DescribeClusterInstancesResponseParams) *resource.RetryError {
	if resp == nil || len(resp.InstanceSet) != 1 {
		return resource.NonRetryableError(fmt.Errorf("query cvm instance error."))
	}

	instanceDetial := resp.InstanceSet[0]
	insId := *instanceDetial.InstanceId
	insState := *instanceDetial.InstanceState

	if insState == "failed" {
		return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster fail, reason: %s", insId, insState))
	}

	if insState != "running" {
		return resource.RetryableError(fmt.Errorf("cvm instance %s in tke status is %s, retry...", insId, insState))
	}

	return nil
}

func resourceTencentCloudKubernetesClusterMasterAttachmentDeletePostFillRequest0(ctx context.Context, req *tkev20180525.ScaleInClusterMasterRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	var (
		instanceId string
		nodeRole   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("node_role"); ok {
		nodeRole = v.(string)
	}

	req.ScaleInMasters = []*tkev20180525.ScaleInMaster{
		{
			InstanceId:         helper.String(instanceId),
			NodeRole:           helper.String(nodeRole),
			InstanceDeleteMode: helper.String("retain"),
		},
	}

	return nil
}

func resourceTencentCloudKubernetesClusterMasterAttachmentDeleteRequestOnError0(ctx context.Context, e error) *resource.RetryError {
	if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
		if sdkErr.GetCode() == "ResourceNotFound" {
			return nil
		}

		if sdkErr.GetCode() == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), `is not exist`) {
			return nil
		}
	}

	return tccommon.RetryError(e, tccommon.InternalError)
}
