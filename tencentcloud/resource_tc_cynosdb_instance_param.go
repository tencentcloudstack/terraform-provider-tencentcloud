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
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
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

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var instanceId string
	if v, ok := d.GetOk("cluster_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(clusterId + FILED_SP + instanceId)

	return resourceTencentCloudCynosdbInstanceParamUpdate(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	instanceParam, err := service.DescribeCynosdbInstanceParamById(ctx, clusterId, instanceId)
	if err != nil {
		return err
	}

	if instanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_ids", instanceId)

	if instanceParam.ParamsItems != nil {
		currentParamMap := make(map[string]*cynosdb.ParamItemDetail)
		for _, param := range instanceParam.ParamsItems {
			currentParamMap[*param.ParamName] = param
		}
		instanceParamListList := []interface{}{}
		if v, ok := d.GetOk("param_items"); ok {
			paramItems := v.([]interface{})
			for _, paramItem := range paramItems {
				item := paramItem.(map[string]interface{})
				name := item["param_name"].(string)
				oldValue := item["old_value"].(string)
				currentParamItem := make(map[string]string)
				currentParamItem["param_name"] = name
				currentParamItem["current_value"] = *currentParamMap[name].CurrentValue
				if oldValue != "" {
					currentParamItem["old_value"] = oldValue
				}
				instanceParamListList = append(instanceParamListList, currentParamItem)
			}
		}

		_ = d.Set("instance_param_list", instanceParamListList)

	}

	return nil
}

func resourceTencentCloudCynosdbInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := cynosdb.NewModifyInstanceParamRequest()
	response := cynosdb.NewModifyInstanceParamResponse()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]
	request.ClusterId = &clusterId
	request.InstanceIds = []*string{&instanceId}

	if v, ok := d.GetOk("instance_param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			modifyParamItem := cynosdb.ModifyParamItem{}
			if v, ok := dMap["param_name"]; ok {
				modifyParamItem.ParamName = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				modifyParamItem.CurrentValue = helper.String(v.(string))
			}
			if v, ok := dMap["old_value"]; ok {
				modifyParamItem.OldValue = helper.String(v.(string))
			}
			request.InstanceParamList = append(request.InstanceParamList, &modifyParamItem)
		}
	}

	if v, ok := d.GetOk("is_in_maintain_period"); ok {
		request.IsInMaintainPeriod = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyInstanceParam(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceParam failed, reason:%+v", logId, err)
		return err
	}

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("update cynosdb instanceParam is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceParam fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbInstanceParamRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_instance_param.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
