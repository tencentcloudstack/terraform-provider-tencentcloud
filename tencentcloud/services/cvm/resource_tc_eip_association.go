package cvm

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipAssociationCreate,
		Read:   resourceTencentCloudEipAssociationRead,
		Delete: resourceTencentCloudEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 25),
				Description:  "The ID of EIP.",
			},
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"network_interface_id",
					"private_ip",
				},
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 25),
				Description:  "The CVM or CLB instance id going to bind with the EIP. This field is conflict with `network_interface_id` and `private_ip fields`.",
			},
			"network_interface_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 25),
				ConflictsWith: []string{
					"instance_id",
				},
				Description: "Indicates the network interface id like `eni-xxxxxx`. This field is conflict with `instance_id`.",
			},
			"private_ip": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(7, 25),
				ConflictsWith: []string{
					"instance_id",
				},
				Description: "Indicates an IP belongs to the `network_interface_id`. This field is conflict with `instance_id`.",
			},
		},
	}
}

func resourceTencentCloudEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_association.create")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
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

func resourceTencentCloudEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_association.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id         = d.Id()
	)

	association, err := ParseEipAssociationId(id)
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, association.EipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}

		if eip == nil {
			d.SetId("")
		}

		return nil
	})

	if err != nil {
		return err
	}

	_ = d.Set("eip_id", association.EipId)
	// associate with instance
	if len(association.InstanceId) > 0 {
		_ = d.Set("instance_id", association.InstanceId)
		return nil
	}

	_ = d.Set("network_interface_id", association.NetworkInterfaceId)
	_ = d.Set("private_ip", association.PrivateIp)
	return nil
}

func resourceTencentCloudEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip_association.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService = svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id         = d.Id()
	)

	association, err := ParseEipAssociationId(id)
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := vpcService.UnattachEip(ctx, association.EipId)
		if e != nil {
			return tccommon.RetryError(e, "DesOperation.MutexTaskRunning")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type EipAssociationId struct {
	EipId              string
	InstanceId         string
	NetworkInterfaceId string
	PrivateIp          string
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
