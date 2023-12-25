package vpc

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcPeerConnectRejectOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPeerConnectRejectOperationCreate,
		Read:   resourceTencentCloudVpcPeerConnectRejectOperationRead,
		Delete: resourceTencentCloudVpcPeerConnectRejectOperationDelete,
		Schema: map[string]*schema.Schema{
			"peering_connection_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Peer connection unique ID.",
			},
		},
	}
}

func resourceTencentCloudVpcPeerConnectRejectOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request             = vpc.NewRejectVpcPeeringConnectionRequest()
		peeringConnectionId string
	)
	if v, ok := d.GetOk("peering_connection_id"); ok {
		peeringConnectionId = v.(string)
		request.PeeringConnectionId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().RejectVpcPeeringConnection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc PeerConnectRejectOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(peeringConnectionId)

	return resourceTencentCloudVpcPeerConnectRejectOperationRead(d, meta)
}

func resourceTencentCloudVpcPeerConnectRejectOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcPeerConnectRejectOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
