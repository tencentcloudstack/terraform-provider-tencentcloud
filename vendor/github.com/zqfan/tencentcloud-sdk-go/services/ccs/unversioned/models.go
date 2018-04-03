package ccs

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const (
	ClusterVipCreate = "Create"
	ClusterVipDelete = "Delete"
)

type DescribeClusterRequest struct {
	*common.BaseRequest
	ClusterIds  []*string `name:"clusterIds" list`
	ClusterName *string   `name:"clusterName"`
	Status      *string   `name:"status"`
	OrderField  *string   `name:"orderField"`
	OrderType   *string   `name:"orderType"`
	Offset      *int      `name:"offset"`
	Limit       *int      `name:"limit"`
}

type Cluster struct {
	ClusterCIDR             *string `json:"clusterCIDR"`
	ClusterExternalEndpoint *string `json:"clusterExternalEndpoint"`
	ClusterId               *string `json:"clusterId"`
	ClusterName             *string `json:"clusterName"`
	CreatedAt               *string `json:"createdAt"`
	Description             *string `json:"description"`
	K8sVersion              *string `json:"k8sVersion"`
	MasterLbSubnetId        *string `json:"masterLbSubnetId"`
	NodeNum                 *int    `json:"nodeNum"`
	NodeStatus              *string `json:"nodeStatus"`
	OpenHttps               *int    `json:"openHttps"`
	OS                      *string `json:"os"`
	ProjectId               *int    `json:"projectId"`
	Region                  *string `json:"region"`
	RegionId                *int    `json:"regionId"`
	Status                  *string `json:"status"`
	TotalCPU                *int    `json:"totalCpu"`
	TotalMem                *int    `json:"totalMem"`
	UnVpcId                 *string `json:"unVpcId"`
	UpdatedAt               *string `json:"updatedAt"`
	VpcId                   *int    `json:"vpcId"`
}

type DescribeClusterResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		TotalCount *int       `json:"totalCount"`
		Clusters   []*Cluster `json:"clusters"`
	} `json:"data"`
}

type CreateClusterRequest struct {
	*common.BaseRequest
	ClusterName               *string `name:"clusterName"`
	ClusterDesc               *string `name:"clusterDesc"`
	ClusterCIDR               *string `name:"clusterCIDR"`
	IgnoreClusterCIDRConflict *int    `name:"ignoreClusterCIDRConflict"`
	ZoneId                    *string `name:"zoneId"`
	GoodsNum                  *int    `name:"goodsNum"`
	CPU                       *int    `name:"cpu"`
	Mem                       *int    `name:"mem"`
	OSName                    *string `name:"osName"`
	InstanceType              *string `name:"instanceType"`
	CVMType                   *string `name:"cvmType"`
	BandwidthType             *string `name:"bandwidthType"`
	Bandwidth                 *int    `name:"bandwidth"`
	WanIp                     *int    `name:"wanIp"`
	VpcId                     *string `name:"vpcId"`
	SubnetId                  *string `name:"subnetId"`
	IsVpcGateway              *int    `name:"isVpcGateway"`
	RootSize                  *int    `name:"rootSize"`
	StorageSize               *int    `name:"storageSize"`
	Password                  *string `name:"password"`
	KeyId                     *string `name:"keyId"`
	Period                    *int    `name:"period"`
	SgId                      *string `name:"sgId"`
	MountTarget		  *string  `name:"mountTarget"`
	DockerGraphPath		  *string  `name:"dockerGraphPath"`
	InstanceName		  *string  `name:"instanceName"`
	ClusterVersion		  *string  `name:"clusterVersion"`
	ProjectId		  *int     `name:"projectId"`
}

type CreateClusterResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		RequestId *int    `json:"requestId"`
		ClusterId *string `json:"clusterId"`
	} `json:"data"`
}

type DeleteClusterRequest struct {
	*common.BaseRequest
	ClusterId      *string `name:"clusterId"`
	NodeDeleteMode *string `name:"nodeDeleteMode"`
}

type DeleteClusterResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		RequestId *int `json:"requestId"`
	} `json:"data"`
}

type DescribeClusterInstancesRequest struct {
	*common.BaseRequest
	ClusterId *string `name:"clusterId"`
	Offset    *int    `name:"offset"`
	Limit     *int    `name:"limit"`
	Namespace *string `name:"namespace"`
}

