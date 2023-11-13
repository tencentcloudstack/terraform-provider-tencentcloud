/*
Provides a resource to create a cdb instance_param

Example Usage

```hcl
resource "tencentcloud_cdb_instance_param" "instance_param" {
  instance_ids = &lt;nil&gt;
  param_list {
		name = &lt;nil&gt;
		current_value = &lt;nil&gt;

  }
  template_id = &lt;nil&gt;
  wait_switch = 0
  not_sync_ro = false
  not_sync_dr = false
}
```

Import

cdb instance_param can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_instance_param.instance_param instance_param_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbInstanceParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbInstanceParamCreate,
		Read:   resourceTencentCloudCdbInstanceParamRead,
		Update: resourceTencentCloudCdbInstanceParamUpdate,
		Delete: resourceTencentCloudCdbInstanceParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id list of instances.",
			},

			"param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A list of parameters to modify.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of parameter.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of parameter.",
						},
					},
				},
			},

			"template_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template id, ParamList and TemplateId must pass at least one of them.",
			},

			"wait_switch": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The way to perform parameter adjustment tasks, the default is 0. Supported values include: 0 - execute immediately, 1 - execute in time window; when the value is 1, only one instance can be passed at a time (the number of InstanceIds is 1).",
			},

			"not_sync_ro": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the parameter is synchronized to the read-only instance under the master instance. default to false.",
			},

			"not_sync_dr": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the parameters are synchronized to the disaster recovery instance under the primary instance. default to false.",
			},
		},
	}
}

func resourceTencentCloudCdbInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_param.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbInstanceParamUpdate(d, meta)
}

func resourceTencentCloudCdbInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_param.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceParamId := d.Id()

	instanceParam, err := service.DescribeCdbInstanceParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceParam.InstanceIds != nil {
		_ = d.Set("instance_ids", instanceParam.InstanceIds)
	}

	if instanceParam.ParamList != nil {
		paramListList := []interface{}{}
		for _, paramList := range instanceParam.ParamList {
			paramListMap := map[string]interface{}{}

			if instanceParam.ParamList.Name != nil {
				paramListMap["name"] = instanceParam.ParamList.Name
			}

			if instanceParam.ParamList.CurrentValue != nil {
				paramListMap["current_value"] = instanceParam.ParamList.CurrentValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	if instanceParam.TemplateId != nil {
		_ = d.Set("template_id", instanceParam.TemplateId)
	}

	if instanceParam.WaitSwitch != nil {
		_ = d.Set("wait_switch", instanceParam.WaitSwitch)
	}

	if instanceParam.NotSyncRo != nil {
		_ = d.Set("not_sync_ro", instanceParam.NotSyncRo)
	}

	if instanceParam.NotSyncDr != nil {
		_ = d.Set("not_sync_dr", instanceParam.NotSyncDr)
	}

	return nil
}

func resourceTencentCloudCdbInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_param.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyInstanceParamRequest()

	instanceParamId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_ids", "param_list", "template_id", "wait_switch", "not_sync_ro", "not_sync_dr"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_ids") {
		if v, ok := d.GetOk("instance_ids"); ok {
			instanceIdsSet := v.(*schema.Set).List()
			for i := range instanceIdsSet {
				instanceIds := instanceIdsSet[i].(string)
				request.InstanceIds = append(request.InstanceIds, &instanceIds)
			}
		}
	}

	if d.HasChange("param_list") {
		if v, ok := d.GetOk("param_list"); ok {
			for _, item := range v.([]interface{}) {
				parameter := cdb.Parameter{}
				if v, ok := dMap["name"]; ok {
					parameter.Name = helper.String(v.(string))
				}
				if v, ok := dMap["current_value"]; ok {
					parameter.CurrentValue = helper.String(v.(string))
				}
				request.ParamList = append(request.ParamList, &parameter)
			}
		}
	}

	if d.HasChange("template_id") {
		if v, ok := d.GetOkExists("template_id"); ok {
			request.TemplateId = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("wait_switch") {
		if v, ok := d.GetOkExists("wait_switch"); ok {
			request.WaitSwitch = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("not_sync_ro") {
		if v, ok := d.GetOkExists("not_sync_ro"); ok {
			request.NotSyncRo = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("not_sync_dr") {
		if v, ok := d.GetOkExists("not_sync_dr"); ok {
			request.NotSyncDr = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyInstanceParam(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb instanceParam failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCdbInstanceParamRead(d, meta)
}

func resourceTencentCloudCdbInstanceParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_instance_param.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
