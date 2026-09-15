package bdrc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBdrcSecurityGroupMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBdrcSecurityGroupMappingCreate,
		Read:   resourceTencentCloudBdrcSecurityGroupMappingRead,
		Update: resourceTencentCloudBdrcSecurityGroupMappingUpdate,
		Delete: resourceTencentCloudBdrcSecurityGroupMappingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"site_pair_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Security group mapping belongs to the site pair ID.",
			},

			"src_security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Production end instance bound to the security group ID.",
			},

			"target_security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Disaster recovery end instance bound to the security group ID.",
			},

			"security_group_mapping_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group mapping ID.",
			},

			"source_security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Production end security group ID.",
			},

			"life_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Life state of the security group mapping; NORMAL: normal.",
			},
		},
	}
}

func resourceTencentCloudBdrcSecurityGroupMappingCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_security_group_mapping.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewCreateSecurityGroupMappingRequest()
	)

	if v, ok := d.GetOk("site_pair_id"); ok {
		request.SitePairId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("src_security_group_id"); ok {
		request.SrcSecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_security_group_id"); ok {
		request.TargetSecurityGroupId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().CreateSecurityGroupMappingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bdrc security_group_mapping failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create bdrc security_group_mapping failed, reason:%+v", logId, err)
		return err
	}

	service := BdrcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	sitePairId := d.Get("site_pair_id").(string)
	srcSecurityGroupId := d.Get("src_security_group_id").(string)
	targetSecurityGroupId := d.Get("target_security_group_id").(string)

	respData, err := service.DescribeSecurityGroupMappingByFilter(ctx, sitePairId, srcSecurityGroupId, targetSecurityGroupId)
	if err != nil {
		log.Printf("[CRITAL]%s Describe bdrc security_group_mapping by filter failed after create, reason:%+v", logId, err)
		return err
	}

	if respData == nil || respData.SecurityGroupMappingId == nil {
		log.Printf("[CRITAL]%s bdrc security_group_mapping not found after create, logId=%s, d.Id()=%s", logId, logId, d.Id())
		return fmt.Errorf("bdrc security_group_mapping not found after create")
	}

	securityGroupMappingId := *respData.SecurityGroupMappingId
	d.SetId(strings.Join([]string{sitePairId, securityGroupMappingId}, tccommon.FILED_SP))
	return resourceTencentCloudBdrcSecurityGroupMappingRead(d, meta)
}

func resourceTencentCloudBdrcSecurityGroupMappingRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_security_group_mapping.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BdrcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sitePairId := idSplit[0]
	securityGroupMappingId := idSplit[1]

	respData, err := service.DescribeSecurityGroupMappingById(ctx, sitePairId, securityGroupMappingId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] bdrc security_group_mapping id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.SitePairId != nil {
		_ = d.Set("site_pair_id", respData.SitePairId)
	}

	if respData.SourceSecurityGroupId != nil {
		_ = d.Set("source_security_group_id", respData.SourceSecurityGroupId)
	}

	if respData.TargetSecurityGroupId != nil {
		_ = d.Set("target_security_group_id", respData.TargetSecurityGroupId)
	}

	if respData.SecurityGroupMappingId != nil {
		_ = d.Set("security_group_mapping_id", respData.SecurityGroupMappingId)
	}

	if respData.LifeState != nil {
		_ = d.Set("life_state", respData.LifeState)
	}

	return nil
}

func resourceTencentCloudBdrcSecurityGroupMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_security_group_mapping.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"src_security_group_id", "target_security_group_id", "site_pair_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed, please delete and recreate", v)
		}
	}

	return resourceTencentCloudBdrcSecurityGroupMappingRead(d, meta)
}

func resourceTencentCloudBdrcSecurityGroupMappingDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_security_group_mapping.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewDeleteSecurityGroupMappingRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sitePairId := idSplit[0]
	securityGroupMappingId := idSplit[1]

	request.SitePairId = helper.String(sitePairId)
	request.SecurityGroupMappingIds = []*string{helper.String(securityGroupMappingId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DeleteSecurityGroupMappingWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete bdrc security_group_mapping failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete bdrc security_group_mapping failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
