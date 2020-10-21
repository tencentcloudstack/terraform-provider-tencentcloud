/*
Use this data source to query API gateway domain list.

Example Usage

```hcl
resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-ohxqslqe"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test","/root#release"]
}

data "tencentcloud_api_gateway_customer_domains" "id" {
	service_id = tencentcloud_api_gateway_custom_domain.foo.service_id
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func dataSourceTencentCloudAPIGatewayCustomerDomains() *schema.Resource {
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
						"status": {
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
	defer logElapsed("data_source.tencentcloud_api_gateway_customer_domains.read")

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		infos             []*apigateway.DomainSetList
		list              []map[string]interface{}
		err               error
	)
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeServiceSubDomains(ctx, serviceId)
		if err != nil {
			return retryError(err, InternalError)
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
			"status":             status,
			"certificate_id":     info.CertificateId,
			"is_default_mapping": info.IsDefaultMapping,
			"protocol":           info.Protocol,
			"net_type":           info.NetType,
			"path_mappings":      pathMapping,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
	}

	d.SetId(serviceId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil
}
