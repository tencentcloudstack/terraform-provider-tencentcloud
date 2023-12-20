package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcPeerConnectRejectOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request             = vpc.NewRejectVpcPeeringConnectionRequest()
		peeringConnectionId string
	)
	if v, ok := d.GetOk("peering_connection_id"); ok {
		peeringConnectionId = v.(string)
		request.PeeringConnectionId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().RejectVpcPeeringConnection(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcPeerConnectRejectOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_reject_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
