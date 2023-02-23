package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ChdfsService struct {
	client *connectivity.TencentCloudClient
}

func (me *ChdfsService) DescribeChdfsAccessGroupById(ctx context.Context, accessGroupId string) (accessGroup *chdfs.AccessGroup, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeAccessGroupRequest()
	request.AccessGroupId = &accessGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeAccessGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.AccessGroup == nil {
		return
	}

	accessGroup = response.Response.AccessGroup
	return
}

func (me *ChdfsService) DeleteChdfsAccessGroupById(ctx context.Context, accessGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDeleteAccessGroupRequest()
	request.AccessGroupId = &accessGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DeleteAccessGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ChdfsService) DescribeChdfsFileSystemById(ctx context.Context, fileSystemId string) (fileSystem *chdfs.FileSystem, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeFileSystemRequest()
	request.FileSystemId = &fileSystemId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeFileSystem(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.FileSystem == nil {
		return
	}

	fileSystem = response.Response.FileSystem
	return
}

func (me *ChdfsService) DeleteChdfsFileSystemById(ctx context.Context, fileSystemId string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDeleteFileSystemRequest()
	request.FileSystemId = &fileSystemId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DeleteFileSystem(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ChdfsService) ChdfsFileSystemStateRefreshFunc(fileSystemId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeChdfsFileSystemById(ctx, fileSystemId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.UInt64ToStr(*object.Status), nil
	}
}
