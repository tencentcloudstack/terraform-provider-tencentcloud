package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapDomainErrorPageInfo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapDomainErrorPageInfoCreate,
		Read:   resourceTencentCloudGaapDomainErrorPageInfoRead,
		Delete: resourceTencentCloudGaapDomainErrorPageInfoDelete,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the layer7 listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HTTP domain.",
			},
			"error_codes": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
				Description: "Original error codes.",
			},
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "New response body.",
			},
			"new_error_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "New error code.",
			},
			"clear_headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Response headers to be removed.",
			},
			"set_headers": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Response headers to be set.",
			},
		},
	}
}

func resourceTencentCloudGaapDomainErrorPageInfoCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_domain_error_page.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)
	body := d.Get("body").(string)

	var errorCodes []int
	for _, errorNo := range d.Get("error_codes").(*schema.Set).List() {
		errorCodes = append(errorCodes, errorNo.(int))
	}

	var (
		newErrorCode *int64
		clearHeaders []string
	)

	if raw, ok := d.GetOk("new_error_code"); ok {
		newErrorCode = helper.IntInt64(raw.(int))
	}

	if raw, ok := d.GetOk("clear_headers"); ok {
		for _, clearHeader := range raw.(*schema.Set).List() {
			clearHeaders = append(clearHeaders, clearHeader.(string))
		}
	}

	setHeaders := helper.GetTags(d, "set_headers")

	id, err := service.CreateDomainErrorPageInfo(ctx, listenerId, domain, body, newErrorCode, errorCodes, clearHeaders, setHeaders)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapDomainErrorPageInfoRead(d, m)
}

func resourceTencentCloudGaapDomainErrorPageInfoRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_domain_error_page.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	info, err := service.DescribeDomainErrorPageInfo(ctx, listenerId, domain, id)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("listener_id", info.ListenerId)
	_ = d.Set("domain", info.Domain)
	_ = d.Set("error_codes", info.ErrorNos)
	_ = d.Set("body", info.Body)
	_ = d.Set("new_error_code", info.NewErrorNo)
	_ = d.Set("clear_headers", info.ClearHeaders)

	setHeaders := make(map[string]string, len(info.SetHeaders))
	for _, header := range info.SetHeaders {
		setHeaders[*header.HeaderName] = *header.HeaderValue
	}

	_ = d.Set("set_headers", setHeaders)

	return nil
}

func resourceTencentCloudGaapDomainErrorPageInfoDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_domain_error_page.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	return service.DeleteDomainErrorPageInfo(ctx, id)
}
