package tke

import (
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

const (
	TKE_CLUSTER_OS_CENTOS72 = "centos7.2x86_64"
	TKE_CLUSTER_OS_CENTOS76 = "centos7.6.0_x64"
	TKE_CLUSTER_OS_UBUNTU18 = "ubuntu18.04.1x86_64"
	TKE_CLUSTER_OS_LINUX24  = "tlinux2.4x86_64"
	TKE_CLUSTER_OS_LINUX22  = "tlinux2.2(tkernel3)x86_64"
	TKE_CLUSTER_OS_LINUXF22 = "Tencent tlinux release 2.2 (Final)"
)

// 兼容旧的 cluster_os 定义
const (
	TkeClusterOsCentOS76 = "centos7.6x86_64"
	//TkeClusterOsUbuntu16 = "ubuntu16.04.1 LTSx86_64"
	TkeClusterOsUbuntu18 = "ubuntu18.04.1 LTSx86_64"
)

var TKE_CLUSTER_OS = []string{TKE_CLUSTER_OS_CENTOS76, TKE_CLUSTER_OS_UBUNTU18, TKE_CLUSTER_OS_LINUX24}

var tkeClusterOsMap = map[string]string{TKE_CLUSTER_OS_CENTOS72: TKE_CLUSTER_OS_CENTOS72,
	TKE_CLUSTER_OS_CENTOS76: TKE_CLUSTER_OS_CENTOS76,
	TKE_CLUSTER_OS_UBUNTU18: TKE_CLUSTER_OS_UBUNTU18,
	TKE_CLUSTER_OS_LINUX24:  TKE_CLUSTER_OS_LINUX24,
	TKE_CLUSTER_OS_LINUX22:  TKE_CLUSTER_OS_LINUX22,
	TKE_CLUSTER_OS_LINUXF22: TKE_CLUSTER_OS_LINUXF22,
}

func tkeToShowClusterOs(apiOs string) string {
	for showName, apiName := range tkeClusterOsMap {
		if apiName == apiOs {
			return showName
		}
	}
	return apiOs
}

const (
	TKE_DEPLOY_TYPE_MANAGED     = "MANAGED_CLUSTER"
	TKE_DEPLOY_TYPE_INDEPENDENT = "INDEPENDENT_CLUSTER"
)

var TKE_DEPLOY_TYPES = []string{TKE_DEPLOY_TYPE_MANAGED, TKE_DEPLOY_TYPE_INDEPENDENT}

const (
	TKE_RUNTIME_DOCKER     = "docker"
	TKE_RUNTIME_CONTAINERD = "containerd"
)

var TKE_RUNTIMES = []string{TKE_RUNTIME_DOCKER, TKE_RUNTIME_CONTAINERD}

const (
	TKE_ROLE_MASTER_ETCD = "MASTER_ETCD"
	TKE_ROLE_WORKER      = "WORKER"
)

var TKE_INSTANCE_CHARGE_TYPE = []string{svccvm.CVM_CHARGE_TYPE_PREPAID, svccvm.CVM_CHARGE_TYPE_POSTPAID}

const (
	TKE_CLUSTER_OS_TYPE_DOCKER_CUSTOMIZE = "DOCKER_CUSTOMIZE"
	TKE_CLUSTER_OS_TYPE_GENERAL          = "GENERAL"
)

var TKE_CLUSTER_OS_TYPES = []string{TKE_CLUSTER_OS_TYPE_GENERAL}

const (
	TkeInternetStatusCreating      = "Creating"
	TkeInternetStatusCreateFailed  = "CreateFailed"
	TkeInternetStatusCreated       = "Created"
	TkeInternetStatusDeleting      = "Deleting"
	TkeInternetStatusDeleted       = "Deleted"
	TkeInternetStatusDeletedFailed = "DeletedFailed"
	TkeInternetStatusNotfound      = "NotFound"
)

const (
	TKE_CLUSTER_NETWORK_TYPE_GR             = "GR"
	TKE_CLUSTER_NETWORK_TYPE_VPC_CNI        = "VPC-CNI"
	TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY = "CiliumOverlay"
)

const (
	TKE_CLUSTER_VPC_CNI_STATUS_RUNNING = "Running"
	TKE_CLUSTER_VPC_CNI_STATUS_SUCCEED = "Succeed"
	TKE_CLUSTER_VPC_CNI_STATUS_FAILED  = "Failed"
)

var TKE_CLUSTER_NETWORK_TYPE = []string{TKE_CLUSTER_NETWORK_TYPE_GR, TKE_CLUSTER_NETWORK_TYPE_VPC_CNI, TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY}

const (
	TKE_CLUSTER_NODE_NAME_TYPE_LAN_IP   = "lan-ip"
	TKE_CLUSTER_NODE_NAME_TYPE_HOSTNAME = "hostname"
)

var TKE_CLUSTER_NODE_NAME_TYPE = []string{TKE_CLUSTER_NODE_NAME_TYPE_LAN_IP, TKE_CLUSTER_NODE_NAME_TYPE_HOSTNAME}

const (
	TKE_CLUSTER_KUBE_PROXY_MODE_BPF = "kube-proxy-bpf"
)

var TKE_CLUSTER_KUBE_PROXY_MODE = []string{TKE_CLUSTER_KUBE_PROXY_MODE_BPF}

type OverrideSettings struct {
	Master []tke.InstanceAdvancedSettings
	Work   []tke.InstanceAdvancedSettings
}

const (
	DefaultDesiredPodNum = 0
)

const (
	DefaultAuthenticationOptionsIssuer = "https://kubernetes.default.svc.cluster.local"
)

// This use to filter default values the addon returns.
var TKE_ADDON_DEFAULT_VALUES_KEY = []string{
	"global.image.host",
	"global.cluster.id",
	"global.cluster.appid",
	"global.cluster.uin",
	"global.cluster.subuin",
	"global.cluster.type",
	"global.cluster.clustertype",
	"global.cluster.kubeversion",
	"global.cluster.kubeminor",
	"global.cluster.region",
	"global.cluster.longregion",
	"global.testenv",
}

const (
	InstallSecurityAgentCommandId = "cmd-d8jj2skv"
)

const (
	TKE_CLUSTER_INTERNET = true
	TKE_CLUSTER_INTRANET = false

	TKE_CLUSTER_OPEN_ACCESS  = true
	TKE_CLUSTER_CLOSE_ACCESS = false
)

const (
	backupStorageLocationStateAvailable = "Available"
)

// Content automatically added by the backend of cloud products
var tkeNativeNodePoolAnnotationsMap = map[string]string{
	"cluster-autoscaler.kubernetes.io/scale-down-disabled": "cluster-autoscaler.kubernetes.io/scale-down-disabled",
	"node.tke.cloud.tencent.com/security-agent":            "node.tke.cloud.tencent.com/security-agent",
	"node.tke.cloud.tencent.com/security-groups-spread":    "node.tke.cloud.tencent.com/security-groups-spread",
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

	if v, ok := dMap["taints"]; ok {
		taints := v.([]interface{})
		setting.Taints = make([]*tke.Taint, len(taints))
		for i, d := range taints {
			taint := d.(map[string]interface{})
			var value, key, effect string
			if v, ok := taint["key"].(string); ok {
				key = v
			}
			if v, ok := taint["value"].(string); ok {
				value = v
			}
			if v, ok := taint["effect"].(string); ok {
				effect = v
			}
			taintItem := &tke.Taint{
				Key:    &key,
				Value:  &value,
				Effect: &effect,
			}
			setting.Taints[i] = taintItem
		}
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

