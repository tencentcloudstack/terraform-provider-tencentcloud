package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type TcmqService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
// tcmq instance

func (me *TcmqService) DescribeTcmqInstanceById(ctx context.Context,
	clusterId string) (info *tcmq.Cluster, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDescribeClustersRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterIdList = []*string{&clusterId}

	var response *tcmq.DescribeClustersResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeClusters(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tcmq failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.ClusterSet) < 1 {
		return
	}
	has = true
	info = response.Response.ClusterSet[0]
	return
}

func (me *TcmqService) ModifyTcmqInstanceAttribute(ctx context.Context, clusterId, clusterName string,
	remark string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewModifyClusterRequest()
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
		log.Printf("[CRITAL]%s modify tcmq failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TcmqService) DeleteTcmqInstance(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDeleteClusterRequest()
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

// CreateTcmqNamespace tcmq namespace
func (me *TcmqService) CreateTcmqNamespace(ctx context.Context, environName string, msgTtl uint64, clusterId string,
	remark string, retentionPolicy tcmq.RetentionPolicy) (environId string, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewCreateEnvironmentRequest()
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

	var response *tcmq.CreateEnvironmentResponse
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().CreateEnvironment(request)
		if err != nil {
			return retryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tcmq namespace failed, reason: %v", logId, err)
		errRet = err
		return
	}
	environId = *response.Response.EnvironmentId
	return
}

func (me *TcmqService) DescribeTcmqNamespaceById(ctx context.Context,
	environId string, clusterId string) (info *tcmq.Environment, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDescribeEnvironmentsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.ClusterId = &clusterId

	var response *tcmq.DescribeEnvironmentsResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeEnvironments(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tcmq failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.EnvironmentSet) < 1 {
		return
	}
	has = true
	info = response.Response.EnvironmentSet[0]
	return
}

func (me *TcmqService) ModifyTcmqNamespaceAttribute(ctx context.Context, environId string, msgTtl uint64,
	remark string, clusterId string, retentionPolicy *tcmq.RetentionPolicy) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewModifyEnvironmentAttributesRequest()
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
		log.Printf("[CRITAL]%s modify tcmq namespace failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TcmqService) DeleteTcmqNamespace(ctx context.Context, environId string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDeleteEnvironmentsRequest()
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
		log.Printf("[CRITAL]%s delete tcmq namespace failed, reason: %v", logId, err)
		return err
	}
	return
}

// tcmq topic
func (me *TcmqService) CreateTcmqTopic(ctx context.Context, environId string, topicName string, partitions uint64,
	topicType uint64, remark string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewCreateTopicRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.EnvironmentId = &environId
	request.TopicName = &topicName
	request.Partitions = &partitions
	request.TopicType = &topicType
	request.Remark = &remark
	request.ClusterId = &clusterId

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().CreateTopic(request)
		if err != nil {
			return retryError(err)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tcmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	return
}

func (me *TcmqService) DescribeTcmqTopicById(ctx context.Context, topicName string) (topic *tcmq.CmqTopic, errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDescribeCmqTopicDetailRequest()
	request.TopicName = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeCmqTopicDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TopicDescribe == nil {
		return
	}

	topic = response.Response.TopicDescribe
	return
}

func (me *TcmqService) ModifyTcmqTopicAttribute(ctx context.Context, environId string, topicName string,
	partitions uint64, remark string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewModifyTopicRequest()
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

func (me *TcmqService) DeleteTcmqTopicById(ctx context.Context, topicName string) (errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDeleteCmqTopicRequest()
	request.TopicName = &topicName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteCmqTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

//tcmq role
func (me *TcmqService) CreateTcmqRole(ctx context.Context, roleName string, clusterId string,
	remark string) (roleId string, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewCreateRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.RoleName = &roleName
	request.ClusterId = &clusterId
	request.Remark = &remark

	var response *tcmq.CreateRoleResponse
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().CreateRole(request)
		if err != nil {
			return retryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tcmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	roleId = *response.Response.RoleName
	return
}

func (me *TcmqService) DescribeTcmqRoleById(ctx context.Context,
	roleName string, clusterId string) (info *tcmq.Role, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDescribeRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.RoleName = &roleName
	request.ClusterId = &clusterId

	var response *tcmq.DescribeRolesResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tcmq role failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.RoleSets) < 1 {
		return
	}
	has = true
	info = response.Response.RoleSets[0]
	return
}

func (me *TcmqService) ModifyTcmqRoleAttribute(ctx context.Context, roleName string, clusterId string,
	remark string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewModifyRoleRequest()
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
		log.Printf("[CRITAL]%s modify tcmq role failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TcmqService) DeleteTcmqRole(ctx context.Context, roleName string, cluserId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDeleteRolesRequest()
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
		log.Printf("[CRITAL]%s delete tcmq roles failed, reason: %v", logId, err)
		return err
	}
	return
}

//tcmq role
func (me *TcmqService) CreateTcmqNamespaceRoleAttachment(ctx context.Context, environId string,
	roleName string, permissions []*string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewCreateEnvironmentRoleRequest()
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
		log.Printf("[CRITAL]%s create tcmq topic failed, reason: %v", logId, err)
		errRet = err
		return
	}
	return
}

func (me *TcmqService) DescribeTcmqNamespaceRoleAttachment(ctx context.Context,
	environId string, roleName string, clusterId string) (info *tcmq.EnvironmentRole, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDescribeEnvironmentRolesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.EnvironmentId = &environId
	request.RoleName = &roleName
	request.ClusterId = &clusterId

	var response *tcmq.DescribeEnvironmentRolesResponse

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().DescribeEnvironmentRoles(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s read tcmq environment role failed, reason: %v", logId, err)
		return nil, false, err
	}

	if len(response.Response.EnvironmentRoleSets) < 1 {
		return
	}
	has = true
	info = response.Response.EnvironmentRoleSets[0]
	return
}

func (me *TcmqService) ModifyTcmqNamespaceRoleAttachment(ctx context.Context,
	environId string, roleName string, permissions []*string, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewModifyEnvironmentRoleRequest()
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
		log.Printf("[CRITAL]%s modify tcmq environment role failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TcmqService) DeleteTcmqNamespaceRoleAttachment(ctx context.Context, environId string,
	roleName string, cluserId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcmq.NewDeleteEnvironmentRolesRequest()
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
		log.Printf("[CRITAL]%s delete tcmq environments roles failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *TcmqService) DescribeTcmqQueueById(ctx context.Context, queueName string) (queue *tcmq.CmqQueue, errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDescribeCmqQueueDetailRequest()
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

	queue = response.Response.QueueDescribe
	return
}

func (me *TcmqService) DeleteTcmqQueueById(ctx context.Context, queueName string) (errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDeleteCmqQueueRequest()
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

func (me *TcmqService) DescribeTcmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (subscribe *tcmq.CmqSubscription, errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDescribeCmqSubscriptionDetailRequest()
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

	if len(response.Response.SubscriptionSet) < 1 {
		return
	}

	subscribe = response.Response.SubscriptionSet[0]
	return
}

func (me *TcmqService) DeleteTcmqSubscribeById(ctx context.Context, topicName string, subscriptionName string) (errRet error) {
	logId := getLogId(ctx)

	request := tcmq.NewDeleteCmqSubscribeRequest()
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
