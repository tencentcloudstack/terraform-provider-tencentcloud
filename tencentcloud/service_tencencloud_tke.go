package tencentcloud

import (
	"context"
	"fmt"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
	"strings"
)

type ClusterBasicSetting struct {
	ClusterId          string
	ClusterOs          string
	ClusterVersion     string
	ClusterName        string
	ClusterDescription string
	VpcId              string
	ProjectId          int64
	ClusterNodeNum     int64
}

type ClusterAdvancedSettings struct {
	Ipvs             bool
	ContainerRuntime string
}

type RunInstancesForNode struct {
	Master []string
	Work   []string
}

type ClusterCidrSettings struct {
	ClusterCidr               string
	IgnoreClusterCidrConflict bool
	MaxNodePodNum             int64
	MaxClusterServiceNum      int64
}

type ClusterInfo struct {
	ClusterBasicSetting
	ClusterCidrSettings
	ClusterAdvancedSettings

	DeployType string
}

type InstanceInfo struct {
	InstanceId    string
	InstanceRole  string
	InstanceState string
	FailedReason  string
}

type TkeService struct {
	client *connectivity.TencentCloudClient
}

func (me *TkeService) DescribeClusterInstances(ctx context.Context, id string) (masters []InstanceInfo, workers []InstanceInfo, errRet error) {
	logId := getLogId(ctx)
	request := tke.NewDescribeClusterInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &id
	masters = make([]InstanceInfo, 0, 100)
	workers = make([]InstanceInfo, 0, 100)
	var offset int64 = 0
	var limit int64 = 20
	var has = map[string]bool{}
	var total int64 = -1

getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = int64(*response.Response.TotalCount)
	}

	if len(response.Response.InstanceSet) > 0 {
		offset += limit
	} else {
		//get empty set,we're done
		return
	}

	for _, item := range response.Response.InstanceSet {
		if has[*item.InstanceId] {
			errRet = fmt.Errorf("get repeated instance_id[%s] when doing DescribeClusterInstances", *item.InstanceId)
			return
		}
		has[*item.InstanceId] = true
		instanceInfo := InstanceInfo{
			InstanceId:    *item.InstanceId,
			InstanceRole:  *item.InstanceRole,
			InstanceState: *item.InstanceState,
			FailedReason:  *item.FailedReason,
		}
		if instanceInfo.InstanceRole == TKE_ROLE_MASTER_WORKER {
			workers = append(workers, instanceInfo)
		} else {
			masters = append(masters, instanceInfo)
		}
	}
	goto getMoreData

}

func (me *TkeService) DescribeClusters(ctx context.Context, id string, name string) (clusterInfos []ClusterInfo, errRet error) {

	logId := getLogId(ctx)
	request := tke.NewDescribeClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if id != "" && name != "" {
		errRet = fmt.Errorf("cluster_id ,cluster_name only one can be set one")
		return
	}

	if id != "" {
		request.ClusterIds = []*string{&id}
	}

	if name != "" {
		request.Filters = []*tke.Filter{&tke.Filter{
			Name:   stringToPointer("ClusterName"),
			Values: []*string{&name},
		}}
	}

	response, err := me.client.UseTkeClient().DescribeClusters(request)

	if err != nil {
		errRet = err
		return
	}

	lenClusters := len(response.Response.Clusters)

	if lenClusters == 0 {
		return
	}
	clusterInfos = make([]ClusterInfo, 0, lenClusters)

	for index := range response.Response.Clusters {
		cluster := response.Response.Clusters[index]
		var clusterInfo ClusterInfo

		clusterInfo.ClusterId = *cluster.ClusterId
		clusterInfo.ClusterOs = *cluster.ClusterOs
		clusterInfo.ClusterVersion = *cluster.ClusterVersion
		clusterInfo.ClusterDescription = *cluster.ClusterDescription
		clusterInfo.ClusterName = *cluster.ClusterName

		clusterInfo.ProjectId = int64(*cluster.ProjectId)
		clusterInfo.VpcId = *cluster.ClusterNetworkSettings.VpcId
		clusterInfo.ClusterNodeNum = int64(*cluster.ClusterNodeNum)

		clusterInfo.IgnoreClusterCidrConflict = *cluster.ClusterNetworkSettings.IgnoreClusterCIDRConflict
		clusterInfo.ClusterCidr = *cluster.ClusterNetworkSettings.ClusterCIDR
		clusterInfo.MaxClusterServiceNum = int64(*cluster.ClusterNetworkSettings.MaxClusterServiceNum)

		clusterInfo.MaxNodePodNum = int64(*cluster.ClusterNetworkSettings.MaxNodePodNum)
		clusterInfo.DeployType = strings.ToUpper(*cluster.ClusterType)
		clusterInfo.Ipvs = *cluster.ClusterNetworkSettings.Ipvs
		clusterInfos = append(clusterInfos, clusterInfo)
	}
	return
}

