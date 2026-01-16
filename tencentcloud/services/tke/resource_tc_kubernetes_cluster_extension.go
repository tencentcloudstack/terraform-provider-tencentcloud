package tke

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
)

var importClsFlag = false

func customResourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importClsFlag = true
	err := resourceTencentCloudKubernetesClusterRead(d, m)
	if err != nil {
		return nil, fmt.Errorf("failed to import resource")
	}
	return []*schema.ResourceData{d}, nil
}

func resourceTencentCloudKubernetesClusterCreatePostFillRequest0(ctx context.Context, req *tke.CreateClusterRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	var (
		basic              ClusterBasicSetting
		advanced           ClusterAdvancedSettings
		cvms               RunInstancesForNode
		iAdvanced          InstanceAdvancedSettings
		iDiskMountSettings []*tke.InstanceDataDiskMountSetting
		cidrSet            ClusterCidrSettings
		clusterInternet    = d.Get("cluster_internet").(bool)
		clusterIntranet    = d.Get("cluster_intranet").(bool)
		intranetSubnetId   = d.Get("cluster_intranet_subnet_id").(string)
	)

	clusterDeployType := d.Get("cluster_deploy_type").(string)

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}
	if !clusterIntranet && intranetSubnetId != "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` can only set when `cluster_intranet` is true")
	}

	_, workerConfigOk := d.GetOk("worker_config")
	if !workerConfigOk && clusterInternet {
		return fmt.Errorf("when creating a cluster, if `cluster_internet` is true, you need to configure the `worker_config` field to ensure that there are available nodes in the cluster.")
	}

	vpcId := d.Get("vpc_id").(string)
	if vpcId != "" {
		basic.VpcId = vpcId
	}

	cluster_os := d.Get("cluster_os").(string)

	if v, ok := tkeClusterOsMap[cluster_os]; ok {
		basic.ClusterOs = v
	} else {
		basic.ClusterOs = cluster_os
	}

	if tkeClusterOsMap[cluster_os] != "" {
		basic.ClusterOs = tkeClusterOsMap[cluster_os]
	} else {
		basic.ClusterOs = cluster_os
	}

	advanced.NetworkType = d.Get("network_type").(string)

	if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
		if v, ok := d.GetOk("vpc_cni_type"); ok {
			advanced.VpcCniType = v.(string)
		} else {
			advanced.VpcCniType = "tke-route-eni"
		}
	}

	cidrSet.ClusterCidr = d.Get("cluster_cidr").(string)
	cidrSet.ServiceCIDR = d.Get("service_cidr").(string)

	if ClaimExpiredSeconds, ok := d.GetOk("claim_expired_seconds"); ok {
		cidrSet.ClaimExpiredSeconds = int64(ClaimExpiredSeconds.(int))
	} else {
		cidrSet.ClaimExpiredSeconds = int64(300)

		if err := d.Set("claim_expired_seconds", 300); err != nil {
			return fmt.Errorf("error setting claim_expired_seconds: %s", err)
		}
	}

	if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
		// VPC-CNI cluster need to set eni subnet and service cidr.
		eniSubnetIdList := d.Get("eni_subnet_ids").([]interface{})
		for index := range eniSubnetIdList {
			subnetId := eniSubnetIdList[index].(string)
			cidrSet.EniSubnetIds = append(cidrSet.EniSubnetIds, subnetId)
		}
		if cidrSet.ServiceCIDR == "" || len(cidrSet.EniSubnetIds) == 0 {
			return fmt.Errorf("`service_cidr` must be set and `eni_subnet_ids` must be set when cluster `network_type` is VPC-CNI.")
		}
	} else {
		// GR cluster
		if cidrSet.ClusterCidr == "" {
			return fmt.Errorf("`cluster_cidr` must be set when cluster `network_type` is GR")
		}
		items := strings.Split(cidrSet.ClusterCidr, "/")
		if len(items) != 2 {
			return fmt.Errorf("`cluster_cidr` must be network segment ")
		}

		bitNumber, err := strconv.ParseInt(items[1], 10, 64)

		if err != nil {
			return fmt.Errorf("`cluster_cidr` must be network segment ")
		}

		if math.Pow(2, float64(32-bitNumber)) <= float64(cidrSet.MaxNodePodNum) {
			return fmt.Errorf("`cluster_cidr` Network segment range is too small, can not cover cluster_max_service_num")
		}

		if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY && d.Get("cluster_subnet_id").(string) == "" {
			return fmt.Errorf("`cluster_subnet_id` must be set ")
		}
	}

	overrideSettings := &OverrideSettings{
		Master: make([]tke.InstanceAdvancedSettings, 0),
		Work:   make([]tke.InstanceAdvancedSettings, 0),
	}

	cdc_id := d.Get("cdc_id").(string)

	if cdc_id == "" {
		if masters, ok := d.GetOk("master_config"); ok {
			if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED {
				return fmt.Errorf("if `cluster_deploy_type` is `MANAGED_CLUSTER` , You don't need define the master yourself")
			}
			var masterCount int64 = 0
			masterList := masters.([]interface{})
			for index := range masterList {
				master := masterList[index].(map[string]interface{})
				paraJson, count, err := tkeGetCvmRunInstancesPara(master, meta, vpcId, basic.ProjectId)
				if err != nil {
					return err
				}

				cvms.Master = append(cvms.Master, paraJson)
				masterCount += count

				if v, ok := master["desired_pod_num"]; ok {
					dpNum := int64(v.(int))
					if dpNum != DefaultDesiredPodNum {
						overrideSettings.Master = append(overrideSettings.Master, tke.InstanceAdvancedSettings{DesiredPodNumber: helper.Int64(dpNum)})
					}
				}
			}
			if masterCount < 3 {
				return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` len(master_config) should >=3")
			}
		} else if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` , You need define the master yourself")
		}
	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, vpcId, basic.ProjectId)
			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)

			if v, ok := worker["desired_pod_num"]; ok {
				dpNum := int64(v.(int))
				if dpNum != DefaultDesiredPodNum {
					overrideSettings.Work = append(overrideSettings.Work, tke.InstanceAdvancedSettings{DesiredPodNumber: helper.Int64(dpNum)})
				}
			}

			if v, ok := worker["data_disk"]; ok {
				var (
					instanceType = worker["instance_type"].(string)
					zone         = worker["availability_zone"].(string)
				)
				iDiskMountSetting := &tke.InstanceDataDiskMountSetting{
					InstanceType: &instanceType,
					Zone:         &zone,
				}

				diskList := v.([]interface{})
				for _, d := range diskList {
					var (
						disk               = d.(map[string]interface{})
						diskType           = disk["disk_type"].(string)
						diskSize           = int64(disk["disk_size"].(int))
						fileSystem         = disk["file_system"].(string)
						autoFormatAndMount = disk["auto_format_and_mount"].(bool)
						mountTarget        = disk["mount_target"].(string)
						diskPartition      = disk["disk_partition"].(string)
					)

					dataDisk := &tke.DataDisk{
						DiskType:           &diskType,
						DiskSize:           &diskSize,
						AutoFormatAndMount: &autoFormatAndMount,
					}

					if fileSystem != "" {
						dataDisk.FileSystem = &fileSystem
					}

					if mountTarget != "" {
						dataDisk.MountTarget = &mountTarget
					}

					if diskPartition != "" {
						dataDisk.DiskPartition = &diskPartition
					}

					iDiskMountSetting.DataDisks = append(iDiskMountSetting.DataDisks, dataDisk)
				}

				iDiskMountSettings = append(iDiskMountSettings, iDiskMountSetting)
			}
		}
	}

	tags := helper.GetTags(d, "tags")

	iAdvanced.Labels = GetTkeLabels(d, "labels")

	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for i := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArgs[i])
		}
	}

	if temp, ok := d.GetOk("docker_graph_path"); ok {
		iAdvanced.DockerGraphPath = temp.(string)
	} else {
		iAdvanced.DockerGraphPath = "/var/lib/docker"
	}

	if preStartUserScript, ok := d.GetOk("pre_start_user_script"); ok {
		iAdvanced.PreStartUserScript = preStartUserScript.(string)
	}

	// ExistedInstancesForNode
	existInstances := make([]*tke.ExistedInstancesForNode, 0)
	if instances, ok := d.GetOk("exist_instance"); ok {
		instanceList := instances.(*schema.Set).List()
		for index := range instanceList {
			instance := instanceList[index].(map[string]interface{})
			existedInstance, _ := tkeGetCvmExistInstancesPara(instance)
			existInstances = append(existInstances, &existedInstance)
		}
	}

	// RunInstancesForNode（master_config+worker_config) 和 ExistedInstancesForNode 不能同时存在
	if len(cvms.Master)+len(cvms.Work) > 0 && len(existInstances) > 0 {
		return fmt.Errorf("master_config+worker_config and exist_instance can not exist at the same time")
	}

	//request传参
	req.ClusterBasicSettings.ClusterOs = &basic.ClusterOs
	req.ClusterBasicSettings.VpcId = &basic.VpcId
	for k, v := range tags {
		if len(req.ClusterBasicSettings.TagSpecification) == 0 {
			req.ClusterBasicSettings.TagSpecification = []*tke.TagSpecification{{
				ResourceType: helper.String("cluster"),
			}}
		}

		req.ClusterBasicSettings.TagSpecification[0].Tags = append(req.ClusterBasicSettings.TagSpecification[0].Tags, &tke.Tag{
			Key:   helper.String(k),
			Value: helper.String(v),
		})
	}

	req.ClusterAdvancedSettings.VpcCniType = &advanced.VpcCniType

	req.InstanceAdvancedSettings.DockerGraphPath = &iAdvanced.DockerGraphPath
	req.InstanceAdvancedSettings.PreStartUserScript = &iAdvanced.PreStartUserScript
	req.InstanceAdvancedSettings.UserScript = &iAdvanced.UserScript

	if len(iAdvanced.DataDisks) > 0 {
		req.InstanceAdvancedSettings.DataDisks = iAdvanced.DataDisks
	}

	if overrideSettings != nil {
		if len(overrideSettings.Master)+len(overrideSettings.Work) > 0 &&
			len(overrideSettings.Master)+len(overrideSettings.Work) != (len(cvms.Master)+len(cvms.Work)) {
			return fmt.Errorf("len(overrideSettings) != (len(cvms.Master)+len(cvms.Work))")
		}
	}

	req.RunInstancesForNode = []*tke.RunInstancesForNode{}

	if cdc_id != "" {
		req.ClusterType = helper.String(clusterDeployType)
	} else if len(cvms.Master) != 0 {

		var node tke.RunInstancesForNode
		node.NodeRole = helper.String(TKE_ROLE_MASTER_ETCD)
		node.RunInstancesPara = []*string{}
		req.ClusterType = helper.String(TKE_DEPLOY_TYPE_INDEPENDENT)
		for v := range cvms.Master {
			node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Master[v])
			if overrideSettings != nil && len(overrideSettings.Master) != 0 {
				node.InstanceAdvancedSettingsOverrides = append(node.InstanceAdvancedSettingsOverrides, &overrideSettings.Master[v])
			}
		}
		req.RunInstancesForNode = append(req.RunInstancesForNode, &node)

	} else {
		req.ClusterType = helper.String(TKE_DEPLOY_TYPE_MANAGED)
	}

	if len(cvms.Work) != 0 {
		var node tke.RunInstancesForNode
		node.NodeRole = helper.String(TKE_ROLE_WORKER)
		node.RunInstancesPara = []*string{}
		for v := range cvms.Work {
			node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Work[v])
			if overrideSettings != nil && len(overrideSettings.Work) != 0 {
				node.InstanceAdvancedSettingsOverrides = append(node.InstanceAdvancedSettingsOverrides, &overrideSettings.Work[v])
			}
		}
		req.RunInstancesForNode = append(req.RunInstancesForNode, &node)
	}

	if len(iDiskMountSettings) != 0 {
		req.InstanceDataDiskMountSettings = iDiskMountSettings
	}

	req.ClusterCIDRSettings.EniSubnetIds = common.StringPtrs(cidrSet.EniSubnetIds)
	req.ClusterCIDRSettings.ClaimExpiredSeconds = &cidrSet.ClaimExpiredSeconds

	if len(existInstances) > 0 {
		req.ExistedInstancesForNode = existInstances
	}
	return nil
}

