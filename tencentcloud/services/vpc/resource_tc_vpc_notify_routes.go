package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcNotifyRoutes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcNotifyRoutesCreate,
		Read:   resourceTencentCloudVpcNotifyRoutesRead,
		Delete: resourceTencentCloudVpcNotifyRoutesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the routing table.",
			},

			"route_item_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The unique ID of the routing policy.",
			},

			"expected_published_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Set the desired publication status: true: published; false: not published.",
			},

			"published_to_vbc": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to publish policies to vbc.",
			},
		},
	}
}

func resourceTencentCloudVpcNotifyRoutesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_notify_routes.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = vpc.NewNotifyRoutesRequest()
		routeTableId string
		routeItemId  string
	)

	if v, ok := d.GetOk("route_table_id"); ok {
		routeTableId = v.(string)
		request.RouteTableId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_item_ids"); ok {
		routeItemIdsSet := v.(*schema.Set).List()
		for i := range routeItemIdsSet {
			routeItemId = routeItemIdsSet[i].(string)
			request.RouteItemIds = append(request.RouteItemIds, &routeItemId)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().NotifyRoutes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate vpc notifyRoutes failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(routeTableId + tccommon.FILED_SP + routeItemId)
	return resourceTencentCloudVpcNotifyRoutesRead(d, meta)
}

func resourceTencentCloudVpcNotifyRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_notify_routes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routeTableId := idSplit[0]
	routeItemId := idSplit[1]

	notifyRoutes, err := service.DescribeVpcNotifyRoutesById(ctx, routeTableId, routeItemId)
	if err != nil {
		return err
	}

	if notifyRoutes == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_notify_routes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if notifyRoutes.RouteTableId != nil {
		_ = d.Set("route_table_id", notifyRoutes.RouteTableId)
	}

	if notifyRoutes.RouteItemId != nil {
		_ = d.Set("route_item_ids", []*string{notifyRoutes.RouteItemId})
	}

	if notifyRoutes.PublishedToVbc != nil {
		_ = d.Set("expected_published_status", notifyRoutes.PublishedToVbc)
	}

	if notifyRoutes.PublishedToVbc != nil {
		_ = d.Set("published_to_vbc", notifyRoutes.PublishedToVbc)
	}

	return nil
}

func resourceTencentCloudVpcNotifyRoutesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_notify_routes.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routeTableId := idSplit[0]
	routeItemId := idSplit[1]

	if err := service.DeleteVpcNotifyRoutesById(ctx, routeTableId, routeItemId); err != nil {
		return err
	}

	return nil
}
