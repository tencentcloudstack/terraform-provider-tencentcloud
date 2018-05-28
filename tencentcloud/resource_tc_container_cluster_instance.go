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
	INSTANCE_NOT_FOUND_CODE = -1
)

func resourceTencentCloudContainerClusterInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudContainerClusterInstancesCreate,
		Read:   resourceTencentCloudContainerClusterInstancesRead,
		Update: resourceTencentCloudContainerClusterInstancesUpdate,
		Delete: resourceTencentCloudContainerClusterInstancesDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
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
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": &schema.Schema{
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
			"unschedulable": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"abnormal_reason": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_normal": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wan_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lan_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, m interface{}) error {
	instanceId := d.Id()
	client := m.(*TencentCloudClient).ccsConn
	describeClusterInstancesReq := ccs.NewDescribeClusterInstancesRequest()

	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		describeClusterInstancesReq.ClusterId = common.StringPtr(clusterId.(string))
	} else {
		return fmt.Errorf("data_source_tencent_cloud_container_cluster_instances read action needs param cluster_id")
	}

	response, err := client.DescribeClusterInstances(describeClusterInstancesReq)
	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("data_source_tencent_cloud_container_cluster_instances got error, no code response")
	}

	if *response.Code != 0 {
		return fmt.Errorf("data_source_tencent_cloud_container_cluster_instances got error, code %v , message %v", *response.Code, *response.CodeDesc)
	}

	found := false
	for _, node := range response.Data.Nodes {
		if *node.InstanceId == instanceId {
			found = true

			if node.AbnormalReason != nil {
				d.Set("abnormal_reason", *node.AbnormalReason)
			}
			if node.CPU != nil {
				d.Set("cpu", *node.CPU)
			}
			if node.Mem != nil {
				d.Set("mem", *node.Mem)
			}
			if node.InstanceId != nil {
				d.Set("instance_id", *node.InstanceId)
			}
			if node.IsNormal != nil {
				d.Set("is_normal", *node.IsNormal)
			}
			if node.WanIp != nil {
				d.Set("wan_ip", *node.WanIp)
			}
			if node.LanIp != nil {
				d.Set("lan_ip", *node.LanIp)
			}
		}
	}

	if found == false {
		d.SetId("")
	}

	return nil
}

func resourceTencentCloudContainerClusterInstancesUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("the container cluster instances resource doesn't support update")
}

// CreateClusterNode one node per time
func resourceTencentCloudContainerClusterInstancesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).ccsConn

	createInstanceReq := ccs.NewAddClusterInstancesRequest()

	clusterId := d.Get("cluster_id").(string)
	createInstanceReq.ClusterId = &clusterId
	createInstanceReq.GoodsNum = common.IntPtr(1)
	if cpuRaw, ok := d.GetOkExists("cpu"); ok {
		cpu := cpuRaw.(int)
		if cpu > 0 {
			createInstanceReq.CPU = common.IntPtr(cpu)
		}
	}

	if memRaw, ok := d.GetOkExists("mem"); ok {
		mem := memRaw.(int)
		if mem > 0 {
			createInstanceReq.Mem = common.IntPtr(mem)
		}
	}

	if bandwidthRaw, ok := d.GetOkExists("bandwidth"); ok {
		bandwidth := bandwidthRaw.(int)
		createInstanceReq.Bandwidth = common.IntPtr(bandwidth)
	}

	if bandwidthTypeRaw, ok := d.GetOkExists("bandwidth_type"); ok {
		bandwidthType := bandwidthTypeRaw.(string)
		if len(bandwidthType) > 0 {
			createInstanceReq.BandwidthType = common.StringPtr(bandwidthType)
		}
	}

	if wanIpRaw, ok := d.GetOkExists("require_wan_ip"); ok {
		wanIp := wanIpRaw.(int)
		createInstanceReq.WanIp = common.IntPtr(wanIp)
	}

	if subnetIdRaw, ok := d.GetOkExists("subnet_id"); ok {
		subnetId := subnetIdRaw.(string)
		if len(subnetId) > 0 {
			createInstanceReq.SubnetId = common.StringPtr(subnetId)
		}
	}

	if isVpcGatewayRaw, ok := d.GetOkExists("is_vpc_gateway"); ok {
		isVpcGateway := isVpcGatewayRaw.(int)
		createInstanceReq.IsVpcGateway = common.IntPtr(isVpcGateway)
	}

	if storageSizeRaw, ok := d.GetOkExists("storage_size"); ok {
		storageSize := storageSizeRaw.(int)
		createInstanceReq.StorageSize = common.IntPtr(storageSize)
	}

	if storageTypeRaw, ok := d.GetOkExists("storage_type"); ok {
		storageType := storageTypeRaw.(string)
		if len(storageType) > 0 {
			createInstanceReq.StorageType = common.StringPtr(storageType)
		}
	}

	if rootSizeRaw, ok := d.GetOkExists("root_size"); ok {
		rootSize := rootSizeRaw.(int)
		createInstanceReq.RootSize = common.IntPtr(rootSize)
	}

	if rootTypeRaw, ok := d.GetOkExists("root_type"); ok {
		rootType := rootTypeRaw.(string)
		if len(rootType) > 0 {
			createInstanceReq.RootType = common.StringPtr(rootType)
		}
	}

	if passwordRaw, ok := d.GetOkExists("password"); ok {
		password := passwordRaw.(string)
		createInstanceReq.Password = common.StringPtr(password)
	}

	if keyIdRaw, ok := d.GetOkExists("key_id"); ok {
		keyId := keyIdRaw.(string)
		if len(keyId) > 0 {
			createInstanceReq.KeyId = common.StringPtr(keyId)
		}
	}

	if cvmTypeRaw, ok := d.GetOkExists("cvm_type"); ok {
		cvmType := cvmTypeRaw.(string)
		if len(cvmType) > 0 {
			createInstanceReq.CvmType = common.StringPtr(cvmType)
		}
	}

	if periodRaw, ok := d.GetOkExists("period"); ok {
		period := periodRaw.(int)
		createInstanceReq.Period = common.IntPtr(period)
	}

	if zoneIdRaw, ok := d.GetOkExists("zone_id"); ok {
		zoneId := zoneIdRaw.(string)
		if len(zoneId) > 0 {
			createInstanceReq.ZoneId = common.StringPtr(zoneId)
		}
	}

	if instanceTypeRaw, ok := d.GetOkExists("instance_type"); ok {
		instanceType := instanceTypeRaw.(string)
		if len(instanceType) > 0 {
			createInstanceReq.InstanceType = common.StringPtr(instanceType)
		}
	}

	if sgIdRaw, ok := d.GetOkExists("sg_id"); ok {
		sgId := sgIdRaw.(string)
		if len(sgId) > 0 {
			createInstanceReq.SgId = common.StringPtr(sgId)
		}
	}

	if mountTargetRaw, ok := d.GetOkExists("mount_target"); ok {
		mountTarget := mountTargetRaw.(string)
		if len(mountTarget) > 0 {
			createInstanceReq.MountTarget = common.StringPtr(mountTarget)
		}
	}

	if dockerGraphPathRaw, ok := d.GetOkExists("docker_graph_path"); ok {
		dockerGraphPath := dockerGraphPathRaw.(string)
		if len(dockerGraphPath) > 0 {
			createInstanceReq.DockerGraphPath = common.StringPtr(dockerGraphPath)
		}
	}

	if unschedulableRaw, ok := d.GetOkExists("unschedulable"); ok {
		unschedulable := unschedulableRaw.(int)
		createInstanceReq.Unschedulable = common.IntPtr(unschedulable)
	}

	if userScriptRaw, ok := d.GetOkExists("user_script"); ok {
		userScript := userScriptRaw.(string)
		if len(userScript) > 0 {
			createInstanceReq.UserScript = common.StringPtr(userScript)
		}
	}

	response, err := client.AddClusterInstances(createInstanceReq)
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

	if len(response.Data.InstanceIds) == 0 {
		return fmt.Errorf("tencentcloud_container_cluster no clusterInstanceId id returned")
	}

	nodeId := response.Data.InstanceIds[0]
	d.SetId(*nodeId)

	if err := waitClusterInstanceRunning(client, clusterId, *nodeId); err != nil {
		return fmt.Errorf("Cluster Instance %s is abnormal, create fail", *nodeId)
	}

	return resourceTencentCloudContainerClusterInstancesRead(d, m)
}

