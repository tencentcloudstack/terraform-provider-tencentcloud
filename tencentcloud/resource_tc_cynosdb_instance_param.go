/*
Provides a resource to create a cynosdb instance_param

Example Usage

```hcl
resource "tencentcloud_cynosdb_instance_param" "instance_param" {
  cluster_id = ""
  instance_ids =
  cluster_param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  instance_param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  is_in_maintain_period = ""
}
```

Import

cynosdb instance_param can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_instance_param.instance_param instance_param_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"log"
	"time"
)

func resourceTencentCloudCynosdbInstanceParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstanceParamCreate,
		Read:   resourceTencentCloudCynosdbInstanceParamRead,
		Update: resourceTencentCloudCynosdbInstanceParamUpdate,
		Delete: resourceTencentCloudCynosdbInstanceParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID.",
			},

			"cluster_param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Cluster parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Current value of parameter.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter old value (only useful when generating parameters) Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"instance_param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Instance parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Current value of parameter.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter old value (only useful when generating parameters) Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"is_in_maintain_period": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Yes: modify within the operation and maintenance time window, no: execute immediately (default value).",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCynosdbInstanceParamUpdate(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceParamId := d.Id()

	instanceParam, err := service.DescribeCynosdbInstanceParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceParam.ClusterId != nil {
		_ = d.Set("cluster_id", instanceParam.ClusterId)
	}

	if instanceParam.InstanceIds != nil {
		_ = d.Set("instance_ids", instanceParam.InstanceIds)
	}

	if instanceParam.ClusterParamList != nil {
		clusterParamListList := []interface{}{}
		for _, clusterParamList := range instanceParam.ClusterParamList {
			clusterParamListMap := map[string]interface{}{}

			if instanceParam.ClusterParamList.ParamName != nil {
				clusterParamListMap["param_name"] = instanceParam.ClusterParamList.ParamName
			}

			if instanceParam.ClusterParamList.CurrentValue != nil {
				clusterParamListMap["current_value"] = instanceParam.ClusterParamList.CurrentValue
			}

			if instanceParam.ClusterParamList.OldValue != nil {
				clusterParamListMap["old_value"] = instanceParam.ClusterParamList.OldValue
			}

			clusterParamListList = append(clusterParamListList, clusterParamListMap)
		}

		_ = d.Set("cluster_param_list", clusterParamListList)

	}

	if instanceParam.InstanceParamList != nil {
		instanceParamListList := []interface{}{}
		for _, instanceParamList := range instanceParam.InstanceParamList {
			instanceParamListMap := map[string]interface{}{}

			if instanceParam.InstanceParamList.ParamName != nil {
				instanceParamListMap["param_name"] = instanceParam.InstanceParamList.ParamName
			}

			if instanceParam.InstanceParamList.CurrentValue != nil {
				instanceParamListMap["current_value"] = instanceParam.InstanceParamList.CurrentValue
			}

			if instanceParam.InstanceParamList.OldValue != nil {
				instanceParamListMap["old_value"] = instanceParam.InstanceParamList.OldValue
			}

			instanceParamListList = append(instanceParamListList, instanceParamListMap)
		}

		_ = d.Set("instance_param_list", instanceParamListList)

	}

	if instanceParam.IsInMaintainPeriod != nil {
		_ = d.Set("is_in_maintain_period", instanceParam.IsInMaintainPeriod)
	}

	return nil
}

func resourceTencentCloudCynosdbInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cynosdb.NewModifyInstanceParamRequest()

	instanceParamId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"cluster_id", "instance_ids", "cluster_param_list", "instance_param_list", "is_in_maintain_period"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyInstanceParam(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceParam failed, reason:%+v", logId, err)
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbInstanceParamStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbInstanceParamRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