func resourceTencentCloudKubernetesClusterCreatePostHandleResponse0(ctx context.Context, resp *tke.CreateClusterResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := *resp.Response.ClusterId

	var (
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
	)

	_, _, err := service.DescribeClusterInstances(ctx, id)

	if err != nil {
		// create often cost more than 20 Minutes.
		err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, _, err = service.DescribeClusterInstances(ctx, id)

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

	err = service.CheckOneOfClusterNodeReady(ctx, id, clusterInternet || clusterIntranet)

	if err != nil {
		return err
	}

	//intranet
	if clusterIntranet {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, intranetSubnetId, clusterInternetSecurityGroup, false, clusterIntranetDomain, "")
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id, false)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create intranet cluster endpoint status still is %s", id, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create intranet cluster endpoint error ,status is %s,message is %s", id, status, message))
		})
		if err != nil {
			return err
		}
	}

	if clusterInternet {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, "", clusterInternetSecurityGroup, true, clusterInternetDomain, "")
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id, true)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create cluster internet endpoint status still is %s", id, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create cluster internet endpoint error ,status is %s,message is %s", id, status, message))
		})
		if err != nil {
			return err
		}
	}

	//Modify node pool global config(sync)
	if _, ok := d.GetOk("node_pool_global_config"); ok {
		request := tkeGetNodePoolGlobalConfig(d)
		request.ClusterId = &id
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := service.ModifyClusterNodePoolGlobalConfig(ctx, request)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	// sync
	if v, ok := d.GetOk("acquire_cluster_admin_role"); ok && v.(bool) {
		_, err := service.AcquireClusterAdminRole(ctx, id)
		if err != nil {
			return err
		}
	}

	// async
	if _, ok := d.GetOk("auth_options"); ok {
		request := tkeGetAuthOptions(d, id)
		if err := service.ModifyClusterAuthenticationOptions(ctx, request); err != nil {
			return err
		}

		// wait
		err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			resp, inErr := service.DescribeKubernetesAuthAttachmentById(ctx, id)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			if resp == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe cluster aauthentication options failed, Response is nil."))
			}

			if resp.LatestOperationState != nil || *resp.LatestOperationState == "Success" {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("Modify auth options running..."))
		})

		if err != nil {
			return err
		}
	}

	// async
	if v, ok := helper.InterfacesHeadMap(d, "log_agent"); ok {
		enabled := v["enabled"].(bool)
		rootDir := v["kubelet_root_dir"].(string)

		if enabled {
			err := service.SwitchLogAgent(ctx, id, rootDir, enabled)
			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				resp, inErr := service.DescribeLogSwitches(ctx, id)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}

				if resp == nil || len(resp) < 1 {
					return resource.NonRetryableError(fmt.Errorf("Describe log switches failed, Response is nil."))
				}

				ret := resp[0]
				if ret.Log != nil && ret.Log.Status != nil && *ret.Log.Status == "opened" {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("Modify log agent running..."))
			})

			if err != nil {
				return err
			}
		}
	}

	// async
	if v, ok := helper.InterfacesHeadMap(d, "event_persistence"); ok {
		enabled := v["enabled"].(bool)
		logSetId := v["log_set_id"].(string)
		topicId := v["topic_id"].(string)
		if enabled {
			err := service.SwitchEventPersistence(ctx, id, logSetId, topicId, enabled, false)
			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				resp, inErr := service.DescribeLogSwitches(ctx, id)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}

				if resp == nil || len(resp) < 1 {
					return resource.NonRetryableError(fmt.Errorf("Describe event persistence failed, Response is nil."))
				}

				ret := resp[0]
				if ret.Event != nil && ret.Event.Status != nil && *ret.Event.Status == "opened" {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("Modify event persistence running..."))
			})

			if err != nil {
				return err
			}
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "cluster_audit"); ok {
		enabled := v["enabled"].(bool)
		logSetId := v["log_set_id"].(string)
		topicId := v["topic_id"].(string)
		if enabled {
			err := service.SwitchClusterAudit(ctx, id, logSetId, topicId, enabled, false)
			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				resp, inErr := service.DescribeLogSwitches(ctx, id)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}

				if resp == nil || len(resp) < 1 {
					return resource.NonRetryableError(fmt.Errorf("Describe cluster audit failed, Response is nil."))
				}

				ret := resp[0]
				if ret.Audit != nil && ret.Audit.Status != nil && *ret.Audit.Status == "opened" {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("Modify cluster audit running..."))
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceTencentCloudKubernetesClusterReadPostHandleResponse0(ctx context.Context, resp *tke.Cluster) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	cluster := resp

	var clusterInfo ClusterInfo
	clusterInfo.AutoUpgradeClusterLevel = cluster.AutoUpgradeClusterLevel
	if cluster.ClusterNetworkSettings != nil {
		clusterInfo.KubeProxyMode = helper.PString(cluster.ClusterNetworkSettings.KubeProxyMode)
		clusterInfo.ServiceCIDR = helper.PString(cluster.ClusterNetworkSettings.ServiceCIDR)
		clusterInfo.IsDualStack = helper.PBool(cluster.ClusterNetworkSettings.IsDualStack)
	}
	clusterInfo.ContainerRuntime = helper.PString(cluster.ContainerRuntime)
	clusterInfo.OsCustomizeType = helper.PString(cluster.OsCustomizeType)
	clusterInfo.RuntimeVersion = helper.PString(cluster.RuntimeVersion)
	clusterInfo.Property = helper.PString(cluster.Property)
	clusterInfo.DeployType = strings.ToUpper(*cluster.ClusterType)
	projectMap, err := helper.JsonToMap(*cluster.Property)
	if err != nil {
		return err
	}
	if projectMap != nil {
		if projectMap["VpcCniType"] != nil {
			vpcCniType := projectMap["VpcCniType"].(string)
			clusterInfo.VpcCniType = vpcCniType
		}
		if projectMap["NetworkType"] != nil {
			networkType := projectMap["NetworkType"].(string)
			clusterInfo.NetworkType = networkType
		}
	}
	if len(cluster.TagSpecification) > 0 {
		clusterInfo.Tags = make(map[string]string)
		for _, tag := range cluster.TagSpecification[0].Tags {
			clusterInfo.Tags[*tag.Key] = *tag.Value
		}
	}

	// 兼容旧的 cluster_os 的 key, 由于 cluster_os有默认值，所以不大可能为空
	oldOs := d.Get("cluster_os").(string)
	newOs := tkeToShowClusterOs(*cluster.ClusterOs)

	if (oldOs == TkeClusterOsCentOS76 && newOs == TKE_CLUSTER_OS_CENTOS76) ||
		(oldOs == TkeClusterOsUbuntu18 && newOs == TKE_CLUSTER_OS_UBUNTU18) {
		newOs = oldOs
	}
	_ = d.Set("cluster_os", newOs)
	// When ImageId is not empty, cluster_os is ImageId. When ImageId is empty, cluster_os displays ClusterOs
	if cluster.ImageId != nil && *cluster.ImageId != "" {
		_ = d.Set("cluster_os", *cluster.ImageId)
	}

	_ = d.Set("tags", clusterInfo.Tags)

	_ = d.Set("vpc_cni_type", clusterInfo.VpcCniType)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(clusterInfo.Property), &data)
	if err != nil {
		return fmt.Errorf("error:%v", err)
	}

	if importClsFlag && clusterInfo.DeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
		var masters []InstanceInfo
		var errRet error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			masters, _, errRet = service.DescribeClusterInstancesByRole(ctx, d.Id(), "MASTER_OR_ETCD")
			if e, ok := errRet.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if errRet != nil {
				return resource.RetryableError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		var instances []*cvm.Instance
		instanceIds := make([]*string, 0)
		for _, instance := range masters {
			instanceIds = append(instanceIds, helper.String(instance.InstanceId))
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instances, errRet = cvmService.DescribeInstanceByFilter(ctx, instanceIds, nil)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}

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

			mapping["data_disk"] = dataDisks
			instanceList = append(instanceList, mapping)
		}

		_ = d.Set("master_config", instanceList)
	}

	if importClsFlag {
		_ = d.Set("is_dual_stack", clusterInfo.IsDualStack)
		networkType, _ := data["NetworkType"].(string)
		_ = d.Set("network_type", networkType)

		nodeNameType, _ := data["NodeNameType"].(string)
		_ = d.Set("node_name_type", nodeNameType)

		enableCustomizedPodCIDR, _ := data["EnableCustomizedPodCIDR"].(bool)
		_ = d.Set("enable_customized_pod_cidr", enableCustomizedPodCIDR)

		basePodNumber, _ := data["BasePodNumber"].(int)
		_ = d.Set("base_pod_num", basePodNumber)

		isNonStaticIpMode, _ := data["IsNonStaticIpMode"].(bool)
		_ = d.Set("is_non_static_ip_mode", isNonStaticIpMode)

		_ = d.Set("runtime_version", clusterInfo.RuntimeVersion)
		_ = d.Set("cluster_os_type", clusterInfo.OsCustomizeType)
		_ = d.Set("container_runtime", clusterInfo.ContainerRuntime)
		_ = d.Set("kube_proxy_mode", clusterInfo.KubeProxyMode)
		_ = d.Set("service_cidr", clusterInfo.ServiceCIDR)
		_ = d.Set("upgrade_instances_follow_cluster", false)

		switchSet, err := service.DescribeLogSwitches(ctx, d.Id())
		if err != nil {
			return err
		}
		logAgents := make([]map[string]interface{}, 0)
		events := make([]map[string]interface{}, 0)
		audits := make([]map[string]interface{}, 0)
		for _, switchItem := range switchSet {
			if switchItem.Log != nil && switchItem.Log.Enable != nil && helper.PBool(switchItem.Log.Enable) {
				logAgent := map[string]interface{}{
					"enabled": helper.PBool(switchItem.Log.Enable),
				}
				logAgents = append(logAgents, logAgent)
			}
			if switchItem.Event != nil && switchItem.Event.Enable != nil && helper.PBool(switchItem.Event.Enable) {
				event := map[string]interface{}{
					"enabled":    helper.PBool(switchItem.Event.Enable),
					"log_set_id": helper.PString(switchItem.Event.LogsetId),
					"topic_id":   helper.PString(switchItem.Event.TopicId),
				}
				events = append(events, event)
			}
			if switchItem.Audit != nil && switchItem.Audit.Enable != nil && helper.PBool(switchItem.Audit.Enable) {
				audit := map[string]interface{}{
					"enabled":    helper.PBool(switchItem.Audit.Enable),
					"log_set_id": helper.PString(switchItem.Audit.LogsetId),
					"topic_id":   helper.PString(switchItem.Audit.TopicId),
				}
				audits = append(audits, audit)
			}
		}
		if len(logAgents) > 0 {
			_ = d.Set("log_agent", logAgents)
		}
		if len(events) > 0 {
			_ = d.Set("event_persistence", events)
		}
		if len(audits) > 0 {
			_ = d.Set("cluster_audit", audits)
		}

		resp, err := service.DescribeClusterExtraArgs(ctx, d.Id())
		if err != nil {
			return err
		}
		fmt.Println(&resp)
		flag := false
		extraArgs := make(map[string]interface{}, 0)
		if len(resp.KubeAPIServer) > 0 {
			flag = true
			extraArgs["kube_apiserver"] = resp.KubeAPIServer
		}
		if len(resp.KubeControllerManager) > 0 {
			flag = true
			extraArgs["kube_controller_manager"] = resp.KubeControllerManager
		}
		if len(resp.KubeScheduler) > 0 {
			flag = true
			extraArgs["kube_scheduler"] = resp.KubeScheduler
		}

		if flag {
			_ = d.Set("cluster_extra_args", []map[string]interface{}{extraArgs})
		}

		if networkType == TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY {
			resp, err := service.DescribeExternalNodeSupportConfig(ctx, d.Id())
			if err != nil {
				return err
			}
			_ = d.Set("cluster_subnet_id", resp.SubnetId)
		}

		if networkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
			if cluster.ClusterNetworkSettings != nil && cluster.ClusterNetworkSettings.SubnetId != nil {
				_ = d.Set("cluster_subnet_id", cluster.ClusterNetworkSettings.SubnetId)
			}

			resp, err := service.DescribeIPAMD(ctx, d.Id())
			if err != nil {
				return err
			}
			_ = d.Set("eni_subnet_ids", helper.PStrings(resp.SubnetIds))

			duration, err := time.ParseDuration(helper.PString(resp.ClaimExpiredDuration))
			if err != nil {
				return err
			}
			seconds := int(duration.Seconds())
			if seconds > 0 {
				_ = d.Set("claim_expired_seconds", seconds)
			}
		}

		if clusterInfo.DeployType == TKE_DEPLOY_TYPE_MANAGED {
			options, state, _, err := service.DescribeClusterAuthenticationOptions(ctx, d.Id())
			if err != nil {
				return err
			}
			if state == "Success" {
				authOptions := make(map[string]interface{}, 0)
				if helper.PBool(options.UseTKEDefault) {
					authOptions["use_tke_default"] = helper.PBool(options.UseTKEDefault)
				} else {
					authOptions["jwks_uri"] = helper.PString(options.JWKSURI)
					authOptions["issuer"] = helper.PString(options.Issuer)
				}
				authOptions["auto_create_discovery_anonymous_auth"] = helper.PBool(options.AutoCreateDiscoveryAnonymousAuth)
				_ = d.Set("auth_options", []map[string]interface{}{authOptions})
			}
		}
	}

	if _, ok := d.GetOkExists("auto_upgrade_cluster_level"); ok {
		_ = d.Set("auto_upgrade_cluster_level", clusterInfo.AutoUpgradeClusterLevel)
	} else if importClsFlag {
		_ = d.Set("auto_upgrade_cluster_level", clusterInfo.AutoUpgradeClusterLevel)
		importClsFlag = false
	}

	err = checkClusterEndpointStatus(ctx, &service, d, false)
	if err != nil {
		return fmt.Errorf("get internet failed, %s", err.Error())
	}

	err = checkClusterEndpointStatus(ctx, &service, d, true)
	if err != nil {
		return fmt.Errorf("get intranet failed, %s\n", err.Error())
	}
	return nil
}

