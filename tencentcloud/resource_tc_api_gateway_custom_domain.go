/*
Use this resource to create custom domain of API gateway.

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
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayCustomDomainCreate,
		Read:   resourceTencentCloudAPIGatewayCustomDomainRead,
		Update: resourceTencentCloudAPIGatewayCustomDomainUpdate,
		Delete: resourceTencentCloudAPIGatewayCustomDomainDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				ForceNew:     true,
				Description:  "Unique service ID.",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom domain name to be bound.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Protocol supported by service. Valid values: `http`, `https`, `http&https`.",
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Network type. Valid values: `OUTER`, `INNER`.",
			},
			"is_default_mapping": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Whether the default path mapping is used. The default value is `true`. When it is `false`, it means custom path mapping. In this case, the `path_mappings` attribute is required.",
			},
			"default_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Default domain name.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unique certificate ID of the custom domain name to be bound. You can choose to upload for the `protocol` attribute value `https` or `http&https`.",
			},
			"path_mappings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom domain name path mapping. The data format is: `path#environment`. Optional values for the environment are `test`, `prepub`, and `release`.",
			},
			//compute
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Domain name resolution status. `1` means normal analysis, `0` means parsing failed.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceId         = d.Get("service_id").(string)
		subDomain         = d.Get("sub_domain").(string)
		protocol          = d.Get("protocol").(string)
		netType           = d.Get("net_type").(string)
		defaultDomain     = d.Get("default_domain").(string)
		isDefaultMapping  = d.Get("is_default_mapping").(bool)
		certificateId     string
		pathMappings      []string
		err               error
	)

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
	}
	if v, ok := d.GetOk("path_mappings"); ok {
		pathMappings = helper.InterfacesStrings(v.(*schema.Set).List())
	}

	err = apiGatewayService.BindSubDomainService(ctx, serviceId, subDomain, protocol, netType, defaultDomain, isDefaultMapping, certificateId, pathMappings)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{serviceId, subDomain}, FILED_SP))

	return resourceTencentCloudAPIGatewayCustomDomainRead(d, meta)
}

func resourceTencentCloudAPIGatewayCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.read")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		err               error
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. id:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]
	resultList, err := apiGatewayService.DescribeServiceSubDomainsService(ctx, serviceId, subDomain)
	if err != nil {
		return err
	}

	if len(resultList) == 0 {
		d.SetId("")
		return nil
	}

	resultInfo := resultList[0]
	info, err := apiGatewayService.DescribeServiceSubDomainMappings(ctx, serviceId, *resultInfo.DomainName)
	if err != nil {
		return fmt.Errorf("DescribeServiceSubDomainMappings err: %s", err.Error())
	}
	pathMap := make([]string, 0, len(info.PathMappingSet))
	for _, v := range info.PathMappingSet {
		pathMap = append(pathMap, strings.Join([]string{*v.Path, *v.Environment}, FILED_SP))
	}

	_ = d.Set("path_mappings", pathMap)
	_ = d.Set("status", resultInfo.Status)
	_ = d.Set("certificate_id", resultInfo.CertificateId)
	_ = d.Set("is_default_mapping", resultInfo.IsDefaultMapping)
	_ = d.Set("protocol", resultInfo.Protocol)
	_ = d.Set("net_type", resultInfo.NetType)
	_ = d.Set("service_id", serviceId)

	return nil
}

func resourceTencentCloudAPIGatewayCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.update")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		id                = d.Id()
		subDomain         string
		isDefaultMapping  bool
		certificateId     string
		protocol          string
		netType           string
		pathMappings      []string
		hasChange         bool
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]

	subDomain = d.Get("sub_domain").(string)
	if d.HasChange("sub_domain") {
		hasChange = true
	}

	isDefaultMapping = d.Get("is_default_mapping").(bool)
	if d.HasChange("is_default_mapping") {
		hasChange = true
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
	}
	if d.HasChange("certificate_id") {
		hasChange = true
	}

	netType = d.Get("net_type").(string)
	if d.HasChange("net_type") {
		hasChange = true
	}

	protocol = d.Get("protocol").(string)
	if d.HasChange("protocol") {
		hasChange = true
	}

	if v, ok := d.GetOk("path_mappings"); ok {
		pathMappings = helper.InterfacesStrings(v.(*schema.Set).List())
	}
	if d.HasChange("path_mappings") {
		hasChange = true
	}

	if hasChange {
		err := apiGatewayService.ModifySubDomainService(ctx, serviceId, subDomain, isDefaultMapping, certificateId, protocol, netType, pathMappings)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayCustomDomainRead(d, meta)
}

func resourceTencentCloudAPIGatewayCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_custom_domain.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		id                = d.Id()
		apigatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	results := strings.Split(id, FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]

	return apigatewayService.UnBindSubDomainService(ctx, serviceId, subDomain)
}
