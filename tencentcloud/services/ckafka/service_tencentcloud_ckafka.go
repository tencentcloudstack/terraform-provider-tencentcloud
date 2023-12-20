package ckafka

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCkafkaService(client *connectivity.TencentCloudClient) CkafkaService {
	return CkafkaService{client: client}
}

type CkafkaService struct {
	client *connectivity.TencentCloudClient
}

func (me *CkafkaService) CheckCkafkaInstanceReady(ctx context.Context,
	instanceId string) (has bool, ready bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	var (
		request  = ckafka.NewDescribeInstancesDetailRequest()
		response = ckafka.NewDescribeInstancesDetailResponse()
		info     *ckafka.InstanceDetail
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, err := me.client.UseCkafkaClient().DescribeInstancesDetail(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read ckafka instance failed, reason: %v", logId, err)
		return false, false, err
	}
	if len(response.Response.Result.InstanceList) < 1 {
		return
	}
	has = true
	info = response.Response.Result.InstanceList[0]
	if *info.Status == 1 {
		ready = true
	}
	return
}

func (me *CkafkaService) ModifyCkafkaInstanceAttributes(ctx context.Context,
	request *ckafka.ModifyInstanceAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCkafkaClient().ModifyInstanceAttributes(request)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId,
			request.GetAction(), request.ToJsonString(), err.Error())
	}
	return
}

func (me *CkafkaService) DescribeCkafkaInstanceById(ctx context.Context,
	instanceId string) (info *ckafka.InstanceDetail, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	var (
		request  = ckafka.NewDescribeInstancesDetailRequest()
		response = ckafka.NewDescribeInstancesDetailResponse()
	)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, err := me.client.UseCkafkaClient().DescribeInstancesDetail(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read ckafka instance failed, reason: %v", logId, err)
		return nil, false, err
	}
	if len(response.Response.Result.InstanceList) < 1 {
		return
	}
	has = true
	info = response.Response.Result.InstanceList[0]
	return
}

func (me *CkafkaService) CreateUser(ctx context.Context, instanceId, user, password string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreateUserRequest()
	request.InstanceId = &instanceId
	request.Name = &user
	request.Password = &password

	var response *ckafka.CreateUserResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().CreateUser(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
	}
	return nil
}

func (me *CkafkaService) OperateStatusCheck(ctx context.Context, result *ckafka.JgwOperateResponse) (isSucceed bool) {
	logId := tccommon.GetLogId(ctx)
	if result == nil {
		log.Printf("[CRITAL]%s OperateStatusCheck fail, result is nil", logId)
		return false
	}

	if result != nil && *result.ReturnCode == "0" {
		return true
	} else {
		return false
	}
}

func (me *CkafkaService) DescribeUserByUserId(ctx context.Context, userId string) (userInfo *ckafka.User, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	items := strings.Split(userId, tccommon.FILED_SP)
	if len(items) != 2 {
		errRet = fmt.Errorf("id of resource.tencentcloud_ckafka_user is wrong")
		return
	}
	instanceId, user := items[0], items[1]

	if _, has, _ = me.DescribeInstanceById(ctx, instanceId); !has {
		return
	}

	request := ckafka.NewDescribeUserRequest()
	request.InstanceId = &instanceId
	request.SearchWord = &user

	var response *ckafka.DescribeUserResponse
	var err error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().DescribeUser(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}

	if response != nil && response.Response != nil && response.Response.Result != nil && response.Response.Result.Users != nil {
		if len(response.Response.Result.Users) < 1 {
			has = false
			return
		} else if len(response.Response.Result.Users) > 1 {
			errRet = fmt.Errorf("[CRITAL]%s dumplicated users found", logId)
			return
		}

		userInfo = response.Response.Result.Users[0]
		has = true
		return
	}

	return
}

func (me *CkafkaService) ModifyPassword(ctx context.Context, instanceId, user, oldPasswd, newPasswd string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewModifyPasswordRequest()
	request.InstanceId = &instanceId
	request.Name = &user
	request.Password = &oldPasswd
	request.PasswordNew = &newPasswd

	var response *ckafka.ModifyPasswordResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().ModifyPassword(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
	}
	return nil
}