func resourceTencentCloudKubernetesClusterReadRequestOnError1(ctx context.Context, resp *tke.DescribeClusterInstancesResponseParams, e error) *resource.RetryError {
	if e, ok := e.(*errors.TencentCloudSDKError); ok {
		if e.GetCode() == "InternalError.ClusterNotFound" {
			return nil
		}
	}
	return nil
}

func resourceTencentCloudKubernetesClusterReadRequestOnError2(ctx context.Context, resp *tke.DescribeClusterSecurityResponseParams, e error) *resource.RetryError {
	if e, ok := e.(*errors.TencentCloudSDKError); ok {
		if e.GetCode() == "InternalError.ClusterNotFound" {
			return nil
		}
	}
	return nil
}

func resourceTencentCloudKubernetesClusterReadPostHandleResponse2(ctx context.Context, resp *tke.DescribeClusterSecurityResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	policies := make([]string, 0, len(resp.SecurityPolicy))
	for _, v := range resp.SecurityPolicy {
		policies = append(policies, *v)
	}
	_ = d.Set("security_policy", policies)

	var globalConfig *tke.ClusterAsGroupOption
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var err error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		globalConfig, err = service.DescribeClusterNodePoolGlobalConfig(ctx, d.Id())
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
	if err != nil {
		return err
	}

	if globalConfig != nil {
		temp := make(map[string]interface{})
		temp["is_scale_in_enabled"] = globalConfig.IsScaleDownEnabled
		temp["expander"] = globalConfig.Expander
		temp["max_concurrent_scale_in"] = globalConfig.MaxEmptyBulkDelete
		temp["scale_in_delay"] = globalConfig.ScaleDownDelay
		temp["scale_in_unneeded_time"] = globalConfig.ScaleDownUnneededTime
		temp["scale_in_utilization_threshold"] = globalConfig.ScaleDownUtilizationThreshold
		temp["ignore_daemon_sets_utilization"] = globalConfig.IgnoreDaemonSetsUtilization
		temp["skip_nodes_with_local_storage"] = globalConfig.SkipNodesWithLocalStorage
		temp["skip_nodes_with_system_pods"] = globalConfig.SkipNodesWithSystemPods

		_ = d.Set("node_pool_global_config", []map[string]interface{}{temp})
	}
	return nil
}

