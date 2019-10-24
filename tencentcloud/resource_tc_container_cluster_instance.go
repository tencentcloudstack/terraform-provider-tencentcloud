package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudContainerClusterInstance() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use 'tencentcloud_kubernetes_scale_worker' instead.",
		Create:             resourceTencentCloudContainerClusterInstancesCreate,
		Read:               resourceTencentCloudContainerClusterInstancesRead,
		Update:             resourceTencentCloudContainerClusterInstancesUpdate,
		Delete:             resourceTencentCloudContainerClusterInstancesDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
			},
			"mem": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"require_wan_ip": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_vpc_gateway": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"root_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cvm_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sg_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"docker_graph_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"unschedulable": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"abnormal_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_normal": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wan_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lan_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudContainerClusterInstancesUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_container_cluster_instance.update")()

	return fmt.Errorf("the container cluster instance resource doesn't support update")
}

func resourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_container_cluster_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var workers []InstanceInfo
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		_, workers, e = service.DescribeClusterInstances(ctx, clusterId)
		if e, ok := e.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				d.SetId("")
				return nil
			}
		}
		if e != nil {
			return resource.RetryableError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	found := false
	for _, v := range workers {
		if v.InstanceId == d.Id() {
			found = true
			d.Set("instance_id", v.InstanceId)
			d.Set("abnormal_reason", v.FailedReason)
			d.Set("wan_ip", "")
			d.Set("lan_ip", "")
			d.Set("is_normal", true)
			if v.InstanceState == "failed" {
				d.Set("is_normal", false)
			}

			describeInstancesreq := cvm.NewDescribeInstancesRequest()
			describeInstancesreq.InstanceIds = []*string{common.StringPtr(v.InstanceId)}
			var describeInstancesResponse *cvm.DescribeInstancesResponse
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().DescribeInstances(describeInstancesreq)
				if e != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, describeInstancesreq.GetAction(), describeInstancesreq.ToJsonString(), e.Error())
					return retryError(e)
				}
				describeInstancesResponse = result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s DescribeInstances failed, reason:%s\n ", logId, err.Error())
				return err
			}

			if len(describeInstancesResponse.Response.InstanceSet) > 0 {
				if len(describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses) > 0 {
					d.Set("wan_ip", *describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses[0])
				}
				if len(describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses) > 0 {
					d.Set("lan_ip", *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0])
				}
			}
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceTencentCloudContainerClusterInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_container_cluster_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	runInstancesPara := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	if v, ok := d.GetOkExists("zone_id"); ok {
		var zones []*cvm.ZoneInfo
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			var e error
			zones, e = cvmService.DescribeZones(ctx)
			if e != nil {
				return retryError(e, "InternalError")
			}
			return nil
		})
		if err != nil {
			return err
		}
		zone := ""
		for _, z := range zones {
			if *z.ZoneId == v.(string) {
				zone = *z.Zone
				break
			}
		}
		place.Zone = stringToPointer(zone)
	}
	runInstancesPara.Placement = &place
	runInstancesPara.InstanceCount = common.Int64Ptr(1)

	var iAdvanced InstanceAdvancedSettings
	var cvms RunInstancesForNode

	if v, ok := d.GetOkExists("vpc_id"); ok {
		if len(v.(string)) > 0 {
			vpcId := v.(string)
			subnetId := ""
			asVpcGateway := false
			if subnetIdRaw, ok := d.GetOkExists("subnet_id"); ok {
				subnetId = subnetIdRaw.(string)
			}
			if isVpcGatewayRaw, ok := d.GetOkExists("is_vpc_gateway"); ok {
				if isVpcGatewayRaw.(int) == 1 {
					asVpcGateway = true
				}
			}
			runInstancesPara.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
				VpcId:        common.StringPtr(vpcId),
				SubnetId:     common.StringPtr(subnetId),
				AsVpcGateway: common.BoolPtr(asVpcGateway),
			}
		}
	}

	if instanceTypeRaw, ok := d.GetOkExists("instance_type"); ok {
		if len(instanceTypeRaw.(string)) > 0 {
			runInstancesPara.InstanceType = common.StringPtr(instanceTypeRaw.(string))
		}
	}

	if v, ok := d.GetOkExists("require_wan_ip"); ok {
		publicIpAssigned := false
		internetMaxBandwidthOut := int64(0)
		internetChargeType := ""
		if v.(int) == 1 {
			publicIpAssigned = true
			if v, ok := d.GetOkExists("bandwidth"); ok {
				internetMaxBandwidthOut = int64(v.(int))
			}
			if v, ok := d.GetOkExists("bandwidth_type"); ok {
				bandwidthTypes := map[string]string{
					"PayByMonth":   "BANDWIDTH_PREPAID",
					"PayByTraffic": "TRAFFIC_POSTPAID_BY_HOUR",
					"PayByHour":    "TRAFFIC_POSTPAID_BY_HOUR",
				}
				if v, ok := bandwidthTypes[v.(string)]; ok {
					internetChargeType = v
				}
			}
		}
		runInstancesPara.InternetAccessible = &cvm.InternetAccessible{
			PublicIpAssigned:        common.BoolPtr(publicIpAssigned),
			InternetMaxBandwidthOut: common.Int64Ptr(internetMaxBandwidthOut),
			InternetChargeType:      common.StringPtr(internetChargeType),
		}
	}

	if v, ok := d.GetOkExists("user_script"); ok {
		runInstancesPara.UserData = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("sg_id"); ok {
		runInstancesPara.SecurityGroupIds = []*string{common.StringPtr(v.(string))}
	}

	if v, ok := d.GetOkExists("password"); ok {
		runInstancesPara.LoginSettings = &cvm.LoginSettings{
			Password: common.StringPtr(v.(string)),
		}
	}

	if v, ok := d.GetOkExists("key_id"); ok {
		runInstancesPara.LoginSettings = &cvm.LoginSettings{
			KeyIds: []*string{common.StringPtr(v.(string))},
		}
	}

	runInstancesPara.SystemDisk = &cvm.SystemDisk{}
	if v, ok := d.GetOkExists("root_size"); ok {
		runInstancesPara.SystemDisk.DiskSize = common.Int64Ptr(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("root_type"); ok {
		runInstancesPara.SystemDisk.DiskType = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		if v.(int) > 0 {
			dataDisk := &cvm.DataDisk{
				DiskSize: common.Int64Ptr(int64(v.(int))),
				DiskType: common.StringPtr("CLOUD_PREMIUM"),
			}
			if v, ok := d.GetOkExists("storage_type"); ok {
				if len(v.(string)) > 0 {
					dataDisk.DiskType = common.StringPtr(v.(string))
				}
			}
			runInstancesPara.DataDisks = []*cvm.DataDisk{dataDisk}
		}
	}

	if v, ok := d.GetOkExists("cvm_type"); ok {
		cvmTypes := map[string]string{
			"PayByHour":  "POSTPAID_BY_HOUR",
			"PayByMonth": "PREPAID",
		}
		if vv, ok := cvmTypes[v.(string)]; ok {
			runInstancesPara.InstanceChargeType = common.StringPtr(vv)
		}
	}

	if v, ok := d.GetOkExists("period"); ok {
		runInstancesPara.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
			Period: common.Int64Ptr(int64(v.(int))),
		}
	}

	if v, ok := d.GetOkExists("instance_name"); ok {
		runInstancesPara.InstanceName = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("mount_target"); ok {
		iAdvanced.MountTarget = v.(string)
	}

	if v, ok := d.GetOkExists("docker_graph_path"); ok {
		iAdvanced.DockerGraphPath = v.(string)
	}

	if v, ok := d.GetOkExists("user_script"); ok {
		iAdvanced.UserScript = v.(string)
	}

	if v, ok := d.GetOkExists("unschedulable"); ok {
		iAdvanced.Unschedulable = int64(v.(int))
	}

	runInstancesParas := runInstancesPara.ToJsonString()
	cvms.Work = []string{runInstancesParas}

	instanceIds, err := service.CreateClusterInstances(ctx, clusterId, cvms.Work[0], iAdvanced)
	if err != nil {
		return err
	}

	if len(instanceIds) != 1 {
		return fmt.Errorf("no instance return")
	}

	err = resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		_, workers, e := service.DescribeClusterInstances(ctx, clusterId)
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		for _, v := range workers {
			if v.InstanceId != instanceIds[0] {
				continue
			}
			if v.InstanceState != "running" {
				return resource.RetryableError(fmt.Errorf("instance state not ready"))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(instanceIds[0])

	return resourceTencentCloudContainerClusterInstancesRead(d, meta)
}

func resourceTencentCloudContainerClusterInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_container_cluster_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var workers []InstanceInfo
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		_, workers, e = service.DescribeClusterInstances(ctx, clusterId)
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if e != nil {
			return resource.RetryableError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	found := false
	var node InstanceInfo
	for _, v := range workers {
		if v.InstanceId == d.Id() {
			found = true
			node = v
		}
	}
	if !found {
		return nil
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := service.DeleteClusterInstances(ctx, clusterId, []string{node.InstanceId})
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return retryError(err)
		}
		return nil
	})

	return err
}
