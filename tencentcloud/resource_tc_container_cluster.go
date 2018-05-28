package tencentcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	ccs "github.com/zqfan/tencentcloud-sdk-go/services/ccs/unversioned"
)

const (
	CLUSTER_NOT_FOUND_CODE    = -16009
	CLUSTER_LIFESTATE_RUNNING = "Running"
)

func resourceTencentCloudContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudContainerClusterCreate,
		Read:   resourceTencentCloudContainerClusterRead,
		Update: resourceTencentCloudContainerClusterUpdate,
		Delete: resourceTencentCloudContainerClusterDelete,

		Schema: map[string]*schema.Schema{
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"mem": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"os_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"require_wan_ip": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_vpc_gateway": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"root_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"goods_num": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_desc": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cvm_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sg_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_target": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"docker_graph_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kubernetes_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes_num": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nodes_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_mem": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudContainerClusterUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("the container cluster resource doesn't support update")
}

func resourceTencentCloudContainerClusterCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).ccsConn

	createReq := ccs.NewCreateClusterRequest()

	createReq.ClusterName = common.StringPtr(d.Get("cluster_name").(string))

	if cpuRaw, ok := d.GetOkExists("cpu"); ok {
		cpu := cpuRaw.(int)
		if cpu > 0 {
			createReq.CPU = common.IntPtr(cpu)
		}
	}

	if memRaw, ok := d.GetOkExists("mem"); ok {
		mem := memRaw.(int)
		if mem > 0 {
			createReq.Mem = common.IntPtr(mem)
		}
	}

	if osNameRaw, ok := d.GetOkExists("os_name"); ok {
		osName := osNameRaw.(string)
		if len(osName) > 0 {
			createReq.OSName = common.StringPtr(osName)
		}
	}

	if bandwidthRaw, ok := d.GetOkExists("bandwidth"); ok {
		bandwidth := bandwidthRaw.(int)
		createReq.Bandwidth = common.IntPtr(bandwidth)
	}

	if bandwidthTypeRaw, ok := d.GetOkExists("bandwidth_type"); ok {
		bandwidthType := bandwidthTypeRaw.(string)
		if len(bandwidthType) > 0 {
			createReq.BandwidthType = common.StringPtr(bandwidthType)
		}
	}

	if wanIpRaw, ok := d.GetOkExists("require_wan_ip"); ok {
		wanIp := wanIpRaw.(int)
		createReq.WanIp = common.IntPtr(wanIp)
	}

	if subnetIdRaw, ok := d.GetOkExists("subnet_id"); ok {
		subnetId := subnetIdRaw.(string)
		if len(subnetId) > 0 {
			createReq.SubnetId = common.StringPtr(subnetId)
		}
	}

	if isVpcGatewayRaw, ok := d.GetOkExists("is_vpc_gateway"); ok {
		isVpcGateway := isVpcGatewayRaw.(int)
		createReq.IsVpcGateway = common.IntPtr(isVpcGateway)
	}

	if storageSizeRaw, ok := d.GetOkExists("storage_size"); ok {
		storageSize := storageSizeRaw.(int)
		createReq.StorageSize = common.IntPtr(storageSize)
	}

	if storageTypeRaw, ok := d.GetOkExists("storage_type"); ok {
		storageType := storageTypeRaw.(string)
		if len(storageType) > 0 {
			createReq.StorageType = common.StringPtr(storageType)
		}
	}

	if rootSizeRaw, ok := d.GetOkExists("root_size"); ok {
		rootSize := rootSizeRaw.(int)
		createReq.RootSize = common.IntPtr(rootSize)
	}

	if rootTypeRaw, ok := d.GetOkExists("root_type"); ok {
		rootType := rootTypeRaw.(string)
		if len(rootType) > 0 {
			createReq.RootType = common.StringPtr(rootType)
		}
	}

	if userScriptRaw, ok := d.GetOkExists("user_script"); ok {
		userScript := userScriptRaw.(string)
		if len(userScript) > 0 {
			createReq.UserScript = common.StringPtr(userScript)
		}
	}

	if goodsNumRaw, ok := d.GetOkExists("goods_num"); ok {
		goodsNum := goodsNumRaw.(int)
		createReq.GoodsNum = common.IntPtr(goodsNum)
	}

	if passwordRaw, ok := d.GetOkExists("password"); ok {
		password := passwordRaw.(string)
		createReq.Password = common.StringPtr(password)
	}

	if keyIdRaw, ok := d.GetOkExists("key_id"); ok {
		keyId := keyIdRaw.(string)
		if len(keyId) > 0 {
			createReq.KeyId = common.StringPtr(keyId)
		}
	}

	if vpcIdRaw, ok := d.GetOkExists("vpc_id"); ok {
		vpcId := vpcIdRaw.(string)
		if len(vpcId) > 0 {
			createReq.VpcId = common.StringPtr(vpcId)
		}
	}

	if clusterDescRaw, ok := d.GetOkExists("cluster_desc"); ok {
		clusterDesc := clusterDescRaw.(string)
		if len(clusterDesc) > 0 {
			createReq.ClusterDesc = common.StringPtr(clusterDesc)
		}
	}

	if clusterCIDRRaw, ok := d.GetOkExists("cluster_cidr"); ok {
		clusterCIDR := clusterCIDRRaw.(string)
		if len(clusterCIDR) > 0 {
			createReq.ClusterCIDR = common.StringPtr(clusterCIDR)
		}
	}

	if cvmTypeRaw, ok := d.GetOkExists("cvm_type"); ok {
		cvmType := cvmTypeRaw.(string)
		if len(cvmType) > 0 {
			createReq.CVMType = common.StringPtr(cvmType)
		}
	}

	if periodRaw, ok := d.GetOkExists("period"); ok {
		period := periodRaw.(int)
		createReq.Period = common.IntPtr(period)
	}

	if zoneIdRaw, ok := d.GetOkExists("zone_id"); ok {
		zoneId := zoneIdRaw.(string)
		if len(zoneId) > 0 {
			createReq.ZoneId = common.StringPtr(zoneId)
		}
	}

	if instanceTypeRaw, ok := d.GetOkExists("instance_type"); ok {
		instanceType := instanceTypeRaw.(string)
		if len(instanceType) > 0 {
			createReq.InstanceType = common.StringPtr(instanceType)
		}
	}

	if sgIdRaw, ok := d.GetOkExists("sg_id"); ok {
		sgId := sgIdRaw.(string)
		if len(sgId) > 0 {
			createReq.SgId = common.StringPtr(sgId)
		}
	}

	if mountTargetRaw, ok := d.GetOkExists("mount_target"); ok {
		mountTarget := mountTargetRaw.(string)
		if len(mountTarget) > 0 {
			createReq.MountTarget = common.StringPtr(mountTarget)
		}
	}

	if dockerGraphPathRaw, ok := d.GetOkExists("docker_graph_path"); ok {
		dockerGraphPath := dockerGraphPathRaw.(string)
		if len(dockerGraphPath) > 0 {
			createReq.DockerGraphPath = common.StringPtr(dockerGraphPath)
		}
	}

	if instanceNameRaw, ok := d.GetOkExists("instance_name"); ok {
		instanceName := instanceNameRaw.(string)
		if len(instanceName) > 0 {
			createReq.InstanceName = common.StringPtr(instanceName)
		}
	}

	if clusterVersionRaw, ok := d.GetOkExists("cluster_version"); ok {
		clusterVersion := clusterVersionRaw.(string)
		if len(clusterVersion) > 0 {
			createReq.ClusterVersion = common.StringPtr(clusterVersion)
		}
	}

	response, err := client.CreateCluster(createReq)
	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("tencentcloud_container_cluster get code error")
	}

	if *response.Code != 0 {
		return fmt.Errorf(
			"tencentcloud_container_cluster create error, code:%d, message:%v",
			*response.Code,
			*response.CodeDesc,
		)
	}

	if response.Data.ClusterId == nil || *response.Data.ClusterId == "" {
		return fmt.Errorf("tencentcloud_container_cluster no clusterInstanceId id returned")
	}

	clusterInstanceId := *response.Data.ClusterId

	d.SetId(clusterInstanceId)

	if err := waitClusterStatusReady(client, clusterInstanceId); err != nil {
		return fmt.Errorf("tencentcloud_container_cluster cluster status error")
	}

	return resourceTencentCloudContainerClusterRead(d, m)
}

