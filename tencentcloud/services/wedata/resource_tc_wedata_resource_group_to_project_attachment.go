package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataResourceGroupToProjectAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataResourceGroupToProjectAttachmentCreate,
		Read:   resourceTencentCloudWedataResourceGroupToProjectAttachmentRead,
		Delete: resourceTencentCloudWedataResourceGroupToProjectAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource group ID.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID.",
			},
		},
	}
}

func resourceTencentCloudWedataResourceGroupToProjectAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group_to_project_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = wedatav20250806.NewAssociateResourceGroupToProjectRequest()
		resourceGroupId string
		projectId       string
	)

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = helper.String(v.(string))
		resourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AssociateResourceGroupToProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata resource group to project attachment failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Create wedata resource group to project attachment failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata resource group to project attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{resourceGroupId, projectId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataResourceGroupToProjectAttachmentRead(d, meta)
}

func resourceTencentCloudWedataResourceGroupToProjectAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group_to_project_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	resourceGroupId := idSplit[0]
	projectId := idSplit[1]

	respData, err := service.DescribeWedataResourceGroupToProjectAttachmentById(ctx, resourceGroupId, projectId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_resource_group_to_project_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	_ = d.Set("project_id", respData.ProjectId)
	_ = d.Set("resource_group_id", resourceGroupId)

	return nil
}

func resourceTencentCloudWedataResourceGroupToProjectAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group_to_project_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDissociateResourceGroupFromProjectRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	resourceGroupId := idSplit[0]
	projectId := idSplit[1]

	request.ResourceGroupId = &resourceGroupId
	request.ProjectId = &projectId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DissociateResourceGroupFromProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata resource group to project attachment failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata resource group to project attachment failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata resource group to project attachment failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
