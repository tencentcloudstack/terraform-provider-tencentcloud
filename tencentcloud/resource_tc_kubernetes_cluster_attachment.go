/*
Provide a resource to attach an existing  cvm to kubernetes cluster.

Example Usage

```hcl

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}


data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  cpu_core_count = 8
  memory_size    = 16
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-auto-test-1-1"
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = var.default_instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_cluster_attachment" "test_attach" {
  cluster_id  = tencentcloud_kubernetes_cluster.managed_cluster.id
  instance_id = tencentcloud_instance.foo.id
  password    = "Lo4wbdit"

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func TkeInstanceAdvancedSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mount_target": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Mount target. Default is not mounting.",
		},
		"docker_graph_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "/var/lib/docker",
			Description: "Docker graph path. Default is `/var/lib/docker`.",
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
						Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"file_system": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Default:     "",
						Description: "File system, e.g. `ext3/ext4/xfs`.",
					},
					"auto_format_and_mount": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
						Description: "Indicate whether to auto format and mount or not. Default is `false`.",
					},
					"mount_target": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Default:     "",
						Description: "Mount target.",
					},
				},
			},
		},
		"extra_args": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Custom parameter information related to the node. This is a white-list parameter.",
		},
		"user_data": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Base64-encoded User Data text, the length limit is 16KB.",
		},
		"is_schedule": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "Indicate to schedule the adding node or not. Default is true.",
		},
	}
}

func resourceTencentCloudTkeClusterAttachment() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "ID of the cluster.",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the CVM instance, this cvm will reinstall the system.",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateAsConfigPassword,
			Description:  "Password to access, should be set if `key_ids` not set.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.",
		},
		"hostname": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Description: "The host name of the attached instance. " +
				"Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. " +
				"Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. " +
				"Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
		},
		"worker_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeInstanceAdvancedSetting(),
			},
			Description: "Deploy the machine configuration information of the 'WORKER', commonly used to attach existing instances.",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Labels of tke attachment exits CVM.",
		},
		"unschedulable": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
		},
		//compute
		"security_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "A list of security group IDs after attach to cluster.",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the node.",
		},
	}

	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterAttachmentCreate,
		Read:   resourceTencentCloudTkeClusterAttachmentRead,
		Delete: resourceTencentCloudTkeClusterAttachmentDelete,
		Schema: schemaBody,
	}
}

func tkeGetInstanceAdvancedPara(dMap map[string]interface{}, meta interface{}) (setting tke.InstanceAdvancedSettings) {
	setting = tke.InstanceAdvancedSettings{}
	if v, ok := dMap["mount_target"]; ok {
		setting.MountTarget = helper.String(v.(string))
	}

	if v, ok := dMap["data_disk"]; ok {

		dataDisks := v.([]interface{})
		setting.DataDisks = make([]*tke.DataDisk, 0, len(dataDisks))

		for _, d := range dataDisks {
			var (
				value              = d.(map[string]interface{})
				diskType           = value["disk_type"].(string)
				diskSize           = int64(value["disk_size"].(int))
				fileSystem         = value["file_system"].(string)
				autoFormatAndMount = value["auto_format_and_mount"].(bool)
				mountTarget        = value["mount_target"].(string)
				dataDisk           = tke.DataDisk{
					DiskType:           &diskType,
					DiskSize:           &diskSize,
					FileSystem:         &fileSystem,
					AutoFormatAndMount: &autoFormatAndMount,
					MountTarget:        &mountTarget,
				}
			)
			setting.DataDisks = append(setting.DataDisks, &dataDisk)
		}
	}
	if v, ok := dMap["is_schedule"]; ok {
		setting.Unschedulable = helper.BoolToInt64Ptr(!v.(bool))
	}

	if v, ok := dMap["user_data"]; ok {
		setting.UserScript = helper.String(v.(string))
	}

	if v, ok := dMap["docker_graph_path"]; ok {
		setting.DockerGraphPath = helper.String(v.(string))
	}

	if temp, ok := dMap["extra_args"]; ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		clusterExtraArgs := tke.InstanceExtraArgs{}
		clusterExtraArgs.Kubelet = make([]*string, 0)
		for i := range extraArgs {
			clusterExtraArgs.Kubelet = append(clusterExtraArgs.Kubelet, &extraArgs[i])
		}
		setting.ExtraArgs = &clusterExtraArgs
	}

	return setting
}
func resourceTencentCloudTkeClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	/*tke has been deleted*/
	_, has, err := tkeService.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = tkeService.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err, InternalError)
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

	/*cvm has been deleted*/
	var instance *cvm.Instance
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, err = cvmService.DescribeInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	instanceState := ""
	has = false
	/*attachment has been  deleted*/

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, workers, err := tkeService.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			return retryError(err, InternalError)
		}
		for _, worker := range workers {
			if worker.InstanceId == instanceId {
				has = true
				instanceState = worker.InstanceState
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						worker.InstanceId, clusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						worker.InstanceId, worker.InstanceState))
				}
				_ = d.Set("unschedulable", worker.InstanceAdvancedSettings.Unschedulable)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	}

	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	_ = d.Set("state", instanceState)
	return nil
}

func resourceTencentCloudTkeClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	request := tke.NewAddExistedInstancesRequest()

	instanceId := helper.String(d.Get("instance_id").(string))
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceIds = []*string{instanceId}
	request.LoginSettings = &tke.LoginSettings{}

	var loginSettingsNumbers = 0

	if v, ok := d.GetOk("key_ids"); ok {
		request.LoginSettings.KeyIds = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		loginSettingsNumbers++
	}

	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
		loginSettingsNumbers++
	}

	if loginSettingsNumbers != 1 {
		return fmt.Errorf("parameters `key_ids` and `password` must set and only set one")
	}

	request.InstanceAdvancedSettings = &tke.InstanceAdvancedSettings{}
	if workConfig, ok := d.GetOk("worker_config"); ok {
		workConfigList := workConfig.([]interface{})
		if len(workConfigList) == 1 {
			workConfigPara := workConfigList[0].(map[string]interface{})
			setting := tkeGetInstanceAdvancedPara(workConfigPara, meta)
			request.InstanceAdvancedSettings = &setting
		}
	}

	request.InstanceAdvancedSettings.Labels = GetTkeLabels(d, "labels")
	if hostName, ok := d.GetOk("hostname"); ok {
		hostNameStr := hostName.(string)
		request.HostName = &hostNameStr
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		request.InstanceAdvancedSettings.Unschedulable = helper.Int64(v.(int64))
	}

	/*cvm has been  attached*/
	var err error
	_, workers, err := tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has := false
	for _, worker := range workers {
		if worker.InstanceId == *instanceId {
			has = true
		}
	}
	if has {
		return fmt.Errorf("instance %s has been attached to cluster %s,can not attach again", *instanceId, *request.ClusterId)
	}

	var response *tke.AddExistedInstancesResponse

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = tkeService.client.UseTkeClient().AddExistedInstances(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("add existed instance %s to cluster %s error,reason %v", *instanceId, *request.ClusterId, err)
	}
	var success = false
	for _, v := range response.Response.SuccInstanceIds {
		if *v == *instanceId {
			d.SetId(*instanceId + "_" + *request.ClusterId)
			success = true
		}
	}

	if !success {
		return fmt.Errorf("add existed instance %s to cluster %s error, instance not in success instanceIds", *instanceId, *request.ClusterId)
	}

	/*wait for cvm status*/
	if err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, *instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance %s status is %s, retry...", *instanceId, *instance.InstanceState))
	}); err != nil {
		return err
	}

	/*wait for tke init ok */
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
		if err != nil {
			return retryError(err, InternalError)
		}
		has := false
		for _, worker := range workers {
			if worker.InstanceId == *instanceId {
				has = true
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						*instanceId, *request.ClusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						*instanceId, worker.InstanceState))
				}

			}
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s not exist in tke instance list", *instanceId))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTkeClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	request := tke.NewDeleteClusterInstancesRequest()

	request.ClusterId = &clusterId
	request.InstanceIds = []*string{
		&instanceId,
	}
	request.InstanceDeleteMode = helper.String("retain")

	var err error

	if err = resource.Retry(4*writeRetryTimeout, func() *resource.RetryError {
		_, err := tkeService.client.UseTkeClient().DeleteClusterInstances(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
			if e.GetCode() == "InternalError.Param" &&
				strings.Contains(e.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
				return nil
			}
		}

		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