func (me *TkeService) DescribeCluster(ctx context.Context, id string) (
	clusterInfo ClusterInfo,
	has bool,
	errRet error) {

	logId := getLogId(ctx)
	request := tke.NewDescribeClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterIds = []*string{&id}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusters(request)

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Clusters) == 0 {
		return
	}

	has = true
	cluster := response.Response.Clusters[0]
	clusterInfo.ClusterId = *cluster.ClusterId
	clusterInfo.ClusterOs = *cluster.ClusterOs
	clusterInfo.ClusterVersion = *cluster.ClusterVersion
	clusterInfo.ClusterDescription = *cluster.ClusterDescription
	clusterInfo.ClusterName = *cluster.ClusterName

	clusterInfo.ProjectId = int64(*cluster.ProjectId)
	clusterInfo.VpcId = *cluster.ClusterNetworkSettings.VpcId
	clusterInfo.ClusterNodeNum = int64(*cluster.ClusterNodeNum)

	clusterInfo.IgnoreClusterCidrConflict = *cluster.ClusterNetworkSettings.IgnoreClusterCIDRConflict
	clusterInfo.ClusterCidr = *cluster.ClusterNetworkSettings.ClusterCIDR
	clusterInfo.MaxClusterServiceNum = int64(*cluster.ClusterNetworkSettings.MaxClusterServiceNum)

	clusterInfo.MaxNodePodNum = int64(*cluster.ClusterNetworkSettings.MaxNodePodNum)
	clusterInfo.DeployType = strings.ToUpper(*cluster.ClusterType)
	clusterInfo.Ipvs = *cluster.ClusterNetworkSettings.Ipvs

	return
}

func (me *TkeService) CreateCluster(ctx context.Context,
	basic ClusterBasicSetting,
	advanced ClusterAdvancedSettings,
	cvms RunInstancesForNode,
	cidrSetting ClusterCidrSettings) (id string, errRet error) {

	logId := getLogId(ctx)
	request := tke.NewCreateClusterRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterBasicSettings = &tke.ClusterBasicSettings{}
	request.ClusterBasicSettings.ClusterOs = &basic.ClusterOs
	request.ClusterBasicSettings.ClusterVersion = &basic.ClusterVersion
	request.ClusterBasicSettings.ProjectId = &basic.ProjectId
	request.ClusterBasicSettings.VpcId = &basic.VpcId
	request.ClusterBasicSettings.ClusterDescription = &basic.ClusterDescription
	request.ClusterBasicSettings.ClusterName = &basic.ClusterName

	request.ClusterAdvancedSettings = &tke.ClusterAdvancedSettings{}
	request.ClusterAdvancedSettings.IPVS = &advanced.Ipvs
	request.ClusterAdvancedSettings.ContainerRuntime = &advanced.ContainerRuntime

	request.RunInstancesForNode = []*tke.RunInstancesForNode{}

	if len(cvms.Master) != 0 {

		var node tke.RunInstancesForNode
		node.NodeRole = stringToPointer(TKE_ROLE_MASTER_ETCD)
		node.RunInstancesPara = []*string{}
		request.ClusterType = stringToPointer(TKE_DEPLOY_TYPE_INDEPENDENT)
		for v := range cvms.Master {
			node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Master[v])
		}
		request.RunInstancesForNode = append(request.RunInstancesForNode, &node)

	} else {
		request.ClusterType = stringToPointer(TKE_DEPLOY_TYPE_MANAGED)
	}

	var node tke.RunInstancesForNode
	node.NodeRole = stringToPointer(TKE_ROLE_MASTER_WORKER)
	node.RunInstancesPara = []*string{}
	for v := range cvms.Work {
		node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Work[v])
	}
	request.RunInstancesForNode = append(request.RunInstancesForNode, &node)

	request.ClusterCIDRSettings = &tke.ClusterCIDRSettings{}

	maxNodePodNum := uint64(cidrSetting.MaxNodePodNum)
	request.ClusterCIDRSettings.MaxNodePodNum = &maxNodePodNum

	maxClusterServiceNum := uint64(cidrSetting.MaxClusterServiceNum)
	request.ClusterCIDRSettings.MaxClusterServiceNum = &maxClusterServiceNum
	request.ClusterCIDRSettings.ClusterCIDR = &cidrSetting.ClusterCidr
	request.ClusterCIDRSettings.IgnoreClusterCIDRConflict = &cidrSetting.IgnoreClusterCidrConflict

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateCluster(request)

	if err != nil {
		errRet = err
		return
	}

	id = *response.Response.ClusterId
	return
}

func (me *TkeService) DeleteCluster(ctx context.Context, id string) (errRet error) {

	logId := getLogId(ctx)
	request := tke.NewDeleteClusterRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.InstanceDeleteMode = stringToPointer("terminate")

	_, err := me.client.UseTkeClient().DeleteCluster(request)

	return err
}
