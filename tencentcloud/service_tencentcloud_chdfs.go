package tencentcloud

import (
	"context"
	"log"

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
