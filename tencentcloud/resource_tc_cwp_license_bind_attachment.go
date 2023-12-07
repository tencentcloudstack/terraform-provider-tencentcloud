package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCwpLicenseBindAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCwpLicenseBindAttachmentCreate,
		Read:   resourceTencentCloudCwpLicenseBindAttachmentRead,
		Delete: resourceTencentCloudCwpLicenseBindAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resource ID.",
			},
			"license_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "License ID.",
			},
			"license_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(LICENSE_TYPE),
				Description:  "LicenseType, 0 CWP Pro - Pay as you go, 1 CWP Pro - Monthly subscription, 2 CWP Ultimate - Monthly subscription. Default is 0.",
			},
			"quuid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Machine quota that needs to be bound.",
			},
			"machine_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "machine name.",
			},
			"machine_wan_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "machine wan ip.",
			},
			"machine_ip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "machine ip.",
			},
			"uuid": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "uuid.",
			},
			"agent_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "agent status.",
			},
			"is_unbind": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Allow unbinding, false does not allow unbinding.",
			},
			"is_switch_bind": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Is it allowed to change the binding, false is not allowed to change the binding.",
			},
		},
	}
}

func resourceTencentCloudCwpLicenseBindAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_bind_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = cwp.NewModifyLicenseBindsRequest()
		response    = cwp.NewModifyLicenseBindsResponse()
		taskRequest = cwp.NewDescribeLicenseBindScheduleRequest()
		resourceId  string
		licenseId   string
		quuid       string
		licenseType string
	)

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
		resourceId = v.(string)
	}

	if v, ok := d.GetOkExists("license_id"); ok {
		licenseIdInt := v.(int)
		licenseId = strconv.Itoa(licenseIdInt)
	}

	if v, ok := d.GetOkExists("license_type"); ok {
		request.LicenseType = helper.IntUint64(v.(int))
		licenseTypeInt := v.(int)
		licenseType = strconv.Itoa(licenseTypeInt)
	}

	if v, ok := d.GetOk("quuid"); ok {
		quuid = v.(string)
		request.QuuidList = append(request.QuuidList, &quuid)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCwpClient().ModifyLicenseBinds(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("cwp licenseBindAttachment not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cwp licenseBindAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{resourceId, licenseId, quuid, licenseType}, FILED_SP))

	// wait
	taskRequest.TaskId = response.Response.TaskId
	err = resource.Retry(writeRetryTimeout*6, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCwpClient().DescribeLicenseBindSchedule(taskRequest)
		if e != nil {
			return retryError(e)
		}

		if result == nil {
			e = fmt.Errorf("license bind schedule failed")
			return resource.NonRetryableError(e)
		}

		if *result.Response.List[0].Status == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("license bind schedule is processing, status: %d", *result.Response.List[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cwp licenseBindAttachment failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCwpLicenseBindAttachmentRead(d, meta)
}

func resourceTencentCloudCwpLicenseBindAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_bind_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CwpService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]
	licenseId := idSplit[1]
	quuid := idSplit[2]
	licenseType := idSplit[3]

	licenseIdInt, _ := strconv.ParseUint(licenseId, 10, 64)
	licenseTypeInt, _ := strconv.ParseUint(licenseType, 10, 64)

	licenseBindAttachment, err := service.DescribeCwpLicenseBindAttachmentById(ctx, resourceId, quuid, licenseIdInt, licenseTypeInt)
	if err != nil {
		return err
	}

	if licenseBindAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CwpLicenseBindAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("resource_id", resourceId)
	_ = d.Set("license_id", licenseIdInt)
	_ = d.Set("license_type", licenseTypeInt)
	_ = d.Set("quuid", quuid)

	if licenseBindAttachment.MachineName != nil {
		_ = d.Set("machine_name", licenseBindAttachment.MachineName)
	}

	if licenseBindAttachment.MachineWanIp != nil {
		_ = d.Set("machine_wan_ip", licenseBindAttachment.MachineWanIp)
	}

	if licenseBindAttachment.MachineIp != nil {
		_ = d.Set("machine_ip", licenseBindAttachment.MachineIp)
	}

	if licenseBindAttachment.Uuid != nil {
		_ = d.Set("uuid", licenseBindAttachment.Uuid)
	}

	if licenseBindAttachment.AgentStatus != nil {
		_ = d.Set("agent_status", licenseBindAttachment.AgentStatus)
	}

	if licenseBindAttachment.IsUnBind != nil {
		_ = d.Set("is_unbind", licenseBindAttachment.IsUnBind)
	}

	if licenseBindAttachment.IsSwitchBind != nil {
		_ = d.Set("is_switch_bind", licenseBindAttachment.IsSwitchBind)
	}

	return nil
}

func resourceTencentCloudCwpLicenseBindAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cwp_license_bind_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CwpService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	resourceId := idSplit[0]
	quuid := idSplit[2]
	licenseType := idSplit[3]

	if err := service.DeleteCwpLicenseBindAttachmentById(ctx, resourceId, quuid, licenseType); err != nil {
		return err
	}

	return nil
}
