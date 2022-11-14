/*
Provides a resource to create a rum taw_instance

Example Usage

```hcl
resource "tencentcloud_rum_taw_instance" "taw_instance" {
  area_id = ""
  charge_type = ""
  data_retention_days = ""
  instance_name = ""
  tags {
			key = ""
			value = ""

  }
  instance_desc = ""
  count_num = ""
  period_retain = ""
  buying_channel = ""
}

```
Import

rum taw_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_taw_instance.taw_instance tawInstance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRumTawInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRumTawInstanceRead,
		Create: resourceTencentCloudRumTawInstanceCreate,
		Update: resourceTencentCloudRumTawInstanceUpdate,
		Delete: resourceTencentCloudRumTawInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"area_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Region ID (at least greater than 0).",
			},

			"charge_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Billing type (1: Pay-as-you-go).",
			},

			"data_retention_days": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Data retention period (at least greater than 0).",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name (up to 255 bytes).",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"instance_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance description (up to 1,024 bytes).",
			},

			"count_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of data entries reported per day.",
			},

			"period_retain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Billing for data storage.",
			},

			"buying_channel": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance purchase channel. Valid value: cdn.",
			},

			"instance_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance status (`1` = creating, `2` = running, `3` = exception, `4` = restarting, `5` = stopping, `6` = stopped, `7` = deleted).",
			},

			"cluster_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Cluster ID.",
			},

			"charge_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Billing status (`1` = in use, `2` = expired, `3` = destroyed, `4` = assigning, `5` = failed).",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},
		},
	}
}

func resourceTencentCloudRumTawInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_taw_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = rum.NewCreateTawInstanceRequest()
		response   *rum.CreateTawInstanceResponse
		instanceId string
	)

	if v, ok := d.GetOk("area_id"); ok {
		request.AreaId = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("data_retention_days"); ok {
		request.DataRetentionDays = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("instance_name"); ok {

		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := rum.Tag{}
			if v, ok := dMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("instance_desc"); ok {
		request.InstanceDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("count_num"); ok {
		request.CountNum = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period_retain"); ok {
		request.PeriodRetain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("buying_channel"); ok {
		request.BuyingChannel = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateTawInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum tawInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId

	d.SetId(instanceId)
	return resourceTencentCloudRumTawInstanceRead(d, meta)
}

func resourceTencentCloudRumTawInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_taw_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	tawInstance, err := service.DescribeRumTawInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if tawInstance == nil {
		d.SetId("")
		return fmt.Errorf("resource `tawInstance` %s does not exist", instanceId)
	}

	if tawInstance.AreaId != nil {
		_ = d.Set("area_id", tawInstance.AreaId)
	}

	if tawInstance.ChargeType != nil {
		_ = d.Set("charge_type", tawInstance.ChargeType)
	}

	if tawInstance.DataRetentionDays != nil {
		_ = d.Set("data_retention_days", tawInstance.DataRetentionDays)
	}

	if tawInstance.InstanceName != nil {
		_ = d.Set("instance_name", tawInstance.InstanceName)
	}

	if tawInstance.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range tawInstance.Tags {
			tagsMap := map[string]interface{}{}
			if tags.Key != nil {
				tagsMap["key"] = tags.Key
			}
			if tags.Value != nil {
				tagsMap["value"] = tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if tawInstance.InstanceDesc != nil {
		_ = d.Set("instance_desc", tawInstance.InstanceDesc)
	}

	// if tawInstance.CountNum != nil {
	// 	_ = d.Set("count_num", tawInstance.CountNum)
	// }

	// if tawInstance.PeriodRetain != nil {
	// 	_ = d.Set("period_retain", tawInstance.PeriodRetain)
	// }

	// if tawInstance.BuyingChannel != nil {
	// 	_ = d.Set("buying_channel", tawInstance.BuyingChannel)
	// }

	if tawInstance.InstanceStatus != nil {
		_ = d.Set("instance_status", tawInstance.InstanceStatus)
	}

	if tawInstance.ClusterId != nil {
		_ = d.Set("cluster_id", tawInstance.ClusterId)
	}

	if tawInstance.ChargeStatus != nil {
		_ = d.Set("charge_status", tawInstance.ChargeStatus)
	}

	if tawInstance.UpdatedAt != nil {
		_ = d.Set("updated_at", tawInstance.UpdatedAt)
	}

	// if tawInstance.DataRetentionDays != nil {
	// 	_ = d.Set("data_retention_days", tawInstance.DataRetentionDays)
	// }

	if tawInstance.CreatedAt != nil {
		_ = d.Set("created_at", tawInstance.CreatedAt)
	}

	return nil
}

func resourceTencentCloudRumTawInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_taw_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := rum.NewModifyInstanceRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("area_id") {
		return fmt.Errorf("`area_id` do not support change now.")
	}

	if d.HasChange("charge_type") {
		return fmt.Errorf("`charge_type` do not support change now.")
	}

	if d.HasChange("data_retention_days") {
		return fmt.Errorf("`data_retention_days` do not support change now.")
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		return fmt.Errorf("`tags` do not support change now.")
	}

	if d.HasChange("instance_desc") {
		if v, ok := d.GetOk("instance_desc"); ok {
			request.InstanceDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("count_num") {
		return fmt.Errorf("`count_num` do not support change now.")
	}

	if d.HasChange("period_retain") {
		return fmt.Errorf("`period_retain` do not support change now.")
	}

	if d.HasChange("buying_channel") {
		return fmt.Errorf("`buying_channel` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().ModifyInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum tawInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRumTawInstanceRead(d, meta)
}

func resourceTencentCloudRumTawInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_taw_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteRumTawInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