func resourceTencentCloudKubernetesClusterUpdatePostFillRequest0(ctx context.Context, req *tke.ModifyClusterAttributeRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	clusterLevel := d.Get("cluster_level").(string)
	autoUpgradeClusterLevel := d.Get("auto_upgrade_cluster_level").(bool)

	ins, _, err := tkeService.DescribeCluster(ctx, id)
	if err != nil {
		return err
	}

	//ignore same cluster level if same
	if *ins.ClusterLevel == clusterLevel {
		clusterLevel = ""
	}

	if clusterLevel != "" {
		req.ClusterLevel = &clusterLevel
	}
	req.AutoUpgradeClusterLevel = &tke.AutoUpgradeClusterLevel{
		IsAutoUpgrade: &autoUpgradeClusterLevel,
	}
	return nil
}

func resourceTencentCloudKubernetesClusterUpdatePostFillRequest1(ctx context.Context, req *tke.UpdateClusterVersionRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	newVersion := d.Get("cluster_version").(string)
	isOk, err := tkeService.CheckClusterVersion(ctx, id, newVersion)
	if err != nil {
		return err
	}
	if !isOk {
		return fmt.Errorf("version %s is unsupported", newVersion)
	}
	extraArgs, ok := d.GetOk("cluster_extra_args")
	if !ok {
		extraArgs = nil
	}

	if extraArgs != nil && len(extraArgs.([]interface{})) > 0 {
		// the first elem is in use
		extraInterface := extraArgs.([]interface{})
		extraMap := extraInterface[0].(map[string]interface{})

		kas := make([]*string, 0)
		if kaArgs, exist := extraMap["kube_apiserver"]; exist {
			args := kaArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kas = append(kas, &str)
			}
		}
		kcms := make([]*string, 0)
		if kcmArgs, exist := extraMap["kube_controller_manager"]; exist {
			args := kcmArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kcms = append(kcms, &str)
			}
		}
		kss := make([]*string, 0)
		if ksArgs, exist := extraMap["kube_scheduler"]; exist {
			args := ksArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kss = append(kss, &str)
			}
		}

		req.ExtraArgs = &tke.ClusterExtraArgs{
			KubeAPIServer:         kas,
			KubeControllerManager: kcms,
			KubeScheduler:         kss,
		}
	}
	return nil
}

