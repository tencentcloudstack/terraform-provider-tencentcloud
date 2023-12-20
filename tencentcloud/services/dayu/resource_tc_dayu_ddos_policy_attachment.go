package dayu

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDayuDdosPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyAttachmentCreate,
		Read:   resourceTencentCloudDayuDdosPolicyAttachmentRead,
		Delete: resourceTencentCloudDayuDdosPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the attached resource.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the policy.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy works for. Valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	policyId := d.Get("policy_id").(string)
	resourceType := d.Get("resource_type").(string)
	dayuService := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	//only one bind relationship supported, check unbind status
	_, has, statusErr := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, "")
	if statusErr != nil {
		statusErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, "")
			if statusErr != nil {
				return tccommon.RetryError(statusErr)
			}
			return nil
		})
	}
	if statusErr != nil {
		return statusErr
	}
	if has {
		return fmt.Errorf("DDoS is already bined by policy")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.BindDdosPolicy(ctx, resourceId, resourceType, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	//check bind status
	_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if statusErr != nil {
		statusErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if statusErr != nil {
				return tccommon.RetryError(statusErr)
			}
			return nil
		})
	}
	if statusErr != nil {
		return statusErr
	}
	if !has {
		return fmt.Errorf("Create DDoS policy attachment faild")
	}

	d.SetId(resourceId + tccommon.FILED_SP + resourceType + tccommon.FILED_SP + policyId)

	return resourceTencentCloudDayuDdosPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of DDoS policy attachment")
	}
	resourceId := items[0]
	resourceType := items[1]
	policyId := items[2]

	dayuService := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	_, has, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		log.Printf("[CRITAL]%s read DDoS Policy attachment failed, reason:%s\n", logId, err)
		return err
	}

	if has {
		_ = d.Set("resource_id", resourceId)
		_ = d.Set("resource_type", resourceType)
		_ = d.Set("policy_id", policyId)
		return nil
	} else {
		d.SetId("")
		return nil
	}
}

func resourceTencentCloudDayuDdosPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of DDoS policy attachment")
	}
	resourceId := items[0]
	resourceType := items[1]
	policyId := items[2]

	dayuService := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.UnbindDdosPolicy(ctx, resourceId, resourceType, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	time.Sleep(1 * time.Minute)
	//check bind status
	_, has, statusErr := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if statusErr != nil {
		statusErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if statusErr != nil {
				return tccommon.RetryError(statusErr)
			}
			return nil
		})
	}
	if statusErr != nil {
		return statusErr
	}
	if has {
		return fmt.Errorf("Delete DDoS policy attachment faild")
	}

	return nil
}
