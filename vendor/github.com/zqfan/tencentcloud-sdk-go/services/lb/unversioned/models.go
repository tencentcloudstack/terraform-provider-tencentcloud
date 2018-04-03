package lb

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const (
	LBNetworkTypePublic  = 2
	LBNetworkTypePrivate = 3

	LBForwardTypeApplication = 1
	LBForwardTypeClassic     = 0
	LBForwardTypeAll         = -1

	LBListenerProtocolHTTP  = 1
	LBListenerProtocolTCP   = 2
	LBListenerProtocolUDP   = 3
	LBListenerProtocolHTTPS = 4

	LBTaskSuccess = 0
	LBTaskFail    = 1
	LBTaskDoing   = 2
)

type DescribeLoadBalancersRequest struct {
	*common.BaseRequest
	LoadBalancerIds  []*string `name:"loadBalancerIds" list`
	LoadBalancerType *int      `name:"loadBalancerType"`
	Forward          *int      `name:"forward"`
	LoadBalancerName *string   `name:"loadBalancerName"`
	Domain           *string   `name:"domain"`
	LoadBalancerVips []*string `name:"loadBalancerVips" list`
	BackendWanIps    []*string `name:"backendWanIps" list`
	BackendLanIps    []*string `name:"backendLanIps" list`
	Offset           *int      `name:"offset"`
	Limit            *int      `name:"limit"`
	OrderBy          *string   `name:"orderBy"`
	OrderType        *int      `name:"orderType"`
	SearchKey        *string   `name:"searchKey"`
	ProjectId        *int      `name:"projectId"`
	WithRs           *int      `name:"withRs"`
}

type LoadBalancer struct {
	CreateTime         *string `json:"createTime"`
	Domain             *string `json:"domain"`
	ExpireTime         *string `json:"expireTime"`
	Forward            *int    `json:"forward"`
	InternetAccessible *struct {
		InternetChargeType      *string `json:"internetChargeType"`
		InternetMaxBandwidthOut *int    `json:"internetMaxBandwidthOut"`
	} `json:"internetAccessible"`
	IsolatedTime    *string `json:"isolatedTime"`
	Isolation       *int    `json:"isolation"`
	LbChargePrepaid *struct {
		Period    *int    `json:"period"`
		RenewFlag *string `json:"renewFlag"`
	} `json"lbChargePrepaid"`
	LbChargeType     *string   `json:"lbChargeType"`
	LoadBalancerId   *string   `json:"loadBalancerId"`
	LoadBalancerName *string   `json:"loadBalancerName"`
	LoadBalancerType *int      `json:"loadBalancerType"`
	LoadBalancerVips []*string `json:"loadBalancerVips"`
	Log              *string   `json:"log"`
	OpenBgp          *int      `json:"openBgp"`
	ProjectId        *int      `json:"projectId"`
	RsRegionInfo     *struct {
		Region *string `json:"region"`
		VpcId  *string `json:"vpcId"`
	} `json:"rsRegionInfo"`
	Snat             *bool   `json:"snat"`
	Status           *int    `json:"status"`
	StatusTime       *string `json:"statusTime"`
	SubnetId         *int    `json:"subnetId"`
	UnLoadBalancerId *string `json:"unLoadBalancerId"`
	UniqVpcId        *string `json:"uniqVpcId"`
	VpcId            *int    `json:"vpcId"`
}

type DescribeLoadBalancersResponse struct {
	*common.BaseResponse
	Code            *int            `json:"code"`
	Message         *string         `json:"message"`
	CodeDesc        *string         `json:"codeDesc"`
	TotalCount      *int            `json:"totalCount"`
	LoadBalancerSet []*LoadBalancer `json:"loadBalancerSet"`
}

type DescribeForwardLBListenersRequest struct {
	*common.BaseRequest
	LoadBalancerId   *string   `name:"loadBalancerId"`
	ListenerIds      []*string `name:"listenerIds" list`
	Protocol         *int      `name:"protocol"`
	LoadBalancerPort *int      `name:"loadBalancerPort"`
}

type ListenerRule struct {
	BAutoCreated      *int    `json:"bAutoCreated"`
	Domain            *string `json:"domain"`
	HealthNum         *int    `json:"healthNum"`
	HealthSwitch      *int    `json:"healthSwitch"`
	HttpCheckDomain   *string `json:"httpCheckDomain"`
	HttpCheckMethod   *string `json:"httpCheckMethod"`
	HttpCheckPath     *string `json:"httpCheckPath"`
	HttpCode          *int    `json:"httpCode"`
	HttpHash          *string `json:"httpHash"`
	IntervalTime      *int    `json:"intervalTime"`
	LocationId        *string `json:"locationId"`
	SessionExpire     *int    `json:"sessionExpire"`
	TargetuListenerId *string `json:"targetuListenerId"`
	TargetuLocationId *string `json:"targetuLocationId"`
	TimeOut           *int    `json:"timeOut"`
	UnhealthNum       *int    `json:"unhealthNum"`
	Url               *string `json:"url"`
}

