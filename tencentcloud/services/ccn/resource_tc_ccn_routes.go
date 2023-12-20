package ccn

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func ResourceTencentCloudCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRoutesCreate,
		Read:   resourceTencentCloudCcnRoutesRead,
		Update: resourceTencentCloudCcnRoutesUpdate,
		Delete: resourceTencentCloudCcnRoutesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "CCN Instance ID.",
			},

			"route_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "CCN Route Id List.",
			},

			"switch": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`on`: Enable, `off`: Disable.",
			},
		},
	}
}

func resourceTencentCloudCcnRoutesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_routes.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	ccnId := d.Get("ccn_id").(string)
	routeId := d.Get("route_id").(string)

	d.SetId(ccnId + tccommon.FILED_SP + routeId)

	return resourceTencentCloudCcnRoutesUpdate(d, meta)
}

func resourceTencentCloudCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_routes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	ccnId := idSplit[0]
	routeId := idSplit[1]

	ccnRoutes, err := service.DescribeVpcCcnRoutesById(ctx, ccnId, routeId)
	if err != nil {
		return err
	}

	if ccnRoutes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcCcnRoutes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	_ = d.Set("route_id", routeId)

	if ccnRoutes.Enabled != nil {
		if *ccnRoutes.Enabled {
			_ = d.Set("switch", "on")
		} else {
			_ = d.Set("switch", "off")
		}
	}

	return nil
}

func resourceTencentCloudCcnRoutesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_routes.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	ccnId := idSplit[0]
	routeId := idSplit[1]

	certSwitch := d.Get("switch").(string)

	if certSwitch == "on" {

		var (
			request = vpc.NewEnableCcnRoutesRequest()
		)
		request.CcnId = &ccnId
		request.RouteIds = []*string{&routeId}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableCcnRoutes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc ccnRoutes failed, reason:%+v", logId, err)
			return err
		}
	} else {

		var (
			request = vpc.NewDisableCcnRoutesRequest()
		)
		request.CcnId = &ccnId
		request.RouteIds = []*string{&routeId}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisableCcnRoutes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc ccnRoutes failed, reason:%+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudCcnRoutesRead(d, meta)
}

func resourceTencentCloudCcnRoutesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_routes.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
