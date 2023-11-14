package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type TdmqService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
// tdmq instance

func (me *TdmqService) DescribeTdmqInstanceById(ctx context.Context,
	clusterId string) (info *tdmq.Cluster, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDescribeClustersRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterIdList = []*string{&clusterId}

	var response *tdmq.DescribeClustersResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeClusters(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tdmq failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.ClusterSet) < 1 {
		return
	}
	has = true
	info = response.Response.ClusterSet[0]
	return
}

func (me *TdmqService) ModifyTdmqInstanceAttribute(ctx context.Context, clusterId, clusterName string,
	remark string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewModifyClusterRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.ClusterName = &clusterName
	request.Remark = &remark

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().ModifyCluster(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify tdmq failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DeleteTdmqInstance(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDeleteClusterRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	response, err := me.client.UseTdmqClient().DeleteCluster(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

// tdmq namespace
func (me *TdmqService) CreateTdmqNamespace(ctx context.Context, environName string, msgTtl uint64, clusterId string,
	remark string, retentionPolicy tdmq.RetentionPolicy) (environId string, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewCreateEnvironmentRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.EnvironmentId = &environName
	request.MsgTTL = &msgTtl
	request.ClusterId = &clusterId
	request.Remark = &remark
	request.RetentionPolicy = &retentionPolicy

	var response *tdmq.CreateEnvironmentResponse
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().CreateEnvironment(request)
		if err != nil {
			return retryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tdmq namespace failed, reason: %v", logId, err)
		errRet = err
		return
	}
	environId = *response.Response.EnvironmentId
	return
}

func (me *TdmqService) DescribeTdmqNamespaceById(ctx context.Context,
	environId string, clusterId string) (info *tdmq.Environment, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDescribeEnvironmentsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.ClusterId = &clusterId

	var response *tdmq.DescribeEnvironmentsResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeEnvironments(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tdmq failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.EnvironmentSet) < 1 {
		return
	}
	has = true
	info = response.Response.EnvironmentSet[0]
	return
}

func (me *TdmqService) ModifyTdmqNamespaceAttribute(ctx context.Context, environId string, msgTtl uint64,
	remark string, clusterId string, retentionPolicy *tdmq.RetentionPolicy) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewModifyEnvironmentAttributesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.MsgTTL = &msgTtl
	request.Remark = &remark
	request.ClusterId = &clusterId
	request.RetentionPolicy = retentionPolicy

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().ModifyEnvironmentAttributes(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify tdmq namespace failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DeleteTdmqNamespace(ctx context.Context, environId string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDeleteEnvironmentsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentIds = []*string{&environId}
	request.ClusterId = &clusterId
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().DeleteEnvironments(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete tdmq namespace failed, reason: %v", logId, err)
		return err
	}
	return
}

// tdmq topic
func (me *TdmqService) CreateTdmqTopic(ctx context.Context, environId string, topicName string, partitions uint64,
	topicType int64, remark string, clusterId string, pulsarTopicType int64) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewCreateTopicRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.EnvironmentId = &environId
	request.TopicName = &topicName
	request.Partitions = &partitions
	if topicType != NoneTopicType {
		request.TopicType = common.Uint64Ptr(uint64(topicType))
	}
	request.Remark = &remark
	request.ClusterId = &clusterId
	if pulsarTopicType != NonePulsarTopicType {
		request.PulsarTopicType = &pulsarTopicType
	}

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().CreateTopic(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tdmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	return
}

func (me *TdmqService) DescribeTdmqTopicById(ctx context.Context,
	environId string, topicName string, clusterId string) (info *tdmq.Topic, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDescribeTopicsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.TopicName = &topicName
	request.ClusterId = &clusterId

	var response *tdmq.DescribeTopicsResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeTopics(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tdmq failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.TopicSets) < 1 {
		return
	}
	has = true
	info = response.Response.TopicSets[0]
	return
}

func (me *TdmqService) ModifyTdmqTopicAttribute(ctx context.Context, environId string, topicName string,
	partitions uint64, remark string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewModifyTopicRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.TopicName = &topicName
	request.Partitions = &partitions
	request.Remark = &remark
	request.ClusterId = &clusterId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().ModifyTopic(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify tdmq topic failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DeleteTdmqTopic(ctx context.Context, environId string, topicName string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDeleteTopicsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	var (
		topicRecord tdmq.TopicRecord
	)
	topicRecord.TopicName = &topicName
	topicRecord.EnvironmentId = &environId
	request.TopicSets = []*tdmq.TopicRecord{&topicRecord}
	request.ClusterId = &clusterId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().DeleteTopics(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete tdmq topic failed, reason: %v", logId, err)
		return err
	}
	return
}

//tdmq role
func (me *TdmqService) CreateTdmqRole(ctx context.Context, roleName string, clusterId string,
	remark string) (roleId string, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewCreateRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.RoleName = &roleName
	request.ClusterId = &clusterId
	request.Remark = &remark

	var response *tdmq.CreateRoleResponse
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().CreateRole(request)
		if err != nil {
			return retryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tdmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	roleId = *response.Response.RoleName
	return
}

func (me *TdmqService) DescribeTdmqRoleById(ctx context.Context,
	roleName string, clusterId string) (info *tdmq.Role, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDescribeRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.RoleName = &roleName
	request.ClusterId = &clusterId

	var response *tdmq.DescribeRolesResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tdmq role failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.RoleSets) < 1 {
		return
	}
	has = true
	info = response.Response.RoleSets[0]
	return
}

func (me *TdmqService) ModifyTdmqRoleAttribute(ctx context.Context, roleName string, clusterId string,
	remark string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewModifyRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.RoleName = &roleName
	request.ClusterId = &clusterId
	request.Remark = &remark

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().ModifyRole(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify tdmq role failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DeleteTdmqRole(ctx context.Context, roleName string, cluserId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDeleteRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.RoleNames = []*string{&roleName}
	request.ClusterId = &cluserId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().DeleteRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete tdmq roles failed, reason: %v", logId, err)
		return err
	}
	return
}

//tdmq role
func (me *TdmqService) CreateTdmqNamespaceRoleAttachment(ctx context.Context, environId string,
	roleName string, permissions []*string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewCreateEnvironmentRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.EnvironmentId = &environId
	request.RoleName = &roleName
	request.Permissions = permissions
	request.ClusterId = &clusterId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().CreateEnvironmentRole(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tdmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	return
}

func (me *TdmqService) DescribeTdmqNamespaceRoleAttachment(ctx context.Context,
	environId string, roleName string, clusterId string) (info *tdmq.EnvironmentRole, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDescribeEnvironmentRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.RoleName = &roleName
	request.ClusterId = &clusterId

	var response *tdmq.DescribeEnvironmentRolesResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeEnvironmentRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tdmq environment role failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.EnvironmentRoleSets) < 1 {
		return
	}
	has = true
	info = response.Response.EnvironmentRoleSets[0]
	return
}

func (me *TdmqService) ModifyTdmqNamespaceRoleAttachment(ctx context.Context,
	environId string, roleName string, permissions []*string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewModifyEnvironmentRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.RoleName = &roleName
	request.ClusterId = &clusterId
	request.Permissions = permissions

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().ModifyEnvironmentRole(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify tdmq environment role failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DeleteTdmqNamespaceRoleAttachment(ctx context.Context, environId string,
	roleName string, cluserId string) (errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewDeleteEnvironmentRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.EnvironmentId = &environId
	request.RoleNames = []*string{&roleName}
	request.ClusterId = &cluserId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().DeleteEnvironmentRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete tdmq environments roles failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TdmqService) DescribeTdmqQueueByFilter(ctx context.Context, param map[string]interface{}) (queue []*tdmq.CmqQueue, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeCmqQueuesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "QueueName" {
			request.QueueName = v.(*string)
		}
		if k == "QueueNameList" {
			request.QueueNameList = v.([]*string)
		}
		if k == "IsTagFilter" {
			request.IsTagFilter = v.(*bool)
		}
		if k == "Filters" {
			request.Filters = v.([]*tdmq.Filter)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "QueueList" {
			request.QueueList = v.([]*tdmq.CmqQueue)
		}
		if k == "RequestId" {
			request.RequestId = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeCmqQueues(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.QueueList) < 1 {
			break
		}
		queue = append(queue, response.Response.QueueList...)
		if len(response.Response.QueueList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqSubscribeByFilter(ctx context.Context, param map[string]interface{}) (subscribe []*tdmq.CmqSubscription, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeCmqSubscriptionDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TopicName" {
			request.TopicName = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "SubscriptionName" {
			request.SubscriptionName = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "SubscriptionSet" {
			request.SubscriptionSet = v.([]*tdmq.CmqSubscription)
		}
		if k == "RequestId" {
			request.RequestId = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeCmqSubscriptionDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SubscriptionSet) < 1 {
			break
		}
		subscribe = append(subscribe, response.Response.SubscriptionSet...)
		if len(response.Response.SubscriptionSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqQueueById(ctx context.Context, queueName string) (queue *tdmq.CmqQueue, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeCmqQueueDetailRequest()
	request.QueueName = &queueName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqQueueDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.CmqQueue) < 1 {
		return
	}

	queue = response.Response.CmqQueue[0]
	return
}

func (me *TdmqService) DeleteTdmqQueueById(ctx context.Context, queueName string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteCmqQueueRequest()
	request.QueueName = &queueName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqQueue(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DescribeTdmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (subscribe *tdmq.CmqSubscription, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeCmqSubscriptionDetailRequest()
	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqSubscriptionDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.CmqSubscription) < 1 {
		return
	}

	subscribe = response.Response.CmqSubscription[0]
	return
}

func (me *TdmqService) DeleteTdmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteCmqSubscribeRequest()
	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqSubscribe(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DescribeTdmqDeadLetterSourceQueueByFilter(ctx context.Context, param map[string]interface{}) (deadLetterSourceQueue []*tdmq.CmqDeadLetterSource, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeCmqDeadLetterSourceQueuesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DeadLetterQueueName" {
			request.DeadLetterQueueName = v.(*string)
		}
		if k == "SourceQueueName" {
			request.SourceQueueName = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeCmqDeadLetterSourceQueues(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.QueueSet) < 1 {
			break
		}
		deadLetterSourceQueue = append(deadLetterSourceQueue, response.Response.QueueSet...)
		if len(response.Response.QueueSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqSubscriptionAttachmentById(ctx context.Context, clusterId string) (subscriptionAttachment *tdmq.Subscription, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeSubscriptionsRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeSubscriptions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Subscription) < 1 {
		return
	}

	subscriptionAttachment = response.Response.Subscription[0]
	return
}

func (me *TdmqService) DeleteTdmqSubscriptionAttachmentById(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteSubscriptionsRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteSubscriptions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DescribeTdmqEnvironmentAttributesByFilter(ctx context.Context, param map[string]interface{}) (environmentAttributes []*tdmq.DescribeEnvironmentAttributesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeEnvironmentAttributesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeEnvironmentAttributes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MsgTTL) < 1 {
			break
		}
		environmentAttributes = append(environmentAttributes, response.Response.MsgTTL...)
		if len(response.Response.MsgTTL) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqPublisherSummaryByFilter(ctx context.Context, param map[string]interface{}) (publisherSummary []*tdmq.DescribePublisherSummaryResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribePublisherSummaryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "Topic" {
			request.Topic = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribePublisherSummary(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MsgRateIn) < 1 {
			break
		}
		publisherSummary = append(publisherSummary, response.Response.MsgRateIn...)
		if len(response.Response.MsgRateIn) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqPublishersByFilter(ctx context.Context, param map[string]interface{}) (publishers []*tdmq.Publisher, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribePublishersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "Topic" {
			request.Topic = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tdmq.Filter)
		}
		if k == "Sort" {
			request.Sort = v.(map[string]interface{})
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
		response, err := me.client.UseTdmqClient().DescribePublishers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Publishers) < 1 {
			break
		}
		publishers = append(publishers, response.Response.Publishers...)
		if len(response.Response.Publishers) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqProInstanceDetailByFilter(ctx context.Context, param map[string]interface{}) (proInstanceDetail []*tdmq.PulsarProClusterInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribePulsarProInstanceDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribePulsarProInstanceDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterInfo) < 1 {
			break
		}
		proInstanceDetail = append(proInstanceDetail, response.Response.ClusterInfo...)
		if len(response.Response.ClusterInfo) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqProInstancesByFilter(ctx context.Context, param map[string]interface{}) (proInstances []*tdmq.PulsarProInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribePulsarProInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*tdmq.Filter)
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
		response, err := me.client.UseTdmqClient().DescribePulsarProInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Instances) < 1 {
			break
		}
		proInstances = append(proInstances, response.Response.Instances...)
		if len(response.Response.Instances) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqRabbitmqNodeListByFilter(ctx context.Context, param map[string]interface{}) (rabbitmqNodeList []*tdmq.RabbitMQPrivateNode, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeRabbitMQNodeListRequest()
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
		if k == "NodeName" {
			request.NodeName = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tdmq.Filter)
		}
		if k == "SortElement" {
			request.SortElement = v.(*string)
		}
		if k == "SortOrder" {
			request.SortOrder = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeRabbitMQNodeList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NodeList) < 1 {
			break
		}
		rabbitmqNodeList = append(rabbitmqNodeList, response.Response.NodeList...)
		if len(response.Response.NodeList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqRabbitmqVipInstanceByFilter(ctx context.Context, param map[string]interface{}) (rabbitmqVipInstance []*tdmq.RabbitMQVipInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeRabbitMQVipInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*tdmq.Filter)
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
		response, err := me.client.UseTdmqClient().DescribeRabbitMQVipInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Instances) < 1 {
			break
		}
		rabbitmqVipInstance = append(rabbitmqVipInstance, response.Response.Instances...)
		if len(response.Response.Instances) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqRabbitmqVirtualHostListByFilter(ctx context.Context, param map[string]interface{}) (rabbitmqVirtualHostList []*tdmq.RabbitMQPrivateVirtualHost, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeRabbitMQVirtualHostListRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTdmqClient().DescribeRabbitMQVirtualHostList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.VirtualHostList) < 1 {
			break
		}
		rabbitmqVirtualHostList = append(rabbitmqVirtualHostList, response.Response.VirtualHostList...)
		if len(response.Response.VirtualHostList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqRabbitmqUserById(ctx context.Context, instanceId string) (rabbitmqUser *tdmq.RabbitMQUser, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeRabbitMQUserRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeRabbitMQUser(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RabbitMQUser) < 1 {
		return
	}

	rabbitmqUser = response.Response.RabbitMQUser[0]
	return
}

func (me *TdmqService) DeleteTdmqRabbitmqUserById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteRabbitMQUserRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteRabbitMQUser(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DescribeTdmqRabbitmqVirtualHostById(ctx context.Context, instanceId string) (rabbitmqVirtualHost *tdmq.RabbitMQVirtualHostInfo, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeRabbitMQVirtualHostRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeRabbitMQVirtualHost(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RabbitMQVirtualHostInfo) < 1 {
		return
	}

	rabbitmqVirtualHost = response.Response.RabbitMQVirtualHostInfo[0]
	return
}

func (me *TdmqService) DeleteTdmqRabbitmqVirtualHostById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteRabbitMQVirtualHostRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteRabbitMQVirtualHost(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DescribeTdmqMessageByFilter(ctx context.Context, param map[string]interface{}) (message []*tdmq.DescribeRocketMQMsgResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeRocketMQMsgRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "EnvironmentId" {
			request.EnvironmentId = v.(*string)
		}
		if k == "TopicName" {
			request.TopicName = v.(*string)
		}
		if k == "MsgId" {
			request.MsgId = v.(*string)
		}
		if k == "QueryDlqMsg" {
			request.QueryDlqMsg = v.(*bool)
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
		response, err := me.client.UseTdmqClient().DescribeRocketMQMsg(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Body) < 1 {
			break
		}
		message = append(message, response.Response.Body...)
		if len(response.Response.Body) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmqService) DescribeTdmqVipInstanceByFilter(ctx context.Context, param map[string]interface{}) (vipInstance []*tdmq.RocketMQClusterInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmq.NewDescribeRocketMQVipInstanceDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
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
		response, err := me.client.UseTdmqClient().DescribeRocketMQVipInstanceDetail(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterInfo) < 1 {
			break
		}
		vipInstance = append(vipInstance, response.Response.ClusterInfo...)
		if len(response.Response.ClusterInfo) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
