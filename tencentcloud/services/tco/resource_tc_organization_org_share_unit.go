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

func ResourceTencentCloudOrganizationOrgShareUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgShareUnitCreate,
		Read:   resourceTencentCloudOrganizationOrgShareUnitRead,
		Update: resourceTencentCloudOrganizationOrgShareUnitUpdate,
		Delete: resourceTencentCloudOrganizationOrgShareUnitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Shared unit name. It only supports a combination of uppercase and lowercase letters, numbers, -, and _, with a length of 3-128 characters.",
			},

			"area": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Shared unit region. The regions that support sharing can be obtained through the DescribeShareAreas interface.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Shared unit description. Up to 128 characters.",
			},

			"unit_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Shared unit region. The regions that support sharing can be obtained through the DescribeShareAreas interface.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgShareUnitCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = organization.NewAddShareUnitRequest()
		response = organization.NewAddShareUnitResponse()
		unitId   string
		area     string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
		area = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddShareUnit(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgShareUnit failed, reason:%+v", logId, err)
		return err
	}

	unitId = *response.Response.UnitId
	d.SetId(area + tccommon.FILED_SP + unitId)

	return resourceTencentCloudOrganizationOrgShareUnitRead(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	split := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(split) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	area := split[0]
	unitId := split[1]
	orgShareUnit, err := service.DescribeOrganizationOrgShareUnitById(ctx, area, unitId)
	if err != nil {
		return err
	}

	if orgShareUnit == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgShareUnit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgShareUnit.Name != nil {
		_ = d.Set("name", orgShareUnit.Name)
	}

	if orgShareUnit.Area != nil {
		_ = d.Set("area", orgShareUnit.Area)
	}

	if orgShareUnit.UnitId != nil {
		_ = d.Set("unit_id", orgShareUnit.UnitId)
	}

	if orgShareUnit.Description != nil {
		_ = d.Set("description", orgShareUnit.Description)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := organization.NewUpdateShareUnitRequest()

	split := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(split) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	unitId := split[1]
	request.UnitId = &unitId
	immutableArgs := []string{"area", "unit_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"name", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateShareUnit(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update organization orgShareUnit failed, reason:%+v", logId, err)
			return err
		}

	}
	return resourceTencentCloudOrganizationOrgShareUnitRead(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	split := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(split) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	unitId := split[1]
	if err := service.DeleteOrganizationOrgShareUnitById(ctx, unitId); err != nil {
		return err
	}

	return nil
}
