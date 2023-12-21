package eb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEbEventBus() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventBusCreate,
		Read:   resourceTencentCloudEbEventBusRead,
		Update: resourceTencentCloudEbEventBusUpdate,
		Delete: resourceTencentCloudEbEventBusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"event_bus_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event set name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event set description, unlimited character type, description within 200 characters.",
			},

			"save_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "EB storage duration.",
			},

			"enable_store": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the EB storage is enabled.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudEbEventBusCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_bus.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = eb.NewCreateEventBusRequest()
		response   = eb.NewCreateEventBusResponse()
		eventBusId string
	)
	if v, ok := d.GetOk("event_bus_name"); ok {
		request.EventBusName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("save_days"); ok {
		request.SaveDays = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_store"); ok {
		request.EnableStore = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().CreateEventBus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb eventBus failed, reason:%+v", logId, err)
		return err
	}

	eventBusId = *response.Response.EventBusId
	d.SetId(eventBusId)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::eb:%s:uin/:eventbusid/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudEbEventBusRead(d, meta)
}

func resourceTencentCloudEbEventBusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_bus.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	eventBusId := d.Id()

	eventBus, err := service.DescribeEbEventBusById(ctx, eventBusId)
	if err != nil {
		return err
	}

	if eventBus == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventBus` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if eventBus.EventBusName != nil {
		_ = d.Set("event_bus_name", eventBus.EventBusName)
	}

	if eventBus.Description != nil {
		_ = d.Set("description", eventBus.Description)
	}

	if eventBus.SaveDays != nil {
		_ = d.Set("save_days", eventBus.SaveDays)
	}

	if eventBus.EnableStore != nil {
		_ = d.Set("enable_store", eventBus.EnableStore)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "eb", "eventbusid", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudEbEventBusUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_bus.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := eb.NewUpdateEventBusRequest()

	eventBusId := d.Id()

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	eventBus, err := service.DescribeEbEventBusById(ctx, eventBusId)
	if err != nil {
		return err
	}

	if eventBus == nil {
		return fmt.Errorf("[ERROR] resource `EbEventBus` [%s] not found, please check if it has been deleted.\n", d.Id())
	}

	request.EventBusId = &eventBusId
	request.LogTopicId = eventBus.LogTopicId

	if v, ok := d.GetOkExists("enable_store"); ok {
		request.EnableStore = helper.Bool(v.(bool))
	} else {
		return fmt.Errorf("[ERROR] When EbEventBus is modified, `enable_store` must be entered.\n")
	}

	if v, ok := d.GetOkExists("save_days"); ok {
		request.SaveDays = helper.IntInt64(v.(int))
	} else {
		return fmt.Errorf("[ERROR] When EbEventBus is modified, `save_days` must be entered.\n")
	}

	if d.HasChange("event_bus_name") {
		if v, ok := d.GetOk("event_bus_name"); ok {
			request.EventBusName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().UpdateEventBus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb eventBus failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("eb", "eventbusid", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudEbEventBusRead(d, meta)
}

func resourceTencentCloudEbEventBusDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_bus.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	eventBusId := d.Id()

	if err := service.DeleteEbEventBusById(ctx, eventBusId); err != nil {
		return err
	}

	return nil
}