func (me *CkafkaService) DeleteUser(ctx context.Context, userId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	items := strings.Split(userId, tccommon.FILED_SP)
	if len(items) != 2 {
		errRet = fmt.Errorf("id of resource.tencentcloud_ckafka_user is wrong")
		return
	}
	instanceId, user := items[0], items[1]

	request := ckafka.NewDeleteUserRequest()
	request.InstanceId = &instanceId
	request.Name = &user

	var response *ckafka.DeleteUserResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().DeleteUser(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
	}
	return nil
}

func (me *CkafkaService) DescribeUserByFilter(ctx context.Context, params map[string]interface{}) (userInfos []*ckafka.User, errRet error) {
	logId := tccommon.GetLogId(ctx)

	instanceId := params["instance_id"].(string)
	if _, has, _ := me.DescribeInstanceById(ctx, instanceId); !has {
		return
	}

	request := ckafka.NewDescribeUserRequest()
	var offset int64 = 0
	var pageSize = int64(CKAFKA_DESCRIBE_LIMIT)
	request.InstanceId = &instanceId
	if user, ok := params["account_name"]; ok {
		request.SearchWord = helper.String(user.(string))
	}
	request.Limit = &pageSize
	request.Offset = &offset

	userInfos = make([]*ckafka.User, 0)
	for {
		var response *ckafka.DescribeUserResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseCkafkaClient().DescribeUser(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			userInfos = append(userInfos, response.Response.Result.Users...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		} else {
			if len(response.Response.Result.Users) < CKAFKA_DESCRIBE_LIMIT {
				break
			} else {
				offset += pageSize
			}
		}
	}
	return
}

func (me *CkafkaService) CreateAcl(ctx context.Context, instanceId, resourceType, resourceName, operation, permissionType, host, principal string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreateAclRequest()
	request.InstanceId = &instanceId
	request.ResourceType = helper.Int64(CKAFKA_ACL_RESOURCE_TYPE[resourceType])
	request.ResourceName = &resourceName
	request.Operation = helper.Int64(CKAFKA_ACL_OPERATION[operation])
	request.PermissionType = helper.Int64(CKAFKA_PERMISSION_TYPE[permissionType])
	request.Host = &host
	request.Principal = helper.String(CKAFKA_ACL_PRINCIPAL_STR + principal)

	var response *ckafka.CreateAclResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().CreateAcl(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
	}
	return nil
}

func (me *CkafkaService) DescribeAclByFilter(ctx context.Context, params map[string]interface{}) (aclInfos []*ckafka.Acl, errRet error) {
	logId := tccommon.GetLogId(ctx)

	instanceId := params["instance_id"].(string)
	if _, has, _ := me.DescribeInstanceById(ctx, instanceId); !has {
		return
	}
	resourceType := params["resource_type"].(string)
	resourceName := params["resource_name"].(string)
	if resourceType == "TOPIC" {
		if _, has, _ := me.DescribeTopicById(ctx, instanceId+tccommon.FILED_SP+resourceName); !has {
			return
		}
	}

	request := ckafka.NewDescribeACLRequest()
	var offset int64 = 0
	var pageSize = int64(CKAFKA_DESCRIBE_LIMIT)
	request.InstanceId = &instanceId
	request.ResourceType = helper.Int64(CKAFKA_ACL_RESOURCE_TYPE[resourceType])
	request.ResourceName = helper.String(resourceName)
	if host, ok := params["host"]; ok {
		request.SearchWord = helper.String(host.(string))
	}
	request.Limit = &pageSize
	request.Offset = &offset

	aclInfos = make([]*ckafka.Acl, 0)
	for {
		var response *ckafka.DescribeACLResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseCkafkaClient().DescribeACL(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			aclInfos = append(aclInfos, response.Response.Result.AclList...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		} else {
			if len(response.Response.Result.AclList) < CKAFKA_DESCRIBE_LIMIT {
				break
			} else {
				offset += pageSize
			}
		}
	}
	return
}

func (me *CkafkaService) DescribeAclByAclId(ctx context.Context, aclId string) (aclInfo *ckafka.Acl, has bool, errRet error) {
	// acl id is organized by "instanceId + tccommon.FILED_SP + permissionType + tccommon.FILED_SP + principal + tccommon.FILED_SP + host + tccommon.FILED_SP + operation + tccommon.FILED_SP + resourceType + tccommon.FILED_SP + resourceName"
	items := strings.Split(aclId, tccommon.FILED_SP)
	if len(items) != 7 {
		errRet = fmt.Errorf("id of resource.tencentcloud_ckafka_acl is wrong")
		return
	}
	instanceId, permission, principal, host, operation, resourceType, resourceName := items[0], items[1], items[2], items[3], items[4], items[5], items[6]

	var params = map[string]interface{}{
		"instance_id":   instanceId,
		"resource_type": resourceType,
		"resource_name": resourceName,
		"host":          host,
	}
	aclInfos, err := me.DescribeAclByFilter(ctx, params)
	if err != nil {
		errRet = err
		return
	}
	for _, acl := range aclInfos {
		if CKAFKA_PERMISSION_TYPE_TO_STRING[*acl.PermissionType] == permission && *acl.Principal == CKAFKA_ACL_PRINCIPAL_STR+principal && CKAFKA_ACL_OPERATION_TO_STRING[*acl.Operation] == operation {
			aclInfo = acl
			has = true
			return
		}
	}
	has = false
	return
}

func (me *CkafkaService) DeleteAcl(ctx context.Context, aclId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	// acl id is organized by "instanceId + tccommon.FILED_SP + permissionType + tccommon.FILED_SP + principal + tccommon.FILED_SP + host + tccommon.FILED_SP + operation + tccommon.FILED_SP + resourceType + tccommon.FILED_SP + resourceName"
	items := strings.Split(aclId, tccommon.FILED_SP)
	if len(items) != 7 {
		errRet = fmt.Errorf("id of resource.tencentcloud_ckafka_acl is wrong")
		return
	}
	instanceId, permission, principal, host, operation, resourceType, resourceName := items[0], items[1], items[2], items[3], items[4], items[5], items[6]

	request := ckafka.NewDeleteAclRequest()
	request.InstanceId = &instanceId
	request.ResourceType = helper.Int64(CKAFKA_ACL_RESOURCE_TYPE[resourceType])
	request.ResourceName = &resourceName
	request.Operation = helper.Int64(CKAFKA_ACL_OPERATION[operation])
	request.PermissionType = helper.Int64(CKAFKA_PERMISSION_TYPE[permission])
	request.Host = &host
	request.Principal = helper.String(CKAFKA_ACL_PRINCIPAL_STR + principal)

	var response *ckafka.DeleteAclResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, err = me.client.UseCkafkaClient().DeleteAcl(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && !me.OperateStatusCheck(ctx, response.Response.Result) {
		return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
	}
	return nil
}

func (me *CkafkaService) DescribeInstanceById(ctx context.Context, instanceId string) (instanceInfo *ckafka.InstanceAttributesResponse, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeInstanceAttributesRequest()
	request.InstanceId = &instanceId
	var response *ckafka.DescribeInstanceAttributesResponse
	var err error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseCkafkaClient().DescribeInstanceAttributes(request)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == CkafkaInstanceNotFound || sdkErr.Code == CkafkaFailedOperation {
					return nil
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}

	if response != nil && response.Response != nil {
		if instanceInfo = response.Response.Result; instanceInfo != nil {
			has = true
			return
		}
	}

	has = false
	return
}

func (me *CkafkaService) DescribeTopicById(ctx context.Context, topicId string) (topicInfo *ckafka.TopicAttributesResponse, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeTopicAttributesRequest()
	items := strings.Split(topicId, tccommon.FILED_SP)
	if len(items) != 2 {
		errRet = fmt.Errorf("id of resource.tencentcloud_ckafka_topic is wrong")
		return
	}
	instanceId, topicName := items[0], items[1]
	request.InstanceId = &instanceId
	request.TopicName = &topicName
	var response *ckafka.DescribeTopicAttributesResponse
	var err error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseCkafkaClient().DescribeTopicAttributes(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}

	if response != nil && response.Response != nil {
		if topicInfo = response.Response.Result; topicInfo != nil {
			has = true
			return
		}
	}

	has = false
	return
}

func (me *CkafkaService) DescribeCkafkaTopics(ctx context.Context, instanceId string, topicName string) (topicList []*ckafka.TopicDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewDescribeTopicDetailRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	if topicName != "" {
		request.SearchWord = &topicName
	}
	var offset, limit int64 = 0, 20
	request.Offset = &offset
	request.Limit = &limit
	//check ckafka exist
	_, ckafkaExist, errRet := me.DescribeInstanceById(ctx, instanceId)
	if errRet != nil {
		return
	}
	if !ckafkaExist {
		return
	}
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCkafkaClient().DescribeTopicDetail(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil || response.Response.Result == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		topicList = append(topicList, response.Response.Result.TopicList...)
		if len(response.Response.Result.TopicList) < int(limit) {
			break
		}
		offset += limit
	}
	return
}

func (me *CkafkaService) CreateCkafkaTopic(ctx context.Context, request *ckafka.CreateTopicRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *ckafka.CreateTopicResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		resp, e := me.client.UseCkafkaClient().CreateTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		return nil
	})
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	if *response.Response.Result.TopicId == "" {
		return fmt.Errorf("TencentCloud SDK returns empty ckafka topic ID, %s", request.GetAction())
	}
	return nil
}

func (me *CkafkaService) DescribeCkafkaTopicByName(ctx context.Context, instanceId string, topicName string) (topic *ckafka.TopicDetail, has bool, errRet error) {
	var topicList []*ckafka.TopicDetail
	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		list, err := me.DescribeCkafkaTopics(ctx, instanceId, topicName)
		if err != nil {
			return tccommon.RetryError(err)
		}
		topicList = list
		return nil
	})
	if errRet != nil {
		return
	}
	for _, v := range topicList {
		if *v.TopicName == topicName {
			has = true
			topic = v
			break
		}
	}
	return
}

