package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoBindSecurityTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoBindSecurityTemplateCreate,
		Read:   resourceTencentCloudTeoBindSecurityTemplateRead,
		Delete: resourceTencentCloudTeoBindSecurityTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID of the policy template to be bound to or unbound from.",
			},

			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "List of domain names to bind to/unbind from a policy template.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the policy template or the site global policy to be bound or unbound.\n<li>To bind to a policy template, or unbind from it, specify the policy template ID.</li>.\n<li>To bind to the site's global policy, or unbind from it, use the @ZoneLevel@domain parameter value.</li>.\n\nNote: After unbinding, the domain name will use an independent policy and rule quota will be calculated separately. Please make sure there is sufficient rule quota before unbinding.",
			},

			"operate": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Unbind operation option. valid values: `unbind-keep-policy`: unbind a domain name from the policy template while retaining the current policy. `unbind-use-default`: unbind a domain name from the policy template and use the default blank policy. default value: `unbind-keep-policy`.",
			},

			"over_write": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "If the passed-in domain is already bound to a policy template (including site-level protection policies), setting this parameter indicates whether to replace that template. The default value is true. Supported values are: `true`: Replace the currently bound template for the domain. `false`: Do not replace the currently bound template for the domain. Note: When set to false, if the passed-in domain is already bound to a policy template, the API will return an error; site-level protection policies are also a type of policy template.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance configuration delivery status, the possible values are: `online`: the configuration has taken effect; `fail`: the configuration failed; `process`: the configuration is being delivered.",
			},
		},
	}
}

func resourceTencentCloudTeoBindSecurityTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	var (
		zoneId     string
		templateId string
		entity     string
	)

	request := teov20220901.NewBindSecurityTemplateToEntityRequest()

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
		request.Entities = append(request.Entities, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(string)
		request.TemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		request.OverWrite = helper.Bool(v.(bool))
	} else {
		request.OverWrite = helper.Bool(true)
	}

	request.Operate = helper.String("bind")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().BindSecurityTemplateToEntityWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo bind security template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceTeoBindSecurityTemplateCreateStateRefreshFunc_0_0(ctx, zoneId, templateId, entity),
		Target:     []string{"online"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, templateId, entity}, tccommon.FILED_SP))

	return resourceTencentCloudTeoBindSecurityTemplateRead(d, meta)
}

func resourceTencentCloudTeoBindSecurityTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	templateId := idSplit[1]
	entity := idSplit[2]

	_ = d.Set("zone_id", zoneId)

	_ = d.Set("template_id", templateId)

	_ = d.Set("entity", entity)

	respData, err := service.DescribeTeoBindSecurityTemplateById(ctx, zoneId, templateId, entity)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_bind_security_template` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if v, ok := d.GetOk("operate"); ok {
		_ = d.Set("operate", v.(string))
	} else {
		_ = d.Set("operate", "unbind-keep-policy")
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		_ = d.Set("over_write", v.(bool))
	} else {
		_ = d.Set("over_write", true)
	}

	return nil
}

func resourceTencentCloudTeoBindSecurityTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_bind_security_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	templateId := idSplit[1]
	entity := idSplit[2]

	request := teov20220901.NewBindSecurityTemplateToEntityRequest()
	request.ZoneId = &zoneId
	request.Entities = append(request.Entities, &entity)
	request.TemplateId = &templateId

	if v, ok := d.GetOk("operate"); ok {
		request.Operate = helper.String(v.(string))
	} else {
		request.Operate = helper.String("unbind-keep-policy")
	}

	if v, ok := d.GetOkExists("over_write"); ok {
		request.OverWrite = helper.Bool(v.(bool))
	} else {
		request.OverWrite = helper.Bool(true)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().BindSecurityTemplateToEntityWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo bind security template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
