package ccn

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCcnAttachmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnAttachmentV2Create,
		Read:   resourceTencentCloudCcnAttachmentV2Read,
		Update: resourceTencentCloudCcnAttachmentV2Update,
		Delete: resourceTencentCloudCcnAttachmentV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of instance is attached.",
			},
			"instance_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region that the instance locates at.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CNN_INSTANCE_TYPE_VPC, CNN_INSTANCE_TYPE_DIRECTCONNECT, CNN_INSTANCE_TYPE_BMVPC, CNN_INSTANCE_TYPE_VPNGW}),
				ForceNew:     true,
				Description:  "Type of attached instance network, and available values include `VPC`, `DIRECTCONNECT`, `BMVPC` and `VPNGW`. Note: `VPNGW` type is only for whitelist customer now.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark of attachment.",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Ccn instance route table ID.",
			},
			"ccn_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Uin of the ccn attached. If not set, which means the uin of this account. This parameter is used with case when attaching ccn of other account to the instance of this account. For now only support instance type `VPC`.",
			},
			// Computed
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
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A network address block of the instance that is attached.",
			},
			"route_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Route id list.",
			},
		},
	}
}

func resourceTencentCloudCcnAttachmentV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request        = vpc.NewAttachCcnInstancesRequest()
		ccnId          string
		instanceId     string
		instanceRegion string
		instanceType   string
	)

	if v, ok := d.GetOk("ccn_id"); ok {
		request.CcnId = helper.String(v.(string))
		ccnId = v.(string)
	}

	ccnInstance := new(vpc.CcnInstance)
	if v, ok := d.GetOk("instance_id"); ok {
		ccnInstance.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_region"); ok {
		ccnInstance.InstanceRegion = helper.String(v.(string))
		instanceRegion = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		ccnInstance.InstanceType = helper.String(v.(string))
		instanceType = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		ccnInstance.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		ccnInstance.RouteTableId = helper.String(v.(string))
	}

	request.Instances = append(request.Instances, ccnInstance)

	if v, ok := d.GetOk("ccn_uin"); ok {
		request.CcnUin = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AttachCcnInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create attach ccn instance failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create attach ccn instance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{ccnId, instanceType, instanceRegion, instanceId}, tccommon.FILED_SP))
	return resourceTencentCloudCcnAttachmentV2Read(d, meta)
}

func resourceTencentCloudCcnAttachmentV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ccnId := idSplit[0]
	instanceType := idSplit[1]
	instanceRegion := idSplit[2]
	instanceId := idSplit[3]

	respData, err := service.DescribeCcnAttachedInstanceByFilter(ctx, ccnId, instanceType, instanceRegion, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ccn_attachment_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("instance_region", instanceRegion)
	_ = d.Set("instance_type", instanceType)

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.RouteTableId != nil {
		_ = d.Set("route_table_id", respData.RouteTableId)
	}

	if respData.CcnUin != nil {
		_ = d.Set("ccn_uin", respData.CcnUin)
	}

	if respData.State != nil {
		_ = d.Set("state", respData.State)
	}

	if respData.AttachedTime != nil {
		_ = d.Set("attached_time", respData.AttachedTime)
	}

	if respData.CidrBlock != nil {
		_ = d.Set("cidr_block", respData.CidrBlock)
	}

	// get route ids
	routeIds := make([]string, 0)
	request := vpc.NewDescribeCcnRoutesRequest()
	request.CcnId = &ccnId
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeCcnRoutes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result != nil && result.Response != nil && len(result.Response.RouteSet) > 0 {
			for _, route := range result.Response.RouteSet {
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

func resourceTencentCloudCcnAttachmentV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = vpc.NewModifyCcnAttachedInstancesAttributeRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ccnId := idSplit[0]
	instanceType := idSplit[1]
	instanceRegion := idSplit[2]
	instanceId := idSplit[3]

	if d.HasChange("description") {
		request.CcnId = &ccnId
		ccnInstance := new(vpc.CcnInstance)
		ccnInstance.InstanceType = &instanceType
		ccnInstance.InstanceRegion = &instanceRegion
		ccnInstance.InstanceId = &instanceId
		if v, ok := d.GetOk("description"); ok {
			ccnInstance.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("route_table_id"); ok {
			ccnInstance.RouteTableId = helper.String(v.(string))
		}

		request.Instances = append(request.Instances, ccnInstance)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyCcnAttachedInstancesAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify ccn instance attribute failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCcnAttachmentV2Read(d, meta)
}

func resourceTencentCloudCcnAttachmentV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_attachment_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = vpc.NewDetachCcnInstancesRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	ccnId := idSplit[0]
	instanceType := idSplit[1]
	instanceRegion := idSplit[2]
	instanceId := idSplit[3]

	request.CcnId = &ccnId
	ccnInstance := new(vpc.CcnInstance)
	ccnInstance.InstanceType = &instanceType
	ccnInstance.InstanceRegion = &instanceRegion
	ccnInstance.InstanceId = &instanceId
	request.Instances = append(request.Instances, ccnInstance)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DetachCcnInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s detach ccn instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
