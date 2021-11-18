/*
Provide a resource to create a emr cluster.

Example Usage

```hcl
resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=4
	display_strategy="clusterList"
	vpc_settings={vpc_id:"vpc-fuwly8x5", subnet_id:"subnet-d830wfso"}
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
	login_settings={password:"tencent@cloud123"}
	time_span=1
	time_unit="m"
    pay_mode=1
    placement={zone:"ap-guangzhou-3", project_id:0}
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrClusterCreate,
		Read:   resourceTencentCloudEmrClusterRead,
		Update: resourceTencentCloudEmrClusterUpdate,
		Delete: resourceTencentCloudEmrClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display strategy of EMR instance.",
			},
			"product_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 30),
				Description:  "The product id of EMR instance.",
			},
			"vpc_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "The private net config of EMR instance.",
			},
			"softwares": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The softwares of a EMR instance.",
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec":         {Type: schema.TypeString, Optional: true},
									"storage_type": {Type: schema.TypeInt, Optional: true},
									"disk_type":    {Type: schema.TypeString, Optional: true},
									"mem_size":     {Type: schema.TypeInt, Optional: true},
									"cpu":          {Type: schema.TypeInt, Optional: true},
									"disk_size":    {Type: schema.TypeInt, Optional: true},
								},
							},
						},
						"core_resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec":         {Type: schema.TypeString, Optional: true},
									"storage_type": {Type: schema.TypeInt, Optional: true},
									"disk_type":    {Type: schema.TypeString, Optional: true},
									"mem_size":     {Type: schema.TypeInt, Optional: true},
									"cpu":          {Type: schema.TypeInt, Optional: true},
									"disk_size":    {Type: schema.TypeInt, Optional: true},
								},
							},
						},
						"task_resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec":         {Type: schema.TypeString, Optional: true},
									"storage_type": {Type: schema.TypeInt, Optional: true},
									"disk_type":    {Type: schema.TypeString, Optional: true},
									"mem_size":     {Type: schema.TypeInt, Optional: true},
									"cpu":          {Type: schema.TypeInt, Optional: true},
									"disk_size":    {Type: schema.TypeInt, Optional: true},
								},
							},
						},
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
						"common_resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec":         {Type: schema.TypeString, Optional: true},
									"storage_type": {Type: schema.TypeInt, Optional: true},
									"disk_type":    {Type: schema.TypeString, Optional: true},
									"mem_size":     {Type: schema.TypeInt, Optional: true},
									"cpu":          {Type: schema.TypeInt, Optional: true},
									"disk_size":    {Type: schema.TypeInt, Optional: true},
								},
							},
						},
						"common_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of common node.",
						},
					},
				},
				Description: "Resource specification of EMR instance.",
			},
			"support_ha": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 1),
				Description:  "The flag whether the instance support high availability.(0=>not support, 1=>support).",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
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

func resourceTencentCloudEmrClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_emr_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	if err := emrService.CreateInstance(ctx, d); err != nil {
		return err
	}
	instanceId := d.Id()
	params := make(map[string]interface{})
	instances := make([]string, 0)
	instances = append(instances, instanceId)
	params["instance_ids"] = instances
	var displayStrategy string
	if v, ok := d.GetOk("display_strategy"); ok {
		displayStrategy = v.(string)
	}
	params["display_strategy"] = displayStrategy
	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstances(ctx, params)

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
	params := make(map[string]interface{})
	instances := make([]string, 0)
	instances = append(instances, instanceId)
	params["instance_ids"] = instances
	params["display_strategy"] = DisplayStrategyIsclusterList
	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstances(ctx, params)

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

func resourceTencentCloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTencentCloudEmrClusterRead(d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	params := make(map[string]interface{})
	var displayStrategy string
	if v, ok := d.GetOk("display_strategy"); ok {
		displayStrategy = v.(string)
	}
	params["display_strategy"] = displayStrategy
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, err := emrService.DescribeInstances(ctx, params)

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
