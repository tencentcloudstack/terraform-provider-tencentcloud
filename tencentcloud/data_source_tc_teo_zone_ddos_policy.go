/*
Use this data source to query zone ddos policy.

Example Usage

```hcl
data "tencentcloud_teo_zone_ddos_policy" "example" {
  zone_id = ""
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
)

func dataSourceTencentCloudTeoZoneDdosPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoZoneDdosPolicyRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site ID.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "App ID.",
			},
			"shield_areas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Shield areas of the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID.",
						},
						"policy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values: `domain`, `application`.",
						},
						"entity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When `Type` is `domain`, this field is `ZoneId`. When `Type` is `application`, this field is `ProxyId`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"entity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When `Type` is `domain`, this field is `ZoneName`. When `Type` is `application`, this field is `ProxyName`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tcp_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP forwarding rule number of layer 4 application.",
						},
						"udp_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "UDP forwarding rule number of layer 4 application.",
						},
						"share": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the resource is shared.",
						},
						"application": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Layer 7 Domain Name Parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subdomain.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the subdomain. Note: This field may return null, indicating that no valid value can be obtained, init: waiting to config NS; offline: waiting to enable site accelerating; process: config deployment processing; online: normal status.",
									},
									"accelerate_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "on: Enable; off: Disable.",
									},
									"security_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "on: Enable; off: Disable.",
									},
								},
							},
						},
					},
				},
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All subdomain info. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdomain.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the subdomain. Note: This field may return null, indicating that no valid value can be obtained, init: waiting to config NS; offline: waiting to enable site accelerating; process: config deployment processing; online: normal status.",
						},
						"accelerate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "on: Enable; off: Disable.",
						},
						"security_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "on: Enable; off: Disable.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoZoneDdosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_zone_ddos_policy.read")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		zoneId     = d.Get("zone_id").(string)
		ddosPolicy *teo.DescribeZoneDDoSPolicyResponseParams
		err        error
	)

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ddosPolicy, err = service.DescribeZoneDDoSPolicy(ctx, zoneId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	appId := strconv.FormatInt(*ddosPolicy.AppId, 10)

	shieldAreasList := make([]map[string]interface{}, 0, len(ddosPolicy.ShieldAreas))
	shieldAreas := ddosPolicy.ShieldAreas
	for _, v := range shieldAreas {
		applications := make([]map[string]interface{}, 0, len(v.Application))
		for _, vv := range v.Application {
			application := map[string]interface{}{
				"host":            vv.Host,
				"status":          vv.Status,
				"accelerate_type": vv.AccelerateType,
				"security_type":   vv.SecurityType,
			}
			applications = append(applications, application)
		}
		shieldArea := map[string]interface{}{
			"zone_id":     v.ZoneId,
			"policy_id":   v.PolicyId,
			"type":        v.Type,
			"entity_name": v.EntityName,
			"application": applications,
			"tcp_num":     v.TcpNum,
			"udp_num":     v.UdpNum,
			"entity":      v.Entity,
			"share":       v.Share,
		}
		shieldAreasList = append(shieldAreasList, shieldArea)
	}
	if err = d.Set("shield_areas", shieldAreasList); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	domainsList := make([]map[string]interface{}, 0, len(ddosPolicy.Domains))
	for _, v := range ddosPolicy.Domains {
		application := map[string]interface{}{
			"host":            v.Host,
			"status":          v.Status,
			"accelerate_type": v.AccelerateType,
			"security_type":   v.SecurityType,
		}
		domainsList = append(domainsList, application)
	}
	if err = d.Set("domains", domainsList); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	if err = d.Set("app_id", appId); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(appId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"app_id":       appId,
			"Shield_areas": shieldAreasList,
			"domains":      domainsList,
		}); e != nil {
			return e
		}
	}
	return nil
}
