package dayuv2

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosSchedulingDomainUserName() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosSchedulingDomainUserNameCreate,
		Read:   resourceTencentCloudAntiddosSchedulingDomainUserNameRead,
		Update: resourceTencentCloudAntiddosSchedulingDomainUserNameUpdate,
		Delete: resourceTencentCloudAntiddosSchedulingDomainUserNameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "user cname.",
			},

			"domain_user_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "domain name.",
			},
		},
	}
}

func resourceTencentCloudAntiddosSchedulingDomainUserNameCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_scheduling_domain_user_name.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(d.Get("domain_name").(string))

	return resourceTencentCloudAntiddosSchedulingDomainUserNameUpdate(d, meta)
}

func resourceTencentCloudAntiddosSchedulingDomainUserNameRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_scheduling_domain_user_name.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	domainName := d.Id()

	schedulingDomainUserName, err := service.DescribeAntiddosSchedulingDomainUserNameById(ctx, domainName)
	if err != nil {
		return err
	}

	if schedulingDomainUserName == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosSchedulingDomainUserName` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if schedulingDomainUserName.Domain != nil {
		_ = d.Set("domain_name", schedulingDomainUserName.Domain)
	}

	if schedulingDomainUserName.UsrDomainName != nil {
		_ = d.Set("domain_user_name", schedulingDomainUserName.UsrDomainName)
	}

	return nil
}

func resourceTencentCloudAntiddosSchedulingDomainUserNameUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_scheduling_domain_user_name.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewModifyDomainUsrNameRequest()

	domainName := d.Id()

	request.DomainName = &domainName

	if d.HasChange("domain_user_name") {
		if v, ok := d.GetOk("domain_user_name"); ok {
			request.DomainUserName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().ModifyDomainUsrName(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update antiddos schedulingDomainUserName failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudAntiddosSchedulingDomainUserNameRead(d, meta)
}

func resourceTencentCloudAntiddosSchedulingDomainUserNameDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_scheduling_domain_user_name.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
