package dnspod

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDnspodLineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodLineGroupCreate,
		Read:   resourceTencentCloudDnspodLineGroupRead,
		Update: resourceTencentCloudDnspodLineGroupUpdate,
		Delete: resourceTencentCloudDnspodLineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Domain name.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Line group name, length 1-17 characters.",
			},

			"lines": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of lines in the group. Maximum 120 lines.",
			},

			"domain_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID.",
			},

			"line_group_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Line group ID.",
			},

			"created_on": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"updated_on": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudDnspodLineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_line_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = dnspod.NewCreateLineGroupRequest()
		response    = dnspod.NewCreateLineGroupResponse()
		lineGroupId uint64
		domain      string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lines"); ok {
		lines := v.([]interface{})
		lineStrs := make([]string, 0, len(lines))
		for _, line := range lines {
			lineStrs = append(lineStrs, line.(string))
		}
		request.Lines = helper.String(strings.Join(lineStrs, ","))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().CreateLineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod line_group failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Data != nil && response.Response.Data.Id != nil {
		lineGroupId = *response.Response.Data.Id
	}

	d.SetId(strings.Join([]string{domain, helper.UInt64ToStr(lineGroupId)}, tccommon.FILED_SP))

	return resourceTencentCloudDnspodLineGroupRead(d, meta)
}

func resourceTencentCloudDnspodLineGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_line_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_line_group id is broken, id is %s", d.Id())
	}
	domain := idSplit[0]
	lineGroupId := helper.StrToUInt64(idSplit[1])

	lineGroup, err := service.DescribeDnspodLineGroupById(ctx, domain, lineGroupId)
	if err != nil {
		return err
	}

	if lineGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodLineGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if lineGroup.DomainId != nil {
		_ = d.Set("domain_id", lineGroup.DomainId)
	}

	if lineGroup.Name != nil {
		_ = d.Set("name", lineGroup.Name)
	}

	if lineGroup.Lines != nil && len(lineGroup.Lines) > 0 {
		_ = d.Set("lines", lineGroup.Lines)
	}

	if lineGroup.Id != nil {
		_ = d.Set("line_group_id", lineGroup.Id)
	}

	if lineGroup.CreatedOn != nil {
		_ = d.Set("created_on", lineGroup.CreatedOn)
	}

	if lineGroup.UpdatedOn != nil {
		_ = d.Set("updated_on", lineGroup.UpdatedOn)
	}

	return nil
}

func resourceTencentCloudDnspodLineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_line_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dnspod.NewModifyLineGroupRequest()
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_line_group id is broken, id is %s", d.Id())
	}
	request.Domain = helper.String(idSplit[0])
	request.LineGroupId = helper.StrToUint64Point(idSplit[1])

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if d.HasChange("lines") {
		if v, ok := d.GetOk("lines"); ok {
			lines := v.([]interface{})
			lineStrs := make([]string, 0, len(lines))
			for _, line := range lines {
				lineStrs = append(lineStrs, line.(string))
			}
			request.Lines = helper.String(strings.Join(lineStrs, ","))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().ModifyLineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod line_group failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDnspodLineGroupRead(d, meta)
}

func resourceTencentCloudDnspodLineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_line_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("tencentcloud_dnspod_line_group id is broken, id is %s", d.Id())
	}

	request := dnspod.NewDeleteLineGroupRequest()
	request.Domain = helper.String(idSplit[0])
	request.LineGroupId = helper.StrToUint64Point(idSplit[1])

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DeleteLineGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete dnspod line_group failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
