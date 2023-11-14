/*
Provides a resource to create a cdb password_complexity

Example Usage

```hcl
resource "tencentcloud_cdb_password_complexity" "password_complexity" {
  instance_ids =
  param_list {
		name = ""
		current_value = ""

  }
}
```

Import

cdb password_complexity can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_password_complexity.password_complexity password_complexity_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
	"time"
)

func resourceTencentCloudCdbPasswordComplexity() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbPasswordComplexityCreate,
		Read:   resourceTencentCloudCdbPasswordComplexityRead,
		Update: resourceTencentCloudCdbPasswordComplexityUpdate,
		Delete: resourceTencentCloudCdbPasswordComplexityDelete,
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
				Description: "Instance ID list.",
			},

			"param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of parameters to be modified. Every element is a combination of `Name` (parameter name) and `CurrentValue` (new value). Valid values for `Name` of version 8.0: `validate_password.policy`, `validate_password.lengt`, `validate_password.mixed_case_coun`, `validate_password.number_coun`, `validate_password.special_char_count`. Valid values for `Name` of version 5.6 and 5.7: `validate_password_polic`, `validate_password_lengt` `validate_password_mixed_case_coun`, `validate_password_number_coun`, `validate_password_special_char_coun`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCdbPasswordComplexityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_password_complexity.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCdbPasswordComplexityUpdate(d, meta)
}

func resourceTencentCloudCdbPasswordComplexityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_password_complexity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	passwordComplexityId := d.Id()

	passwordComplexity, err := service.DescribeCdbPasswordComplexityById(ctx, instanceId)
	if err != nil {
		return err
	}

	if passwordComplexity == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbPasswordComplexity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if passwordComplexity.InstanceIds != nil {
		_ = d.Set("instance_ids", passwordComplexity.InstanceIds)
	}

	if passwordComplexity.ParamList != nil {
		paramListList := []interface{}{}
		for _, paramList := range passwordComplexity.ParamList {
			paramListMap := map[string]interface{}{}

			if passwordComplexity.ParamList.Name != nil {
				paramListMap["name"] = passwordComplexity.ParamList.Name
			}

			if passwordComplexity.ParamList.CurrentValue != nil {
				paramListMap["current_value"] = passwordComplexity.ParamList.CurrentValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	return nil
}

func resourceTencentCloudCdbPasswordComplexityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_password_complexity.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cdb.NewModifyInstancePasswordComplexityRequest()

	passwordComplexityId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_ids", "param_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().ModifyInstancePasswordComplexity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cdb passwordComplexity failed, reason:%+v", logId, err)
		return err
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbPasswordComplexityStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbPasswordComplexityRead(d, meta)
}

func resourceTencentCloudCdbPasswordComplexityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_password_complexity.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
