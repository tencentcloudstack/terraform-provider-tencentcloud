package pls

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcEnableEndPointConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcEnableEndPointConnectCreate,
		Read:   resourceTencentCloudVpcEnableEndPointConnectRead,
		Delete: resourceTencentCloudVpcEnableEndPointConnectDelete,
		Schema: map[string]*schema.Schema{
			"end_point_service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Endpoint service ID.",
			},

			"end_point_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Endpoint ID.",
			},

			"accept_flag": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to accept endpoint connection requests. `true`: Accept automatically. `false`: Do not automatically accept.",
			},
		},
	}
}

func resourceTencentCloudVpcEnableEndPointConnectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_enable_end_point_connect.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request           = vpc.NewEnableVpcEndPointConnectRequest()
		endPointServiceId string
		endPointId        string
	)
	if v, ok := d.GetOk("end_point_service_id"); ok {
		endPointServiceId = v.(string)
		request.EndPointServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_id"); ok {
		endPointIdSet := v.(*schema.Set).List()
		for i := range endPointIdSet {
			endPointId = endPointIdSet[i].(string)
			request.EndPointId = append(request.EndPointId, &endPointId)
		}
	}

	if v, _ := d.GetOk("accept_flag"); v != nil {
		request.AcceptFlag = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableVpcEndPointConnect(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc enableEndPointConnect failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(endPointServiceId + tccommon.FILED_SP + endPointId)

	return resourceTencentCloudVpcEnableEndPointConnectRead(d, meta)
}

func resourceTencentCloudVpcEnableEndPointConnectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_enable_end_point_connect.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcEnableEndPointConnectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_enable_end_point_connect.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
