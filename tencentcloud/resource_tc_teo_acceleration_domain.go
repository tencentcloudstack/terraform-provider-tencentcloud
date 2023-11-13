/*
Provides a resource to create a teo acceleration_domain

Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id = ""
  domain_name = ""
  origin_info {
		origin_type = ""
		origin = ""
		backup_origin = ""
		private_access = ""
		private_parameters {
			name = ""
			value = ""
		}

  }
}
```

Import

teo acceleration_domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.acceleration_domain acceleration_domain_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTeoAcceleration_domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoAcceleration_domainCreate,
		Read:   resourceTencentCloudTeoAcceleration_domainRead,
		Update: resourceTencentCloudTeoAcceleration_domainUpdate,
		Delete: resourceTencentCloudTeoAcceleration_domainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID to which the acceleration domain name belongs.",
			},

			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Accelerated domain name.",
			},

			"origin_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Origin information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The origin type. Values:&amp;lt;li&amp;gt;`IP_DOMAIN`: IPv4/IPv6 address or domain name&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`COS`: COS bucket address &amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`ORIGIN_GROUP`: Origin group &amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`AWS_S3`: AWS S3 bucket address &amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SPACE`: EdgeOne Shield Space &amp;lt;/li&amp;gt;.",
						},
						"origin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The origin address. Enter the origin group ID if `OriginType=ORIGIN_GROUP`.",
						},
						"backup_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the secondary origin group (valid when `OriginType=ORIGIN_GROUP`). If it’s not specified, it indicates that secondary origins are not used.",
						},
						"private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to authenticate access to the private object storage origin (valid when `OriginType=COS/AWS_S3`). Values: &amp;lt;li&amp;gt;`on`: Enable private authentication.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`off`: Disable private authentication.&amp;lt;/li&amp;gt;If this field is not specified, the default value `off` is used.",
						},
						"private_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The private authentication parameters. This field is valid when `PrivateAccess=on`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter name. Values&amp;lt;li&amp;gt;`AccessKeyId`: Access Key ID&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`SecretAccessKey`: Secret Access Key&amp;lt;/li&amp;gt;.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter value.",
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

func resourceTencentCloudTeoAcceleration_domainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = teo.NewCreateAccelerationDomainRequest()
		response   = teo.NewCreateAccelerationDomainResponse()
		zoneId     string
		domainName string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
		originInfo := teo.OriginInfo{}
		if v, ok := dMap["origin_type"]; ok {
			originInfo.OriginType = helper.String(v.(string))
		}
		if v, ok := dMap["origin"]; ok {
			originInfo.Origin = helper.String(v.(string))
		}
		if v, ok := dMap["backup_origin"]; ok {
			originInfo.BackupOrigin = helper.String(v.(string))
		}
		if v, ok := dMap["private_access"]; ok {
			originInfo.PrivateAccess = helper.String(v.(string))
		}
		if v, ok := dMap["private_parameters"]; ok {
			for _, item := range v.([]interface{}) {
				privateParametersMap := item.(map[string]interface{})
				privateParameter := teo.PrivateParameter{}
				if v, ok := privateParametersMap["name"]; ok {
					privateParameter.Name = helper.String(v.(string))
				}
				if v, ok := privateParametersMap["value"]; ok {
					privateParameter.Value = helper.String(v.(string))
				}
				originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
			}
		}
		request.OriginInfo = &originInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateAccelerationDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo acceleration_domain failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, domainName}, FILED_SP))

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"online"}, 1*readRetryTimeout, time.Second, service.TeoAcceleration_domainStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoAcceleration_domainRead(d, meta)
}

func resourceTencentCloudTeoAcceleration_domainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	acceleration_domain, err := service.DescribeTeoAcceleration_domainById(ctx, zoneId, domainName)
	if err != nil {
		return err
	}

	if acceleration_domain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoAcceleration_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if acceleration_domain.ZoneId != nil {
		_ = d.Set("zone_id", acceleration_domain.ZoneId)
	}

	if acceleration_domain.DomainName != nil {
		_ = d.Set("domain_name", acceleration_domain.DomainName)
	}

	if acceleration_domain.OriginInfo != nil {
		originInfoMap := map[string]interface{}{}

		if acceleration_domain.OriginInfo.OriginType != nil {
			originInfoMap["origin_type"] = acceleration_domain.OriginInfo.OriginType
		}

		if acceleration_domain.OriginInfo.Origin != nil {
			originInfoMap["origin"] = acceleration_domain.OriginInfo.Origin
		}

		if acceleration_domain.OriginInfo.BackupOrigin != nil {
			originInfoMap["backup_origin"] = acceleration_domain.OriginInfo.BackupOrigin
		}

		if acceleration_domain.OriginInfo.PrivateAccess != nil {
			originInfoMap["private_access"] = acceleration_domain.OriginInfo.PrivateAccess
		}

		if acceleration_domain.OriginInfo.PrivateParameters != nil {
			privateParametersList := []interface{}{}
			for _, privateParameters := range acceleration_domain.OriginInfo.PrivateParameters {
				privateParametersMap := map[string]interface{}{}

				if privateParameters.Name != nil {
					privateParametersMap["name"] = privateParameters.Name
				}

				if privateParameters.Value != nil {
					privateParametersMap["value"] = privateParameters.Value
				}

				privateParametersList = append(privateParametersList, privateParametersMap)
			}

			originInfoMap["private_parameters"] = []interface{}{privateParametersList}
		}

		_ = d.Set("origin_info", []interface{}{originInfoMap})
	}

	return nil
}

func resourceTencentCloudTeoAcceleration_domainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyAccelerationDomainRequest  = teo.NewModifyAccelerationDomainRequest()
		modifyAccelerationDomainResponse = teo.NewModifyAccelerationDomainResponse()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	request.ZoneId = &zoneId
	request.DomainName = &domainName

	immutableArgs := []string{"zone_id", "domain_name", "origin_info"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
		}
	}

	if d.HasChange("domain_name") {
		if v, ok := d.GetOk("domain_name"); ok {
			request.DomainName = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_info") {
		if dMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
			originInfo := teo.OriginInfo{}
			if v, ok := dMap["origin_type"]; ok {
				originInfo.OriginType = helper.String(v.(string))
			}
			if v, ok := dMap["origin"]; ok {
				originInfo.Origin = helper.String(v.(string))
			}
			if v, ok := dMap["backup_origin"]; ok {
				originInfo.BackupOrigin = helper.String(v.(string))
			}
			if v, ok := dMap["private_access"]; ok {
				originInfo.PrivateAccess = helper.String(v.(string))
			}
			if v, ok := dMap["private_parameters"]; ok {
				for _, item := range v.([]interface{}) {
					privateParametersMap := item.(map[string]interface{})
					privateParameter := teo.PrivateParameter{}
					if v, ok := privateParametersMap["name"]; ok {
						privateParameter.Name = helper.String(v.(string))
					}
					if v, ok := privateParametersMap["value"]; ok {
						privateParameter.Value = helper.String(v.(string))
					}
					originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
				}
			}
			request.OriginInfo = &originInfo
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyAccelerationDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo acceleration_domain failed, reason:%+v", logId, err)
		return err
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"online"}, 1*readRetryTimeout, time.Second, service.TeoAcceleration_domainStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoAcceleration_domainRead(d, meta)
}

func resourceTencentCloudTeoAcceleration_domainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	if err := service.DeleteTeoAcceleration_domainById(ctx, zoneId, domainName); err != nil {
		return err
	}

	return nil
}
