package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRouteTableEntryConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRouteTableEntryConfigCreate,
		Read:   resourceTencentCloudRouteTableEntryConfigRead,
		Update: resourceTencentCloudRouteTableEntryConfigUpdate,
		Delete: resourceTencentCloudRouteTableEntryConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Route table ID.",
			},

			"route_item_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of route table entry.",
			},

			"disabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the entry is disabled.",
			},
		},
	}
}

func resourceTencentCloudRouteTableEntryConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_entry_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		routeTableId string
		routeItemId  string
	)

	if v, ok := d.GetOk("route_table_id"); ok {
		routeTableId = v.(string)
	}

	if v, ok := d.GetOk("route_item_id"); ok {
		routeItemId = v.(string)
	}

	d.SetId(strings.Join([]string{routeTableId, routeItemId}, tccommon.FILED_SP))

	return resourceTencentCloudRouteTableEntryConfigUpdate(d, meta)
}

func resourceTencentCloudRouteTableEntryConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_entry_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routeTableId := items[0]
	routeItemId := items[1]
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeRouteTable(ctx, routeTableId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		if has != 1 {
			e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return resource.NonRetryableError(e)
		}

		for _, v := range info.entryInfos {
			if v.routeItemId == routeItemId {
				_ = d.Set("route_table_id", routeTableId)
				_ = d.Set("route_item_id", routeItemId)
				_ = d.Set("disabled", !v.enabled)
				return nil
			}
		}

		d.SetId("")
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudRouteTableEntryConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_entry_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		disabled bool
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routeTableId := items[0]
	routeItemId := items[1]

	if v, ok := d.GetOkExists("disabled"); ok {
		disabled = v.(bool)
	}

	if disabled {
		request := vpc.NewDisableRoutesRequest()
		request.RouteTableId = &routeTableId
		request.RouteItemIds = helper.Strings([]string{routeItemId})
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisableRoutesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s disable route table entry config failed, reason:%+v", logId, err)
			return err
		}
	} else {
		request := vpc.NewEnableRoutesRequest()
		request.RouteTableId = &routeTableId
		request.RouteItemIds = helper.Strings([]string{routeItemId})
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableRoutesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s enable route table entry config failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudRouteTableEntryConfigRead(d, meta)
}

func resourceTencentCloudRouteTableEntryConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_route_table_entry_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