func waitClusterStatusReady(client *ccs.Client, id string) error {

	describeClusterReq := ccs.NewDescribeClusterRequest()
	describeClusterReq.ClusterIds = []*string{&id}

	err := resource.Retry(15*time.Minute, func() *resource.RetryError {
		response, err := client.DescribeCluster(describeClusterReq)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(response.Data.Clusters) != 1 {
			return resource.NonRetryableError(fmt.Errorf("cluster create failed"))
		}
		if response.Data.Clusters[0].Status == nil {
			return resource.RetryableError(fmt.Errorf("status not found yet"))
		}
		if *response.Data.Clusters[0].Status == CLUSTER_LIFESTATE_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cluster still creating...will retry"))
	})
	return err
}

// Read Cluster Info
// Read Cluster Security Info
func resourceTencentCloudContainerClusterRead(d *schema.ResourceData, m interface{}) error {
	clusterInstanceId := d.Id()

	client := m.(*TencentCloudClient).ccsConn

	describeClusterReq := ccs.NewDescribeClusterRequest()
	describeClusterReq.ClusterIds = []*string{&clusterInstanceId}

	clusterResponse, err := client.DescribeCluster(describeClusterReq)
	if err != nil {
		return err
	}

	if clusterResponse.Code == nil {
		return fmt.Errorf("tencentcloud_container_cluster get code error")
	}

	if *clusterResponse.Code != 0 {
		return fmt.Errorf(
			"tencentcloud_container_cluster read error, code:%d, message:%v",
			*clusterResponse.Code,
			*clusterResponse.CodeDesc,
		)
	}

	if len(clusterResponse.Data.Clusters) == 1 {
		if clusterResponse.Data.Clusters[0].Description != nil {
			d.Set("cluster_desc", *clusterResponse.Data.Clusters[0].Description)
		}
		if clusterResponse.Data.Clusters[0].K8sVersion != nil {
			d.Set("kubernetes_version", *clusterResponse.Data.Clusters[0].K8sVersion)
		}
		if clusterResponse.Data.Clusters[0].NodeNum != nil {
			d.Set("nodes_num", *clusterResponse.Data.Clusters[0].NodeNum)
		}
		if clusterResponse.Data.Clusters[0].NodeStatus != nil {
			d.Set("nodes_status", *clusterResponse.Data.Clusters[0].NodeStatus)
		}
		if clusterResponse.Data.Clusters[0].TotalCPU != nil {
			d.Set("total_cpu", clusterResponse.Data.Clusters[0].TotalCPU)
		}
		if clusterResponse.Data.Clusters[0].TotalMem != nil {
			d.Set("total_mem", *clusterResponse.Data.Clusters[0].TotalMem)
		}
		if clusterResponse.Data.Clusters[0].ClusterName != nil {
			d.Set("cluster_name", *clusterResponse.Data.Clusters[0].ClusterName)
		}
		if clusterResponse.Data.Clusters[0].UnVpcId != nil {
			d.Set("vpc_id", *clusterResponse.Data.Clusters[0].UnVpcId)
		}
		if clusterResponse.Data.Clusters[0].ClusterCIDR != nil {
			d.Set("cluster_cidr", *clusterResponse.Data.Clusters[0].ClusterCIDR)
		}
	} else {
		d.SetId("")
	}

	return nil
}

func resourceTencentCloudContainerClusterDelete(d *schema.ResourceData, m interface{}) error {
	clusterInstanceId := d.Id()
	client := m.(*TencentCloudClient).ccsConn

	deleteClusterReq := ccs.NewDeleteClusterRequest()
	deleteClusterReq.ClusterId = &clusterInstanceId

	if nodeDeleteModeRaw, ok := d.GetOkExists("node_delete_mode"); ok {
		nodeDeleteMode := nodeDeleteModeRaw.(string)
		deleteClusterReq.NodeDeleteMode = &nodeDeleteMode
	}

	response, err := client.DeleteCluster(deleteClusterReq)

	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("tencentcloud_container_cluster get code error")
	}

	// resource not existed, return done
	if *response.Code == CLUSTER_NOT_FOUND_CODE {
		return nil
	}

	if *response.Code != 0 {
		return fmt.Errorf(
			"tencentcloud_container_cluster delete error, code:%d, message:%v",
			*response.Code,
			*response.CodeDesc,
		)
	}

	return nil
}
