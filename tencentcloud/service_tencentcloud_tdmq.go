package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type TdmqNamespaceInfo struct {
	environId        string
	msgTtl           uint64
	remark           string
	retentionSize    uint64
	retentionMinutes uint64
}

type TdmqService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
// tdmq instance
func (me *TdmqService) CreateTdmqInstance(ctx context.Context, clusterName string, bindClusterId uint64,
	remark string) (clusterId string, errRet error) {
	logId := getLogId(ctx)
	request := tdmq.NewCreateClusterRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterName = &clusterName
	request.BindClusterId = &bindClusterId
	request.Remark = &remark

	var response *tdmq.CreateClusterResponse
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseTdmqClient().CreateCluster(request)
		if err != nil {
			return retryError(err)
		}
		response = result
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s create tdmq failed, reason: %v", logId, err)
		errRet = err
		return
	}
	clusterId = *response.Response.ClusterId
	return
}

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
	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseTdmqClient().DeleteCluster(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s delete tdmq failed, reason: %v", logId, err)
		return err
	}
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
	topicType uint64, remark string, clusterId string) (errRet error) {
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
		topicRecord  tdmq.TopicRecord
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
		log.Printf("[CRITAL]%s modify tdmq enviroment role failed, reason: %v", logId, err)
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