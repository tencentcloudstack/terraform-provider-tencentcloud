package tencentcloud

import (
	"context"
	"log"

	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TdmqRocketmqService struct {
	client *connectivity.TencentCloudClient
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqCluster(ctx context.Context, clusterId string) (cluster *tdmqRocketmq.RocketMQClusterInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRocketMQClusterRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId

	response, err := me.client.UseTdmqClient().DescribeRocketMQCluster(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	cluster = response.Response.ClusterInfo
	return
}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqClusterById(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteRocketMQClusterRequest()

	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteRocketMQCluster(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqNamespace(ctx context.Context, namespaceName, clusterId string) (namespace []*tdmqRocketmq.RocketMQNamespace, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRocketMQNamespacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NameKeyword = &namespaceName
	offset := uint64(0)
	limit := uint64(100)
	namespace = make([]*tdmqRocketmq.RocketMQNamespace, 0)
	for {
		request.Limit = &limit
		request.Offset = &offset
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeRocketMQNamespaces(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		namespaces := response.Response.Namespaces
		if len(namespaces) > 0 {
			namespace = append(namespace, namespaces...)
		}
		if len(namespaces) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqNamespaceById(ctx context.Context, namespaceName, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteRocketMQNamespaceRequest()

	request.NamespaceId = &namespaceName
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteRocketMQNamespace(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqRole(ctx context.Context, clusterId, roleName string) (role *tdmqRocketmq.Role, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRolesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId
	request.Filters = append(
		request.Filters,
		&tdmqRocketmq.Filter{
			Name:   helper.String("RoleName"),
			Values: []*string{&roleName},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*tdmqRocketmq.Role, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeRoles(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RoleSets) < 1 {
			break
		}
		instances = append(instances, response.Response.RoleSets...)
		if len(response.Response.RoleSets) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	role = instances[0]

	return

}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqRoleById(ctx context.Context, clusterId, roleName string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteRolesRequest()

	request.ClusterId = &clusterId
	request.RoleNames = []*string{&roleName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteRoles(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqTopic(ctx context.Context, clusterId, namespaceId, topicName string) (result []*tdmqRocketmq.RocketMQTopic, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRocketMQTopicsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.FilterName = &topicName

	var offset uint64 = 0
	var pageSize uint64 = 100
	result = make([]*tdmqRocketmq.RocketMQTopic, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeRocketMQTopics(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Topics) < 1 {
			break
		}
		result = append(result, response.Response.Topics...)
		if len(response.Response.Topics) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	return
}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqTopicById(ctx context.Context, clusterId, namespaceId, topic string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteRocketMQTopicRequest()

	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.Topic = &topic

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteRocketMQTopic(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqGroup(ctx context.Context, clusterId, namespaceId, groupId string) (result []*tdmqRocketmq.RocketMQGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRocketMQGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.FilterGroup = &groupId

	var offset uint64 = 0
	var pageSize uint64 = 100
	result = make([]*tdmqRocketmq.RocketMQGroup, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeRocketMQGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Groups) < 1 {
			break
		}
		result = append(result, response.Response.Groups...)
		if len(response.Response.Groups) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	return
}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqGroupById(ctx context.Context, clusterId, namespaceId, groupId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteRocketMQGroupRequest()

	request.ClusterId = &clusterId
	request.NamespaceId = &namespaceId
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteRocketMQGroup(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
