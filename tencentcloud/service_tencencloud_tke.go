package tencentcloud

import (
	"context"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type ClusterBasicSetting struct {
	ClusterOs          string
	ClusterVersion     string
	ClusterName        string
	ClusterDescription string
	VpcId              string
	ProjectId          int64
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

type TkeService struct {
	client *connectivity.TencentCloudClient
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


	request.ClusterBasicSettings =&tke.ClusterBasicSettings{}
	request.ClusterBasicSettings.ClusterOs = &basic.ClusterOs
	request.ClusterBasicSettings.ClusterVersion = &basic.ClusterVersion
	request.ClusterBasicSettings.ProjectId = &basic.ProjectId
	request.ClusterBasicSettings.VpcId = &basic.VpcId
	request.ClusterBasicSettings.ClusterDescription = &basic.ClusterDescription
	request.ClusterBasicSettings.ClusterName = &basic.ClusterName

	request.ClusterAdvancedSettings =&tke.ClusterAdvancedSettings{}
	request.ClusterAdvancedSettings.IPVS = &advanced.Ipvs
	request.ClusterAdvancedSettings.ContainerRuntime =&advanced.ContainerRuntime

	request.RunInstancesForNode = []*tke.RunInstancesForNode{}

	if len(cvms.Master)!=0{

		var node  tke.RunInstancesForNode
		node.NodeRole =stringToPointer(TKE_ROLE_MASTER_ETCD)
		node.RunInstancesPara =[]*string{}
		request.ClusterType = stringToPointer(TKE_DEPLOY_TYPE_INDEPENDENT)
		for v:=range cvms.Master{
			node.RunInstancesPara = append(node.RunInstancesPara,&cvms.Master[v])
		}
		request.RunInstancesForNode =  append(request.RunInstancesForNode,&node)

	}else{
		request.ClusterType = stringToPointer(TKE_DEPLOY_TYPE_MANAGED)
	}

	var node  tke.RunInstancesForNode
	node.NodeRole =stringToPointer(TKE_ROLE_MASTER_WORKER)
	node.RunInstancesPara =[]*string{}
	for v:=range cvms.Work{
		node.RunInstancesPara = append(node.RunInstancesPara,&cvms.Work[v])
	}
	request.RunInstancesForNode =  append(request.RunInstancesForNode,&node)

	request.ClusterCIDRSettings = &tke.ClusterCIDRSettings{}

	maxNodePodNum := uint64(cidrSetting.MaxNodePodNum)
	request.ClusterCIDRSettings.MaxNodePodNum = &maxNodePodNum


	maxClusterServiceNum := uint64(cidrSetting.MaxClusterServiceNum)
	request.ClusterCIDRSettings.MaxClusterServiceNum = &maxClusterServiceNum
	request.ClusterCIDRSettings.ClusterCIDR = &cidrSetting.ClusterCidr
	request.ClusterCIDRSettings.IgnoreClusterCIDRConflict = &cidrSetting.IgnoreClusterCidrConflict


	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateCluster(request)

	if err!=nil{
		errRet = err
		return
	}

	id = *response.Response.ClusterId
	return
}

func (me *TkeService)DeleteCluster(ctx context.Context,id string)(errRet error)  {

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

	return  err
}
