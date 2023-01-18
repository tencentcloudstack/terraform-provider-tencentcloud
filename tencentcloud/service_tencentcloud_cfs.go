package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CfsService struct {
	client *connectivity.TencentCloudClient
}

func (me *CfsService) DescribeFileSystem(ctx context.Context, fsId, vpcId, subnetId string) (fs []*cfs.FileSystemInfo, errRet error) {
	logId := getLogId(ctx)
	request := cfs.NewDescribeCfsFileSystemsRequest()
	if fsId != "" {
		request.FileSystemId = &fsId
	}
	if vpcId != "" {
		request.VpcId = &vpcId
	}
	if subnetId != "" {
		request.SubnetId = &subnetId
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DescribeCfsFileSystems(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	fs = response.Response.FileSystems
	return
}

func (me *CfsService) DescribeMountTargets(ctx context.Context, fsId string) (targets []*cfs.MountInfo, errRet error) {
	logId := getLogId(ctx)
	request := cfs.NewDescribeMountTargetsRequest()
	request.FileSystemId = &fsId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DescribeMountTargets(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	targets = response.Response.MountTargets
	return
}

func (me *CfsService) ModifyFileSystemName(ctx context.Context, fsId, fsName string) error {
	logId := getLogId(ctx)
	request := cfs.NewUpdateCfsFileSystemNameRequest()
	request.FileSystemId = &fsId
	request.FsName = &fsName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().UpdateCfsFileSystemName(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CfsService) ModifyFileSystemAccessGroup(ctx context.Context, fsId, accessGroupId string) error {
	logId := getLogId(ctx)
	request := cfs.NewUpdateCfsFileSystemPGroupRequest()
	request.FileSystemId = &fsId
	request.PGroupId = &accessGroupId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().UpdateCfsFileSystemPGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CfsService) DeleteFileSystem(ctx context.Context, fsId string) error {
	logId := getLogId(ctx)
	request := cfs.NewDeleteCfsFileSystemRequest()
	request.FileSystemId = &fsId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DeleteCfsFileSystem(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CfsService) CreateAccessGroup(ctx context.Context, name, description string) (id string, errRet error) {
	logId := getLogId(ctx)
	request := cfs.NewCreateCfsPGroupRequest()
	request.Name = &name
	if description != "" {
		request.DescInfo = &description
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().CreateCfsPGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.PGroupId == nil {
		errRet = fmt.Errorf("cfs access group id is nil")
		return
	}
	id = *response.Response.PGroupId
	return
}

func (me *CfsService) DescribeAccessGroup(ctx context.Context, id, name string) (accessGroups []*cfs.PGroupInfo, errRet error) {
	logId := getLogId(ctx)
	request := cfs.NewDescribeCfsPGroupsRequest()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DescribeCfsPGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accessGroups = make([]*cfs.PGroupInfo, 0)
	for _, accessGroup := range response.Response.PGroupList {
		if id != "" && *accessGroup.PGroupId != id {
			continue
		}
		if name != "" && *accessGroup.Name != name {
			continue
		}
		accessGroups = append(accessGroups, accessGroup)
	}
	return
}

func (me *CfsService) DeleteAccessGroup(ctx context.Context, id string) error {
	logId := getLogId(ctx)
	request := cfs.NewDeleteCfsPGroupRequest()
	request.PGroupId = &id
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DeleteCfsPGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CfsService) DescribeAccessRule(ctx context.Context, accessGroupId, accessRuleId string) (accessRules []*cfs.PGroupRuleInfo, errRet error) {
	logId := getLogId(ctx)
	request := cfs.NewDescribeCfsRulesRequest()
	request.PGroupId = &accessGroupId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DescribeCfsRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			if sdkErr.Code == CfsInvalidPgroup {
				return nil, nil
			}
		}
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accessRules = make([]*cfs.PGroupRuleInfo, 0)
	for _, accessRule := range response.Response.RuleList {
		if accessRuleId != "" && *accessRule.RuleId != accessRuleId {
			continue
		}
		accessRules = append(accessRules, accessRule)
	}
	return
}

func (me *CfsService) DeleteAccessRule(ctx context.Context, accessGroupId, accessRuleId string) error {
	logId := getLogId(ctx)
	request := cfs.NewDeleteCfsRuleRequest()
	request.PGroupId = &accessGroupId
	request.RuleId = &accessRuleId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCfsClient().DeleteCfsRule(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CfsService) DescribeCfsAutoSnapshotPolicyById(ctx context.Context, autoSnapshotPolicyId string) (autoSnapshotPolicy *cfs.AutoSnapshotPolicyInfo, errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDescribeAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyId = &autoSnapshotPolicyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*cfs.AutoSnapshotPolicyInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCfsClient().DescribeAutoSnapshotPolicies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutoSnapshotPolicies) < 1 {
			break
		}
		instances = append(instances, response.Response.AutoSnapshotPolicies...)
		if len(response.Response.AutoSnapshotPolicies) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	autoSnapshotPolicy = instances[0]
	return
}

func (me *CfsService) DeleteCfsAutoSnapshotPolicyById(ctx context.Context, autoSnapshotPolicyId string) (errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDeleteAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyId = &autoSnapshotPolicyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().DeleteAutoSnapshotPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfsService) DescribeCfsAutoSnapshotPolicyAttachmentById(ctx context.Context, autoSnapshotPolicyId string, fileSystemIds string) (autoSnapshotPolicyAttachment *cfs.AutoSnapshotPolicyInfo, errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDescribeAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyId = &autoSnapshotPolicyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*cfs.AutoSnapshotPolicyInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCfsClient().DescribeAutoSnapshotPolicies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutoSnapshotPolicies) < 1 {
			break
		}
		instances = append(instances, response.Response.AutoSnapshotPolicies...)
		if len(response.Response.AutoSnapshotPolicies) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	var fileSystemIdsList []string
	autoSnapshotPolicy := instances[0]

	for _, fileSystem := range autoSnapshotPolicy.FileSystems {
		fileSystemIdsList = append(fileSystemIdsList, *fileSystem.FileSystemId)
	}

	res := strings.Split(fileSystemIds, ",")
	sort.Strings(res)
	sort.Strings(fileSystemIdsList)

	if reflect.DeepEqual(res, fileSystemIdsList) {
		autoSnapshotPolicyAttachment = autoSnapshotPolicy
	}

	return
}

func (me *CfsService) DeleteCfsAutoSnapshotPolicyAttachmentById(ctx context.Context, autoSnapshotPolicyId string, fileSystemIds string) (errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewUnbindAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyId = &autoSnapshotPolicyId
	request.FileSystemIds = &fileSystemIds

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().UnbindAutoSnapshotPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfsService) DescribeCfsSnapshotById(ctx context.Context, snapshotId string) (snapshot *cfs.SnapshotInfo, errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDescribeCfsSnapshotsRequest()
	request.SnapshotId = &snapshotId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*cfs.SnapshotInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCfsClient().DescribeCfsSnapshots(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Snapshots) < 1 {
			break
		}
		instances = append(instances, response.Response.Snapshots...)
		if len(response.Response.Snapshots) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	snapshot = instances[0]
	return
}

func (me *CfsService) DeleteCfsSnapshotById(ctx context.Context, snapshotId string) (errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDeleteCfsSnapshotRequest()
	request.SnapshotId = &snapshotId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().DeleteCfsSnapshot(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfsService) CfsSnapshotStateRefreshFunc(snapshotId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeCfsSnapshotById(ctx, snapshotId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CfsService) DescribeCfsMountTargetsById(ctx context.Context, fileSystemId string) (mountTargets []*cfs.MountInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cfs.NewDescribeMountTargetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.FileSystemId = helper.String(fileSystemId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().DescribeMountTargets(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	mountTargets = response.Response.MountTargets

	return
}

func (me *CfsService) DescribeCfsFileSystemClientsById(ctx context.Context, fileSystemId string) (fileSystemClients []*cfs.FileSystemClient, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cfs.NewDescribeCfsFileSystemClientsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.FileSystemId = helper.String(fileSystemId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().DescribeCfsFileSystemClients(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	fileSystemClients = response.Response.ClientList

	return
}

func (me *CfsService) DescribeCfsUserQuotaById(ctx context.Context, fileSystemId string, userType string, userId string) (userQuota *cfs.UserQuota, errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDescribeUserQuotaRequest()
	request.FileSystemId = &fileSystemId

	typeFilter := cfs.Filter{
		Name:   helper.String("UserType"),
		Values: []*string{helper.String(userType)},
	}
	request.Filters = append(request.Filters, &typeFilter)

	idFilter := cfs.Filter{
		Name:   helper.String("UserId"),
		Values: []*string{helper.String(userId)},
	}
	request.Filters = append(request.Filters, &idFilter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*cfs.UserQuota, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCfsClient().DescribeUserQuota(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.UserQuotaInfo) < 1 {
			break
		}
		instances = append(instances, response.Response.UserQuotaInfo...)
		if len(response.Response.UserQuotaInfo) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	userQuota = instances[0]
	return
}

func (me *CfsService) DeleteCfsUserQuotaById(ctx context.Context, fileSystemId string, userType string, userId string) (errRet error) {
	logId := getLogId(ctx)

	request := cfs.NewDeleteUserQuotaRequest()
	request.FileSystemId = &fileSystemId
	request.UserType = &userType
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCfsClient().DeleteUserQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CfsService) DescribeCfsAvailableZoneByFilter(ctx context.Context) (availableZone []*cfs.AvailableRegion, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cfs.NewDescribeAvailableZoneInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseCfsClient().DescribeAvailableZoneInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	availableZone = response.Response.RegionZones

	return
}
