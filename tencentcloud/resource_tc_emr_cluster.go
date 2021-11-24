/*
Provide a resource to create a emr cluster.

Example Usage

```hcl
resource "tencentcloud_emr_cluster" "emrrrr" {
  product_id=4
  display_strategy="clusterList"
  vpc_settings={
    vpc_id="vpc-fuwly8x5"
    subnet_id:"subnet-d830wfso"
  }
  softwares=["hadoop-2.8.4", "zookeeper-3.4.9"]
  support_ha=0
  instance_name="emr-test"
  resource_spec {
    master_resource_spec {
      mem_size=8192
      cpu=4
      disk_size=100
      disk_type="CLOUD_PREMIUM"
      spec="CVM.S2"
      storage_type=5
    }
    core_resource_spec {
      mem_size=8192
      cpu=4
      disk_size=100
      disk_type="CLOUD_PREMIUM"
      spec="CVM.S2"
      storage_type=5
    }
    master_count=1
    core_count=2
  }
  login_settings={
    password="tencent@cloud123"
  }
  time_span=1
  time_unit="m"
  pay_mode=1
  placement={
    zone="ap-guangzhou-3"
    project_id=0
  }
}
```
*/
package tencentcloud

import (
	"context"
	innerErr "errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

func resourceTencentCloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrClusterCreate,
		Read:   resourceTencentCloudEmrClusterRead,
		Delete: resourceTencentCloudEmrClusterDelete,
		Update: resourceTencentCloudEmrClusterUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Display strategy of EMR instance.",
			},
			"product_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 30),
				Description:  "The product id of EMR instance.",
			},
			"vpc_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "The private net config of EMR instance.",
			},
			"softwares": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The softwares of a EMR instance.",
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_resource_spec": buildResourceSpecSchema(),
						"core_resource_spec":   buildResourceSpecSchema(),
						"task_resource_spec":   buildResourceSpecSchema(),
						"master_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of master node.",
						},
						"core_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of core node.",
						},
						"task_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of core node.",
						},
						"common_resource_spec": buildResourceSpecSchema(),
						"common_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The number of common node.",
						},
					},
				},
				Description: "Resource specification of EMR instance.",
			},
			"support_ha": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(0, 1),
				Description:  "The flag whether the instance support high availability.(0=>not support, 1=>support).",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(6, 36),
				Description:  "Name of the instance, which can contain 6 to 36 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			},
			"pay_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 1),
				Description:  "The pay mode of instance. 0 is pay on an annual basis, 1 is pay on a measure basis.",
			},
			"placement": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "The location of the instance.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The length of time the instance was purchased. Use with TimeUnit.When TimeUnit is s, the parameter can only be filled in at 3600, representing a metered instance.\nWhen TimeUnit is m, the number filled in by this parameter indicates the length of purchase of the monthly instance of the package year, such as 1 for one month of purchase.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unit of time in which the instance was purchased. When PayMode is 0, TimeUnit can only take values of s(second). When PayMode is 1, TimeUnit can only take the value m(month).",
			},
			"login_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "Instance login settings.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created EMR instance id.",
			},
		},
	}
}

func resourceTencentCloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_emr_instance.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId := d.Id()
	timeUnit, hasTimeUnit := d.GetOkExists("time_unit")
	timeSpan, hasTimeSpan := d.GetOkExists("time_span")
	payMode, hasPayMode := d.GetOkExists("pay_mode")
	if !hasTimeUnit || !hasTimeSpan || !hasPayMode {
		return innerErr.New("Time_unit, time_span or pay_mode must be set.")
	}
	request := emr.NewScaleOutInstanceRequest()
	request.TimeUnit = common.StringPtr(timeUnit.(string))
	request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
	request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
	request.InstanceId = common.StringPtr(instanceId)

	tmpResourceSpec := d.Get("resource_spec").([]interface{})
	resourceSpec := tmpResourceSpec[0].(map[string]interface{})

	if d.HasChange("resource_spec.0.master_count") {
		request.MasterCount = common.Uint64Ptr((uint64)(resourceSpec["master_count"].(int)))
	}
	if d.HasChange("resource_spec.0.task_count") {
		request.TaskCount = common.Uint64Ptr((uint64)(resourceSpec["task_count"].(int)))
	}
	if d.HasChange("resource_spec.0.core_count") {
		request.CoreCount = common.Uint64Ptr((uint64)(resourceSpec["core_count"].(int)))
	}
	_, err := emrService.UpdateInstance(ctx, request)
	if err != nil {
		return err
	}
	err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusCreated {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudEmrClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_emr_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId, err := emrService.CreateInstance(ctx, d)
	if err != nil {
		return err
	}
	d.SetId(instanceId)
	d.Set("instance_id", instanceId)
	var displayStrategy string
	if v, ok := d.GetOk("display_strategy"); ok {
		displayStrategy = v.(string)
	}
	err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, displayStrategy)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusCreated {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudEmrClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_emr_instance.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	if err := emrService.DeleteInstance(ctx, d); err != nil {
		return err
	}
	instanceId := d.Id()
	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusDeleted {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudEmrClusterRead(d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instanceId := d.Id()
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
