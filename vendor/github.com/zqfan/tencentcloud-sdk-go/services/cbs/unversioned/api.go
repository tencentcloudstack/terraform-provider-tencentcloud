package cbs

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

func NewCreateSnapshotRequest() (request *CreateSnapshotRequest) {
	request = &CreateSnapshotRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("snapshot", APIVersion, "CreateSnapshot")
	return
}

func NewCreateSnapshotResponse() (response *CreateSnapshotResponse) {
	response = &CreateSnapshotResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateSnapshot(request *CreateSnapshotRequest) (response *CreateSnapshotResponse, err error) {
	if request == nil {
		request = NewCreateSnapshotRequest()
	}
	response = NewCreateSnapshotResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteSnapshotRequest() (request *DeleteSnapshotRequest) {
	request = &DeleteSnapshotRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("snapshot", APIVersion, "DeleteSnapshot")
	return
}

func NewDeleteSnapshotResponse() (response *DeleteSnapshotResponse) {
	response = &DeleteSnapshotResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeleteSnapshot(request *DeleteSnapshotRequest) (response *DeleteSnapshotResponse, err error) {
	if request == nil {
		request = NewDeleteSnapshotRequest()
	}
	response = NewDeleteSnapshotResponse()
	err = c.Send(request, response)
	return
}

func NewModifySnapshotRequest() (request *ModifySnapshotRequest) {
	request = &ModifySnapshotRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("snapshot", APIVersion, "ModifySnapshot")
	return
}

func NewModifySnapshotResponse() (response *ModifySnapshotResponse) {
	response = &ModifySnapshotResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) ModifySnapshot(request *ModifySnapshotRequest) (response *ModifySnapshotResponse, err error) {
	if request == nil {
		request = NewModifySnapshotRequest()
	}
	response = NewModifySnapshotResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeSnapshotsRequest() (request *DescribeSnapshotsRequest) {
	request = &DescribeSnapshotsRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("snapshot", APIVersion, "DescribeSnapshots")
	return
}

func NewDescribeSnapshotsResponse() (response *DescribeSnapshotsResponse) {
	response = &DescribeSnapshotsResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeSnapshots(request *DescribeSnapshotsRequest) (response *DescribeSnapshotsResponse, err error) {
	if request == nil {
		request = NewDescribeSnapshotsRequest()
	}
	response = NewDescribeSnapshotsResponse()
	err = c.Send(request, response)
	return
}

func NewDescribeCbsStoragesRequest() (request *DescribeCbsStoragesRequest) {
	request = &DescribeCbsStoragesRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("snapshot", APIVersion, "DescribeCbsStorages")
	return
}

func NewDescribeCbsStoragesResponse() (response *DescribeCbsStoragesResponse) {
	response = &DescribeCbsStoragesResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeCbsStorages(request *DescribeCbsStoragesRequest) (response *DescribeCbsStoragesResponse, err error) {
	if request == nil {
		request = NewDescribeCbsStoragesRequest()
	}
	response = NewDescribeCbsStoragesResponse()
	err = c.Send(request, response)
	return
}