func resourceTencentCloudKubernetesClusterUpdatePostHandleResponse1(ctx context.Context, resp *tke.UpdateClusterVersionResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	// todo check status
	err := resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ins, has, inErr := tkeService.DescribeCluster(ctx, id)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("Cluster %s is not exist", id))
		}
		if ins.ClusterStatus == "Running" {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cluster %s status %s, retry...", id, ins.ClusterStatus))
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesClusterDeletePostFillRequest0(ctx context.Context, req *tke.DeleteClusterRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	deleteEventLogSetAndTopic := false
	enableEventLog := false
	deleteAuditLogSetAndTopic := false
	if v, ok := helper.InterfacesHeadMap(d, "event_persistence"); ok {
		deleteEventLogSetAndTopic = v["delete_event_log_and_topic"].(bool)
		// get cluster current enabled status
		enableEventLog = v["enabled"].(bool)
	}

	if v, ok := helper.InterfacesHeadMap(d, "cluster_audit"); ok {
		deleteAuditLogSetAndTopic = v["delete_audit_log_and_topic"].(bool)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if deleteEventLogSetAndTopic && enableEventLog {
			err := service.SwitchEventPersistence(ctx, d.Id(), "", "", false, true)
			if err != nil {
				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() != "FailedOperation.ClusterNotFound" {
						return tccommon.RetryError(err, tccommon.InternalError)
					}
				}
				return tccommon.RetryError(err, tccommon.InternalError)
			}
		}
		if deleteAuditLogSetAndTopic {
			err := service.SwitchClusterAudit(ctx, d.Id(), "", "", false, true)
			if err != nil {
				if e, ok := err.(*errors.TencentCloudSDKError); ok {
					if e.GetCode() != "ResourceNotFound.ClusterNotFound" {
						return tccommon.RetryError(err, tccommon.InternalError)
					}
				}
				return tccommon.RetryError(err, tccommon.InternalError)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesClusterDeleteRequestOnError0(ctx context.Context, e error) *resource.RetryError {
	if err, ok := e.(*errors.TencentCloudSDKError); ok {
		if err.GetCode() == "InternalError.ClusterNotFound" {
			return nil
		}
	}
	return tccommon.RetryError(e)
}

func resourceTencentCloudKubernetesClusterDeletePostHandleResponse0(ctx context.Context, resp *tke.DeleteClusterResponse) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	_, _, err := service.DescribeClusterInstances(ctx, d.Id())

	if err != nil {
		err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, _, err = service.DescribeClusterInstances(ctx, d.Id())
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InvalidParameter.ClusterNotFound" {
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

func resourceTencentCloudKubernetesClusterUpdateOnStart(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)
	id := d.Id()

	if d.HasChange("cluster_subnet_id") {
		return fmt.Errorf("argument cluster_subnet_id cannot be changed")
	}

	if d.HasChange("tags") {
		if err := modifyClusterTags(ctx); err != nil {
			return err
		}

		// wait for tags ok
		err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			request := tke.NewDescribeBatchModifyTagsStatusRequest()
			request.ClusterId = &id
			resp, errRet := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeBatchModifyTagsStatus(request)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			// TaskFailed = "failed"; TaskRunning = "running"; TaskDone = "done"
			if resp != nil && *resp.Response.Status == "done" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("modify tags status is %s, retry...", *resp.Response.Status))
		})
		if err != nil {
			return err
		}

	}

	var (
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
	)

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}

	if d.HasChange("cluster_intranet_subnet_id") && !d.HasChange("cluster_intranet") {
		return fmt.Errorf("`cluster_intranet_subnet_id` must modified with `cluster_intranet`")
	}

	if d.HasChange("cluster_internet_security_group") && !d.HasChange("cluster_internet") {
		if clusterInternet {
			err := tkeService.ModifyClusterEndpointSG(ctx, id, clusterInternetSecurityGroup)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("cluster_intranet") {
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, clusterIntranet, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}

	}

	if d.HasChange("cluster_internet") {
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, clusterInternet, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
	}

	// situation when only domain changed
	if !d.HasChange("cluster_intranet") && clusterIntranet && d.HasChange("cluster_intranet_domain") {
		// recreate the cluster intranet endpoint using new domain
		// first close
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, TKE_CLUSTER_CLOSE_ACCESS, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}
		// then reopen
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, TKE_CLUSTER_OPEN_ACCESS, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}
	}
	if !d.HasChange("cluster_internet") && clusterInternet && d.HasChange("cluster_internet_domain") {
		// recreate the cluster internet endpoint using new domain
		// first close
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, TKE_CLUSTER_CLOSE_ACCESS, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
		// then reopen
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, TKE_CLUSTER_OPEN_ACCESS, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
	}

	//update VPC-CNI container network capability
	if !d.HasChange("eni_subnet_ids") && (d.HasChange("vpc_cni_type") || d.HasChange("claim_expired_seconds")) {
		err := fmt.Errorf("changing only `vpc_cni_type` or `claim_expired_seconds` is not supported, when turning on or off the vpc-cni container network capability, `eni_subnet_ids` must be changed")
		return err
	}
	if d.HasChange("eni_subnet_ids") {
		eniSubnetIdList := d.Get("eni_subnet_ids").([]interface{})
		if len(eniSubnetIdList) == 0 {
			err := tkeService.DisableVpcCniNetworkType(ctx, id)
			if err != nil {
				return err
			}
			time.Sleep(3 * time.Second)
			err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				ipamdResp, inErr := tkeService.DescribeIPAMD(ctx, id)
				enableIPAMD := *ipamdResp.EnableIPAMD
				disableVpcCniMode := *ipamdResp.DisableVpcCniMode
				phase := *ipamdResp.Phase
				if inErr != nil {
					return resource.NonRetryableError(inErr)
				}
				if !enableIPAMD || (disableVpcCniMode && phase != "upgrading") {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("%s close vpc cni network type task is in progress and waiting to be completed", id))
			})
			if err != nil {
				return err
			}
		} else {
			info, _, err := tkeService.DescribeCluster(ctx, id)
			if err != nil {
				err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					newInfo, _, inErr := tkeService.DescribeCluster(ctx, id)
					if inErr != nil {
						return tccommon.RetryError(inErr)
					}
					info = newInfo
					return nil
				})
				if err != nil {
					return err
				}
			}
			oldSubnets := info.EniSubnetIds
			var subnets []string
			for index := range eniSubnetIdList {
				subnetId := eniSubnetIdList[index].(string)
				subnets = append(subnets, subnetId)
			}
			if len(oldSubnets) > 0 {
				exist, addSubnets := helper.CheckElementsExist(oldSubnets, subnets)
				if !exist {
					err = fmt.Errorf("the `eni_subnet_ids` parameter does not allow modification of existing subnet ID data %v. "+
						"if you want to modify the existing subnet ID, please first set eni_subnet_ids to empty to turn off the VPC-CNI network capability, "+
						"and then fill in the latest subnet ID", oldSubnets)
					return err
				}
				if d.HasChange("vpc_cni_type") || d.HasChange("claim_expired_seconds") {
					err = fmt.Errorf("modifying `vpc_cni_type` and `claim_expired_seconds` is not supported when adding a cluster subnet")
					return err
				}
				if len(addSubnets) > 0 {
					vpcId := d.Get("vpc_id").(string)
					err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						inErr := tkeService.AddVpcCniSubnets(ctx, id, addSubnets, vpcId)
						if inErr != nil {
							return resource.NonRetryableError(inErr)
						}
						return nil
					})
					if err != nil {
						return err
					}
				}
			} else {
				var vpcCniType string
				if v, ok := d.GetOk("vpc_cni_type"); ok {
					vpcCniType = v.(string)
				} else {
					vpcCniType = "tke-route-eni"
				}
				enableStaticIp := !d.Get("is_non_static_ip_mode").(bool)
				expiredSeconds := uint64(d.Get("claim_expired_seconds").(int))

				err = tkeService.EnableVpcCniNetworkType(ctx, id, vpcCniType, enableStaticIp, subnets, expiredSeconds)
				if err != nil {
					return err
				}
				time.Sleep(3 * time.Second)
				err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					ipamdResp, inErr := tkeService.DescribeIPAMD(ctx, id)
					disableVpcCniMode := *ipamdResp.DisableVpcCniMode
					phase := *ipamdResp.Phase
					if inErr != nil {
						return resource.NonRetryableError(inErr)
					}
					if !disableVpcCniMode && phase == "running" {
						return nil
					}
					if !disableVpcCniMode && phase == "initializing" {
						return resource.RetryableError(fmt.Errorf("%s enable vpc cni network type task is in progress and waiting to be completed", id))
					}
					return resource.NonRetryableError(fmt.Errorf("%s enable vpc cni network type task disableVpcCniMode is %v and phase is %s,we won't wait for it finish", id, disableVpcCniMode, phase))
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func resourceTencentCloudKubernetesClusterUpdateOnExit(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)

	tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if d.HasChange("auth_options") {
		request := tkeGetAuthOptions(d, id)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := tkeService.ModifyClusterAuthenticationOptions(ctx, request)
			if inErr != nil {
				return tccommon.RetryError(inErr, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
			return nil
		})
		if err != nil {
			return err
		}
		_, _, err = tkeService.WaitForAuthenticationOptionsUpdateSuccess(ctx, id)
		if err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") {
		enable := d.Get("deletion_protection").(bool)
		if err := tkeService.ModifyDeletionProtection(ctx, id, enable); err != nil {
			return err
		}

	}

	if d.HasChange("acquire_cluster_admin_role") {
		o, n := d.GetChange("acquire_cluster_admin_role")
		if o.(bool) && !n.(bool) {
			return fmt.Errorf("argument `acquire_cluster_admin_role` cannot set to false")
		}
		_, err := tkeService.AcquireClusterAdminRole(ctx, id)
		if err != nil {
			return err
		}
	}

	if d.HasChange("log_agent") {
		v, ok := helper.InterfacesHeadMap(d, "log_agent")
		enabled := false
		rootDir := ""
		if ok {
			rootDir = v["kubelet_root_dir"].(string)
			enabled = v["enabled"].(bool)
		}
		err := tkeService.SwitchLogAgent(ctx, id, rootDir, enabled)
		if err != nil {
			return err
		}
	}

	if d.HasChange("event_persistence") {
		v, ok := helper.InterfacesHeadMap(d, "event_persistence")
		enabled := false
		logSetId := ""
		topicId := ""
		deleteEventLog := false
		if ok {
			enabled = v["enabled"].(bool)
			logSetId = v["log_set_id"].(string)
			topicId = v["topic_id"].(string)
			deleteEventLog = v["delete_event_log_and_topic"].(bool)
		}

		err := tkeService.SwitchEventPersistence(ctx, id, logSetId, topicId, enabled, deleteEventLog)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cluster_audit") {
		v, ok := helper.InterfacesHeadMap(d, "cluster_audit")
		enabled := false
		logSetId := ""
		topicId := ""
		deleteAuditLog := false
		if ok {
			enabled = v["enabled"].(bool)
			logSetId = v["log_set_id"].(string)
			topicId = v["topic_id"].(string)
			deleteAuditLog = v["delete_audit_log_and_topic"].(bool)
		}

		err := tkeService.SwitchClusterAudit(ctx, id, logSetId, topicId, enabled, deleteAuditLog)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return nil
}

func unschedulableDiffSuppressFunc(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if newValue == "0" && oldValue == "" {
		return true
	} else {
		return oldValue == newValue
	}
}

func dockerGraphPathDiffSuppressFunc(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if newValue == "/var/lib/docker" && oldValue == "" || oldValue == "/var/lib/docker" && newValue == "" {
		return true
	} else {
		return oldValue == newValue
	}
}

func clusterCidrValidateFunc(v interface{}, k string) (ws []string, errs []error) {
	value := v.(string)
	if value == "" {
		return
	}
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}
	if ipnet == nil || value != ipnet.String() {
		errs = append(errs, fmt.Errorf("%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
		return
	}
	if !strings.Contains(value, "/") {
		errs = append(errs, fmt.Errorf("%q must be a network segment", k))
		return
	}
	if !strings.HasPrefix(value, "9.") && !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "11.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
		errs = append(errs, fmt.Errorf("%q must in 9. | 10. | 11. | 192.168. | 172.[16-31]", k))
		return
	}

	if strings.HasPrefix(value, "172.") {
		nextNo := strings.Split(value, ".")[1]
		no, _ := strconv.ParseInt(nextNo, 10, 64)
		if no < 16 || no > 31 {
			errs = append(errs, fmt.Errorf("%q must in 9.0 | 10. | 11. | 192.168. | 172.[16-31]", k))
			return
		}
	}
	return
}

func serviceCidrValidateFunc(v interface{}, k string) (ws []string, errs []error) {
	value := v.(string)
	if value == "" {
		return
	}
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errs = append(errs, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}
	if ipnet == nil || value != ipnet.String() {
		errs = append(errs, fmt.Errorf("%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
		return
	}
	if !strings.Contains(value, "/") {
		errs = append(errs, fmt.Errorf("%q must be a network segment", k))
		return
	}
	if !strings.HasPrefix(value, "9.") && !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
		errs = append(errs, fmt.Errorf("%q must in 9. | 10. | 192.168. | 172.[16-31]", k))
		return
	}

	if strings.HasPrefix(value, "172.") {
		nextNo := strings.Split(value, ".")[1]
		no, _ := strconv.ParseInt(nextNo, 10, 64)
		if no < 16 || no > 31 {
			errs = append(errs, fmt.Errorf("%q must in 9. | 10. | 192.168. | 172.[16-31]", k))
			return
		}
	}
	return
}

func claimExpiredSecondsValidateFunc(v interface{}, k string) (ws []string, errs []error) {
	value := v.(int)
	if value < 300 || value > 15768000 {
		errs = append(errs, fmt.Errorf("%q must greater or equal than 300 and less than 15768000", k))
		return
	}
	return
}

func ResourceTkeGetAddonsDiffs(o, n []interface{}) (adds, removes, changes []interface{}) {
	indexByName := func(i interface{}) int {
		v := i.(map[string]interface{})
		return helper.HashString(v["name"].(string))
	}
	indexAll := func(i interface{}) int {
		v := i.(map[string]interface{})
		name := v["name"].(string)
		param := v["param"].(string)
		return helper.HashString(fmt.Sprintf("%s#%s", name, param))
	}

	os := schema.NewSet(indexByName, o)
	ns := schema.NewSet(indexByName, n)

	adds = ns.Difference(os).List()
	removes = os.Difference(ns).List()

	fullIndexedKeeps := schema.NewSet(indexAll, ns.Intersection(os).List())
	fullIndexedOlds := schema.NewSet(indexAll, o)

	changes = fullIndexedKeeps.Difference(fullIndexedOlds).List()
	return
}

func modifyClusterTags(ctx context.Context) error {
	d := tccommon.ResourceDataFromContext(ctx)
	meta := tccommon.ProviderMetaFromContext(ctx)
	logId := tccommon.GetLogId(ctx)

	id := d.Id()
	tags := GetTkeTags(d, "tags")
	body := map[string]interface{}{
		"ClusterId":       id,
		"SyncSubresource": false,
		"Tags":            tags,
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("tke")
	request := tchttp.NewCommonRequest("tke", "2018-05-25", "ModifyClusterTags")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = client.Send(request, response)
	if err != nil {
		fmt.Printf("Modify Cluster Tags failed: %v \n", err)
		return err
	}
	reqBody, _ := request.MarshalJSON()
	respBody := response.GetBody()

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), string(reqBody), string(respBody))
	return nil
}

// upgradeClusterInstances upgrade instances, upgrade type try seq:major, hot.
func upgradeClusterInstances(tkeService TkeService, ctx context.Context, id string) error {
	// get all available instances for upgrade
	upgradeType := "major"
	instanceIds, err := tkeService.CheckInstancesUpgradeAble(ctx, id, upgradeType)
	if err != nil {
		return err
	}
	if len(instanceIds) == 0 {
		upgradeType = "hot"
		instanceIds, err = tkeService.CheckInstancesUpgradeAble(ctx, id, upgradeType)
		if err != nil {
			return err
		}
	}
	log.Println("instancesIds for upgrade:", instanceIds)
	instNum := len(instanceIds)
	if instNum == 0 {
		return nil
	}

	// upgrade instances
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := tkeService.UpgradeClusterInstances(ctx, id, upgradeType, instanceIds)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// check update status: upgrade instance one by one, so timeout depend on instance number.
	timeout := tccommon.ReadRetryTimeout * time.Duration(instNum)
	err = resource.Retry(timeout, func() *resource.RetryError {
		done, inErr := tkeService.GetUpgradeInstanceResult(ctx, id)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if done {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cluster %s, retry...", id))
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func tkeGetCvmRunInstancesPara(dMap map[string]interface{}, meta interface{},
	vpcId string, projectId int64) (cvmJson string, count int64, errRet error) {

	request := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	request.Placement = &place

	place.ProjectId = &projectId

	if v, ok := dMap["availability_zone"]; ok {
		place.Zone = helper.String(v.(string))
	}

	if v, ok := dMap["instance_type"]; ok {
		request.InstanceType = helper.String(v.(string))
	} else {
		errRet = fmt.Errorf("instance_type must be set.")
		return
	}

	subnetId := ""

	if v, ok := dMap["subnet_id"]; ok {
		subnetId = v.(string)
	}

	if (vpcId == "" && subnetId != "") ||
		(vpcId != "" && subnetId == "") {
		errRet = fmt.Errorf("Parameters cvm.`subnet_id` and cluster.`vpc_id` are both set or neither")
		return
	}

	if vpcId != "" {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
			VpcId:    &vpcId,
			SubnetId: &subnetId,
		}
	}

	if v, ok := dMap["system_disk_type"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskSize = helper.Int64(int64(v.(int)))

	}

	if v, ok := dMap["cam_role_name"]; ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := dMap["data_disk"]; ok {

		dataDisks := v.([]interface{})
		request.DataDisks = make([]*cvm.DataDisk, 0, len(dataDisks))

		for _, d := range dataDisks {

			var (
				value      = d.(map[string]interface{})
				diskType   = value["disk_type"].(string)
				diskSize   = int64(value["disk_size"].(int))
				snapshotId = value["snapshot_id"].(string)
				encrypt    = value["encrypt"].(bool)
				kmsKeyId   = value["kms_key_id"].(string)
				dataDisk   = cvm.DataDisk{
					DiskType: &diskType,
				}
			)
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if encrypt {
				dataDisk.Encrypt = &encrypt
			}
			if kmsKeyId != "" {
				dataDisk.KmsKeyId = &kmsKeyId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if v, ok := dMap["internet_charge_type"]; ok {

		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetMaxBandwidthOut = helper.Int64(int64(v.(int)))
	}

	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}

	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	if v, ok := dMap["password"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}

		if v.(string) != "" {
			request.LoginSettings.Password = helper.String(v.(string))
		}
	}

	if v, ok := dMap["instance_name"]; ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := dMap["key_ids"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}
		keyIds := v.([]interface{})

		if len(keyIds) != 0 {
			request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
			for i := range keyIds {
				if keyId, ok := keyIds[i].(string); ok && keyId != "" {
					request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
				}
			}
		}
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters cvm.`key_ids` and cluster.`password` should be set one")
		return
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters cvm.`key_ids` and cluster.`password` can only be supported one")
		return
	}

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroup, ok := securityGroups[i].(string); ok && securityGroup != "" {
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
			}
		}
	}

	if v, ok := dMap["disaster_recover_group_ids"]; ok {
		disasterGroups := v.([]interface{})
		request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterGroups))
		for i := range disasterGroups {
			if disasterGroup, ok := disasterGroups[i].(string); ok && disasterGroup != "" {
				request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterGroup)
			}
		}
	}

	if v, ok := dMap["enhanced_security_service"]; ok {

		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}

		securityService := v.(bool)
		request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}
	if v, ok := dMap["user_data"]; ok {
		request.UserData = helper.String(v.(string))
	}
	if v, ok := dMap["instance_charge_type"]; ok {
		instanceChargeType := v.(string)
		request.InstanceChargeType = &instanceChargeType
		if instanceChargeType == svccvm.CVM_CHARGE_TYPE_PREPAID || instanceChargeType == svccvm.CVM_CHARGE_TYPE_UNDERWRITE {
			request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
			if period, ok := dMap["instance_charge_type_prepaid_period"]; ok {
				periodInt64 := int64(period.(int))
				request.InstanceChargePrepaid.Period = &periodInt64
			} else {
				errRet = fmt.Errorf("instance charge type prepaid period can not be empty when charge type is %s",
					instanceChargeType)
				return
			}
			if renewFlag, ok := dMap["instance_charge_type_prepaid_renew_flag"].(string); ok && renewFlag != "" {
				request.InstanceChargePrepaid.RenewFlag = helper.String(renewFlag)
			}
		}
	}
	if v, ok := dMap["count"]; ok {
		count = int64(v.(int))
	} else {
		count = 1
	}
	request.InstanceCount = &count

	if v, ok := dMap["hostname"]; ok {
		hostname := v.(string)
		if hostname != "" {
			request.HostName = &hostname
		}
	}

	if v, ok := dMap["img_id"]; ok && v.(string) != "" {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := dMap["hpc_cluster_id"]; ok && v.(string) != "" {
		request.HpcClusterId = helper.String(v.(string))
	}

	if v, ok := dMap["tags"].([]interface{}); ok && len(v) != 0 {
		tmpTagSpec := cvm.TagSpecification{}
		tmpTagSpec.ResourceType = helper.String("instance")
		for _, item := range v {
			value := item.(map[string]interface{})
			tmpTag := cvm.Tag{}
			if v, ok := value["key"].(string); ok && v != "" {
				tmpTag.Key = &v
			}

			if v, ok := value["value"].(string); ok && v != "" {
				tmpTag.Value = &v
			}

			tmpTagSpec.Tags = append(tmpTagSpec.Tags, &tmpTag)
		}

		request.TagSpecification = append(request.TagSpecification, &tmpTagSpec)
	}

	if v, ok := dMap["cdc_id"]; ok && v.(string) != "" {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	cvmJson = request.ToJsonString()

	cvmJson = strings.Replace(cvmJson, `"Password":"",`, "", -1)

	return
}

func tkeGetCvmExistInstancesPara(dMap map[string]interface{}) (tke.ExistedInstancesForNode, error) {
	inst := tke.ExistedInstancesForNode{}
	if temp, ok := dMap["node_role"]; ok {
		nodeRole := temp.(string)
		inst.NodeRole = &nodeRole
	}

	if temp, ok := dMap["instances_para"]; ok {
		paras := temp.([]interface{})
		if len(paras) > 0 {
			paraMap := paras[0].(map[string]interface{})
			inst.ExistedInstancesPara = &tke.ExistedInstancesPara{}
			loginSettings := &tke.LoginSettings{}
			enhancedService := &tke.EnhancedService{}

			if v, ok := paraMap["instance_ids"]; ok && len(v.([]interface{})) > 0 {
				insIDs := v.([]interface{})
				inst.ExistedInstancesPara.InstanceIds = make([]*string, 0, len(insIDs))
				for _, v := range insIDs {
					inst.ExistedInstancesPara.InstanceIds = append(inst.ExistedInstancesPara.InstanceIds, helper.String(v.(string)))
				}
			}

			if v, ok := paraMap["security_group_ids"]; ok && len(v.([]interface{})) > 0 {
				sgIds := v.([]interface{})
				inst.ExistedInstancesPara.SecurityGroupIds = make([]*string, 0, len(sgIds))
				for i := range sgIds {
					if sgId, ok := sgIds[i].(string); ok && sgId != "" {
						inst.ExistedInstancesPara.SecurityGroupIds = append(inst.ExistedInstancesPara.SecurityGroupIds, &sgId)
					}
				}
			}

			if v, ok := paraMap["password"]; ok {
				loginSettings.Password = helper.String(v.(string))
				inst.ExistedInstancesPara.LoginSettings = loginSettings
			}

			if v, ok := paraMap["key_ids"]; ok && len(v.([]interface{})) > 0 {
				keyIds := v.([]interface{})
				loginSettings.KeyIds = make([]*string, 0, len(keyIds))
				for i := range keyIds {
					if keyId, ok := keyIds[i].(string); ok && keyId != "" {
						loginSettings.KeyIds = append(loginSettings.KeyIds, &keyId)
					}
				}

				inst.ExistedInstancesPara.LoginSettings = loginSettings
			}

			if v, ok := paraMap["enhanced_security_service"]; ok {
				enhancedService.SecurityService = &tke.RunSecurityServiceEnabled{Enabled: helper.Bool(v.(bool))}
				inst.ExistedInstancesPara.EnhancedService = enhancedService
			}

			if v, ok := paraMap["enhanced_monitor_service"]; ok {
				enhancedService.MonitorService = &tke.RunMonitorServiceEnabled{Enabled: helper.Bool(v.(bool))}
				inst.ExistedInstancesPara.EnhancedService = enhancedService
			}

			if v, ok := paraMap["master_config"]; ok && len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					instanceAdvancedSettingsOverridesMap := item.(map[string]interface{})
					instanceAdvancedSettings := tke.InstanceAdvancedSettings{}
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
							labels := tke.Label{}
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
							dataDisk := tke.DataDisk{}
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
							args := tke.InstanceExtraArgs{}
							if v, ok := extraArgsMap["kubelet"]; ok {
								args.Kubelet = helper.InterfacesStringsPoint(v.([]interface{}))
							}

							instanceAdvancedSettings.ExtraArgs = &args
						}
					}

					if v, ok := instanceAdvancedSettingsOverridesMap["desired_pod_number"]; ok {
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

						instanceAdvancedSettings.GPUArgs = &tkeGpuArgs
					}

					if v, ok := instanceAdvancedSettingsOverridesMap["taints"]; ok && len(v.([]interface{})) > 0 {
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

							instanceAdvancedSettings.Taints = append(instanceAdvancedSettings.Taints, &taint)
						}
					}

					inst.InstanceAdvancedSettingsOverride = &instanceAdvancedSettings
				}
			}
		}
	}

	if temp, ok := dMap["desired_pod_numbers"]; ok {
		inst.DesiredPodNumbers = make([]*int64, 0)
		podNums := temp.([]interface{})
		for _, v := range podNums {
			inst.DesiredPodNumbers = append(inst.DesiredPodNumbers, helper.Int64(int64(v.(int))))
		}
	}

	return inst, nil
}

