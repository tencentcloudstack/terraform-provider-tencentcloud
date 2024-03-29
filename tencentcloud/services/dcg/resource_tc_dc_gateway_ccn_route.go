package dcg

import (
	"context"
	"fmt"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDcGatewayCcnRouteInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayCcnRouteCreate,
		Read:   resourceTencentCloudDcGatewayCcnRouteRead,
		Delete: resourceTencentCloudDcGatewayCcnRouteDelete,
		Schema: map[string]*schema.Schema{
			"dcg_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the DCG.",
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCIDRNetworkAddress,
				Description:  "A network address segment of IDC.",
			},

			//compute
			"as_path": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "As path list of the BGP.",
			},
		},
	}
}

func resourceTencentCloudDcGatewayCcnRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway_ccn_route.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		dcgId     = d.Get("dcg_id").(string)
		cidrBlock = d.Get("cidr_block").(string)
	)

	//Modification of this parameter[as_path] is not yet supported
	routeId, err := service.CreateDirectConnectGatewayCcnRoute(ctx, dcgId, cidrBlock, nil)

	if err != nil {
		return err
	}

	d.SetId(dcgId + "#" + routeId)

	// add sleep protect, either network_instance_id will be set "".
	time.Sleep(1 * time.Second)

	return resourceTencentCloudDcGatewayCcnRouteRead(d, meta)
}

func resourceTencentCloudDcGatewayCcnRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway_ccn_route.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
	}

	dcgId, routeId := items[0], items[1]
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("dcg_id", info.dcgId)
		_ = d.Set("cidr_block", info.cidrBlock)
		_ = d.Set("as_path", info.asPaths)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudDcGatewayCcnRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway_ccn_route.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_dc_gateway_ccn_route is wrong")
	}

	dcgId, routeId := items[0], items[1]
	_, has, err := service.DescribeDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
	if err != nil {
		return err
	}

	if has == 0 {
		return nil
	}

	return service.DeleteDirectConnectGatewayCcnRoute(ctx, dcgId, routeId)
}