type Listener struct {
	AddTimestamp     *string `json:"addTimestamp"`
	ListenerId       *string `json:"listenerId"`
	ListenerName     *string `json:"listenerName"`
	LoadBalancerPort *int    `json:"loadBalancerPort"`
	Protocol         *int    `json:"protocol"`
	ProtocolType     *string `json:"protocolType"`
	// HTTPS type only
	SSLMode  *string `json:"SSLMode"`
	CertId   *string `json:"certId"`
	CertCaId *string `json:"certCaId"`
	// TCP/UDP type only
	HealthNum     *int            `json:"healthNum"`
	HealthSwitch  *int            `json:"healthSwitch"`
	IntervalTime  *int            `json:"intervalTime"`
	Scheduler     *string         `json:"scheduler"`
	SessionExpire *int            `json:"sessionExpire"`
	TimeOut       *int            `json:"timeOut"`
	UnhealthNum   *int            `json:"unhealthNum"`
	Rules         []*ListenerRule `json:"rules"`
}

type DescribeForwardLBListenersResponse struct {
	*common.BaseResponse
	Code        *int        `json:"code"`
	Message     *string     `json:"message"`
	CodeDesc    *string     `json:"codeDesc"`
	ListenerSet []*Listener `json:"listenerSet"`
}

type Backend struct {
	InstanceId *string `name:"instanceId"`
	Port       *int    `name:"port"`
	Weight     *int    `name:"weight"`
}

type DeregisterInstancesFromForwardLBRequest struct {
	*common.BaseRequest
	LoadBalancerId *string    `name:"loadBalancerId"`
	ListenerId     *string    `name:"listenerId"`
	LocationIds    []*string  `name:"locationIds" list`
	Domain         *string    `name:"domain"`
	Url            *string    `name:"url"`
	Backends       []*Backend `name:"backends" list`
}

type DeregisterInstancesFromForwardLBResponse struct {
	*common.BaseResponse
	Code      *int    `json:"code"`
	Message   *string `json:"message"`
	CodeDesc  *string `json:"codeDesc"`
	RequestId *int    `json:"requestId"`
}

type RegisterInstancesWithForwardLBSeventhListenerRequest struct {
	*common.BaseRequest
	LoadBalancerId *string    `name:"loadBalancerId"`
	ListenerId     *string    `name:"listenerId"`
	LocationIds    []*string  `name:"locationIds" list`
	Domain         *string    `name:"domain"`
	Url            *string    `name:"url"`
	Backends       []*Backend `name:"backends" list`
}

type RegisterInstancesWithForwardLBSeventhListenerResponse struct {
	*common.BaseResponse
	Code      *int    `json:"code"`
	Message   *string `json:"message"`
	CodeDesc  *string `json:"codeDesc"`
	RequestId *int    `json:"requestId"`
}

type DescribeLoadBalancersTaskResultRequest struct {
	*common.BaseRequest
	RequestId *int `name:"requestId"`
}

type DescribeLoadBalancersTaskResultResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	Message  *string `json:"message"`
	CodeDesc *string `json:"codeDesc"`
	Data     *struct {
		Status *int `json:"status"`
	} `json:"data"`
}

type DescribeForwardLBBackendsRequest struct {
	*common.BaseRequest
	LoadBalancerId   *string   `name:"loadBalancerId"`
	ListenerIds      []*string `name:"listenerIds" list`
	Protocol         *int      `name:"protocol"`
	LoadBalancerPort *int      `name:"loadBalancerPort"`
}

type DescribeForwardLBBackendsResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	Message  *string `json:"message"`
	CodeDesc *string `json:"codeDesc"`
	Data     []*struct {
		Backends []*struct {
			AddTimestamp   *string   `json:"addTimestamp"`
			InstanceName   *string   `json:"instanceName"`
			InstanceStatus *int      `json:"instanceStatus"`
			LanIp          *string   `json:"lanIp"`
			Port           *int      `json:"port"`
			UnInstanceId   *string   `json:"unInstanceId"`
			Uuid           *string   `json:"uuid"`
			WanIpSet       []*string `json:"wanIpSet"`
			Weight         *int      `json:"weight"`
		} `json:"backends"`
		ListenerId       *string `json:"listenerId"`
		LoadBalancerPort *int    `json:"loadBalancerPort"`
		Protocol         *int    `json:"protocol"`
		ProtocolType     *string `json:"protocolType"`
		Rules            []*struct {
			Backends []*struct {
				AddTimestamp   *string   `json:"addTimestamp"`
				InstanceName   *string   `json:"instanceName"`
				InstanceStatus *int      `json:"instanceStatus"`
				LanIp          *string   `json:"lanIp"`
				Port           *int      `json:"port"`
				UnInstanceId   *string   `json:"unInstanceId"`
				Uuid           *string   `json:"uuid"`
				WanIpSet       []*string `json:"wanIpSet"`
				Weight         *int      `json:"weight"`
			} `json:"backends"`
			Domain     *string `json:"domain"`
			LocationId *string `json:"locationId"`
			Url        *string `json:"url"`
		} `json:"rules"`
	} `json:"data"`
}

type Request struct {
	*common.BaseRequest
}

type Response struct {
	*common.BaseResponse
}