func waitClusterInstanceRunning(conn *ccs.Client, clusterId, nodeId string) error {
	req := ccs.NewDescribeClusterInstancesRequest()
	req.ClusterId = &clusterId
	err := resource.Retry(15*time.Minute, func() *resource.RetryError {
		resp, err := conn.DescribeClusterInstances(req)
		if err != nil {
			return resource.RetryableError(err)
		}
		if _, ok := err.(*common.APIError); ok {
			return resource.NonRetryableError(err)
		}
		for _, node := range resp.Data.Nodes {
			if *node.InstanceId != nodeId {
				continue
			}
			if *node.IsNormal == 0 {
				return resource.NonRetryableError(fmt.Errorf("Instance status = 0"))
			}
			if *node.IsNormal == 1 {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Instance status = %d", *node.IsNormal))
		}
		// not found
		return nil
	})
	return err
}

func resourceTencentCloudContainerClusterInstancesDelete(d *schema.ResourceData, m interface{}) error {
	nodeId := d.Id()

	deleteClusterNodeReq := ccs.NewDeleteClusterInstancesRequest()
	client := m.(*TencentCloudClient).ccsConn

	describeClusterInstancesReq := ccs.NewDescribeClusterInstancesRequest()
	describeClusterInstancesReq.ClusterId = common.StringPtr(d.Get("cluster_id").(string))
	describeClusterInstancesRsp, err := client.DescribeClusterInstances(describeClusterInstancesReq)
	if err != nil {
		return err
	}
	// node no longer exists
	if len(describeClusterInstancesRsp.Data.Nodes) == 0 {
		return nil
	}

	nodeFound := false
	for _, instance := range describeClusterInstancesRsp.Data.Nodes {
		if *instance.InstanceId == nodeId {
			nodeFound = true
			break
		}
	}

	// node no longer exists
	if nodeFound == false {
		return nil
	}

	deleteClusterNodeReq.ClusterId = common.StringPtr(d.Get("cluster_id").(string))
	deleteClusterNodeReq.InstanceIds = common.StringPtrs([]string{nodeId})

	if nodeDeleteModeRaw, ok := d.GetOkExists("node_delete_mode"); ok {
		nodeDeleteMode := nodeDeleteModeRaw.(string)
		deleteClusterNodeReq.NodeDeleteMode = &nodeDeleteMode
	}

	response, err := client.DeleteClusterInstances(deleteClusterNodeReq)

	if err != nil {
		return err
	}

	if response.Code == nil {
		return fmt.Errorf("tencentcloud_container_cluster_instances get code error")
	}

	if *response.Code != 0 {
		return fmt.Errorf(
			"tencentcloud_container_cluster_instances delete error, code:%d, message:%v",
			*response.Code,
			*response.CodeDesc,
		)
	}

	return nil
}