func tkeGetNodePoolGlobalConfig(d *schema.ResourceData) *tke.ModifyClusterAsGroupOptionAttributeRequest {
	request := tke.NewModifyClusterAsGroupOptionAttributeRequest()
	request.ClusterId = helper.String(d.Id())

	clusterAsGroupOption := &tke.ClusterAsGroupOption{}
	if v, ok := d.GetOkExists("node_pool_global_config.0.is_scale_in_enabled"); ok {
		clusterAsGroupOption.IsScaleDownEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.expander"); ok {
		clusterAsGroupOption.Expander = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.max_concurrent_scale_in"); ok {
		clusterAsGroupOption.MaxEmptyBulkDelete = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_delay"); ok {
		clusterAsGroupOption.ScaleDownDelay = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_unneeded_time"); ok {
		clusterAsGroupOption.ScaleDownUnneededTime = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_utilization_threshold"); ok {
		clusterAsGroupOption.ScaleDownUtilizationThreshold = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.ignore_daemon_sets_utilization"); ok {
		clusterAsGroupOption.IgnoreDaemonSetsUtilization = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.skip_nodes_with_local_storage"); ok {
		clusterAsGroupOption.SkipNodesWithLocalStorage = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.skip_nodes_with_system_pods"); ok {
		clusterAsGroupOption.SkipNodesWithSystemPods = helper.Bool(v.(bool))
	}

	request.ClusterAsGroupOption = clusterAsGroupOption
	return request
}

