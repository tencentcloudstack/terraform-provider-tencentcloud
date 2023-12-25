package vpc

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcPeerConnectAcceptOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPeerConnectAcceptOperationCreate,
		Read:   resourceTencentCloudVpcPeerConnectAcceptOperationRead,
		Delete: resourceTencentCloudVpcPeerConnectAcceptOperationDelete,
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

func resourceTencentCloudVpcPeerConnectAcceptOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_accept_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request             = vpc.NewAcceptVpcPeeringConnectionRequest()
		peeringConnectionId string
	)
	if v, ok := d.GetOk("peering_connection_id"); ok {
		peeringConnectionId = v.(string)
		request.PeeringConnectionId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AcceptVpcPeeringConnection(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc PeerConnectAcceptOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(peeringConnectionId)

	return resourceTencentCloudVpcPeerConnectAcceptOperationRead(d, meta)
}

func resourceTencentCloudVpcPeerConnectAcceptOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_accept_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcPeerConnectAcceptOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_peer_connect_accept_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
