package dayu

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuDdosPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the attached resource to be queried.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"dayu_ddos_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of dayu DDoS policy attachments. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the attached resource.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the DDoS policy works for.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_ddos_policy_attachments.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	policyId := d.Get("policy_id").(string)

	dayuService := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	attachments, _, err := dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			attachments, _, err = dayuService.DescribeDdosPolicyAttachments(ctx, resourceId, resourceType, policyId)
			if err != nil {
				return tccommon.RetryError(err, "ClientError.NetworkError")
			}
			return nil
		})
	}
	if err != nil {
		log.Printf("[CRITAL]%s read DDoS Policy attachment failed, reason:%s\n", logId, err)
		return err
	}

	ddosPolicyAttachmentList := make([]map[string]interface{}, 0, len(attachments))
	ids := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		ddosPolicyAttachmentList = append(ddosPolicyAttachmentList, attachment)
		ids = append(ids, resourceId+tccommon.FILED_SP+resourceType+""+tccommon.FILED_SP+policyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("dayu_ddos_policy_attachment_list", ddosPolicyAttachmentList); e != nil {
		log.Printf("[CRITAL]%s provider set DDoS policy attachment list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), ddosPolicyAttachmentList); e != nil {
			return e
		}
	}

	return nil
}
