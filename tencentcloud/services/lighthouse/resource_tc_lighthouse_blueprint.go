package lighthouse

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseBlueprint() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_blueprint.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CreateBlueprint(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 10*tccommon.ReadRetryTimeout, time.Second, service.LighthouseBlueprintStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseBlueprintRead(d, meta)
}

func resourceTencentCloudLighthouseBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_blueprint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_blueprint.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyBlueprintAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_blueprint.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	blueprintId := d.Id()

	if err := service.DeleteLighthouseBlueprintById(ctx, blueprintId); err != nil {
		return err
	}

	return nil
}
