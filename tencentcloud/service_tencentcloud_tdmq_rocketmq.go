package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
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

func (me *TdmqRocketmqService) DescribeTdmqRocketmqEnvironmentRole(ctx context.Context, clusterId, roleName, environmentId string) (environmentRoles []*tdmqRocketmq.EnvironmentRole, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeEnvironmentRolesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.RoleName = &roleName
	request.EnvironmentId = &environmentId
	environmentRoles = make([]*tdmqRocketmq.EnvironmentRole, 0)
	var offset int64 = 0
	var pageSize int64 = 100

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeEnvironmentRoles(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EnvironmentRoleSets) < 1 {
			break
		}
		environmentRoles = append(environmentRoles, response.Response.EnvironmentRoleSets...)
		if len(response.Response.EnvironmentRoleSets) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqRocketmqService) DeleteTdmqRocketmqEnvironmentRoleById(ctx context.Context, clusterId, roleName, environmentId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmqRocketmq.NewDeleteEnvironmentRolesRequest()

	request.ClusterId = &clusterId
	request.RoleNames = []*string{&roleName}
	request.EnvironmentId = &environmentId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdmqClient().DeleteEnvironmentRoles(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqRocketmqService) DescribeRocketmqClusterByFilter(ctx context.Context, param map[string]interface{}) (cluster []*tdmqRocketmq.RocketMQClusterDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdmqRocketmq.NewDescribeRocketMQClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if v, ok := param["id_keyword"]; ok {
		request.IdKeyword = helper.String(v.(string))
	}
	if v, ok := param["name_keyword"]; ok {
		request.NameKeyword = helper.String(v.(string))
	}
	if v, ok := param["cluster_id_list"]; ok {
		request.ClusterIdList = make([]*string, 0)
		for _, cluster := range v.([]interface{}) {
			clusterId := cluster.(string)
			request.ClusterIdList = append(request.ClusterIdList, &clusterId)
		}
	}

	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTdmqClient().DescribeRocketMQClusters(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterList) < 1 {
			break
		}
		cluster = append(cluster, response.Response.ClusterList...)
		if len(response.Response.ClusterList) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqNamespaceByFilter(ctx context.Context, param map[string]interface{}) (namespace []*tdmqRocketmq.RocketMQNamespace, errRet error) {
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

	for k, v := range param {
		if k == "cluster_id" {
			request.ClusterId = helper.String(v.(string))
		}

		if k == "name_keyword" {
			request.NameKeyword = helper.String(v.(string))
		}

	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
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

		if response == nil || len(response.Response.Namespaces) < 1 {
			break
		}
		namespace = append(namespace, response.Response.Namespaces...)
		if len(response.Response.Namespaces) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqTopicByFilter(ctx context.Context, param map[string]interface{}) (topic []*tdmqRocketmq.RocketMQTopic, errRet error) {
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

	for k, v := range param {
		if k == "cluster_id" {
			request.ClusterId = helper.String(v.(string))
		}

		if k == "namespace_id" {
			request.NamespaceId = helper.String(v.(string))
		}

		if k == "filter_type" {
			filterTypes := make([]*string, 0)
			for _, item := range v.([]interface{}) {
				fileterType := item.(string)
				filterTypes = append(filterTypes, &fileterType)
			}
			request.FilterType = filterTypes
		}

		if k == "filter_name" {
			request.FilterName = helper.String(v.(string))
		}

	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

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
		topic = append(topic, response.Response.Topics...)
		if len(response.Response.Topics) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqRoleByFilter(ctx context.Context, param map[string]interface{}) (role []*tdmqRocketmq.Role, errRet error) {
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
	request.ClusterId = helper.String(param["cluster_id"].(string))
	if v, ok := param["role_name"]; ok {
		request.RoleName = helper.String(v.(string))
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

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
		role = append(role, response.Response.RoleSets...)
		if len(response.Response.RoleSets) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqRocketmqService) DescribeTdmqRocketmqGroupByFilter(ctx context.Context, param map[string]interface{}) (group []*tdmqRocketmq.RocketMQGroup, errRet error) {
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

	for k, v := range param {
		if k == "cluster_id" {
			request.ClusterId = helper.String(v.(string))
		}

		if k == "namespace_id" {
			request.NamespaceId = helper.String(v.(string))
		}

		if k == "filter_topic" {
			request.FilterTopic = helper.String(v.(string))
		}

		if k == "filter_group" {
			request.FilterGroup = helper.String(v.(string))
		}

		if k == "filter_one_group" {
			request.FilterOneGroup = helper.String(v.(string))
		}

	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

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
		group = append(group, response.Response.Groups...)
		if len(response.Response.Groups) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TdmqService) DescribeTdmqEnvironmentAttributesByFilter(ctx context.Context, param map[string]interface{}) (environmentAttributes *tdmq.DescribeEnvironmentAttributesResponseParams, errRet error) {
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
		if k == "EnvironmentId" {
			request.EnvironmentId = v.(*string)
		}

		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeEnvironmentAttributes(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	environmentAttributes = response.Response
	return
}

func (me *TdmqService) DescribeTdmqPublisherSummaryByFilter(ctx context.Context, param map[string]interface{}) (publisherSummary *tdmq.DescribePublisherSummaryResponseParams, errRet error) {
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

	response, err := me.client.UseTdmqClient().DescribePublisherSummary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	publisherSummary = response.Response

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
			request.Sort = v.(*tdmq.Sort)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
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

		if response == nil || *response.Response.TotalCount == 0 {
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

func (me *TdmqService) DescribeTdmqSubscriptionAttachmentById(ctx context.Context, environmentId, Topic, subscriptionName, clusterId string) (subscriptionAttachment *tdmq.Subscription, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeSubscriptionsRequest()
	request.EnvironmentId = &environmentId
	request.TopicName = &Topic
	request.SubscriptionName = &subscriptionName
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

	if *response.Response.TotalCount == 0 {
		return
	}

	subscriptionAttachment = response.Response.SubscriptionSets[0]
	return
}

func (me *TdmqService) GetTdmqTopicsAttachmentById(ctx context.Context, environmentId, Topic, subscriptionName, clusterId string) (has bool, errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDescribeTopicsRequest()
	topicRetry := fmt.Sprint(Topic + "-" + subscriptionName + "-" + "RETRY")
	topicDLQ := fmt.Sprint(Topic + "-" + subscriptionName + "-" + "DLQ")

	request.EnvironmentId = &environmentId
	request.ClusterId = &clusterId

	request.Filters = []*tdmq.Filter{
		{
			Name:   common.StringPtr("TopicName"),
			Values: common.StringPtrs([]string{topicRetry, topicDLQ}),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeTopics(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return false, nil
	}

	return true, nil
}

func (me *TdmqService) DeleteTdmqTopicsAttachmentById(ctx context.Context, environmentId, Topic, subscriptionName, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteTopicsRequest()
	topicRetry := fmt.Sprint(Topic + "-" + subscriptionName + "-" + "RETRY")
	topicDLQ := fmt.Sprint(Topic + "-" + subscriptionName + "-" + "DLQ")
	request.TopicSets = []*tdmq.TopicRecord{
		{
			EnvironmentId: &environmentId,
			TopicName:     &topicRetry,
		},
		{
			EnvironmentId: &environmentId,
			TopicName:     &topicDLQ,
		},
	}
	request.ClusterId = &clusterId
	request.EnvironmentId = &environmentId
	request.Force = common.BoolPtr(true)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DeleteTopics(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdmqService) DeleteTdmqSubscriptionAttachmentById(ctx context.Context, environmentId, Topic, subscriptionName, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdmq.NewDeleteSubscriptionsRequest()
	request.SubscriptionTopicSets = []*tdmq.SubscriptionTopic{
		{
			EnvironmentId:    &environmentId,
			TopicName:        &Topic,
			SubscriptionName: &subscriptionName,
		},
	}
	request.ClusterId = &clusterId
	request.EnvironmentId = &environmentId

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
		offset uint64 = 0
		limit  uint64 = 20
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

		if response == nil || *response.Response.TotalCount == 0 {
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
		offset uint64 = 0
		limit  uint64 = 20
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

		if response == nil || *response.Response.TotalCount == 0 {
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
		offset uint64 = 0
		limit  uint64 = 20
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

		if response == nil || *response.Response.TotalCount == 0 {
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

func (me *TdmqService) DescribeTdmqVipInstanceByFilter(ctx context.Context, param map[string]interface{}) (vipInstance *tdmq.DescribeRocketMQVipInstanceDetailResponseParams, errRet error) {
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

	response, err := me.client.UseTdmqClient().DescribeRocketMQVipInstanceDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	vipInstance = response.Response
	return
}

func (me *TdmqService) DescribeTdmqProInstanceDetailByFilter(ctx context.Context, param map[string]interface{}) (proInstanceDetail *tdmq.DescribePulsarProInstanceDetailResponseParams, errRet error) {
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

	response, err := me.client.UseTdmqClient().DescribePulsarProInstanceDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	proInstanceDetail = response.Response

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
		offset uint64 = 0
		limit  uint64 = 20
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

func (me *TdmqService) DescribeTdmqMessageByFilter(ctx context.Context, param map[string]interface{}) (message *tdmq.DescribeRocketMQMsgResponseParams, errRet error) {
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
			request.PulsarMsgId = v.(*string)
		}
		if k == "QueryDlqMsg" {
			request.QueryDlqMsg = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTdmqClient().DescribeRocketMQMsg(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	message = response.Response

	return
}
