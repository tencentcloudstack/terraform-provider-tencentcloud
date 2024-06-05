package cvm

import (
	"context"
	"fmt"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type EipAssociationId struct {
	EipId              string
	InstanceId         string
	NetworkInterfaceId string
	PrivateIp          string
}

func resourceTencentCloudEipAssociationCreateOnExit(ctx context.Context) error {
	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		meta       = tccommon.ProviderMetaFromContext(ctx)
		vpcService = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		d          = tccommon.ResourceDataFromContext(ctx)
		eip        *vpc.Address
		errRet     error
	)

	eipId := d.Get("eip_id").(string)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet = vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if eip == nil {
			return resource.NonRetryableError(fmt.Errorf("eip is not found"))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if *eip.AddressStatus != svcvpc.EIP_STATUS_UNBIND {
		return fmt.Errorf("eip status is illegal %s", *eip.AddressStatus)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId := v.(string)
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := vpcService.AttachEip(ctx, eipId, instanceId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			return err
		}

		associationId := fmt.Sprintf("%v::%v", eipId, instanceId)
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			eip, errRet = vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}

			if eip == nil {
				return resource.NonRetryableError(fmt.Errorf("eip is not found"))
			}

			if *eip.AddressStatus == svcvpc.EIP_STATUS_BIND {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("wait for binding success: %s", *eip.AddressStatus))
		})

		if err != nil {
			return err
		}

		d.SetId(associationId)
		return resourceTencentCloudEipAssociationRead(d, meta)
	}

	needRequest := false
	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eipId
	var networkId string
	var privateIp string
	if v, ok := d.GetOk("network_interface_id"); ok {
		needRequest = true
		networkId = v.(string)
		request.NetworkInterfaceId = &networkId
	}

	if v, ok := d.GetOk("private_ip"); ok {
		needRequest = true
		privateIp = v.(string)
		request.PrivateIpAddress = &privateIp
	}

	if needRequest {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateAddress(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		})

		if err != nil {
			return err
		}

		id := fmt.Sprintf("%v::%v::%v", eipId, networkId, privateIp)

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			eip, errRet = vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}

			if eip == nil {
				return resource.NonRetryableError(fmt.Errorf("eip is not found"))
			}

			if *eip.AddressStatus == svcvpc.EIP_STATUS_BIND_ENI || *eip.AddressStatus == svcvpc.EIP_STATUS_BIND {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("wait for binding success: %s", *eip.AddressStatus))
		})

		if err != nil {
			return err
		}

		d.SetId(id)
		return resourceTencentCloudEipAssociationRead(d, meta)
	}

	return nil
}

func resourceTencentCloudEipAssociationReadPreRequest0(ctx context.Context, req *vpc.DescribeAddressesRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	id := d.Id()
	association, err := ParseEipAssociationId(id)
	if err != nil {
		return err
	}

	_ = d.Set("eip_id", association.EipId)
	if len(association.InstanceId) > 0 {
		_ = d.Set("instance_id", association.InstanceId)
		return nil
	}

	_ = d.Set("network_interface_id", association.NetworkInterfaceId)
	_ = d.Set("private_ip", association.PrivateIp)

	req.AddressIds = []*string{&association.EipId}
	return nil
}

func ParseEipAssociationId(associationId string) (association EipAssociationId, errRet error) {
	ids := strings.Split(associationId, "::")
	if len(ids) < 2 || len(ids) > 3 {
		errRet = fmt.Errorf("Invalid eip association ID: %v", associationId)
		return
	}
	association.EipId = ids[0]

	// associate with instance
	if len(ids) == 2 {
		association.InstanceId = ids[1]
		return
	}

	// associate with network interface
	association.NetworkInterfaceId = ids[1]
	association.PrivateIp = ids[2]
	return
}
