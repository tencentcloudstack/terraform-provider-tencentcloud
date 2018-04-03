package cvm

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

type Filter struct {
	Name   *string   `name:"Name"`
	Values []*string `name:"Values"`
}

type DescribeAddressesRequest struct {
	*common.BaseRequest
	AddressIds []*string `name:"AddressIds" list`
	Filters    []*Filter `name:"Filters" list`
	Offset     *int      `name:"Offset" type:"int"`
	Limit      *int      `name:"Limit" type:"int"`
}

type Address struct {
	AddressId             *string `json:"AddressId"`
	AddressIp             *string `json:"AddressIp"`
	AddressName           *string `json:"AddressName"`
	AddressStatus         *string `json:"AddressStatus"`
	BindedResourceId      *string `json:"BindedResourceId"`
	CreatedTime           *string `json:"CreatedTime"`
	InstanceId            *string `json:"InstanceId"`
	IsArrears             *bool   `json:"IsArrears"`
	IsBlocked             *bool   `json:"IsBlocked"`
	IsEipDirectConnection *bool   `json:"IsEipDirectConnection"`
	NetworkInterfaceId    *string `json:"NetworkInterfaceId"`
	PrivateAddressIp      *string `json:"PrivateAddressIp"`
}

type DescribeAddressesResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId  *string    `json:"RequestId"`
		TotalCount *int       `json:"TotalCount"`
		AddressSet []*Address `json:"AddressSet`
	} `json:"Response"`
}

type AllocateAddressesRequest struct {
	*common.BaseRequest
	AddressCount *int `name:"AddressCount" type:"int"`
}

type AllocateAddressesResponse struct {
	*common.BaseResponse
	Response *struct {
		AddressSet []*string `json:"AddressSet"`
		RequestId  *string   `json:"RequestId"`
	}
}

type ReleaseAddressesRequest struct {
	*common.BaseRequest
	AddressIds []*string `name:"AddressIds" list`
}

type ReleaseAddressesResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId *string `json:"RequestId"`
	}
}

type ModifyAddressAttributeRequest struct {
	*common.BaseRequest
	AddressId   *string `name:"AddressId"`
	AddressName *string `name:"AddressName"`
}

type ModifyAddressAttributeResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId *string `json:"RequestId"`
	}
}

type AssociateAddressRequest struct {
	*common.BaseRequest
	AddressId          *string `name:"AddressId"`
	InstanceId         *string `name:"InstanceId"`
	NetworkInterfaceId *string `name:"NetworkInterfaceId"`
	PrivateIpAddress   *string `name:"PrivateIpAddress"`
}

type AssociateAddressResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId *string `json:"RequestId"`
	}
}

type DisassociateAddressRequest struct {
	*common.BaseRequest
	AddressId                *string `name:"AddressId"`
	ReallocateNormalPublicIp *string `name:"ReallocateNormalPublicIp"`
}

type DisassociateAddressResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId *string `json:"RequestId"`
	}
}

type DescribeInstancesRequest struct {
	*common.BaseRequest
	InstanceIds []*string `name:"InstanceIds" list`
	Filters     []*Filter `name:"Filters" list`
	Offset      *string   `name:"Offset"`
	Limit       *string   `name:"Limit"`
}

type Placement struct {
	Zone      *string   `json:"Zone" name:"Zone"`
	ProjectId *int      `json:"ProjectId" name:"ProjectId" type:"int"`
	HostIds   []*string `json:"HostIds" name:"HostIds" list`
}

type SystemDisk struct {
	DiskType *string `json:"DiskType" name:"DiskType"`
	DiskId   *string `json:"DiskId" name:"DiskId"`
	DiskSize *int    `json:"DiskSize" name:"DiskSize" type:"int"`
}

type DataDisks struct {
	DiskType *string `json:"DiskType" name:"DiskType"`
	DiskId   *string `json:"DiskId" name:"DiskId"`
	DiskSize *int    `json:"DiskSize" name:"DiskSize" type:"int"`
}

type InternetAccessible struct {
	InternetChargeType      *string `json:"InternetChargeType" name:"InternetChargeType"`
	InternetMaxBandwidthOut *int    `json:"InternetMaxBandwidthOut" name:"InternetMaxBandwidthOut" type:"int"`
	PublicIpAssigned        *string `json:"PublicIpAssigned" name:"PublicIpAssigned"`
}

