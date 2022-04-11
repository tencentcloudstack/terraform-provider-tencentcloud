package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type EksService struct {
	client *connectivity.TencentCloudClient
}

type EksClusterInfo struct {
	ClusterId        string
	ClusterName      string
	ClusterDesc      string
	K8SVersion       string
	VpcId            string
	SubnetIds        []string
	Status           string
	CreatedTime      string
	ServiceSubnetId  string
	ExtraParam       string
	DnsServers       []map[string]interface{}
	NeedDeleteCbs    bool
	EnableVpcCoreDNS bool
	Tags             map[string]string
}

type EksClusterCredentialResponse struct {
	Addresses []*tke.IPAddress `json:"Addresses,omitempty" name:"Addresses"`

	Credential *tke.ClusterCredential `json:"Credential,omitempty" name:"Credential"`

	PublicLB *tke.ClusterPublicLB `json:"PublicLB,omitempty" name:"PublicLB"`

	InternalLB *tke.ClusterInternalLB `json:"InternalLB,omitempty" name:"InternalLB"`

	KubeConfig string

	ProxyLB bool
}

var versionSuffix = regexp.MustCompile(`-eks\.\d*$`)

func getClusterInfo(cluster *tke.EksCluster) EksClusterInfo {
	var clusterInfo EksClusterInfo

	clusterInfo.ClusterId = *cluster.ClusterId
	clusterInfo.ClusterName = *cluster.ClusterName
	clusterInfo.ClusterDesc = *cluster.ClusterDesc
	clusterInfo.K8SVersion = versionSuffix.ReplaceAllString(*cluster.K8SVersion, "")
	clusterInfo.VpcId = *cluster.VpcId
	clusterInfo.Status = *cluster.Status
	clusterInfo.CreatedTime = *cluster.CreatedTime
	clusterInfo.ServiceSubnetId = *cluster.ServiceSubnetId
	clusterInfo.NeedDeleteCbs = *cluster.NeedDeleteCbs
	clusterInfo.EnableVpcCoreDNS = *cluster.EnableVpcCoreDNS

	if len(cluster.SubnetIds) > 0 {
		for _, i := range cluster.SubnetIds {
			clusterInfo.SubnetIds = append(clusterInfo.SubnetIds, *i)
		}
	}

	if len(cluster.DnsServers) > 0 {
		for _, i := range cluster.DnsServers {
			server := make(map[string]interface{})
			server["domain"] = *i.Domain
			var servers []string
			for _, s := range i.DnsServers {
				servers = append(servers, *s)
			}
			server["servers"] = servers
			clusterInfo.DnsServers = append(clusterInfo.DnsServers, server)
		}
	}

	if cluster.TagSpecification != nil && len(cluster.TagSpecification) > 0 {
		clusterInfo.Tags = make(map[string]string)
		for _, tag := range cluster.TagSpecification[0].Tags {
			clusterInfo.Tags[*tag.Key] = *tag.Value
		}
	}

	return clusterInfo
}

func (me *EksService) DescribeEKSClusters(ctx context.Context, id string, name string) (eksClusters []EksClusterInfo, errRet error) {

	logId := getLogId(ctx)
	request := tke.NewDescribeEKSClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if id != "" && name != "" {
		errRet = fmt.Errorf("cluster_id, cluster_name only one can be set one")
		return
	}

	if id != "" {
		request.ClusterIds = []*string{&id}
	}

	if name != "" {
		filter := &tke.Filter{
			Name:   helper.String("ClusterName"),
			Values: []*string{&name},
		}
		request.Filters = []*tke.Filter{filter}
	}

	response, err := me.client.UseTkeClient().DescribeEKSClusters(request)

	if err != nil {
		errRet = err
		return
	}

	lenClusters := len(response.Response.Clusters)

	if lenClusters == 0 {
		return
	}
	eksClusters = make([]EksClusterInfo, 0, lenClusters)

	for index := range response.Response.Clusters {
		cluster := response.Response.Clusters[index]
		eksClusters = append(eksClusters, getClusterInfo(cluster))
	}

	return eksClusters, nil
}

func (me *EksService) DescribeEksCluster(ctx context.Context, id string) (clusterInfo EksClusterInfo, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tke.NewDescribeEKSClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterIds = []*string{&id}

	ratelimit.Check("DescribeEksCluster")
	response, err := me.client.UseTkeClient().DescribeEKSClusters(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) == 0 {
		return
	}

	has = true
	clusterInfo = getClusterInfo(response.Response.Clusters[0])

	return
}

func (me *EksService) DescribeEksContainerInstanceById(ctx context.Context, id string) (instance *tke.EksCi, has bool, errRet error) {
	logId := getLogId(ctx)

	request := tke.NewDescribeEKSContainerInstancesRequest()
	request.EksCiIds = []*string{&id}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeEKSContainerInstances(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil && len(response.Response.EksCis) > 0 {
		has = true
		instance = response.Response.EksCis[0]
	}
	return
}

func (me *EksService) DescribeEksContainerInstancesByFilter(ctx context.Context, filters []*tke.Filter, limit uint64, offset uint64) (instances []*tke.EksCi, errRet error) {
	logId := getLogId(ctx)

	request := tke.NewDescribeEKSContainerInstancesRequest()
	if limit > 0 {
		request.Limit = &limit
	}

	if offset > 0 {
		request.Offset = &offset
	}

	if len(filters) > 0 {
		request.Filters = filters
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeEKSContainerInstances(request)

	if err != nil {
		errRet = err
		return
	}

	instances = response.Response.EksCis

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EksService) DescribeEKSClusterCredentialById(ctx context.Context, id string) (info *EksClusterCredentialResponse, errRet error) {
	request := tke.NewDescribeEKSClusterCredentialRequest()
	request.ClusterId = &id
	info, errRet = me.DescribeEKSClusterCredential(ctx, request)
	return
}

func (me *EksService) DescribeEKSClusterCredential(ctx context.Context, request *tke.DescribeEKSClusterCredentialRequest) (info *EksClusterCredentialResponse, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeEKSClusterCredential(request)

	if err != nil {
		errRet = err
		return
	}

	if body := response.Response; body != nil {
		info = &EksClusterCredentialResponse{
			Credential: body.Credential,
			Addresses:  body.Addresses,
			PublicLB:   body.PublicLB,
			InternalLB: body.InternalLB,
			ProxyLB:    *body.ProxyLB,
			KubeConfig: *body.Kubeconfig,
		}
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EksService) CreateEksCluster(ctx context.Context, request *tke.CreateEKSClusterRequest) (id string, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().CreateEKSCluster(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return *response.Response.ClusterId, nil
}

func (me *EksService) CreateEksContainerInstances(ctx context.Context, request *tke.CreateEKSContainerInstancesRequest) (id string, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateEKSContainerInstances(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil && len(response.Response.EksCiIds) > 0 {
		id = *response.Response.EksCiIds[0]
	}

	return
}

func (me *EksService) UpdateEksContainerInstances(ctx context.Context, request *tke.UpdateEKSContainerInstanceRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().UpdateEKSContainerInstance(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EksService) UpdateEksCluster(ctx context.Context, request *tke.UpdateEKSClusterRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().UpdateEKSCluster(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *EksService) DeleteEksCluster(ctx context.Context, request *tke.DeleteEKSClusterRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DeleteEKSCluster(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *EksService) DeleteEksContainerInstance(ctx context.Context, request *tke.DeleteEKSContainerInstancesRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DeleteEKSContainerInstances(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
