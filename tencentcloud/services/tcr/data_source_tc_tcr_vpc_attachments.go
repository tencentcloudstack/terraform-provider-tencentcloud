package tcr

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTCRVPCAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRVPCAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the instance to query.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of VPC to query.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of subnet to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values
			"vpc_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated TCR namespaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of this VPC access.",
						},
						"access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address of this VPC access.",
						},
						"enable_public_domain_dns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable public domain dns.",
						},
						"enable_vpc_domain_dns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable vpc domain dns.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRVPCAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_vpc_attachments.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var vpcId, subnetId, instanceId string
	instanceId = d.Get("instance_id").(string)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error
	vpcAccesses, outErr := tcrService.DescribeTCRVPCAttachments(ctx, instanceId, vpcId, subnetId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			vpcAccesses, outErr = tcrService.DescribeTCRVPCAttachments(ctx, instanceId, vpcId, subnetId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(vpcAccesses))
	vpcAccessList := make([]map[string]interface{}, 0, len(vpcAccesses))
	for _, vpcAccess := range vpcAccesses {
		mapping := map[string]interface{}{
			"vpc_id":    vpcAccess.VpcId,
			"subnet_id": vpcAccess.SubnetId,
			"status":    vpcAccess.Status,
			"access_ip": vpcAccess.AccessIp,
		}
		if *vpcAccess.AccessIp != "" {
			publicDomainDnsStatus, err := GetDnsStatus(ctx, tcrService, instanceId, *vpcAccess.VpcId, *vpcAccess.AccessIp, true)
			if err != nil {
				return err
			}
			mapping["enable_public_domain_dns"] = *publicDomainDnsStatus.Status == TCR_VPC_DNS_STATUS_ENABLED

			vpcDomainDnsStatus, err := GetDnsStatus(ctx, tcrService, instanceId, *vpcAccess.VpcId, *vpcAccess.AccessIp, false)
			if err != nil {
				return err
			}
			mapping["enable_vpc_domain_dns"] = *vpcDomainDnsStatus.Status == TCR_VPC_DNS_STATUS_ENABLED
		}

		vpcAccessList = append(vpcAccessList, mapping)
		ids = append(ids, instanceId+tccommon.FILED_SP+*vpcAccess.VpcId+tccommon.FILED_SP+*vpcAccess.SubnetId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("vpc_attachment_list", vpcAccessList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR VPC access list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), vpcAccessList); e != nil {
			return e
		}
	}

	return nil

}