func (me *CkafkaService) DescribeCkafkaTopicAttributes(ctx context.Context, instanceId string, topicName string) (topicInfo *ckafka.TopicAttributesResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewDescribeTopicAttributesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.TopicName = &topicName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCkafkaClient().DescribeTopicAttributes(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	topicInfo = response.Response.Result
	return
}

func (me *CkafkaService) AddCkafkaTopicIpWhiteList(ctx context.Context, instanceId string, topicName string, whiteIpList []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreateTopicIpWhiteListRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.TopicName = &topicName
	request.InstanceId = &instanceId
	request.IpWhiteList = whiteIpList
	var response *ckafka.CreateTopicIpWhiteListResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		resp, e := me.client.UseCkafkaClient().CreateTopicIpWhiteList(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		return nil
	})
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return nil
}

func (me *CkafkaService) AddCkafkaTopicPartition(ctx context.Context, instanceId string, topicName string, partitionNum int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreatePartitionRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.TopicName = &topicName
	request.InstanceId = &instanceId
	request.PartitionNum = &partitionNum
	var response *ckafka.CreatePartitionResponse
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		resp, e := me.client.UseCkafkaClient().CreatePartition(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		return nil
	})
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return nil
}

func (me *CkafkaService) RemoveCkafkaTopicIpWhiteList(ctx context.Context, instaneId string, topicName string, whiteIpList []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewDeleteTopicIpWhiteListRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.TopicName = &topicName
	request.InstanceId = &instaneId
	request.IpWhiteList = whiteIpList
	ratelimit.Check(request.GetAction())
	var response *ckafka.DeleteTopicIpWhiteListResponse
	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resp, e := me.client.UseCkafkaClient().DeleteTopicIpWhiteList(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = resp
		return nil
	})
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return nil
}

