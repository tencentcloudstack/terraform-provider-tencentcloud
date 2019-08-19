/*
Provide a resource to increase instance to cluster

Example Usage

```hcl

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "subnet" {
  default = "subnet-pqfek0t8"
}

variable "scale_instance_type" {
  default = "S2.LARGE16"
}


resource tencentcloud_kubernetes_scale_worker test_scale {
  cluster_id = "cls-godovr32"

  worker_config {
    count                      = 3
    availability_zone          = "${var.availability_zone}"
    instance_type              = "${var.scale_instance_type}"
    subnet_id                  = "${var.subnet}"
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 50
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "AABBccdd1122"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceTencentCloudTkeScaleWorker() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeScaleWorkerCreate,
		Read:   resourceTencentCloudTkeScaleWorkerRead,
		Delete: resourceTencentCloudTkeScaleWorkerDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				Description:"ID of the cluster.",
			},
			"worker_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: TkeCvmCreateInfo(),
				},
				Description:"Deploy the machine configuration information of the 'WORK' service, and create <=20 units for common users. ",
			},
			// Computed values
			"worker_instances_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: TkeCvmState(),
				},
				Description:"An information list of kubernetes cluster 'WORKER'. Each element contains the following attributes:",
			},
		},
	}
}

func resourceTencentCloudTkeScaleWorkerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_scale_worker.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var cvms RunInstancesForNode

	cvms.Work = []string{}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	if clusterId == "" {
		return fmt.Errorf("`cluster_id` is empty.")
	}

	info, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, clusterId)
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
		return fmt.Errorf("cluster [%s]  is not exist.", clusterId)
	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, info.VpcId, info.ProjectId)
			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)
		}
	}

	if len(cvms.Work) != 1 {
		return fmt.Errorf("only one additional configuration of virtual machines is now supported now, " +
			"so len(cvms.Work) should be 1")
	}

	instanceIds, err := service.CreateClusterInstances(ctx, clusterId, cvms.Work[0])

	if err != nil {
		return err
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(instanceIds))

	for _, v := range instanceIds {
		if v == "" {
			return fmt.Errorf("CreateClusterInstances return one instanceId is empty")
		}
		infoMap := make(map[string]interface{})
		infoMap["instance_id"] = v
		infoMap["instance_role"] = TKE_ROLE_WORKER
		workerInstancesList = append(workerInstancesList, infoMap)
	}

	if err = d.Set("worker_instances_list", workerInstancesList); err != nil {
		return err
	}

	md := md5.New()

	if _, err = md.Write([]byte(clusterId)); err != nil {
		return err
	}

	instanceIdJoin := strings.Join(instanceIds, "#")

	if _, err = md.Write([]byte(instanceIdJoin)); err != nil {
		return err
	}

	id := fmt.Sprintf("TkeScaleWorker.%x", md.Sum(nil))

	d.SetId(id)

	return nil
}

func resourceTencentCloudTkeScaleWorkerRead(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_kubernetes_scale_worker.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)
	if clusterId == "" {
		return fmt.Errorf("tke.`cluster_id` is empty.")
	}

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}
	//The cluster has been deleted
	if !has {
		d.SetId("")
		return nil
	}

	oldWorkerInstancesList := d.Get("worker_instances_list").([]interface{})


	instanceMap := make(map[string]bool)

	for _, v := range oldWorkerInstancesList {

		infoMap, ok := v.(map[string]interface{})

		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}

		if instanceMap[instanceId] {
			log.Printf("[WARN]The same instance id exists in the list")
		}

		instanceMap[instanceId] = true

	}

	_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	newWorkerInstancesList := make([]map[string]interface{}, 0, len(workers))

	for _, cvm := range workers {
		if _, ok := instanceMap[cvm.InstanceId]; !ok {
			continue
		}
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = cvm.InstanceId
		tempMap["instance_role"] = cvm.InstanceRole
		tempMap["instance_state"] = cvm.InstanceState
		tempMap["failed_reason"] = cvm.FailedReason
		newWorkerInstancesList = append(newWorkerInstancesList, tempMap)
	}

	//The machines I generated was deleted by others.
	if len(newWorkerInstancesList) == 0 {
		d.SetId("")
		return nil
	}

	return d.Set("worker_instances_list", newWorkerInstancesList)
}
func resourceTencentCloudTkeScaleWorkerDelete(d *schema.ResourceData, meta interface{}) error {

	defer logElapsed("resource.tencentcloud_kubernetes_scale_worker.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Get("cluster_id").(string)

	if clusterId == "" {
		return fmt.Errorf("`cluster_id` is empty.")
	}

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}
	//The cluster has been deleted
	if !has {
		return nil
	}
	workerInstancesList := d.Get("worker_instances_list").([]interface{})

	instanceMap := make(map[string]bool)

	for _, v := range workerInstancesList {

		infoMap, ok := v.(map[string]interface{})

		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}

		if instanceMap[instanceId] {
			log.Printf("[WARN]The same instance id exists in the list")
		}

		instanceMap[instanceId] = true

	}

	_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	needDeletes := []string{}
	for _, cvm := range workers {
		if _, ok := instanceMap[cvm.InstanceId]; ok {
			needDeletes = append(needDeletes, cvm.InstanceId)
		}
	}
	//The machines I generated was deleted by others.
	if len(needDeletes) == 0 {
		return nil
	}

	err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)
	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	return err
}
