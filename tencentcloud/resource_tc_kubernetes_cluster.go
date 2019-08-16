package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

func TkeCvmState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"failed_reason": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func TkeCvmCreateInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"count": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  1,
		},
		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "sub machine of tke",
			Description: "Indicates which availability zone will be used.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Specified types of CVM instance.",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := strings.ToUpper(v.(string))
				if !strings.Contains(value, "LARGE") {
					errors = append(errors, fmt.Errorf(
						"%q  has to be `LARGE` type", k))
				}
				return
			},
		},
		"subnet_id": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(4, 100),
			Description:  "Private network ID. If vpc_id is set, this value is required.",
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(50, 500),
			Description:  "Volume of system disk in GB. Default is 50.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			MaxItems:    11,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:         schema.TypeString,
						ForceNew:     true,
						Optional:     true,
						Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk，available values: CLOUD_PREMIUM and CLOUD_SSD.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is 0.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Data disk snapshot ID.",
					},
				},
			},
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Available values include TRAFFIC_POSTPAID_BY_HOUR.",
		},
		"internet_max_bandwidth_out": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      0,
			ValidateFunc: validateIntegerInRange(0, 100),
			Description:  "Max bandwidth of Internet access in Mbps. Default is 0.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Specify whether to assign an Internet IP address.",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateAsConfigPassword,
			Description:  "Password to access.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "ID list of keys.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
		"user_data": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "ase64-encoded User Data text, the length limit is 16KB.",
		},
	}
}

func resourceTencentCloudTkeCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterCreate,
		Read:   resourceTencentCloudTkeClusterRead,
		Delete: resourceTencentCloudTkeClusterDelete,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_desc": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"cluster_os": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      TKE_CLUSTER_OS_UBUNTU,
				ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_OS),
			},
			"container_runtime": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      TKE_RUNTIME_DOCKER,
				ValidateFunc: validateAllowedStringValue(TKE_RUNTIMES),
			},
			"cluster_deploy_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      TKE_DEPLOY_TYPE_MANAGED,
				ValidateFunc: validateAllowedStringValue(TKE_DEPLOY_TYPES),
			},
			"cluster_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "1.10.5",
			},
			"cluster_ipvs": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(4, 100),
			},
			"project_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "Project ID, default value is 0.",
			},
			"cluster_cidr": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,

				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					_, ipnet, err := net.ParseCIDR(value)
					if err != nil {
						errors = append(errors, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %s", k, err))
						return
					}
					if ipnet == nil || value != ipnet.String() {
						errors = append(errors, fmt.Errorf("%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
						return
					}
					if !strings.Contains(value, "/") {
						errors = append(errors, fmt.Errorf("%q must be a network segment", k))
						return
					}
					if !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
						errors = append(errors, fmt.Errorf("%q must in  10.  |  192.168. |  172.[16-31]", k))
						return
					}

					if strings.HasPrefix(value, "172.") {
						nextNo := strings.Split(value, ".")[1]
						no, _ := strconv.ParseInt(nextNo, 10, 64)
						if no < 16 || no > 31 {
							errors = append(errors, fmt.Errorf("%q must in  10.  |  192.168. |  172.[16-31]", k))
							return
						}
					}
					return
				},
			},
			"ignore_cluster_cidr_conflict": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"cluster_max_pod_num": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  256,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value%16 != 0 {
						errors = append(errors, fmt.Errorf(
							"%q  has to be a multiple of 16 ", k))
					}
					if value < 32 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be lower than %d: %d", k, 32, value))
					}
					return
				},
			},
			"cluster_max_service_num": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  256,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value%16 != 0 {
						errors = append(errors, fmt.Errorf(
							"%q  has to be a multiple of 16 ", k))
					}
					return
				},
			},
			"master_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: TkeCvmCreateInfo(),
				},
			},
			"worker_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: TkeCvmCreateInfo(),
				},
			},
			// Computed values
			"cluster_node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"worker_instances_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: TkeCvmState(),
				},
			},
		},
	}
}

