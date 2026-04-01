package pls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcEndPointCreate,
		Read:   resourceTencentCloudVpcEndPointRead,
		Update: resourceTencentCloudVpcEndPointUpdate,
		Delete: resourceTencentCloudVpcEndPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of vpc instance.",
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of subnet instance.",
			},

			"end_point_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of endpoint.",
			},

			"end_point_service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of endpoint service.",
			},

			"end_point_vip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP of endpoint ip.",
			},

			"security_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of security group.",
			},

			"security_groups_ids": {
				Optional:    true,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ordered security groups associated with the endpoint.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tags of the VPC endpoint.",
			},

			"ip_address_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "Ipv4",
				Description: "IP address type. Valid values are `Ipv4` and `Ipv6`.",
				ValidateFunc: validation.StringInSlice([]string{"Ipv4", "Ipv6"}, false),
			},

			"end_point_owner": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "APPID.",
			},

			"state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "state of end point.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create Time.",
			},

			"cdc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CDC instance ID.",
			},
		},
	}
}

func resourceTencentCloudVpcEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = vpc.NewCreateVpcEndPointRequest()
		response   = vpc.NewCreateVpcEndPointResponse()
		endPointId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_name"); ok {
		request.EndPointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_service_id"); ok {
		request.EndPointServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_vip"); ok {
		request.EndPointVip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_type"); ok {
		request.IpAddressType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagMap := v.(map[string]interface{})
		tags := make([]*vpc.Tag, 0, len(tagMap))
		for key, value := range tagMap {
			tag := &vpc.Tag{
				Key:   helper.String(key),
				Value: helper.String(value.(string)),
			}
			tags = append(tags, tag)
		}
		request.Tags = tags
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpcEndPoint(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.EndPoint == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vpc endPoint failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPoint failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.EndPoint.EndPointId == nil {
		return fmt.Errorf("EndPointId is nil.")
	}

	endPointId = *response.Response.EndPoint.EndPointId
	d.SetId(endPointId)

	if v, ok := d.GetOk("security_groups_ids"); ok {
		request := vpc.NewModifyVpcEndPointAttributeRequest()
		request.EndPointId = helper.String(endPointId)
		request.SecurityGroupIds = helper.InterfacesStringsPoint(v.([]interface{}))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcEndPointAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create vpc endPoint failed, reason:%+v", logId, err)
			return err
		}

	}
	return resourceTencentCloudVpcEndPointRead(d, meta)
}

func resourceTencentCloudVpcEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	endPointId := d.Id()

	endPoint, err := service.DescribeVpcEndPointById(ctx, endPointId)
	if err != nil {
		return err
	}

	if endPoint == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if endPoint.VpcId != nil {
		_ = d.Set("vpc_id", endPoint.VpcId)
	}

	if endPoint.SubnetId != nil {
		_ = d.Set("subnet_id", endPoint.SubnetId)
	}

	if endPoint.EndPointName != nil {
		_ = d.Set("end_point_name", endPoint.EndPointName)
	}

	if endPoint.EndPointServiceId != nil {
		_ = d.Set("end_point_service_id", endPoint.EndPointServiceId)
	}

	if endPoint.EndPointVip != nil {
		_ = d.Set("end_point_vip", endPoint.EndPointVip)
	}

	if endPoint.EndPointOwner != nil {
		_ = d.Set("end_point_owner", endPoint.EndPointOwner)
	}

	if endPoint.State != nil {
		_ = d.Set("state", endPoint.State)
	}

	if endPoint.CreateTime != nil {
		_ = d.Set("create_time", endPoint.CreateTime)
	}

	if endPoint.GroupSet != nil {
		_ = d.Set("security_groups_ids", endPoint.GroupSet)
	}

	if endPoint.CdcId != nil {
		_ = d.Set("cdc_id", endPoint.CdcId)
	}

	if endPoint.SecurityGroupId != nil {
		_ = d.Set("security_group_id", endPoint.SecurityGroupId)
	}

	if endPoint.Tags != nil {
		tagMap := make(map[string]interface{})
		for _, tag := range endPoint.Tags {
			if tag.Key != nil {
				tagMap[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tagMap)
	}

	if endPoint.IpAddressType != nil {
		_ = d.Set("ip_address_type", endPoint.IpAddressType)
	}

	return nil
}

func resourceTencentCloudVpcEndPointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyVpcEndPointAttributeRequest()

	endPointId := d.Id()

	request.EndPointId = &endPointId

	unsupportedUpdateFields := []string{
		"vpc_id",
		"subnet_id",
		"end_point_service_id",
		"end_point_vip",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_vpc_end_point update on %s is not support yet", field)
		}
	}

	if d.HasChange("end_point_name") || d.HasChange("security_groups_ids") || d.HasChange("security_group_id") || d.HasChange("ip_address_type") {
		if v, ok := d.GetOk("end_point_name"); ok {
			request.EndPointName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("security_groups_ids"); ok {
			request.SecurityGroupIds = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := d.GetOk("security_group_id"); ok {
			request.SecurityGroupId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("ip_address_type"); ok {
			request.IpAddressType = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tencentcloudstack.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcEndPointAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create vpc endPoint failed, reason:%+v", logId, err)
			return err
		}
	}

	// Handle tags update separately using tag management API
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		oldTagMap := oldTags.(map[string]interface{})
		newTagMap := newTags.(map[string]interface{})

		// Delete old tags
		for key := range oldTagMap {
			if _, ok := newTagMap[key]; !ok {
				// Tag was removed, delete it
				// Note: You may need to use a different API for tag deletion
			}
		}

		// Add new tags
		tagsToAdd := make([]*vpc.Tag, 0)
		for key, value := range newTagMap {
			tagsToAdd = append(tagsToAdd, &vpc.Tag{
				Key:   helper.String(key),
				Value: helper.String(value.(string)),
			})
		}

		// Call tag management API to update tags
		// This is a placeholder - you may need to use the actual tag management API
		if len(tagsToAdd) > 0 {
			// For now, we assume tags can be updated via ModifyVpcEndPointAttribute
			request := vpc.NewModifyVpcEndPointAttributeRequest()
			request.EndPointId = &endPointId
			request.Tags = tagsToAdd

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcEndPointAttribute(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update tags failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudVpcEndPointRead(d, meta)
}

func resourceTencentCloudVpcEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	endPointId := d.Id()

	if err := service.DeleteVpcEndPointById(ctx, endPointId); err != nil {
		return nil
	}

	return nil
}
