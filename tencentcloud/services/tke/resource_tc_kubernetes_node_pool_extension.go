package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"

	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

var importFlag = false

func nodePoolCustomResourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importFlag = true
	err := resourceTencentCloudKubernetesNodePoolRead(d, m)
	if err != nil {
		return nil, fmt.Errorf("failed to import resource")
	}
	return []*schema.ResourceData{d}, nil
}

func nodeOsTypeDiffSuppressFunc(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("node_os"); ok {
		if strings.Contains(v.(string), "img-") {
			return true
		}
	}
	return false
}

func resourceTencentCloudKubernetesNodePoolCreatePostFillRequest0(ctx context.Context, req *tke.CreateClusterNodePoolRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		configParas = d.Get("auto_scaling_config").([]interface{})
		iAdvanced   tke.InstanceAdvancedSettings
	)
	if len(configParas) != 1 {
		return fmt.Errorf("need only one auto_scaling_config")
	}

	groupParaStr, err := composeParameterToAsScalingGroupParaSerial(d)
	if err != nil {
		return err
	}
	req.AutoScalingGroupPara = &groupParaStr

	configParaStr, err := composedKubernetesAsScalingConfigParaSerial(configParas[0].(map[string]interface{}), meta)
	if err != nil {
		return err
	}
	req.LaunchConfigurePara = &configParaStr

	labels := GetTkeLabels(d, "labels")
	tags := GetTkeTags(d, "tags")
	if len(labels) > 0 {
		req.Labels = labels
	}
	if len(tags) > 0 {
		req.Tags = tags
	}

	//compose InstanceAdvancedSettings
	if workConfig, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		iAdvanced = tkeGetInstanceAdvancedPara(workConfig, meta)
		req.InstanceAdvancedSettings = &iAdvanced
	}

	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for _, extraArg := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArg)
		}
	}
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.Int64(int64(temp.(int)))
	}
	req.InstanceAdvancedSettings = &iAdvanced

	nodeOs := d.Get("node_os").(string)
	nodeOsType := d.Get("node_os_type").(string)
	//自定镜像不能指定节点操作系统类型
	if strings.Contains(nodeOs, "img-") {
		nodeOsType = ""
	}
	req.NodePoolOs = &nodeOs
	req.OsCustomizeType = &nodeOsType

	return nil
}

func resourceTencentCloudKubernetesNodePoolCreatePostHandleResponse0(ctx context.Context, resp *tke.CreateClusterNodePoolResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	nodePoolId := *resp.Response.NodePoolId

	// todo wait for status ok
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		nodePool, _, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if nodePool != nil && *nodePool.LifeState == "normal" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("node pool status is %s, retry...", *nodePool.LifeState))
	})
	if err != nil {
		return err
	}

	instanceTypes := getNodePoolInstanceTypes(d)

	if len(instanceTypes) != 0 {
		err := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
		if err != nil {
			return err
		}
	}

	//modify os, instanceTypes and image
	//err = resourceTencentCloudKubernetesNodePoolUpdate(d, meta)
	//if err != nil {
	//	return err
	//}
	d.SetId(strings.Join([]string{clusterId, nodePoolId}, tccommon.FILED_SP))
	if err := resourceTencentCloudKubernetesNodePoolUpdateOnStart(ctx); err != nil {
		return err
	}
	if err := resourceTencentCloudKubernetesNodePoolUpdateOnExit(ctx); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolReadRequestOnError1(ctx context.Context, resp *tke.NodePool, e error) *resource.RetryError {
	if e != nil {
		return resource.NonRetryableError(e)
	}
	return nil
}

func resourceTencentCloudKubernetesNodePoolReadRequestOnSuccess1(ctx context.Context, resp *tke.NodePool) *resource.RetryError {
	status := *resp.AutoscalingGroupStatus
	if status == "enabling" || status == "disabling" {
		return resource.RetryableError(fmt.Errorf("node pool status is %s, retrying", status))
	}
	return nil
}

