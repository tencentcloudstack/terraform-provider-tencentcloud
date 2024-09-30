package tke

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var importFlag1 = false

var GlobalClusterId string
var CreateClusterInstancesVpcId string
var CreateClusterInstancesProjectId int64
var WorkersInstanceIds []*string
var WorkersNewWorkerInstancesList []map[string]interface{}
var WorkersLabelsMap map[string]string

func init() {
	// need to support append by multiple calls when the paging occurred
	WorkersNewWorkerInstancesList = make([]map[string]interface{}, 0)
	WorkersLabelsMap = make(map[string]string)
	WorkersInstanceIds = make([]*string, 0)
}

func customScaleWorkerResourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importFlag1 = true
	err := resourceTencentCloudKubernetesScaleWorkerRead(d, m)
	if err != nil {
		return nil, fmt.Errorf("failed to import resource")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPostRequest1(ctx context.Context, req *cvm.DescribeInstancesRequest, resp *cvm.DescribeInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)

	instances := make([]*cvm.Instance, 0)
	instances = append(instances, resp.Response.InstanceSet...)

	instanceList := make([]interface{}, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"count":                               1,
			"instance_charge_type_prepaid_period": 1,
			"instance_type":                       helper.PString(instance.InstanceType),
			"subnet_id":                           helper.PString(instance.VirtualPrivateCloud.SubnetId),
			"availability_zone":                   helper.PString(instance.Placement.Zone),
			"instance_name":                       helper.PString(instance.InstanceName),
			"instance_charge_type":                helper.PString(instance.InstanceChargeType),
			"system_disk_type":                    helper.PString(instance.SystemDisk.DiskType),
			"system_disk_size":                    helper.PInt64(instance.SystemDisk.DiskSize),
			"internet_charge_type":                helper.PString(instance.InternetAccessible.InternetChargeType),
			"bandwidth_package_id":                helper.PString(instance.InternetAccessible.BandwidthPackageId),
			"internet_max_bandwidth_out":          helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut),
			"security_group_ids":                  helper.StringsInterfaces(instance.SecurityGroupIds),
			"img_id":                              helper.PString(instance.ImageId),
		}

		if instance.RenewFlag != nil && helper.PString(instance.InstanceChargeType) == "PREPAID" {
			mapping["instance_charge_type_prepaid_renew_flag"] = helper.PString(instance.RenewFlag)
		} else {
			mapping["instance_charge_type_prepaid_renew_flag"] = ""
		}
		if helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut) > 0 {
			mapping["public_ip_assigned"] = true
		}

		if instance.CamRoleName != nil {
			mapping["cam_role_name"] = instance.CamRoleName
		}
		if instance.LoginSettings != nil {
			if instance.LoginSettings.KeyIds != nil && len(instance.LoginSettings.KeyIds) > 0 {
				mapping["key_ids"] = helper.StringsInterfaces(instance.LoginSettings.KeyIds)
			}
			if instance.LoginSettings.Password != nil {
				mapping["password"] = helper.PString(instance.LoginSettings.Password)
			}
		}
		if instance.DisasterRecoverGroupId != nil && helper.PString(instance.DisasterRecoverGroupId) != "" {
			mapping["disaster_recover_group_ids"] = []string{helper.PString(instance.DisasterRecoverGroupId)}
		}
		if instance.HpcClusterId != nil {
			mapping["hpc_cluster_id"] = helper.PString(instance.HpcClusterId)
		}

		dataDisks := make([]interface{}, 0, len(instance.DataDisks))
		for _, v := range instance.DataDisks {
			dataDisk := map[string]interface{}{
				"disk_type":   helper.PString(v.DiskType),
				"disk_size":   helper.PInt64(v.DiskSize),
				"snapshot_id": helper.PString(v.DiskId),
				"encrypt":     helper.PBool(v.Encrypt),
				"kms_key_id":  helper.PString(v.KmsKeyId),
			}
			dataDisks = append(dataDisks, dataDisk)
		}

		mapping["data_disk"] = dataDisks // worker_config.data_disk
		instanceList = append(instanceList, mapping)
	}
	if importFlag1 {
		_ = d.Set("worker_config", instanceList)
	}

	// The machines I generated was deleted by others.
	if len(WorkersNewWorkerInstancesList) == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", GlobalClusterId)
	_ = d.Set("labels", WorkersLabelsMap)
	_ = d.Set("worker_instances_list", WorkersNewWorkerInstancesList)

	return nil
}
func clusterInstanceParamHandle(ctx context.Context, req *tke.DescribeClusterInstancesRequest, resp *tke.DescribeClusterInstancesResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var has = map[string]bool{}

	workerInstancesList := d.Get("worker_instances_list").([]interface{})
	instanceMap := make(map[string]bool)
	for _, v := range workerInstancesList {
		infoMap, ok := v.(map[string]interface{})
		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}

		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}

		if instanceMap[instanceId] {
			log.Printf("[WARN]The same instance id exists in the list")
		}

		instanceMap[instanceId] = true
	}
	workers := make([]InstanceInfo, 0, 100)
	for _, item := range resp.Response.InstanceSet {
		if has[*item.InstanceId] {
			return fmt.Errorf("get repeated instance_id[%s] when doing DescribeClusterInstances", *item.InstanceId)
		}
		has[*item.InstanceId] = true
		instanceInfo := InstanceInfo{
			InstanceId:               *item.InstanceId,
			InstanceRole:             *item.InstanceRole,
			InstanceState:            *item.InstanceState,
			FailedReason:             *item.FailedReason,
			InstanceAdvancedSettings: item.InstanceAdvancedSettings,
		}
		if item.CreatedTime != nil {
			instanceInfo.CreatedTime = *item.CreatedTime
		}
		if item.NodePoolId != nil {
			instanceInfo.NodePoolId = *item.NodePoolId
		}
		if item.LanIP != nil {
			instanceInfo.LanIp = *item.LanIP
		}
		if instanceInfo.InstanceRole == TKE_ROLE_WORKER {
			workers = append(workers, instanceInfo)
		}
	}

	for sub, cvmInfo := range workers {
		if _, ok := instanceMap[cvmInfo.InstanceId]; !ok {
			continue
		}
		WorkersInstanceIds = append(WorkersInstanceIds, &workers[sub].InstanceId)
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = cvmInfo.InstanceId
		tempMap["instance_role"] = cvmInfo.InstanceRole
		tempMap["instance_state"] = cvmInfo.InstanceState
		tempMap["failed_reason"] = cvmInfo.FailedReason
		tempMap["lan_ip"] = cvmInfo.LanIp

		WorkersNewWorkerInstancesList = append(WorkersNewWorkerInstancesList, tempMap)
		if cvmInfo.InstanceAdvancedSettings != nil {
			if cvmInfo.InstanceAdvancedSettings.Labels != nil {
				for _, v := range cvmInfo.InstanceAdvancedSettings.Labels {
					WorkersLabelsMap[helper.PString(v.Name)] = helper.PString(v.Value)
				}
			}

			_ = d.Set("unschedulable", helper.PInt64(cvmInfo.InstanceAdvancedSettings.Unschedulable))

			if importFlag1 {
				_ = d.Set("docker_graph_path", helper.PString(cvmInfo.InstanceAdvancedSettings.DockerGraphPath))
				_ = d.Set("desired_pod_num", helper.PInt64(cvmInfo.InstanceAdvancedSettings.DesiredPodNumber))
				_ = d.Set("mount_target", helper.PString(cvmInfo.InstanceAdvancedSettings.MountTarget))
			}

			if cvmInfo.InstanceAdvancedSettings.DataDisks != nil && len(cvmInfo.InstanceAdvancedSettings.DataDisks) > 0 {
				dataDisks := make([]interface{}, 0, len(cvmInfo.InstanceAdvancedSettings.DataDisks))
				for i := range cvmInfo.InstanceAdvancedSettings.DataDisks {
					item := cvmInfo.InstanceAdvancedSettings.DataDisks[i]
					disk := make(map[string]interface{})
					disk["disk_type"] = helper.PString(item.DiskType)
					disk["disk_size"] = helper.PInt64(item.DiskSize)
					disk["file_system"] = helper.PString(item.FileSystem)
					disk["auto_format_and_mount"] = helper.PBool(item.AutoFormatAndMount)
					disk["mount_target"] = helper.PString(item.MountTarget)
					disk["disk_partition"] = helper.PString(item.MountTarget)
					dataDisks = append(dataDisks, disk)
				}
				if importFlag1 {
					_ = d.Set("data_disk", dataDisks) // out layer data_disk
				}
			}

			if cvmInfo.InstanceAdvancedSettings.GPUArgs != nil {
				setting := cvmInfo.InstanceAdvancedSettings.GPUArgs

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

				if importFlag1 {
					if driverEmptyFlag || cudaEmptyFlag || cudnnEmptyFlag || customDriverEmptyFlag {
						_ = d.Set("gpu_args", []interface{}{gpuArgs})
					}
				}
			}
		}
	}
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerDeletePostRequest0(ctx context.Context, req *tke.DescribeClustersRequest, resp *tke.DescribeClustersResponse) *resource.RetryError {
	if len(resp.Response.Clusters) == 0 {
		return resource.NonRetryableError(fmt.Errorf("The cluster has been deleted"))
	}
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPostFillRequest0(ctx context.Context, req *tke.DescribeClustersRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	items := strings.Split(d.Id(), tccommon.FILED_SP)

	instanceMap := make(map[string]bool)
	oldWorkerInstancesList := d.Get("worker_instances_list").([]interface{})
	if importFlag1 {
		GlobalClusterId = items[0]
		if len(items[1:]) >= 2 {
			return fmt.Errorf("only one additional configuration of virtual machines is now supported now, " +
				"so should be 1")
		}
		infoMap := map[string]interface{}{
			"instance_id": items[1],
		}
		oldWorkerInstancesList = append(oldWorkerInstancesList, infoMap)
	} else {
		GlobalClusterId = d.Get("cluster_id").(string)
	}

	if GlobalClusterId == "" {
		return fmt.Errorf("tke.`cluster_id` is empty.")
	}

	for _, v := range oldWorkerInstancesList {
		infoMap, ok := v.(map[string]interface{})
		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}
		if instanceMap[instanceId] {
			continue
		}
		instanceMap[instanceId] = true
	}

	return nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPostRequest0(ctx context.Context, req *tke.DescribeClustersRequest, resp *tke.DescribeClustersResponse) error {
	if len(resp.Response.Clusters) == 0 {
		return fmt.Errorf("The cluster has been deleted")
	}
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPostFillRequest2(ctx context.Context, req *cvm.DescribeInstancesRequest) error {
	req.InstanceIds = WorkersInstanceIds
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerCreateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var cvms RunInstancesForNode
	var iAdvanced tke.InstanceAdvancedSettings
	cvms.Work = []string{}

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)
	if clusterId == "" {
		return fmt.Errorf("`cluster_id` is empty.")
	}

	info, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		return fmt.Errorf("cluster [%s] is not exist.", clusterId)
	}

	dMap := make(map[string]interface{}, 5)
	//mount_target, docker_graph_path, data_disk, extra_args, desired_pod_num
	iAdvancedParas := []string{"mount_target", "docker_graph_path", "extra_args", "data_disk", "desired_pod_num", "gpu_args", "taints"}
	for _, k := range iAdvancedParas {
		if v, ok := d.GetOk(k); ok {
			dMap[k] = v
		}
	}
	iAdvanced = tkeGetInstanceAdvancedPara(dMap, meta)

	iAdvanced.Labels = GetTkeLabels(d, "labels")
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.Int64(int64(temp.(int)))
	}

	if v, ok := d.GetOk("pre_start_user_script"); ok {
		iAdvanced.PreStartUserScript = helper.String(v.(string))
	}

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
			iAdvanced.Taints = append(iAdvanced.Taints, &taint)
		}
	}

	if v, ok := d.GetOk("user_script"); ok {
		iAdvanced.UserScript = helper.String(v.(string))
	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, info.VpcId, info.ProjectId)
			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)
		}
	}
	if len(cvms.Work) != 1 {
		return fmt.Errorf("only one additional configuration of virtual machines is now supported now, " +
			"so len(cvms.Work) should be 1")
	}

	instanceIds, err := service.CreateClusterInstances(ctx, clusterId, cvms.Work[0], iAdvanced)
	if err != nil {
		return err
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(instanceIds))
	for _, v := range instanceIds {
		if v == "" {
			return fmt.Errorf("CreateClusterInstances return one instanceId is empty")
		}
		infoMap := make(map[string]interface{})
		infoMap["instance_id"] = v
		infoMap["instance_role"] = TKE_ROLE_WORKER
		workerInstancesList = append(workerInstancesList, infoMap)
	}

	if err = d.Set("worker_instances_list", workerInstancesList); err != nil {
		return err
	}

	//修改id设置,不符合id规则
	id := clusterId + tccommon.FILED_SP + strings.Join(instanceIds, tccommon.COMMA_SP)
	d.SetId(id)

	//wait for LANIP
	time.Sleep(tccommon.ReadRetryTimeout)

	// wait for all instances status running
	waitRequest := tke.NewDescribeClusterInstancesRequest()
	waitRequest.ClusterId = &clusterId
	waitRequest.InstanceIds = helper.Strings(instanceIds)
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		var (
			offset         int64 = 0
			limit          int64 = 100
			tmpInstanceSet []*tke.Instance
		)

		// get all instances
		for {
			waitRequest.Limit = &limit
			waitRequest.Offset = &offset
			ratelimit.Check(waitRequest.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusterInstances(waitRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG] api[%s] success, request body [%s], response body [%s]\n", waitRequest.GetAction(), waitRequest.ToJsonString(), result.ToJsonString())
			}

			if result == nil || len(result.Response.InstanceSet) == 0 {
				break
			}

			tmpInstanceSet = append(tmpInstanceSet, result.Response.InstanceSet...)

			if len(result.Response.InstanceSet) < int(limit) {
				break
			}

			offset += limit
		}

		// check instances status
		if len(tmpInstanceSet) == 0 {
			return resource.NonRetryableError(fmt.Errorf("there is no instances in set"))
		} else {
			var (
				stop int
				flag bool
			)

			for _, v := range instanceIds {
				for _, instance := range tmpInstanceSet {
					if v == *instance.InstanceId {
						if *instance.InstanceState == "running" {
							stop += 1
							flag = true
						} else if *instance.InstanceState == "failed" {
							stop += 1
							log.Printf("instance: %s status is failed.", v)
						} else {
							continue
						}
					}
				}
			}

			if stop == len(instanceIds) && flag {
				return nil
			} else if stop == len(instanceIds) && !flag {
				return resource.NonRetryableError(fmt.Errorf("The instances being created have all failed."))
			} else {
				return resource.RetryableError(fmt.Errorf("cluster instances is still initializing."))
			}
		}
	})

	if err != nil {
		log.Printf("[CRITAL] kubernetes scale worker instances status error, reason:%+v", err)
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesScaleWorkerDeleteOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	clusterId := idSplit[0]
	workerInstancesList := d.Get("worker_instances_list").([]interface{})

	instanceMap := make(map[string]bool)
	for _, v := range workerInstancesList {

		infoMap, ok := v.(map[string]interface{})

		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}

		if instanceMap[instanceId] {
			log.Printf("[WARN]The same instance id exists in the list")
		}

		instanceMap[instanceId] = true

	}

	_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}

			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	needDeletes := []string{}
	for _, cvm := range workers {
		if _, ok := instanceMap[cvm.InstanceId]; ok {
			needDeletes = append(needDeletes, cvm.InstanceId)
		}
	}
	// The machines I generated was deleted by others.
	if len(needDeletes) == 0 {
		return nil
	}

	err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)
	if err != nil {
		err = resource.Retry(3*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)

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
		})
	}
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPostFillRequest1(ctx context.Context, req *cvm.DescribeInstancesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := tccommon.GetLogId(ctx)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]

	request := tke.NewDescribeClusterInstancesRequest()
	request.ClusterId = helper.String(clusterId)

	instanceSet := make([]*tke.Instance, 0)

	var offset int64 = 0
	var pageSize int64 = 100
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusterInstances(request)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if err := clusterInstanceParamHandle(ctx, request, response); err != nil {
			return err
		}

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		count := len(response.Response.InstanceSet)
		instanceSet = append(instanceSet, response.Response.InstanceSet...)

		if count < int(pageSize) {
			break
		}
		offset += pageSize
	}
	if instanceSet == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_scale_worker` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	return nil
}

func resourceTencentCloudKubernetesScaleWorkerReadPreRequest1(ctx context.Context, req *cvm.DescribeInstancesRequest) error {
	req.InstanceIds = WorkersInstanceIds

	return nil
}
