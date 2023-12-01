/*
Provides a resource to create a lighthouse blueprint

Example Usage

```hcl
resource "tencentcloud_lighthouse_blueprint" "blueprint" {
  blueprint_name = "blueprint_name_test"
  description = "blueprint_description_test"
  instance_id = "lhins-xxxxxx"
}
```

Import

lighthouse blueprint can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_blueprint.blueprint blueprint_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseBlueprint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseBlueprintCreate,
		Read:   resourceTencentCloudLighthouseBlueprintRead,
		Update: resourceTencentCloudLighthouseBlueprintUpdate,
		Delete: resourceTencentCloudLighthouseBlueprintDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"blueprint_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Blueprint name, which can contain up to 60 characters.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Blueprint description, which can contain up to 60 characters.",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to make a blueprint.",
			},
		},
	}
}

func resourceTencentCloudLighthouseBlueprintCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_blueprint.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = lighthouse.NewCreateBlueprintRequest()
		response    = lighthouse.NewCreateBlueprintResponse()
		blueprintId string
	)
	if v, ok := d.GetOk("blueprint_name"); ok {
		request.BlueprintName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateBlueprint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse blueprint failed, reason:%+v", logId, err)
		return err
	}

	blueprintId = *response.Response.BlueprintId
	d.SetId(blueprintId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"NORMAL"}, 10*readRetryTimeout, time.Second, service.LighthouseBlueprintStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseBlueprintRead(d, meta)
}

func resourceTencentCloudLighthouseBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_blueprint.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	blueprintId := d.Id()

	blueprint, err := service.DescribeLighthouseBlueprintById(ctx, blueprintId)
	if err != nil {
		return err
	}

	if blueprint == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseBlueprint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if blueprint.BlueprintName != nil {
		_ = d.Set("blueprint_name", blueprint.BlueprintName)
	}

	if blueprint.Description != nil {
		_ = d.Set("description", blueprint.Description)
	}

	return nil
}

func resourceTencentCloudLighthouseBlueprintUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_blueprint.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifyBlueprintAttributeRequest()

	blueprintId := d.Id()

	request.BlueprintId = &blueprintId

	if d.HasChange("blueprint_name") {
		if v, ok := d.GetOk("blueprint_name"); ok {
			request.BlueprintName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifyBlueprintAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse blueprint failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseBlueprintRead(d, meta)
}

func resourceTencentCloudLighthouseBlueprintDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_blueprint.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	blueprintId := d.Id()

	if err := service.DeleteLighthouseBlueprintById(ctx, blueprintId); err != nil {
		return err
	}

	return nil
}
