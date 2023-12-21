package dnspod

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func ResourceTencentCloudDnspodDomainInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDomainInstanceCreate,
		Read:   resourceTencentCloudDnspodDomainInstanceRead,
		Update: resourceTencentCloudDnspodDomainInstanceUpdate,
		Delete: resourceTencentCloudDnspodDomainInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Domain.",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The Group Id of Domain.",
			},
			"is_mark": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DNSPOD_DOMAIN_MARK_TYPE),
				Description:  "Whether to Mark the Domain.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DNSPOD_DOMAIN_STATUS_TYPE),
				Description:  "The status of Domain.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remark of Domain.",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the domain.",
			},
		},
	}
}

func resourceTencentCloudDnspodDomainInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := dnspod.NewCreateDomainRequest()
	var (
		domain   string
		groupId  uint64
		isMark   string
		response *dnspod.CreateDomainResponse
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(uint64)
	}
	if v, ok := d.GetOk("is_mark"); ok {
		isMark = v.(string)
	}

	request.Domain = &domain
	request.GroupId = &groupId
	request.IsMark = &isMark

	result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().CreateDomain(request)
	response = result

	if err != nil {
		log.Printf("[CRITAL]%s create DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(*response.Response.DomainInfo.Domain)

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if v, ok := d.GetOk("status"); ok {
		domainId := response.Response.DomainInfo.Domain
		status := v.(string)
		err := service.ModifyDnsPodDomainStatus(ctx, *domainId, status)
		if err != nil {
			log.Printf("[CRITAL]%s set DnsPod Domain status failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if v, ok := d.GetOk("remark"); ok {
		domainId := response.Response.DomainInfo.Domain
		remark := v.(string)
		err := service.ModifyDnsPodDomainRemark(ctx, *domainId, remark)
		if err != nil {
			log.Printf("[CRITAL]%s set DnsPod Domain remark failed, reason:%s\n", logId, err.Error())
			return err
		}
	}
	return resourceTencentCloudDnspodDomainInstanceRead(d, meta)
}

func resourceTencentCloudDnspodDomainInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()

	request := dnspod.NewDescribeDomainRequest()
	request.Domain = helper.String(id)

	var response *dnspod.DescribeDomainResponse

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeDomain(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		response = result
		info := response.Response.DomainInfo

		d.SetId(*response.Response.DomainInfo.Domain)

		_ = d.Set("domain", info.Domain)
		_ = d.Set("create_time", info.CreatedOn)
		_ = d.Set("is_mark", info.IsMark)

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudDnspodDomainInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		err := service.ModifyDnsPodDomainStatus(ctx, id, status)
		if err != nil {
			log.Printf("[CRITAL]%s modify DnsPod Domain status failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("remark") {
		remark := d.Get("remark").(string)
		err := service.ModifyDnsPodDomainRemark(ctx, id, remark)
		if err != nil {
			log.Printf("[CRITAL]%s modify DnsPod Domain remark failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	d.Partial(false)
	return resourceTencentCloudDnspodDomainInstanceRead(d, meta)
}

func resourceTencentCloudDnspodDomainInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dnspod_domain.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dnspod.NewDeleteDomainRequest()
	request.Domain = helper.String(d.Id())

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DeleteDomain(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
