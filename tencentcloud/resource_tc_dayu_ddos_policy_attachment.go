/*
Provides a resource to create a dayu DDoS policy attachment.

Example Usage

```hcl
resource "tencentcloud_dayu_ddos_policy_attachment" "dayu_ddos_policy_attachment_basic" {
  resource_type = tencentcloud_dayu_ddos_policy.test_policy.resource_type
  resource_id = "bgpip-00000294"
  policy_id = tencentcloud_dayu_ddos_policy.test_policy.policy_id
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudDayuDdosPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyAttachmentCreate,
		Read:   resourceTencentCloudDayuDdosPolicyAttachmentRead,
		Delete: resourceTencentCloudDayuDdosPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the attached resource.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the policy.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip`, `net`.",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	resourceId := d.Get("resource_id").(string)
	policyId := d.Get("policy_id").(string)
	resourceType := d.Get("resource_type").(string)
	dayuService := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.BindDdosPolicy(ctx, resourceId, resourceType, policyId)
		if e != nil {
			return retryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	//check bind status
	_, has, statusErr := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if statusErr != nil {
		statusErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if statusErr != nil {
				return retryError(statusErr)
			}
			return nil
		})
	}
	if statusErr != nil {
		return statusErr
	}
	if !has {
		return fmt.Errorf("Create DDos policy attachment faild")
	}

	d.SetId(resourceId + FILED_SP + resourceType + FILED_SP + policyId)

	return resourceTencentCloudDayuDdosPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), "#")
	if len(items) < 3 {
		return fmt.Errorf("broken ID of DDos policy attachment")
	}
	resourceId := items[0]
	resourceType := items[1]
	policyId := items[2]

	dayuService := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	_, has, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if err != nil {
				return retryError(err)
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
	defer logElapsed("resource.tencentcloud_dayu_ddos_policy_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of DDos policy attachment")
	}
	resourceId := items[0]
	resourceType := items[1]
	policyId := items[2]

	dayuService := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := dayuService.UnbindDdosPolicy(ctx, resourceId, resourceType, policyId)
		if e != nil {
			return retryError(e)
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
		statusErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, statusErr = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if statusErr != nil {
				return retryError(statusErr)
			}
			return nil
		})
	}
	if statusErr != nil {
		return statusErr
	}
	if has {
		return fmt.Errorf("Delete DDos policy attachment faild")
	}

	return nil
}
