package ccn

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCcnAttachment() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.81.198. Please use 'tencentcloud_ccn_attachment_v2' instead.",
		Create:             resourceTencentCloudCcnAttachmentCreate,
		Read:               resourceTencentCloudCcnAttachmentRead,
		Update:             resourceTencentCloudCcnAttachmentUpdate,
		Delete:             resourceTencentCloudCcnAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CNN_INSTANCE_TYPE_VPC, CNN_INSTANCE_TYPE_DIRECTCONNECT, CNN_INSTANCE_TYPE_BMVPC, CNN_INSTANCE_TYPE_VPNGW}),
				ForceNew:     true,
				Description:  "Type of attached instance network, and available values include `VPC`, `DIRECTCONNECT`, `BMVPC` and `VPNGW`. Note: `VPNGW` type is only for whitelist customer now.",
			},
			"instance_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region that the instance locates at.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of instance is attached.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark of attachment.",
			},
			"ccn_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Uin of the ccn attached. If not set, which means the uin of this account. This parameter is used with case when attaching ccn of other account to the instance of this account. For now only support instance type `VPC`.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Ccn instance route table ID.",
			},
			// Computed values
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "States of instance is attached. Valid values: `PENDING`, `ACTIVE`, `EXPIRED`, `REJECTED`, `DELETED`, `FAILED`, `ATTACHING`, `DETACHING` and `DETACHFAILED`. `FAILED` means asynchronous forced disassociation after 2 hours. `DETACHFAILED` means asynchronous forced disassociation after 2 hours.",
			},
			"attached_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time of attaching.",
			},
			"cidr_block": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A network address block of the instance that is attached.",
			},
			"route_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Route id list.",
			},
		},
	}
}

func resourceTencentCloudCcnAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment.create")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
		description    = ""
		ccnUin         = ""
		routeTableId   = ""
	)

	if len(ccnId) < 4 || len(instanceRegion) < 3 || len(instanceId) < 3 {
		return fmt.Errorf("param ccn_id or instance_region or instance_id  error")
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}

	if v, ok := d.GetOk("ccn_uin"); ok {
		ccnUin = v.(string)
		if ccnUin != "" && instanceType != CNN_INSTANCE_TYPE_VPC {
			return fmt.Errorf("Other ccn account attachment %s only support instance type of `VPC`.", ccnId)
		}
	} else {
		_, has, err := service.DescribeCcn(ctx, ccnId)
		if err != nil {
			return err
		}

		if has == 0 {
			return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
		}
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		routeTableId = v.(string)
	}

	if err := service.AttachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId, ccnUin, description, routeTableId); err != nil {
		return err
	}

	m := md5.New()
	_, err := m.Write([]byte(ccnId + instanceType + instanceRegion + instanceId))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	return resourceTencentCloudCcnAttachmentRead(d, meta)
}

func resourceTencentCloudCcnAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	var (
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	if _, ok := d.GetOk("ccn_uin"); !ok {
		onlineHas := true

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, e := service.DescribeCcn(ctx, ccnId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if has == 0 {
				d.SetId("")
				onlineHas = false
				return nil
			}

			return nil
		})

		if err != nil {
			return err
		}

		if !onlineHas {
			return nil
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeCcnAttachedInstance(ctx, ccnId, instanceRegion, instanceType, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		if v, ok := d.GetOk("ccn_uin"); ok && v.(string) != info.ccnUin {
			d.SetId("")
			return nil
		}

		_ = d.Set("description", info.description)
		_ = d.Set("route_table_id", info.routeTableId)
		_ = d.Set("state", strings.ToUpper(info.state))
		_ = d.Set("attached_time", info.attachedTime)
		_ = d.Set("cidr_block", info.cidrBlock)
		return nil
	})

	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := vpc.NewDescribeCcnRoutesRequest()
		request.CcnId = &ccnId
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeCcnRoutes(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		routeIds := make([]string, 0)
		if response != nil && response.Response != nil && len(response.Response.RouteSet) > 0 {
			for _, route := range response.Response.RouteSet {
				routeIds = append(routeIds, *route.RouteId)
			}
		}

		_ = d.Set("route_ids", routeIds)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudCcnAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment.update")()

	if d.HasChange("description") {
		var (
			logId          = tccommon.GetLogId(tccommon.ContextNil)
			request        = vpc.NewModifyCcnAttachedInstancesAttributeRequest()
			ccnInstance    vpc.CcnInstance
			ccnId          = d.Get("ccn_id").(string)
			instanceType   = d.Get("instance_type").(string)
			instanceRegion = d.Get("instance_region").(string)
			instanceId     = d.Get("instance_id").(string)
			description    = d.Get("description").(string)
			routeTableId   = d.Get("route_table_id").(string)
		)

		request.CcnId = &ccnId
		ccnInstance.InstanceId = &instanceId
		ccnInstance.InstanceRegion = &instanceRegion
		ccnInstance.InstanceType = &instanceType
		ccnInstance.Description = &description
		ccnInstance.RouteTableId = &routeTableId
		request.Instances = []*vpc.CcnInstance{&ccnInstance}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyCcnAttachedInstancesAttribute(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return tccommon.RetryError(err)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		})

		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCcnAttachmentRead(d, meta)
}

func resourceTencentCloudCcnAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment.delete")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId          = d.Get("ccn_id").(string)
		instanceType   = d.Get("instance_type").(string)
		instanceRegion = d.Get("instance_region").(string)
		instanceId     = d.Get("instance_id").(string)
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeCcn(ctx, ccnId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := service.DetachCcnInstances(ctx, ccnId, instanceRegion, instanceType, instanceId); err != nil {
		return err
	}

	return resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, err := service.DescribeCcnAttachedInstance(ctx, ccnId, instanceRegion, instanceType, instanceId)
		if err != nil {
			return resource.RetryableError(err)
		}

		if has == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
