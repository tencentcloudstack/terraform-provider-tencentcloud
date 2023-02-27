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

func (me *ChdfsService) DescribeChdfsAccessRulesById(ctx context.Context, accessGroupId string, accessRuleId string) (accessRule *chdfs.AccessRule, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeAccessRulesRequest()
	request.AccessGroupId = &accessGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeAccessRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AccessRules) < 1 {
		return
	}

	for _, rule := range response.Response.AccessRules {
		if *rule.AccessRuleId == helper.StrToUInt64(accessRuleId) {
			accessRule = rule
			break
		}
	}
	return
}

func (me *ChdfsService) DeleteChdfsAccessRulesById(ctx context.Context, accessRuleId string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDeleteAccessRulesRequest()
	request.AccessRuleIds = []*uint64{helper.StrToUint64Point(accessRuleId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DeleteAccessRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ChdfsService) DescribeChdfsLifeCycleRuleById(ctx context.Context, fileSystemId string, lifeCycleRuleId string) (lifeCycleRule *chdfs.LifeCycleRule, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeLifeCycleRulesRequest()
	request.FileSystemId = &fileSystemId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeLifeCycleRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LifeCycleRules) < 1 {
		return
	}

	for _, rule := range response.Response.LifeCycleRules {
		if *rule.LifeCycleRuleId == helper.StrToUInt64(lifeCycleRuleId) {
			lifeCycleRule = rule
			break
		}
	}
	return
}

func (me *ChdfsService) DescribeChdfsLifeCycleRuleByPath(ctx context.Context, fileSystemId string, path string) (lifeCycleRule *chdfs.LifeCycleRule, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeLifeCycleRulesRequest()
	request.FileSystemId = &fileSystemId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeLifeCycleRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LifeCycleRules) < 1 {
		return
	}

	for _, rule := range response.Response.LifeCycleRules {
		if *rule.Path == path {
			lifeCycleRule = rule
			break
		}
	}
	return
}

func (me *ChdfsService) DeleteChdfsLifeCycleRuleById(ctx context.Context, lifeCycleRuleId string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDeleteLifeCycleRulesRequest()
	request.LifeCycleRuleIds = []*uint64{helper.StrToUint64Point(lifeCycleRuleId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DeleteLifeCycleRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ChdfsService) DescribeChdfsMountPointById(ctx context.Context, mountPointId string) (mountPoint *chdfs.MountPoint, errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDescribeMountPointRequest()
	request.MountPointId = &mountPointId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DescribeMountPoint(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.MountPoint == nil {
		return
	}

	mountPoint = response.Response.MountPoint
	return
}

func (me *ChdfsService) DeleteChdfsMountPointById(ctx context.Context, mountPointId string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDeleteMountPointRequest()
	request.MountPointId = &mountPointId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DeleteMountPoint(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ChdfsService) DeleteChdfsMountPointAttachmentById(ctx context.Context, mountPointId string, accessGroupIds []*string) (errRet error) {
	logId := getLogId(ctx)

	request := chdfs.NewDisassociateAccessGroupsRequest()
	request.MountPointId = &mountPointId
	request.AccessGroupIds = accessGroupIds

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseChdfsClient().DisassociateAccessGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
