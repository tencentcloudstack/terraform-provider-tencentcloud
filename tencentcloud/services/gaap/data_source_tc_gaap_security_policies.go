package gaap

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudGaapSecurityPolices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapSecurityPoliciesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the security policy to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"proxy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the GAAP proxy.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the security policy.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default policy.",
			},
		},
	}
}

func dataSourceTencentCloudGaapSecurityPoliciesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_security_policies.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Get("id").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	proxyId, status, action, exist, err := service.DescribeSecurityPolicy(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		d.SetId("")
		return nil
	}

	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("status", status)
	_ = d.Set("action", action)

	d.SetId(id)

	m = map[string]interface{}{
		"id":       id,
		"proxy_id": proxyId,
		"status":   status,
		"action":   action,
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), m); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
