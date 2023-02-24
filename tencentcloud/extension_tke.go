package tencentcloud

import tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

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

var TKE_INSTANCE_CHARGE_TYPE = []string{CVM_CHARGE_TYPE_PREPAID, CVM_CHARGE_TYPE_POSTPAID}

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
	TKE_CLUSTER_NETWORK_TYPE_GR      = "GR"
	TKE_CLUSTER_NETWORK_TYPE_VPC_CNI = "VPC-CNI"
)

var TKE_CLUSTER_NETWORK_TYPE = []string{TKE_CLUSTER_NETWORK_TYPE_GR, TKE_CLUSTER_NETWORK_TYPE_VPC_CNI}

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