func (me *CkafkaService) DescribeCkafkaById(ctx context.Context, instanceId string) (instance *ckafka.InstanceDetail, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewDescribeInstancesDetailRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())
	resp, err := me.client.UseCkafkaClient().DescribeInstancesDetail(request)
	if err != nil {
		has = false
		return
	}
	for _, cKafkaInstance := range resp.Response.Result.InstanceList {
		if *cKafkaInstance.InstanceId == instanceId {
			has = true
			instance = cKafkaInstance
			break
		}
	}
	return
}

func (me *CkafkaService) ModifyCkafkaTopicAttribute(ctx context.Context, request *ckafka.ModifyTopicAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, errRet := me.client.UseCkafkaClient().ModifyTopicAttributes(request)
	if errRet != nil {
		return errRet
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		return fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *CkafkaService) DeleteCkafkaTopic(ctx context.Context, instanceId string, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewDeleteTopicRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.TopicName = &name

	ratelimit.Check(request.GetAction())
	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := me.client.UseCkafkaClient().DeleteTopic(request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if errRet != nil {
		return
	}
	//重试超时时间
	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		topicList, err := me.DescribeCkafkaTopics(ctx, instanceId, name)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if len(topicList) != 0 {
			return resource.RetryableError(fmt.Errorf("this Topic %s Delete Failed", name))
		}
		return nil
	})

	if errRet != nil {
		return errRet
	}
	return
}