func resourceTencentCloudKubernetesNodePoolReadPostHandleResponse1(ctx context.Context, resp *tke.NodePool) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		asService = svcas.NewAsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)
	nodePool := resp

	AutoscalingAddedTotal := *nodePool.NodeCountSummary.AutoscalingAdded.Total
	ManuallyAddedTotal := *nodePool.NodeCountSummary.ManuallyAdded.Total
	_ = d.Set("node_count", AutoscalingAddedTotal+ManuallyAddedTotal)
	if _, ok := d.GetOkExists("unschedulable"); !ok && importFlag {
		_ = d.Set("unschedulable", nodePool.Unschedulable)
	}
	if nodePool.AutoscalingGroupStatus != nil {
		_ = d.Set("enable_auto_scale", *nodePool.AutoscalingGroupStatus == "enabled")
	}
	//修复自定义镜像返回信息的不一致
	if nodePool.ImageId != nil && *nodePool.ImageId != "" {
		_ = d.Set("node_os", nodePool.ImageId)
	} else {
		if nodePool.NodePoolOs != nil {
			_ = d.Set("node_os", nodePool.NodePoolOs)
		}
		if nodePool.OsCustomizeType != nil {
			_ = d.Set("node_os_type", nodePool.OsCustomizeType)
		}
	}

	if tags := nodePool.Tags; tags != nil {
		tagMap := make(map[string]string)
		for i := range tags {
			tag := tags[i]
			tagMap[*tag.Key] = *tag.Value
		}
		_ = d.Set("tags", tagMap)
	}

	//set composed struct
	lables := make(map[string]interface{}, len(nodePool.Labels))
	for _, v := range nodePool.Labels {
		lables[*v.Name] = *v.Value
	}
	_ = d.Set("labels", lables)

	// set launch config
	launchCfg, hasLC, err := asService.DescribeLaunchConfigurationById(ctx, *nodePool.LaunchConfigurationId)

	if hasLC > 0 {
		launchConfig := make(map[string]interface{})
		if launchCfg.InstanceTypes != nil {
			insTypes := launchCfg.InstanceTypes
			launchConfig["instance_type"] = insTypes[0]
			backupInsTypes := insTypes[1:]
			if len(backupInsTypes) > 0 {
				launchConfig["backup_instance_types"] = helper.StringsInterfaces(backupInsTypes)
			}
		} else {
			launchConfig["instance_type"] = launchCfg.InstanceType
		}
		if launchCfg.SystemDisk.DiskType != nil {
			launchConfig["system_disk_type"] = launchCfg.SystemDisk.DiskType
		}
		if launchCfg.SystemDisk.DiskSize != nil {
			launchConfig["system_disk_size"] = launchCfg.SystemDisk.DiskSize
		}
		if launchCfg.InternetAccessible.InternetChargeType != nil {
			launchConfig["internet_charge_type"] = launchCfg.InternetAccessible.InternetChargeType
		}
		if launchCfg.InternetAccessible.InternetMaxBandwidthOut != nil {
			launchConfig["internet_max_bandwidth_out"] = launchCfg.InternetAccessible.InternetMaxBandwidthOut
		}
		if launchCfg.InternetAccessible.BandwidthPackageId != nil {
			launchConfig["bandwidth_package_id"] = launchCfg.InternetAccessible.BandwidthPackageId
		}
		if launchCfg.InternetAccessible.PublicIpAssigned != nil {
			launchConfig["public_ip_assigned"] = launchCfg.InternetAccessible.PublicIpAssigned
		}
		if launchCfg.InstanceChargeType != nil {
			launchConfig["instance_charge_type"] = launchCfg.InstanceChargeType
			if *launchCfg.InstanceChargeType == svcas.INSTANCE_CHARGE_TYPE_SPOTPAID && launchCfg.InstanceMarketOptions != nil {
				launchConfig["spot_instance_type"] = launchCfg.InstanceMarketOptions.SpotOptions.SpotInstanceType
				launchConfig["spot_max_price"] = launchCfg.InstanceMarketOptions.SpotOptions.MaxPrice
			}
			if *launchCfg.InstanceChargeType == svcas.INSTANCE_CHARGE_TYPE_PREPAID && launchCfg.InstanceChargePrepaid != nil {
				launchConfig["instance_charge_type_prepaid_period"] = launchCfg.InstanceChargePrepaid.Period
				launchConfig["instance_charge_type_prepaid_renew_flag"] = launchCfg.InstanceChargePrepaid.RenewFlag
			}
		}
		if len(launchCfg.DataDisks) > 0 {
			dataDisks := make([]map[string]interface{}, 0, len(launchCfg.DataDisks))
			for i := range launchCfg.DataDisks {
				item := launchCfg.DataDisks[i]
				disk := make(map[string]interface{})
				disk["disk_type"] = *item.DiskType
				disk["disk_size"] = *item.DiskSize
				if item.SnapshotId != nil {
					disk["snapshot_id"] = *item.SnapshotId
				}
				if item.DeleteWithInstance != nil {
					disk["delete_with_instance"] = *item.DeleteWithInstance
				}
				if item.Encrypt != nil {
					disk["encrypt"] = *item.Encrypt
				}
				if item.ThroughputPerformance != nil {
					disk["throughput_performance"] = *item.ThroughputPerformance
				}
				dataDisks = append(dataDisks, disk)
			}
			launchConfig["data_disk"] = dataDisks
		}
		if launchCfg.LoginSettings != nil {
			launchConfig["key_ids"] = helper.StringsInterfaces(launchCfg.LoginSettings.KeyIds)
		}
		// keep existing password in new launchConfig object
		if v, ok := d.GetOk("auto_scaling_config.0.password"); ok {
			launchConfig["password"] = v.(string)
		}

		if launchCfg.SecurityGroupIds != nil {
			launchConfig["security_group_ids"] = helper.StringsInterfaces(launchCfg.SecurityGroupIds)
			launchConfig["orderly_security_group_ids"] = helper.StringsInterfaces(launchCfg.SecurityGroupIds)
		}

		enableSecurity := launchCfg.EnhancedService.SecurityService.Enabled
		enableMonitor := launchCfg.EnhancedService.MonitorService.Enabled
		// Only declared or diff from exist will set.
		if _, ok := d.GetOk("enhanced_security_service"); ok || enableSecurity != nil {
			launchConfig["enhanced_security_service"] = *enableSecurity
		}
		if _, ok := d.GetOk("enhanced_monitor_service"); ok || enableMonitor != nil {
			launchConfig["enhanced_monitor_service"] = *enableMonitor
		}
		if _, ok := d.GetOk("cam_role_name"); ok || launchCfg.CamRoleName != nil {
			launchConfig["cam_role_name"] = launchCfg.CamRoleName
		}
		if launchCfg.InstanceNameSettings != nil {
			if launchCfg.InstanceNameSettings.InstanceName != nil {
				launchConfig["instance_name"] = launchCfg.InstanceNameSettings.InstanceName
			}

			if launchCfg.InstanceNameSettings.InstanceNameStyle != nil {
				launchConfig["instance_name_style"] = launchCfg.InstanceNameSettings.InstanceNameStyle
			}
		}
		if launchCfg.HostNameSettings != nil && launchCfg.HostNameSettings.HostName != nil {
			launchConfig["host_name"] = launchCfg.HostNameSettings.HostName
		}
		if launchCfg.HostNameSettings != nil && launchCfg.HostNameSettings.HostNameStyle != nil {
			launchConfig["host_name_style"] = launchCfg.HostNameSettings.HostNameStyle
		}

		asgConfig := make([]interface{}, 0, 1)
		asgConfig = append(asgConfig, launchConfig)
		if err := d.Set("auto_scaling_config", asgConfig); err != nil {
			return err
		}
	}

	nodeConfig := make(map[string]interface{})
	nodeConfigs := make([]interface{}, 0, 1)

	if nodePool.DataDisks != nil && len(nodePool.DataDisks) > 0 {
		dataDisks := make([]interface{}, 0, len(nodePool.DataDisks))
		for i := range nodePool.DataDisks {
			item := nodePool.DataDisks[i]
			disk := make(map[string]interface{})
			disk["disk_type"] = helper.PString(item.DiskType)
			disk["disk_size"] = helper.PInt64(item.DiskSize)
			disk["file_system"] = helper.PString(item.FileSystem)
			disk["auto_format_and_mount"] = helper.PBool(item.AutoFormatAndMount)
			disk["mount_target"] = helper.PString(item.MountTarget)
			disk["disk_partition"] = helper.PString(item.MountTarget)
			dataDisks = append(dataDisks, disk)
		}
		nodeConfig["data_disk"] = dataDisks
	}

	if helper.PInt64(nodePool.DesiredPodNum) != 0 {
		nodeConfig["desired_pod_num"] = helper.PInt64(nodePool.DesiredPodNum)
	}

	if helper.PInt64(nodePool.Unschedulable) != 0 {
		nodeConfig["is_schedule"] = false
	} else {
		nodeConfig["is_schedule"] = true
	}

	if helper.PString(nodePool.DockerGraphPath) != "" {
		nodeConfig["docker_graph_path"] = helper.PString(nodePool.DockerGraphPath)
	} else {
		nodeConfig["docker_graph_path"] = "/var/lib/docker"
	}

	if helper.PString(nodePool.PreStartUserScript) != "" {
		nodeConfig["pre_start_user_script"] = helper.PString(nodePool.PreStartUserScript)
	}

	if importFlag {
		if nodePool.ExtraArgs != nil && len(nodePool.ExtraArgs.Kubelet) > 0 {
			extraArgs := make([]string, 0)
			for i := range nodePool.ExtraArgs.Kubelet {
				extraArgs = append(extraArgs, helper.PString(nodePool.ExtraArgs.Kubelet[i]))
			}
			nodeConfig["extra_args"] = extraArgs
		}

		if helper.PString(nodePool.UserScript) != "" {
			nodeConfig["user_data"] = helper.PString(nodePool.UserScript)
		}

		if nodePool.GPUArgs != nil {
			setting := nodePool.GPUArgs
			var driverEmptyFlag, cudaEmptyFlag, cudnnEmptyFlag, customDriverEmptyFlag bool
			gpuArgs := map[string]interface{}{
				"mig_enable": helper.PBool(setting.MIGEnable),
			}

			if !isDriverEmpty(setting.Driver) {
				driverEmptyFlag = true
				driver := map[string]interface{}{
					"version": helper.PString(setting.Driver.Version),
					"name":    helper.PString(setting.Driver.Name),
				}
				gpuArgs["driver"] = driver
			}

			if !isCUDAEmpty(setting.CUDA) {
				cudaEmptyFlag = true
				cuda := map[string]interface{}{
					"version": helper.PString(setting.CUDA.Version),
					"name":    helper.PString(setting.CUDA.Name),
				}
				gpuArgs["cuda"] = cuda
			}

			if !isCUDNNEmpty(setting.CUDNN) {
				cudnnEmptyFlag = true
				cudnn := map[string]interface{}{
					"version":  helper.PString(setting.CUDNN.Version),
					"name":     helper.PString(setting.CUDNN.Name),
					"doc_name": helper.PString(setting.CUDNN.DocName),
					"dev_name": helper.PString(setting.CUDNN.DevName),
				}
				gpuArgs["cudnn"] = cudnn
			}

			if !isCustomDriverEmpty(setting.CustomDriver) {
				customDriverEmptyFlag = true
				customDriver := map[string]interface{}{
					"address": helper.PString(setting.CustomDriver.Address),
				}
				gpuArgs["custom_driver"] = customDriver
			}
			if driverEmptyFlag || cudaEmptyFlag || cudnnEmptyFlag || customDriverEmptyFlag {
				nodeConfig["gpu_args"] = []map[string]interface{}{gpuArgs}
			}
		}
		nodeConfigs = append(nodeConfigs, nodeConfig)
		_ = d.Set("node_config", nodeConfigs)
		importFlag = false
	}

	// Relative scaling group status
	asg, hasAsg, err := asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			asg, hasAsg, err = asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if hasAsg > 0 {
		_ = d.Set("scaling_group_name", asg.AutoScalingGroupName)
		_ = d.Set("zones", asg.ZoneSet)
		_ = d.Set("scaling_group_project_id", asg.ProjectId)
		_ = d.Set("default_cooldown", asg.DefaultCooldown)
		_ = d.Set("termination_policies", helper.StringsInterfaces(asg.TerminationPolicySet))
		_ = d.Set("vpc_id", asg.VpcId)
		_ = d.Set("retry_policy", asg.RetryPolicy)
		_ = d.Set("subnet_ids", helper.StringsInterfaces(asg.SubnetIdSet))
		if v, ok := d.GetOk("scaling_mode"); ok {
			if asg.ServiceSettings != nil && asg.ServiceSettings.ScalingMode != nil {
				_ = d.Set("scaling_mode", helper.PString(asg.ServiceSettings.ScalingMode))
			} else {
				_ = d.Set("scaling_mode", v.(string))
			}
		}
		// If not check, the diff between computed and default empty value leads to force replacement
		if _, ok := d.GetOk("multi_zone_subnet_policy"); ok {
			_ = d.Set("multi_zone_subnet_policy", asg.MultiZoneSubnetPolicy)
		}
	}
	if v, ok := d.GetOkExists("delete_keep_instance"); ok {
		_ = d.Set("delete_keep_instance", v.(bool))
	} else {
		_ = d.Set("delete_keep_instance", true)
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolDeletePostFillRequest0(ctx context.Context, req *tke.DeleteClusterNodePoolRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)

	var deletionProtection = d.Get("deletion_protection").(bool)

	if deletionProtection {
		return fmt.Errorf("deletion protection was enabled, please set `deletion_protection` to `false` and apply first")
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolDeleteRequestOnError0(ctx context.Context, e error) *resource.RetryError {
	if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
		if sdkErr.Code == "InternalError.Param" && strings.Contains(sdkErr.Message, "Not Found") {
			return nil
		}
	}
	return tccommon.RetryError(e)
}

func resourceTencentCloudKubernetesNodePoolDeletePostHandleResponse0(ctx context.Context, resp *tke.DeleteClusterNodePoolResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		items   = strings.Split(d.Id(), tccommon.FILED_SP)
	)
	clusterId := items[0]
	nodePoolId := items[1]

	// todo wait for delete ok
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		nodePool, has, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			errCode := errRet.(*sdkErrors.TencentCloudSDKError).Code
			if errCode == "InternalError.UnexpectedInternal" {
				return nil
			}
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("node pool %s still alive, status %s", nodePoolId, *nodePool.LifeState))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolUpdateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		service    = TkeService{client: client}
		cvmService = svccvm.NewCvmService(client)
		items      = strings.Split(d.Id(), tccommon.FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	d.Partial(true)

	nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		return err
	}
	oldDesiredCapacity := *nodePool.DesiredNodesNum
	oldMinSize := *nodePool.MinNodesNum
	oldMaxSize := *nodePool.MaxNodesNum

	desiredCapacity := int64(d.Get("desired_capacity").(int))
	minSize := int64(d.Get("min_size").(int))
	maxSize := int64(d.Get("max_size").(int))
	if desiredCapacity != oldDesiredCapacity && (minSize != oldMinSize || maxSize != oldMaxSize) {
		log.Printf("[CRITAL]%s modification of min_size[%v] or max_size[%v] at the same time as desired_capacity[%v] failed\n", logId, minSize, maxSize, desiredCapacity)
		return fmt.Errorf("`min_size` or `max_size` cannot be modified at the same time as `desired_capacity`, please modify `min_size` or `max_size` first, and then modify `desired_capacity`")
	}

	// LaunchConfig
	if d.HasChange("auto_scaling_config") {
		launchConfigId := *nodePool.LaunchConfigurationId
		//  change as config here
		request, composeError := composeAsLaunchConfigModifyRequest(d, launchConfigId)
		if composeError != nil {
			return composeError
		}
		_, err = client.UseAsClient().ModifyLaunchConfigurationAttributes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return err
		}

		// change existed cvm security service if necessary
		if err := ModifySecurityServiceOfCvmInNodePool(ctx, d, &service, &cvmService, client, clusterId, *nodePool.NodePoolId); err != nil {
			return err
		}

	}

	var capacityHasChanged = false
	// assuming
	// min 1 max 6 desired 2
	// to
	// min 3 max 6 desired 5
	// modify min/max first will cause error, this case must upgrade desired first
	if d.HasChange("desired_capacity") || !desiredCapacityOutRange(d) {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolDesiredCapacity(ctx, clusterId, nodePoolId, desiredCapacity)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		capacityHasChanged = true
	}

	// ModifyClusterNodePool
	if d.HasChanges(
		"labels",
		"tags",
	) {
		var body map[string]interface{}
		nodeOs := d.Get("node_os").(string)
		nodeOsType := d.Get("node_os_type").(string)
		//自定镜像不能指定节点操作系统类型
		if strings.Contains(nodeOs, "img-") {
			nodeOsType = ""
		}

		labels := GetTkeLabels(d, "labels")
		body = map[string]interface{}{
			"ClusterId":       clusterId,
			"NodePoolId":      nodePoolId,
			"OsName":          nodeOs,
			"OsCustomizeType": nodeOsType,
			"Labels":          labels,
		}

		tags := helper.GetTags(d, "tags")
		if len(tags) > 0 {
			var tmpTags []*tke.Tag
			for k, v := range tags {
				key := k
				val := v
				tmpTags = append(tmpTags, &tke.Tag{
					Key:   &key,
					Value: &val,
				})
			}

			body["Tags"] = tmpTags
		}

		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("tke")
		request := tchttp.NewCommonRequest("tke", "2018-05-25", "ModifyClusterNodePool")
		err := request.SetActionParameters(body)
		if err != nil {
			return err
		}

		response := tchttp.NewCommonResponse()
		err = client.Send(request, response)
		if err != nil {
			fmt.Printf("update kubernetes node pool taints failed: %v \n", err)
			return err
		}

		// todo wait for status ok
		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			nodePool, _, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if nodePool != nil && *nodePool.LifeState == "normal" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("node pool status is %s, retry...", *nodePool.LifeState))
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("desired_capacity") && !capacityHasChanged {
		desiredCapacity := int64(d.Get("desired_capacity").(int))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolDesiredCapacity(ctx, clusterId, nodePoolId, desiredCapacity)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		client    = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		service   = TkeService{client: client}
		asService = svcas.NewAsService(client)
		items     = strings.Split(d.Id(), tccommon.FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	err := resourceTencentCloudKubernetesNodePoolUpdateTaints(ctx, clusterId, nodePoolId)
	if err != nil {
		return err
	}

	// ModifyScalingGroup
	if d.HasChange("scaling_group_name") ||
		d.HasChange("zones") ||
		d.HasChange("scaling_group_project_id") ||
		d.HasChange("multi_zone_subnet_policy") ||
		d.HasChange("default_cooldown") ||
		d.HasChange("termination_policies") {

		nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return err
		}

		var (
			request               = as.NewModifyAutoScalingGroupRequest()
			scalingGroupId        = *nodePool.AutoscalingGroupId
			name                  = d.Get("scaling_group_name").(string)
			projectId             = d.Get("scaling_group_project_id").(int)
			defaultCooldown       = d.Get("default_cooldown").(int)
			multiZoneSubnetPolicy = d.Get("multi_zone_subnet_policy").(string)
		)

		request.AutoScalingGroupId = &scalingGroupId

		if name != "" {
			request.AutoScalingGroupName = &name
		}

		if multiZoneSubnetPolicy != "" {
			request.MultiZoneSubnetPolicy = &multiZoneSubnetPolicy
		}

		// It is safe to use Get() with default value 0.
		request.ProjectId = helper.IntUint64(projectId)

		if defaultCooldown != 0 {
			request.DefaultCooldown = helper.IntUint64(defaultCooldown)
		}

		if v, ok := d.GetOk("zones"); ok {
			request.Zones = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := d.GetOk("termination_policies"); ok {
			request.TerminationPolicies = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := asService.ModifyAutoScalingGroup(ctx, request)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	if d.HasChange("auto_scaling_config.0.backup_instance_types") {
		instanceTypes := getNodePoolInstanceTypes(d)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		_ = d.Set("auto_scaling_config.0.backup_instance_types", instanceTypes)
	}
	d.Partial(false)

	return nil
}

// merge `instance_type` to `backup_instance_types` as param `instance_types`
func getNodePoolInstanceTypes(d *schema.ResourceData) []*string {
	configParas := d.Get("auto_scaling_config").([]interface{})
	dMap := configParas[0].(map[string]interface{})
	instanceType := dMap["instance_type"]
	currInsType := instanceType.(string)
	v, ok := dMap["backup_instance_types"]
	backupInstanceTypes := v.([]interface{})
	instanceTypes := make([]*string, 0)
	if !ok || len(backupInstanceTypes) == 0 {
		instanceTypes = append(instanceTypes, &currInsType)
		return instanceTypes
	}
	headType := backupInstanceTypes[0].(string)
	if headType != currInsType {
		instanceTypes = append(instanceTypes, &currInsType)
	}
	for i := range backupInstanceTypes {
		insType := backupInstanceTypes[i].(string)
		instanceTypes = append(instanceTypes, &insType)
	}

	return instanceTypes
}

// this function composes every single parameter to an as scale parameter with json string format
func composeParameterToAsScalingGroupParaSerial(d *schema.ResourceData) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateAutoScalingGroupRequest()

	//this is an empty string
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))

	if *request.MinSize > *request.MaxSize {
		return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
	}

	request.VpcId = helper.String(d.Get("vpc_id").(string))

	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
		if *request.DesiredCapacity > *request.MaxSize ||
			*request.DesiredCapacity < *request.MinSize {
			return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
		}

	}

	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = helper.InterfacesStringsPoint(subnetIds)
	}

	if v, ok := d.GetOk("scaling_mode"); ok {
		request.ServiceSettings = &as.ServiceSettings{ScalingMode: helper.String(v.(string))}
	}

	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok {
		request.MultiZoneSubnetPolicy = helper.String(v.(string))
	}

	result = request.ToJsonString()

	return result, errRet
}