type VirtualPrivateCloud struct {
	VpcId              *string   `json:"VpcId" name:"VpcId"`
	SubnetId           *string   `json:"SubnetId" name:"SubnetId"`
	AsVpcGateway       *string   `json:"AsVpcGateway" name:"AsVpcGateway"`
	PrivateIpAddresses []*string `json:"PrivateIpAddresses" name:"PrivateIpAddresses" list`
}

type Instance struct {
	Placement           *Placement           `json:"Placement"`
	InstanceId          *string              `json:"InstanceId"`
	InstanceType        *string              `json:"InstanceType"`
	CPU                 *int                 `json:"CPU"`
	Memory              *int                 `json:"Memory"`
	RestrictState       *string              `json:"RestrictState"`
	InstanceName        *string              `json:"InstanceName"`
	InstanceChargeType  *string              `json:"InstanceChargeType"`
	SystemDisk          *SystemDisk          `json:"SystemDisk"`
	DataDisks           []*DataDisks         `json:"DataDisks"`
	PrivateIpAddresses  []*string            `json:"PrivateIpAddresses"`
	PublicIpAddresses   []*string            `json:"PublicIpAddresses"`
	InternetAccessible  *InternetAccessible  `json:"InternetAccessible"`
	VirtualPrivateCloud *VirtualPrivateCloud `json:"VirtualPrivateCloud"`
	ImageId             *string              `json:"ImageId"`
	RenewFlag           *string              `json:"RenewFlag"`
	CreatedTime         *string              `json:"CreatedTime"`
	ExpiredTime         *string              `json:"ExpiredTime"`
}

type DescribeInstancesResponse struct {
	*common.BaseResponse
	Response *struct {
		TotalCount  *int        `json:"TotalCount"`
		InstanceSet []*Instance `json:"InstanceSet"`
		RequestId   *string     `json:"RequestId"`
	}
}

type InstanceChargePrepaid struct {
	Period    *int    `name:"Period" type:"int"`
	RenewFlag *string `name:"RenewFlag"`
}

type LoginSettings struct {
	Password       *string   `name:"Password"`
	KeyIds         []*string `name:"KeyIds" list`
	KeepImageLogin *string   `name:"KeepImageLogin"`
}

type EnhancedService struct {
	UnitPrice     *float64 `name:"UnitPrice" type:"float64"`
	ChargeUnit    *string  `name:"ChargeUnit"`
	OriginalPrice *float64 `name:"OriginalPrice" type:"float64"`
	DiscountPrice *float64 `name:"DiscountPrice" type:"float64"`
}

type RunInstancesRequest struct {
	*common.BaseRequest
	InstanceChargeType    *string                `name:"InstanceChargeType"`
	InstanceChargePrepaid *InstanceChargePrepaid `name:"InstanceChargePrepaid"`
	Placement             *Placement             `name:"Placement"`
	InstanceType          *string                `name:"InstanceType"`
	ImageId               *string                `name:"ImageId"`
	SystemDisk            *SystemDisk            `name:"SystemDisk"`
	DataDisks             []*DataDisks           `name:"DataDisks" list`
	VirtualPrivateCloud   *VirtualPrivateCloud   `name:"VirtualPrivateCloud"`
	InternetAccessible    *InternetAccessible    `name:"InternetAccessible"`
	InstanceCount         *int                   `name:"InstanceCount" type:"int"`
	InstanceName          *string                `name:"InstanceName"`
	LoginSettings         *LoginSettings         `name:"LoginSettings"`
	SecurityGroupIds      []*string              `name:"SecurityGroupIds" list`
	EnhancedService       *EnhancedService       `name:"EnhancedService"`
	ClientToken           *string                `name:"ClientToken"`
}

type RunInstancesResponse struct {
	*common.BaseResponse
	Response *struct {
		InstanceIdSet []*string `json:"InstanceIdSet"`
		RequestId     *string   `json:"RequestId"`
	}
}

type TerminateInstancesRequest struct {
	*common.BaseRequest
	InstanceIds []*string `name:"InstanceIds" list`
}

type TerminateInstancesResponse struct {
	*common.BaseResponse
	Response *struct {
		RequestId *string `json:"RequestId"`
	}
}

type Request struct {
}

type Response struct {
}
