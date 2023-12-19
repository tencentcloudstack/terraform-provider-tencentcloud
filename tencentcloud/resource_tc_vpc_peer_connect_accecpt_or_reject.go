package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcPeerConnectAccecptOrReject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPeerConnectAccecptOrRejectCreate,
		Read:   resourceTencentCloudVpcPeerConnectAccecptOrRejectRead,
		Delete: resourceTencentCloudVpcPeerConnectAccecptOrRejectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceTencentCloudVpcPeerConnectAccecptOrRejectCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_accecpt_or_reject.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request             = vpc.NewAcceptVpcPeeringConnectionRequest()
		service             = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		peeringConnectionId string
		ctx                 = context.WithValue(context.TODO(), logIdKey, logId)
	)
	if v, ok := d.GetOk("peering_connection_id"); ok {
		peeringConnectionId = v.(string)
		request.PeeringConnectionId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AcceptVpcPeeringConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc PeerConnectAccecptOrReject failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {

		PeerConnectAcceptOrReject, errRet := service.DescribeVpcPeerConnectManagerById(ctx, peeringConnectionId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if PeerConnectAcceptOrReject.State == nil {
			return resource.RetryableError(fmt.Errorf("waiting for PeerConnectAcceptOrReject operation "))
		}
		if *PeerConnectAcceptOrReject.State == "FAILED" || *PeerConnectAcceptOrReject.State == "DELETED" ||
			*PeerConnectAcceptOrReject.State == "EXPIRED" {
			return resource.NonRetryableError(fmt.Errorf("failed operation"))
		}
		if *PeerConnectAcceptOrReject.State != "ACTIVE" {
			return resource.RetryableError(fmt.Errorf("waiting for PeerConnectAcceptOrReject %s operation", peeringConnectionId))
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(peeringConnectionId)

	return resourceTencentCloudVpcPeerConnectAccecptOrRejectRead(d, meta)
}

func resourceTencentCloudVpcPeerConnectAccecptOrRejectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_accecpt_or_reject.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	peeringConnectionId := d.Id()

	PeerConnectAcceptOrReject, err := service.DescribeVpcPeerConnectManagerById(ctx, peeringConnectionId)
	if err != nil {
		return err
	}

	if PeerConnectAcceptOrReject == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcPeerConnectAccecptOrReject` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if PeerConnectAcceptOrReject.PeeringConnectionId != nil {
		_ = d.Set("peering_connection_id", PeerConnectAcceptOrReject.PeeringConnectionId)
	}

	return nil
}

func resourceTencentCloudVpcPeerConnectAccecptOrRejectDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peer_connect_accecpt_or_reject.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	peeringConnectionId := d.Id()

	if err := service.DeleteVpcPeerConnectAccecptOrRejectById(ctx, peeringConnectionId); err != nil {
		return err
	}

	return nil
}
