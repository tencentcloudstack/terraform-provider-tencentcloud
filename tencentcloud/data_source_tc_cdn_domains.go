/*
Use this data source to query the detail information of CDN domain.

Example Usage

```hcl
data "tencentcloud_cdn_domains" "foo" {
  domain         	   = "xxxx.com"
  service_type   	   = "web"
  full_url_cache 	   = false
  origin_pull_protocol = "follow"
  status			   = "online"
  https_switch		   = "on"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
)

func dataSourceTencentCloudCdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the acceleration domain.",
			},
			"service_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Service type of Acceleration domain name. The available value include `web`, `download` and `media`.",
			},
			"full_url_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable full-path cache.",
			},
			"origin_pull_protocol": {
				Type:         schema.TypeBool,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
				Description:  "Origin-pull protocol configuration. The available value include `http`, `https` and `follow`.",
			},
			"https_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CDN_HTTPS_SWITCH),
				Description:  "HTTPS configuration. The available value include `on`, `off` and `processing`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"cdn_domain_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of cdn domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the acceleration domain.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service type of Acceleration domain name.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name acceleration region.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project CDN belongs to.",
						},
						"full_url_cache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable full-path cache.",
						},
						"origin": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin server configuration. It's a list and consist of at most one item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Master origin server type.",
									},
									"origin_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Master origin server list.",
									},
									"backup_origin_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backup origin server type.",
									},
									"backup_origin_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Backup origin server list.",
									},
									"backup_server_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.",
									},
									"server_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.",
									},
									"cos_private_access": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "When OriginType is COS, you can specify if access to private buckets is allowed.",
									},
									"origin_pull_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin-pull protocol configuration.",
									},
								},
							},
						},
						"https_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "HTTPS acceleration configuration. It's a list and consist of at most one item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"https_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTPS configuration switch.",
									},
									"http2_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP2 configuration switch.",
									},
									"ocsp_stapling_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OCSP configuration switch.",
									},
									"spdy_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Spdy configuration switch.",
									},
									"verify_client": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client certificate authentication feature.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of cdn domain.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCdnDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdn_domain.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	region := client.Region
	cdnService := CdnService{client: client}
	tagService := TagService{client: client}

	var domainFilterMap = make(map[string]interface{}, 5)

	if v, ok := d.GetOk(CDN_DATASOURCE_NAME_DOMAIN); ok {
		domainFilterMap[CDN_DATASOURCE_NAME_DOMAIN] = v.(string)
	}
	if v, ok := d.GetOk("service_type"); ok {
		domainFilterMap["service_type"] = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		domainFilterMap["status"] = v.(string)
	}
	if v, ok := d.GetOk("https_switch"); ok {
		domainFilterMap["https_switch"] = v.(string)
	}
	if v, ok := d.GetOk("full_url_cache"); ok {
		var value string
		if v.(bool) {
			value = "on"
		} else {
			value = "off"
		}

		domainFilterMap["full_url_cache"] = value
	}

	if len(domainFilterMap) == 0 {
		return nil
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByFilters(ctx, domainFilterMap)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describeDomainsConfigByFilters fail, reason:%s\n ", logId, err.Error())
		return err
	}
	if domainConfig == nil {
		return nil
	}

	domain := domainFilterMap[CDN_DATASOURCE_NAME_DOMAIN]
	_ = d.Set(CDN_DATASOURCE_NAME_DOMAIN, domain)
	_ = d.Set("service_type", domainConfig.ServiceType)
	_ = d.Set("project_id", domainConfig.ProjectId)
	_ = d.Set("area", domainConfig.Area)
	_ = d.Set("status", domainConfig.Status)

	if *domainConfig.CacheKey.FullUrlCache == CDN_SWITCH_OFF {
		_ = d.Set("full_url_cache", false)
	} else {
		_ = d.Set("full_url_cache", true)
	}

	origins := make([]map[string]interface{}, 0, 1)
	origin := make(map[string]interface{}, 8)
	origin["origin_type"] = domainConfig.Origin.OriginType
	origin["origin_list"] = domainConfig.Origin.Origins
	origin["server_name"] = domainConfig.Origin.ServerName
	origin["cos_private_access"] = domainConfig.Origin.CosPrivateAccess
	origin["origin_pull_protocol"] = domainConfig.Origin.OriginPullProtocol
	origin["backup_origin_type"] = domainConfig.Origin.BackupOriginType
	origin["backup_origin_list"] = domainConfig.Origin.BackupOrigins
	origin["backup_server_name"] = domainConfig.Origin.BackupServerName
	origins = append(origins, origin)
	_ = d.Set("origin", origins)

	httpsConfig := make(map[string]interface{}, 7)
	httpsConfig["https_switch"] = domainConfig.Https.Switch
	httpsConfig["http2_switch"] = domainConfig.Https.Http2
	httpsConfig["ocsp_stapling_switch"] = domainConfig.Https.OcspStapling
	httpsConfig["spdy_switch"] = domainConfig.Https.Spdy
	httpsConfig["verify_client"] = domainConfig.Https.VerifyClient
	_ = d.Set("https_config", httpsConfig)

	tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain.(string))
	if errRet != nil {
		return errRet
	}
	_ = d.Set("tags", tags)

	return nil
}