func (me *CkafkaService) DescribeCkafkaDatahubTopicById(ctx context.Context, topicName string) (datahubTopic *ckafka.DescribeDatahubTopicResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeDatahubTopicRequest()
	request.Name = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeDatahubTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	datahubTopic = response.Response.Result
	return
}

func (me *CkafkaService) DeleteCkafkaDatahubTopicById(ctx context.Context, topicName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteDatahubTopicRequest()
	request.Name = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteDatahubTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CkafkaService) DescribeCkafkaConnectResourceById(ctx context.Context, resourceId string) (connectResource *ckafka.DescribeConnectResourceResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeConnectResourceRequest()
	request.ResourceId = &resourceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeConnectResource(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	connectResource = response.Response.Result
	return
}

func (me *CkafkaService) DeleteCkafkaConnectResourceById(ctx context.Context, resourceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteConnectResourceRequest()
	request.ResourceId = &resourceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteConnectResource(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CkafkaService) CkafkaConnectResourceStateRefreshFunc(resourceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeCkafkaConnectResourceById(ctx, resourceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.Status), nil
	}
}

func (me *CkafkaService) DescribeCkafkaConnectResourceByFilter(ctx context.Context, params map[string]interface{}) (describeConnectResourceResp *ckafka.DescribeConnectResourcesResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeConnectResourcesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	offset := 0
	limit := 20
	for k, v := range params {
		if k == "type" {
			request.Type = helper.String(v.(string))
		}
		if k == "search_word" {
			request.SearchWord = helper.String(v.(string))
		}
		if k == "resource_region" {
			request.ResourceRegion = helper.String(v.(string))
		}
		if k == "offset" {
			offset = v.(int)
		}
		if k == "limit" {
			limit = v.(int)
		}
	}

	request.Offset = helper.IntInt64(offset)
	request.Limit = helper.IntInt64(limit)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeConnectResources(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("Response is null")
		return
	}

	describeConnectResourceResp = response.Response.Result
	return
}

func (me *CkafkaService) DescribeCkafkaDatahubTopicByFilter(ctx context.Context, paramMap map[string]interface{}) (result *ckafka.DescribeDatahubTopicsResp, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeDatahubTopicsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	offset := 0
	limit := 50
	for k, v := range paramMap {
		if k == "search_word" {
			request.SearchWord = helper.String(v.(string))
		}
		if k == "offset" {
			offset = v.(int)
		}
		if k == "limit" {
			limit = v.(int)
		}
	}
	request.Limit = helper.IntUint64(limit)
	request.Offset = helper.IntUint64(offset)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCkafkaClient().DescribeDatahubTopics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response is null")
		return
	}

	result = response.Response.Result
	return
}

