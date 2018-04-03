package cbs

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

type DescribeSnapshotsRequest struct {
	*common.BaseRequest
	DiskType    *string   `name:"diskType"`
	ProjectId   *int      `name:"projectId"`
	StorageIds  []*string `name:"storageIds" list`
	SnapshotIds []*string `name:"snapshotIds" list`
	Offset      *int      `name:"offset"`
	Limit       *int      `name:"limit"`
}

type Snapshot struct {
	CreateTime     *string `json:"createTime"`
	SnapshotId     *string `json:"snapshotId"`
	StorageId      *string `json:"storageId"`
	SnapshotName   *string `json:"snapshotName"`
	SnapshotStatus *string `json:"snapshotStatus"`
	ProjectId      *int    `json:"projectId"`
	Percent        *int    `json:"percent"`
	StorageSize    *int    `json:"storageSize"`
	DiskType       *string `json:"diskType"`
	ZoneId         *int    `json:"zoneId"`
	ZoneName       *string `json:"zoneName"`
	Deadline       *string `json:"deadline"`
	Encrypt        *string `json:"encrypt"`
}

type DescribeSnapshotsResponse struct {
	*common.BaseResponse
	Code        *int        `json:"code"`
	Message     *string     `json:"message"`
	CodeDesc    *string     `json:"codeDesc"`
	SnapshotSet []*Snapshot `json:"snapshotSet"`
	TotalCount  *int        `json:"totalCount"`
}

type CreateSnapshotRequest struct {
	*common.BaseRequest
	StorageId    *string `name:"storageId"`
	SnapshotName *string `name:"snapshotName"`
}

type CreateSnapshotResponse struct {
	*common.BaseResponse
	Code       *int    `json:"code"`
	Message    *string `json:"message"`
	CodeDesc   *string `json:"codeDesc"`
	SnapshotId *string `json:"snapshotId"`
}

type DeleteSnapshotRequest struct {
	*common.BaseRequest
	SnapshotIds []*string `name:"snapshotIds" list`
}

type DeleteSnapshotResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	Message  *string `json:"message"`
	CodeDesc *string `json:"codeDesc"`
	Detail   *map[string]*struct {
		Code *int    `json:"code"`
		Msg  *string `json:"msg"`
		// the data field is an array when success
		// and an map when fail
		// so comment it out for now
		//Data *map[string]int `json:"data"`
	} `json:"detail"`
}

type ModifySnapshotRequest struct {
	*common.BaseRequest
	SnapshotId   *string `name:"snapshotId"`
	SnapshotName *string `name:"snapshotName"`
}

type ModifySnapshotResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	Message  *string `json:"message"`
	CodeDesc *string `json:"codeDesc"`
}

type DescribeCbsStoragesRequest struct {
	*common.BaseRequest
	DiskType      *string   `name:"diskType"`
	PayMode       *string   `name:"payMode"`
	Portable      *int      `name:"portable"`
	ProjectId     *int      `name:"projectId"`
	StorageIds    []*string `name:"storageIds" list`
	StorageType   *string   `name:"storageType"`
	StorageStatus []*string `name:"storageStatus" list`
	UInstanceIds  []*string `name:"uInstanceIds" list`
	Zone          *string   `name:"zone"`
	Offset        *int      `name:"offset"`
	Limit         *int      `name:"limit"`
}

type Storage struct {
	StorageId        *string `json:"storageId"`
	Attached         *int    `json:"attached"`
	CreateTime       *string `json:"createTime"`
	DeadlineTime     *string `json:"deadlineTime"`
	DiskType         *string `json:"diskType"`
	PayMode          *string `json:"payMode"`
	Portable         *int    `json:"portable"`
	ProjectId        *int    `json:"projectId"`
	SnapshotAbility  *int    `json:"snapshotAbility"`
	StorageName      *string `json:"storageName"`
	StorageSize      *int    `json:"storageSize"`
	StorageStatus    *string `json:"storageStatus"`
	StorageType      *string `json:"storageType"`
	UInstanceId      *string `json:"uInstanceId"`
	Zone             *string `json:"zone"`
	ZoneId           *int    `json:"zoneId"`
	Rollbacking      *int    `json:"rollbacking"`
	RrollbackPercent *int    `json:"rollbackPercent"`
}

type DescribeCbsStoragesResponse struct {
	*common.BaseResponse
	Code       *int       `json:"code"`
	Message    *string    `json:"message"`
	CodeDesc   *string    `json:"codeDesc"`
	TotalCount *int       `json:"totalCount"`
	StorageSet []*Storage `json:"storageSet"`
}

type Request struct {
}

type Response struct {
}