// This function is used to specify tke as group launch config, similar to kubernetesAsScalingConfigParaSerial, but less parameter
func composedKubernetesAsScalingConfigParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateLaunchConfigurationRequest()

	instanceType := dMap["instance_type"].(string)
	request.InstanceType = &instanceType

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		//request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			encrypt, eOk := value["encrypt"].(bool)
			throughputPerformance := value["throughput_performance"].(int)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
			}
			if eOk {
				dataDisk.Encrypt = &encrypt
			}
			if throughputPerformance > 0 {
				dataDisk.ThroughputPerformance = helper.IntUint64(throughputPerformance)
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}

	if v, ok := dMap["password"]; ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := dMap["key_ids"]; ok {
		keyIds := v.([]interface{})
		//request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}

	if request.LoginSettings.Password != nil && *request.LoginSettings.Password == "" {
		request.LoginSettings.Password = nil
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` should be set one")
		return result, errRet
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` can only be supported one")
		return result, errRet
	}

	if v, ok := dMap["security_group_ids"]; ok {
		if list := v.(*schema.Set).List(); len(list) > 0 {
			errRet = fmt.Errorf("The parameter `security_group_ids` has an issue that the actual order of the security group may be inconsistent with the order of your tf code, which will cause your service to be inaccessible. Please use `orderly_security_group_ids` instead.")
			return result, errRet
		}
	}

	if v, ok := dMap["orderly_security_group_ids"]; ok {
		if list := v.([]interface{}); len(list) > 0 {
			request.SecurityGroupIds = helper.InterfacesStringsPoint(list)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	chargeType, ok := dMap["instance_charge_type"].(string)
	if !ok || chargeType == "" {
		chargeType = svcas.INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == svcas.INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == svcas.INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	request.InstanceChargeType = &chargeType

	if v, ok := dMap["cam_role_name"]; ok {
		request.CamRoleName = helper.String(v.(string))
	}

	tmpInstanceNameSettings := &as.InstanceNameSettings{}
	if v, ok := dMap["instance_name"]; ok && v != "" {
		tmpInstanceNameSettings.InstanceName = helper.String(v.(string))
	}

	if v, ok := dMap["instance_name_style"]; ok && v != "" {
		tmpInstanceNameSettings.InstanceNameStyle = helper.String(v.(string))
	}

	if tmpInstanceNameSettings.InstanceName != nil || tmpInstanceNameSettings.InstanceNameStyle != nil {
		request.InstanceNameSettings = tmpInstanceNameSettings
	}

	if v, ok := dMap["host_name"]; ok && v != "" {
		if request.HostNameSettings == nil {
			request.HostNameSettings = &as.HostNameSettings{
				HostName: helper.String(v.(string)),
			}
		} else {
			request.HostNameSettings.HostName = helper.String(v.(string))
		}
	}

	if v, ok := dMap["host_name_style"]; ok && v != "" {
		if request.HostNameSettings != nil {
			request.HostNameSettings.HostNameStyle = helper.String(v.(string))
		} else {
			request.HostNameSettings = &as.HostNameSettings{
				HostNameStyle: helper.String(v.(string)),
			}
		}
	}
	result = request.ToJsonString()
	return result, errRet
}

func isCUDNNEmpty(cudnn *tke.CUDNN) bool {
	return cudnn == nil || (helper.PString(cudnn.Version) == "" && helper.PString(cudnn.Name) == "" && helper.PString(cudnn.DocName) == "" && helper.PString(cudnn.DevName) == "")
}

func isCUDAEmpty(cuda *tke.DriverVersion) bool {
	return cuda == nil || (helper.PString(cuda.Version) == "" && helper.PString(cuda.Name) == "")
}

func isDriverEmpty(driver *tke.DriverVersion) bool {
	return driver == nil || (helper.PString(driver.Version) == "" && helper.PString(driver.Name) == "")
}

func isCustomDriverEmpty(customDriver *tke.CustomDriver) bool {
	return customDriver == nil || helper.PString(customDriver.Address) == ""
}

func composeAsLaunchConfigModifyRequest(d *schema.ResourceData, launchConfigId string) (*as.ModifyLaunchConfigurationAttributesRequest, error) {
	launchConfigRaw := d.Get("auto_scaling_config").([]interface{})
	dMap := launchConfigRaw[0].(map[string]interface{})
	request := as.NewModifyLaunchConfigurationAttributesRequest()
	request.LaunchConfigurationId = &launchConfigId

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		//request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			encrypt, eOk := value["encrypt"].(bool)
			throughputPerformance := value["throughput_performance"].(int)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
			}
			if eOk {
				dataDisk.Encrypt = &encrypt
			}
			if throughputPerformance > 0 {
				dataDisk.ThroughputPerformance = helper.IntUint64(throughputPerformance)
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	if d.HasChange("auto_scaling_config.0.security_group_ids") {
		if v, ok := dMap["security_group_ids"]; ok {
			if list := v.(*schema.Set).List(); len(list) > 0 {
				errRet := fmt.Errorf("The parameter `security_group_ids` has an issue that the actual order of the security group may be inconsistent with the order of your tf code, which will cause your service to be inaccessible. You can check whether the order of your current security groups meets your expectations through the TencentCloud Console, then use `orderly_security_group_ids` field to update them.")
				return nil, errRet
			}
		}
	}

	if d.HasChange("auto_scaling_config.0.orderly_security_group_ids") {
		if v, ok := dMap["orderly_security_group_ids"]; ok {
			if list := v.([]interface{}); len(list) > 0 {
				request.SecurityGroupIds = helper.InterfacesStringsPoint(list)
			}
		}
	}

	chargeType, ok := dMap["instance_charge_type"].(string)

	if !ok || chargeType == "" {
		chargeType = svcas.INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == svcas.INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == svcas.INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	tmpInstanceNameSettings := &as.InstanceNameSettings{}
	if v, ok := dMap["instance_name"]; ok && v != "" {
		tmpInstanceNameSettings.InstanceName = helper.String(v.(string))
	}

	if v, ok := dMap["instance_name_style"]; ok && v != "" {
		tmpInstanceNameSettings.InstanceNameStyle = helper.String(v.(string))
	}

	if tmpInstanceNameSettings.InstanceName != nil || tmpInstanceNameSettings.InstanceNameStyle != nil {
		request.InstanceNameSettings = tmpInstanceNameSettings
	}

	if v, ok := dMap["host_name"]; ok && v != "" {
		if request.HostNameSettings == nil {
			request.HostNameSettings = &as.HostNameSettings{
				HostName: helper.String(v.(string)),
			}
		} else {
			request.HostNameSettings.HostName = helper.String(v.(string))
		}
	}

	if v, ok := dMap["host_name_style"]; ok && v != "" {
		if request.HostNameSettings != nil {
			request.HostNameSettings.HostNameStyle = helper.String(v.(string))
		} else {
			request.HostNameSettings = &as.HostNameSettings{
				HostNameStyle: helper.String(v.(string)),
			}
		}
	}

	// set enhanced_security_service if necessary
	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		if request.EnhancedService != nil {
			request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
				Enabled: helper.Bool(securityService),
			}
		} else {
			request.EnhancedService = &as.EnhancedService{
				SecurityService: &as.RunSecurityServiceEnabled{
					Enabled: helper.Bool(securityService),
				},
			}
		}

	}

	request.InstanceChargeType = &chargeType

	return request, nil
}

func desiredCapacityOutRange(d *schema.ResourceData) bool {
	capacity := d.Get("desired_capacity").(int)
	minSize := d.Get("min_size").(int)
	maxSize := d.Get("max_size").(int)
	return capacity > maxSize || capacity < minSize
}

func resourceTencentCloudKubernetesNodePoolUpdateTaints(ctx context.Context, clusterId string, nodePoolId string) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)

	if d.HasChange("taints") {
		_, n := d.GetChange("taints")

		// clean taints
		if len(n.([]interface{})) == 0 {
			body := map[string]interface{}{
				"ClusterId":  clusterId,
				"NodePoolId": nodePoolId,
				"Taints":     []interface{}{},
			}

			client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("tke")
			request := tchttp.NewCommonRequest("tke", "2018-05-25", "ModifyClusterNodePool")
			err := request.SetActionParameters(body)
			if err != nil {
				return err
			}

			response := tchttp.NewCommonResponse()
			err = client.Send(request, response)
			if err != nil {
				fmt.Printf("update kubernetes node pool taints failed: %v \n", err)
				return err
			}
		} else {
			request := tke.NewModifyClusterNodePoolRequest()
			request.ClusterId = helper.String(clusterId)
			request.NodePoolId = helper.String(nodePoolId)

			if v, ok := d.GetOk("taints"); ok {
				for _, item := range v.([]interface{}) {
					taintsMap := item.(map[string]interface{})
					taint := tke.Taint{}
					if v, ok := taintsMap["key"]; ok {
						taint.Key = helper.String(v.(string))
					}

					if v, ok := taintsMap["value"]; ok {
						taint.Value = helper.String(v.(string))
					}

					if v, ok := taintsMap["effect"]; ok {
						taint.Effect = helper.String(v.(string))
					}

					request.Taints = append(request.Taints, &taint)
				}
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterNodePoolWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update kubernetes node pool taints failed, reason:%+v", logId, err)
				return err
			}
		}

		service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			nodePool, _, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if nodePool != nil && *nodePool.LifeState == "normal" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("node pool status is %s, retry...", *nodePool.LifeState))
		})

		if err != nil {
			return err
		}

		return nil
	}
	return nil
}