func (me *CkafkaService) DescribeCkafkaDatahubGroupOffsetsByFilter(ctx context.Context, param map[string]interface{}) (groupOffsetTopics []*ckafka.GroupOffsetTopic, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeDatahubGroupOffsetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "name" {
			request.Name = v.(*string)
		}
		if k == "group" {
			request.Group = v.(*string)
		}
		if k == "search_word" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeDatahubGroupOffsets(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil {
			break
		}
		groupOffsetTopics = append(groupOffsetTopics, response.Response.Result.TopicList...)
		if len(response.Response.Result.TopicList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CkafkaService) DescribeCkafkaDatahubTaskByFilter(ctx context.Context, param map[string]interface{}) (datahubTaskInfos []*ckafka.DatahubTaskInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeDatahubTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "search_word" {
			request.SearchWord = v.(*string)
		}
		if k == "target_type" {
			request.TargetType = v.(*string)
		}
		if k == "task_type" {
			request.TaskType = v.(*string)
		}
		if k == "source_type" {
			request.SourceType = v.(*string)
		}
		if k == "resource" {
			request.Resource = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeDatahubTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.TaskList) < 1 {
			break
		}
		datahubTaskInfos = append(datahubTaskInfos, response.Response.Result.TaskList...)
		if len(response.Response.Result.TaskList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CkafkaService) DescribeCkafkaGroupByFilter(ctx context.Context, param map[string]interface{}) (groups []*ckafka.DescribeGroup, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeGroupRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "search_word" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeGroup(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Result.GroupList) < 1 {
			break
		}
		groups = append(groups, response.Response.Result.GroupList...)
		if len(response.Response.Result.GroupList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CkafkaService) DescribeCkafkaGroupOffsetsByFilter(ctx context.Context, param map[string]interface{}) (groupOffsetTopics []*ckafka.GroupOffsetTopic, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeGroupOffsetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "group" {
			request.Group = v.(*string)
		}
		if k == "topics" {
			request.Topics = v.([]*string)
		}
		if k == "search_word" {
			request.SearchWord = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeGroupOffsets(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.TopicList) < 1 {
			break
		}
		groupOffsetTopics = append(groupOffsetTopics, response.Response.Result.TopicList...)
		if len(response.Response.Result.TopicList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CkafkaService) DescribeCkafkaGroupInfoByFilter(ctx context.Context, param map[string]interface{}) (groupInfo []*ckafka.GroupInfoResponse, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeGroupInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "group_list" {
			request.GroupList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeGroupInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	groupInfo = response.Response.Result

	return
}

func (me *CkafkaService) DescribeCkafkaTaskStatusByFilter(ctx context.Context, flowId int) (taskStatus *ckafka.TaskStatusResponse, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeTaskStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.FlowId = helper.IntInt64(flowId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeTaskStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	taskStatus = response.Response.Result

	return
}

func (me *CkafkaService) DescribeCkafkaTopicFlowRankingByFilter(ctx context.Context, param map[string]interface{}) (topicFlowRanking *ckafka.TopicFlowRankingResult, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeTopicFlowRankingRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "ranking_type" {
			request.RankingType = v.(*string)
		}
		if k == "begin_date" {
			request.BeginDate = v.(*string)
		}
		if k == "end_date" {
			request.EndDate = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeTopicFlowRanking(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	topicFlowRanking = response.Response.Result

	return
}

func (me *CkafkaService) DescribeCkafkaTopicProduceConnectionByFilter(ctx context.Context, param map[string]interface{}) (topicProduceConnection []*ckafka.DescribeConnectInfoResultDTO, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeTopicProduceConnectionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "topic_name" {
			request.TopicName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeTopicProduceConnection(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	topicProduceConnection = response.Response.Result

	return
}

func (me *CkafkaService) DescribeCkafkaTopicSubscribeGroupByFilter(ctx context.Context, param map[string]interface{}) (groupInfos []*ckafka.GroupInfoResponse, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeTopicSubscribeGroupRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "TopicName" {
			request.TopicName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeTopicSubscribeGroup(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.GroupsInfo) < 1 {
			break
		}
		groupInfos = append(groupInfos, response.Response.Result.GroupsInfo...)
		if len(response.Response.Result.GroupsInfo) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CkafkaService) DescribeCkafkaTopicSyncReplicaByFilter(ctx context.Context, param map[string]interface{}) (topicInSyncReplicaInfos []*ckafka.TopicInSyncReplicaInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeTopicSyncReplicaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "topic_name" {
			request.TopicName = v.(*string)
		}
		if k == "out_of_sync_replica_only" {
			request.OutOfSyncReplicaOnly = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  int64  = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCkafkaClient().DescribeTopicSyncReplica(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.TopicInSyncReplicaList) < 1 {
			break
		}
		topicInSyncReplicaInfos = append(topicInSyncReplicaInfos, response.Response.Result.TopicInSyncReplicaList...)
		if len(response.Response.Result.TopicInSyncReplicaList) < int(limit) {
			break
		}

		offset += uint64(limit)
	}

	return
}

func (me *CkafkaService) DescribeCkafkaAclRuleById(ctx context.Context, instanceId string, ruleName string) (aclRule *ckafka.AclRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeAclRuleRequest()
	request.InstanceId = &instanceId
	request.RuleName = &ruleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeAclRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.AclRuleList) < 1 {
		return
	}

	aclRule = response.Response.Result.AclRuleList[0]
	return
}

func (me *CkafkaService) DeleteCkafkaAclRuleById(ctx context.Context, instanceId string, ruleName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteAclRuleRequest()
	request.InstanceId = &instanceId
	request.RuleName = &ruleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteAclRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CkafkaService) DescribeCkafkaConsumerGroupById(ctx context.Context, instanceId string, groupName string) (consumerGroup *ckafka.ConsumerGroupResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeConsumerGroupRequest()
	request.InstanceId = &instanceId
	request.GroupName = &groupName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeConsumerGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Result == nil {
		return
	}

	consumerGroup = response.Response.Result
	return
}

func (me *CkafkaService) DeleteCkafkaConsumerGroupById(ctx context.Context, instanceId string, groupName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteGroupRequest()
	request.InstanceId = &instanceId
	request.Group = &groupName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CkafkaService) DeleteCkafkaDatahubTaskById(ctx context.Context, taskId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteDatahubTaskRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteDatahubTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CkafkaService) DescribeDatahubTask(ctx context.Context, taskId string) (result *ckafka.DescribeDatahubTaskRes, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeDatahubTaskRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeDatahubTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil || response.Response.Result == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	result = response.Response.Result
	return
}

func (me *CkafkaService) CkafkaDatahubTaskStateRefreshFunc(taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeDatahubTask(ctx, taskId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.Status), nil
	}
}

func (me *CkafkaService) DescribeCkafkaCkafkaZoneByFilter(ctx context.Context, param map[string]interface{}) (ckafkaZone *ckafka.ZoneResponse, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = ckafka.NewDescribeCkafkaZoneRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CdcId" {
			request.CdcId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeCkafkaZone(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	ckafkaZone = response.Response.Result

	return
}

func (me *CkafkaService) DescribeCkafkaRouteById(ctx context.Context, instanceId string, routeId int64) (route *ckafka.Route, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDescribeRouteRequest()
	request.InstanceId = &instanceId
	request.RouteId = &routeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DescribeRoute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil || len(response.Response.Result.Routers) < 1 {
		return
	}

	route = response.Response.Result.Routers[0]
	return
}

func (me *CkafkaService) DeleteCkafkaRouteById(ctx context.Context, instanceId string, routeId int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := ckafka.NewDeleteRouteRequest()
	request.InstanceId = &instanceId
	request.RouteId = &routeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCkafkaClient().DeleteRoute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	// response.Response.Result.Data.FlowId
	return
}

func (me *CkafkaService) CkafkaRouteStateRefreshFunc(flowId int64, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := ckafka.NewDescribeTaskStatusRequest()
		request.FlowId = helper.Int64(flowId)
		object, err := me.client.UseCkafkaClient().DescribeTaskStatus(request)

		if err != nil {
			return nil, "", err
		}
		status := strconv.FormatInt(*object.Response.Result.Status, 10)
		return object, status, nil
	}
}