type ClusterInstance struct {
	AbnormalReason       *string            `json:"abnormalReason"`
	AutoScalingGroupId   *string            `json:"autoScalingGroupId"`
	CPU                  *int               `json:"cpu"`
	CreatedAt            *string            `json:"createdAt"`
	CvmPayMode           *int               `json:"cvmPayMode"`
	CvmState             *int               `json:"cvmState"`
	InstanceCreateTime   *string            `json:"instanceCreateTime"`
	InstanceDeadlineTime *string            `json:"instanceDeadlineTime"`
	InstanceId           *string            `json:"instanceId"`
	InstanceName         *string            `json:"instanceName"`
	InstanceType         *string            `json:"instanceType"`
	IsNormal             *int               `json:"isNormal"`
	KernelVersion        *string            `json:"kernelVersion"`
	//Labels               *map[string]string `json:"labels,omitempty"`
	LanIp                *string            `json:"lanIp"`
	Mem                  *int               `json:"mem"`
	NetworkPayMode       *int               `json:"networkPayMode"`
	OSImage              *string            `json:"osImage"`
	PodCidr              *string            `json:"podCidr"`
	Unschedulable        *bool              `json:"unschedulable"`
	WanIp                *string            `json:"wanIp"`
	Zone                 *string            `json:"zone"`
	ZoneId               *int               `json:"zoneId"`
}

type DescribeClusterInstancesResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		TotalCount *int               `json:"totalCount"`
		Nodes      []*ClusterInstance `json:"nodes"`
	} `json:"data"`
}

type AddClusterInstancesRequest struct {
	*common.BaseRequest
	ClusterId     *string `name:"clusterId"`
	ClusterDesc   *string `name:"clusterDesc"`
	ZoneId        *string `name:"zoneId"`
	CPU           *int    `name:"cpu"`
	Mem           *int    `name:"mem"`
	InstanceType  *string `name:"instanceType"`
	CvmType       *string `name:"cvmType"`
	BandwidthType *string `name:"bandwidthType"`
	Bandwidth     *int    `name:"bandwidth"`
	WanIp         *int    `name:"wanIp"`
	SubnetId      *string `name:"subnetId"`
	IsVpcGateway  *int    `name:"isVpcGateway"`
	StorageSize   *int    `name:"storageSize"`
	RootSize      *int    `name:"rootSize"`
	GoodsNum      *int    `name:"goodsNum"`
	Password      *string `name:"password"`
	KeyId         *string `name:"keyId"`
	Period        *int    `name:"period"`
	SgId          *string `name:"sgId"`
	MountTarget   *string `name:"mountTarget"`
	DockerGraphPath  *string `name:"dockerGraphPath"`
}

type AddClusterInstancesResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Data     *struct {
		InstanceIds []*string `json:"instanceIds"`
		RequestId   *int      `json:"requestId"`
	}
	Message *string `json:"message"`
}

type AddClusterInstancesFromExistedCvmRequest struct {
	*common.BaseRequest
	ClusterId   *string   `name:"clusterId"`
	InstanceIds []*string `name:"instanceIds" list`
	Password    *string   `name:"password"`
	KeyId       *string   `name:"keyId"`
}

type AddClusterInstancesFromExistedCvmResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Data     *struct {
		SuccInstanceIds []*string `json:"succInstanceIds"`
		FaliInstanceIds []*struct {
			InstanceId *string `json:"instanceId"`
			Message    *string `json:"message"`
		} `json:"faliInstanceIds"`
	} `json:"data"`
	Message *string `json:"message"`
}

type DeleteClusterInstancesRequest struct {
	*common.BaseRequest
	ClusterId      *string   `name:"clusterId"`
	InstanceIds    []*string `name:"instanceIds" list`
	NodeDeleteMode *string   `name:"nodeDeleteMode"`
}

type DeleteClusterInstancesResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		RequestId *int `json:"requestId"`
	} `json:"data"`
}

type DescribeClusterTaskResultRequest struct {
	*common.BaseRequest
	RequestId *int `name:"requestId"`
}

type DescribeClusterTaskResultResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Data     *struct {
		Status *string `json:"status"`
	} `json:"data"`
	Message *string `json:"message"`
}

type DescribeClusterSecurityInfoRequest struct {
	*common.BaseRequest
	ClusterId *string `name:"clusterId"`
}

type DescribeClusterSecurityInfoResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Data     *struct {
		CertificationAuthority  *string `json:"certificationAuthority"`
		ClusterExternalEndpoint *string `json:"clusterExternalEndpoint"`
		Password                *string `json:"password"`
		UserName                *string `json:"userName"`
	} `json:"data"`
	Message *string `json:"message"`
}

type OperateClusterVipRequest struct {
	*common.BaseRequest
	ClusterId *string `name:"clusterId"`
	Operation *string `name:"operation"`
}

type OperateClusterVipResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		RequestId *int `json:"requestId"`
	} `json:"data"`
}

type Request struct {
	*common.BaseRequest
}

type Response struct {
	*common.BaseResponse
}
