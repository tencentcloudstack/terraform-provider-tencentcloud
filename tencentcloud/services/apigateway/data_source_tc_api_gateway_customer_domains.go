package apigateway

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayCustomerDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayCustomerDomainRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			//Computed
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service custom domain name list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},
						"is_status_on": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Domain name resolution status. Valid values: `true`, `false`. `true` means normal parsing, `false` means parsing failed.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate ID.",
						},
						"is_default_mapping": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to use default path mapping. Valid values: `true`, `false`. `true` means to use default path mapping, `false` means to use custom path mapping.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom domain name agreement type.",
						},
						"net_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network type.",
						},
						"path_mappings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain name mapping path and environment list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The domain mapping path.",
									},
									"environment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release environment.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayCustomerDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_customer_domains.read")

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		infos             []*apigateway.DomainSetList
		list              []map[string]interface{}
		err               error
	)
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeServiceSubDomains(ctx, serviceId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var (
			pathMapping []map[string]interface{}
			status      bool
		)
		if !*info.IsDefaultMapping && *info.DomainName != "" {
			var mappings *apigateway.ServiceSubDomainMappings
			mappings, err = apiGatewayService.DescribeServiceSubDomainMappings(ctx, serviceId, *info.DomainName)
			if err != nil {
				return err
			}

			for _, v := range mappings.PathMappingSet {
				pathMapping = append(pathMapping, map[string]interface{}{
					"path":        v.Path,
					"environment": v.Environment,
				})
			}
		}
		if *info.Status == 1 {
			status = true
		}
		list = append(list, map[string]interface{}{
			"domain_name":        info.DomainName,
			"is_status_on":       status,
			"certificate_id":     info.CertificateId,
			"is_default_mapping": info.IsDefaultMapping,
			"protocol":           info.Protocol,
			"net_type":           info.NetType,
			"path_mappings":      pathMapping,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(serviceId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
