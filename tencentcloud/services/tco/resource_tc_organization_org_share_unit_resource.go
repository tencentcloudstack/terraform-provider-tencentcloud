package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationOrgShareUnitResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgShareUnitResourceCreate,
		Read:   resourceTencentCloudOrganizationOrgShareUnitResourceRead,
		Delete: resourceTencentCloudOrganizationOrgShareUnitResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Shared unit ID.",
			},

			"area": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Shared unit area.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Shared resource type.",
			},

			"product_resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product Resource ID.",
			},

			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Shared resource ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},

			"shared_member_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of shared unit members.",
			},

			"shared_member_use_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of shared unit members in use.",
			},

			"share_manager_uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Sharing administrator OwnerUin.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgShareUnitResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		unitId            string
		area              string
		shareResourceType string
		productResourceId string
	)
	var (
		request  = organization.NewAddShareUnitResourcesRequest()
		response = organization.NewAddShareUnitResourcesResponse()
	)

	if v, ok := d.GetOk("unit_id"); ok {
		unitId = v.(string)
		request.UnitId = helper.String(unitId)
	}

	if v, ok := d.GetOk("area"); ok {
		area = v.(string)
		request.Area = helper.String(area)
	}

	if v, ok := d.GetOk("type"); ok {
		shareResourceType = v.(string)
		request.Type = helper.String(shareResourceType)
	}

	productResource := organization.ProductResource{}
	if v, ok := d.GetOk("product_resource_id"); ok {
		productResourceId = v.(string)
		productResource.ProductResourceId = helper.String(productResourceId)
	}
	request.Resources = []*organization.ProductResource{&productResource}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddShareUnitResourcesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization org share unit resource failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(strings.Join([]string{unitId, area, shareResourceType, productResourceId}, tccommon.FILED_SP))

	return resourceTencentCloudOrganizationOrgShareUnitResourceRead(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	unitId := idSplit[0]
	area := idSplit[1]
	shareResourceType := idSplit[2]
	productResourceId := idSplit[3]

	_ = d.Set("unit_id", unitId)
	_ = d.Set("area", area)
	_ = d.Set("type", shareResourceType)
	_ = d.Set("product_resource_id", productResourceId)

	respData, err := service.DescribeOrganizationOrgShareUnitResourceById(ctx, unitId, area, shareResourceType, productResourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `organization_org_share_unit_resource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.ResourceId != nil {
		_ = d.Set("resource_id", respData.ResourceId)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ProductResourceId != nil {
		_ = d.Set("product_resource_id", respData.ProductResourceId)
	}

	if respData.SharedMemberNum != nil {
		_ = d.Set("shared_member_num", respData.SharedMemberNum)
	}

	if respData.SharedMemberUseNum != nil {
		_ = d.Set("shared_member_use_num", respData.SharedMemberUseNum)
	}

	if respData.ShareManagerUin != nil {
		_ = d.Set("share_manager_uin", respData.ShareManagerUin)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	unitId := idSplit[0]
	area := idSplit[1]
	shareResourceType := idSplit[2]
	productResourceId := idSplit[3]

	var (
		request  = organization.NewDeleteShareUnitResourcesRequest()
		response = organization.NewDeleteShareUnitResourcesResponse()
	)

	request.UnitId = helper.String(unitId)

	request.Area = helper.String(area)

	request.Type = helper.String(shareResourceType)

	shareResource := organization.ShareResource{}
	shareResource.ProductResourceId = helper.String(productResourceId)
	request.Resources = []*organization.ShareResource{&shareResource}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteShareUnitResourcesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete organization org share unit resource failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