func tkeGetCvmRunInstancesPara(dMap map[string]interface{}, meta interface{},
	vpcId string, projectId int64) (cvmJson string, count int64, errRet error) {

	request := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	request.Placement = &place

	place.ProjectId = &projectId

	configRegion := meta.(*TencentCloudClient).apiV3Conn.Region
	if v, ok := dMap["availability_zone"]; ok {
		if !strings.Contains(v.(string), configRegion) {
			errRet = fmt.Errorf("availability_zone[%s] should in [%s]", v.(string), configRegion)
			return
		}
		place.Zone = stringToPointer(v.(string))
	}

	if v, ok := dMap["instance_type"]; ok {
		request.InstanceType = stringToPointer(v.(string))
	} else {
		errRet = fmt.Errorf("instance_type must be set.")
		return
	}

	subnetId := ""

	if v, ok := dMap["subnet_id"]; ok {
		subnetId = v.(string)
	}

	if (vpcId == "" && subnetId != "") ||
		(vpcId != "" && subnetId == "") {
		errRet = fmt.Errorf("Parameters cvm.`subnet_id` and cluster.`vpc_id` are both set or neither")
		return
	}

	if vpcId != "" {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
			VpcId:    &vpcId,
			SubnetId: &subnetId,
		}
	}

	if v, ok := dMap["system_disk_type"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskType = stringToPointer(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskSize = int64Pt(int64(v.(int)))

	}

	if v, ok := dMap["data_disk"]; ok {

		dataDisks := v.([]interface{})
		request.DataDisks = make([]*cvm.DataDisk, 0, len(dataDisks))

		for _, d := range dataDisks {

			var (
				value      = d.(map[string]interface{})
				diskType   = value["disk_type"].(string)
				diskSize   = int64(value["disk_size"].(int))
				snapshotId = value["snapshot_id"].(string)
				dataDisk   = cvm.DataDisk{
					DiskType: &diskType,
					DiskSize: &diskSize,
				}
			)
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if v, ok := dMap["internet_charge_type"]; ok {

		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetChargeType = stringToPointer(v.(string))
	}

	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetMaxBandwidthOut = int64Pt(int64(v.(int)))
	}

	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	if v, ok := dMap["password"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}
		request.LoginSettings.Password = stringToPointer(v.(string))
	}

	if v, ok := dMap["instance_name"]; ok {
		request.InstanceName = stringToPointer(v.(string))
	}

	if v, ok := dMap["key_ids"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	if v, ok := dMap["enhanced_security_service"]; ok {

		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}

		securityService := v.(bool)
		request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}
	if v, ok := dMap["user_data"]; ok {
		request.UserData = stringToPointer(v.(string))
	}

	chargeType := INSTANCE_CHARGE_TYPE_POSTPAID
	request.InstanceChargeType = &chargeType

	if v, ok := dMap["count"]; ok {
		count = int64(v.(int))
	} else {
		count = 1
	}
	request.InstanceCount = &count

	cvmJson = request.ToJsonString()

	cvmJson = strings.Replace(cvmJson, `"Password":"",`, "", -1)

	return
}

func resourceTencentCloudTkeClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var basic ClusterBasicSetting
	var advanced ClusterAdvancedSettings
	var cvms RunInstancesForNode
	var cidrSet ClusterCidrSettings
	cvms.Master = []string{}
	cvms.Work = []string{}

	clusterDeployType := d.Get("cluster_deploy_type").(string)

	vpcId := d.Get("vpc_id").(string)
	if vpcId != "" {
		basic.VpcId = vpcId
	}

	basic.ProjectId = int64(d.Get("project_id").(int))
	basic.ClusterOs = d.Get("cluster_os").(string)
	basic.ClusterVersion = d.Get("cluster_version").(string)
	if v, ok := d.GetOk("cluster_name"); ok {
		basic.ClusterName = v.(string)
	}
	if v, ok := d.GetOk("cluster_desc"); ok {
		basic.ClusterDescription = v.(string)
	}

	advanced.ContainerRuntime = d.Get("container_runtime").(string)
	advanced.Ipvs = d.Get("cluster_ipvs").(bool)

	cidrSet.ClusterCidr = d.Get("cluster_cidr").(string)
	cidrSet.IgnoreClusterCidrConflict = d.Get("ignore_cluster_cidr_conflict").(bool)
	cidrSet.MaxClusterServiceNum = int64(d.Get("cluster_max_service_num").(int))
	cidrSet.MaxNodePodNum = int64(d.Get("cluster_max_pod_num").(int))

	items := strings.Split(cidrSet.ClusterCidr, "/")
	if len(items) != 2 {
		return fmt.Errorf("`cluster_cidr` must be network segment ")
	}

	bitNumber, err := strconv.ParseInt(items[1], 10, 64)

	if err != nil {
		return fmt.Errorf("`cluster_cidr` must be network segment ")
	}

	if math.Pow(2, float64(32-bitNumber)) <= float64(cidrSet.MaxNodePodNum) {
		return fmt.Errorf("`cluster_cidr` Network segment range is too small, can not cover cluster_max_service_num")
	}

	if masters, ok := d.GetOk("master_config"); ok {
		if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED {
			return fmt.Errorf("if `cluster_deploy_type` is `MANAGED_CLUSTER` , You don't need define the master yourself")
		}
		var masterCount int64 = 0
		masterList := masters.([]interface{})
		for index := range masterList {
			master := masterList[index].(map[string]interface{})
			paraJson, count, err := tkeGetCvmRunInstancesPara(master, meta, vpcId, int64(basic.ProjectId))
			if err != nil {
				return err
			}

			cvms.Master = append(cvms.Master, paraJson)
			masterCount += count
		}
		if masterCount < 3 {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` len(master_config) should  >=3 ")
		}
	} else {
		if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` , You need define the master yourself")
		}

	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, vpcId, int64(basic.ProjectId))

			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)
		}
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	id, err := service.CreateCluster(ctx, basic, advanced, cvms, cidrSet)

	if err != nil {
		return err
	}

	d.SetId(id)

	_, _, err = service.DescribeClusterInstances(ctx, d.Id())

	if err != nil {
		//create often cost more than half an hour.
		err = resource.Retry(time.Hour, func() *resource.RetryError {
			_, _, err = service.DescribeClusterInstances(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		log.Printf("[WARN]%s resource.kubernetes_cluster DescribeClusterInstances fail , %s", logId, err.Error())
	}

	if err = resourceTencentCloudTkeClusterRead(d, meta); err != nil {
		log.Printf("[WARN]%s resource.kubernetes_cluster.read after create fail , %s", logId, err.Error())
	}

	return nil
}

func resourceTencentCloudTkeClusterRead(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_kubernetes_cluster.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	info, has, err := service.DescribeCluster(ctx, d.Id())
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	d.Set("cluster_name", info.ClusterName)
	d.Set("cluster_desc", info.ClusterDescription)
	d.Set("cluster_os", info.ClusterOs)
	d.Set("cluster_deploy_type", info.DeployType)
	d.Set("cluster_version", info.ClusterVersion)
	d.Set("cluster_ipvs", info.Ipvs)
	d.Set("vpc_id", info.VpcId)
	d.Set("project_id", info.ProjectId)
	d.Set("cluster_cidr", info.ClusterCidr)
	d.Set("ignore_cluster_cidr_conflict", info.IgnoreClusterCidrConflict)
	d.Set("cluster_max_pod_num", info.MaxClusterServiceNum)
	d.Set("cluster_max_service_num", info.MaxClusterServiceNum)
	d.Set("cluster_node_num", info.ClusterNodeNum)

	_, workers, err := service.DescribeClusterInstances(ctx, d.Id())
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, d.Id())
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(workers))
	for _, worker := range workers {
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = worker.InstanceId
		tempMap["instance_role"] = worker.InstanceRole
		tempMap["instance_state"] = worker.InstanceState
		tempMap["failed_reason"] = worker.FailedReason
		workerInstancesList = append(workerInstancesList, tempMap)
	}

	return d.Set("worker_instances_list", workerInstancesList)
}

func resourceTencentCloudTkeClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.DeleteCluster(ctx, d.Id())
		if err != nil {
			return retryError(err)
		}
		return nil
	})
}