func tkeGetAuthOptions(d *schema.ResourceData, clusterId string) *tke.ModifyClusterAuthenticationOptionsRequest {
	raw, ok := d.GetOk("auth_options")
	options := raw.([]interface{})

	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = helper.String(clusterId)
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{
		AutoCreateDiscoveryAnonymousAuth: helper.Bool(false),
	}

	if !ok || len(options) == 0 {
		request.ServiceAccounts.JWKSURI = helper.String("")
		return request
	}

	option := options[0].(map[string]interface{})

	if v, ok := option["auto_create_discovery_anonymous_auth"]; ok {
		request.ServiceAccounts.AutoCreateDiscoveryAnonymousAuth = helper.Bool(v.(bool))
	}

	if v, ok := option["use_tke_default"]; ok && v.(bool) {
		request.ServiceAccounts.UseTKEDefault = helper.Bool(true)
	} else {
		if v, ok := option["issuer"]; ok {
			request.ServiceAccounts.Issuer = helper.String(v.(string))
		}

		if v, ok := option["jwks_uri"]; ok {
			request.ServiceAccounts.JWKSURI = helper.String(v.(string))
		}
	}

	return request
}

func checkClusterEndpointStatus(ctx context.Context, service *TkeService, d *schema.ResourceData, isInternet bool) (err error) {
	var status, config string
	var response tke.DescribeClusterEndpointsResponseParams
	var isOpened bool
	var errRet error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		status, _, errRet = service.DescribeClusterEndpointStatus(ctx, d.Id(), isInternet)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if status == TkeInternetStatusCreating || status == TkeInternetStatusDeleting {
			return resource.RetryableError(
				fmt.Errorf("%s create cluster internet endpoint status still is %s", d.Id(), status))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted {
		isOpened = false
	}
	if status == TkeInternetStatusCreated {
		isOpened = true
	}
	if isInternet {
		_ = d.Set("cluster_internet", isOpened)
	} else {
		_ = d.Set("cluster_intranet", isOpened)
	}

	if isOpened {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			config, errRet = service.DescribeClusterConfig(ctx, d.Id(), isInternet)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, errRet = service.DescribeClusterEndpoints(ctx, d.Id())
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if isInternet {
			_ = d.Set("kube_config", config)
			_ = d.Set("cluster_internet_domain", helper.PString(response.ClusterExternalDomain))
			_ = d.Set("cluster_internet_security_group", helper.PString(response.SecurityGroup))
		} else {
			_ = d.Set("kube_config_intranet", config)
			_ = d.Set("cluster_intranet_domain", helper.PString(response.ClusterIntranetDomain))
			_ = d.Set("cluster_intranet_subnet_id", helper.PString(response.ClusterIntranetSubnetId))
		}

	} else {
		if isInternet {
			_ = d.Set("kube_config", "")
		} else {
			_ = d.Set("kube_config_intranet", "")
		}
	}
	return nil
}

func tkeCvmState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the cvm.",
		},
		"instance_role": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Role of the cvm.",
		},
		"instance_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the cvm.",
		},
		"failed_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Information of the cvm when it is failed.",
		},
		"lan_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "LAN IP of the cvm.",
		},
	}
}

//func tkeSecurityInfo() map[string]*schema.Schema {
//	return map[string]*schema.Schema{
//		"user_name": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "User name of account.",
//		},
//		"password": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "Password of account.",
//		},
//		"certification_authority": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "The certificate used for access.",
//		},
//		"cluster_external_endpoint": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "External network address to access.",
//		},
//		"domain": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "Domain name for access.",
//		},
//		"pgw_endpoint": {
//			Type:        schema.TypeString,
//			Computed:    true,
//			Description: "The Intranet address used for access.",
//		},
//		"security_policy": {
//			Type:        schema.TypeList,
//			Computed:    true,
//			Elem:        &schema.Schema{Type: schema.TypeString},
//			Description: "Access policy.",
//		},
//	}
//}

func TkeCvmCreateInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"count": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1,
			Description: "Number of cvm.",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "sub machine of tke",
			Description: "Name of the CVMs.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Specified types of CVM instance.",
		},
		// payment
		"instance_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      svccvm.CVM_CHARGE_TYPE_POSTPAID,
			ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
			Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.",
		},
		"instance_charge_type_prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1,
			ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
			Description:  "The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
		},
		"instance_charge_type_prepaid_renew_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
			Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
			Description:  "Private network ID.",
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      svcas.SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
			Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      50,
			ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
			Description:  "Volume of system disk in GB. Default is `50`.",
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
						Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Data disk snapshot ID.",
					},
					"encrypt": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Indicates whether to encrypt data disk, default `false`.",
					},
					"kms_key_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "ID of the custom CMK in the format of UUID or `kms-abcd1234`. This parameter is used to encrypt cloud disks.",
					},
					"file_system": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "File system, e.g. `ext3/ext4/xfs`.",
					},
					"auto_format_and_mount": {
						Type:        schema.TypeBool,
						ForceNew:    true,
						Optional:    true,
						Default:     false,
						Description: "Indicate whether to auto format and mount or not. Default is `false`.",
					},
					"mount_target": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Mount target.",
					},
					"disk_partition": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "The name of the device or partition to mount.",
					},
				},
			},
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      svcas.INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is 0.",
		},
		"bandwidth_package_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Specify whether to assign an Internet IP address.",
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
			Description: "ID list of keys, should be set if `password` not set.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
		"user_data": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "ase64-encoded User Data text, the length limit is 16KB.",
		},
		"cam_role_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "CAM role name authorized to access.",
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
		"disaster_recover_group_ids": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Disaster recover groups to which a CVM instance belongs. Only support maximum 1.",
		},
		"img_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: tccommon.ValidateImageID,
			Description:  "The valid image id, format of img-xxx.",
		},
		// InstanceAdvancedSettingsOverrides
		"desired_pod_num": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
			Default:  DefaultDesiredPodNum,
			Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, " +
				"and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.",
		},
		"hpc_cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of cvm hpc cluster.",
		},
	}
}
