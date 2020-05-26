/*
Provides a resource to create a CDN domain.

Example Usage

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }

  tags = {
    hello = "world"
  }
}
```

Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdnDomainCreate,
		Read:   resourceTencentCloudCdnDomainRead,
		Update: resourceTencentCloudCdnDomainUpdate,
		Delete: resourceTencentCloudCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the acceleration domain.",
			},
			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Service type of Acceleration domain name. Valid values are `web`, `download` and `media`.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project CDN belongs to, default to 0.",
			},
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CDN_AREA),
				Description:  "Domain name acceleration region.  Valid values are `mainland`, `overseas` and `global`.",
			},
			"full_url_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable full-path cache. Default value is `true`.",
			},
			"origin": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Origin server configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(CDN_ORIGIN_TYPE),
							Description:  "Master origin server type. Valid values are `domain`, `cos`, `ip`, `ipv6` and `ip_ipv6`.",
						},
						"origin_list": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Master origin server list. Valid values can be ip or doamin name. When modifying the origin server, you need to enter the corresponding `origin_type`.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.",
						},
						"cos_private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`, and default value is `off`.",
						},
						"origin_pull_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
							ValidateFunc: validateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
							Description:  "Origin-pull protocol configuration. Valid values are `http`, `https` and `follow`, and default value is `http`.",
						},
						"backup_origin_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(CDN_BACKUP_ORIGIN_TYPE),
							Description:  "Backup origin server type. Valid values are `domain` and `ip`.",
						},
						"backup_origin_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Backup origin server list. Valid values can be ip or doamin name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.",
						},
						"backup_server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.",
						},
					},
				},
			},
			"https_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "HTTPS acceleration configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTPS configuration switch. Valid values are `on` and `off`.",
						},
						"http2_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTP2 configuration switch. Valid values are `on` and `off`, and default value is `off`.",
						},
						"ocsp_stapling_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "OCSP configuration switch. Valid values are `on` and `off`, and default value is `off`.",
						},
						"spdy_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Spdy configuration switch. Valid values are `on` and `off`, and default value is `off`.",
						},
						"verify_client": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Client certificate authentication feature. Valid values are `on` and `off`, and default value is `off`.",
						},
						"server_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Server certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate ID.",
									},
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.",
									},
									"private_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server key information. This is required when uploading an external certificate.",
									},
									"message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Certificate remarks.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of server certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of server certificate.",
									},
								},
							},
						},
						"client_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Client certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Client Certificate PEM format, requires Base64 encoding.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of client certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of client certificate.",
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of cdn domain.",
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration service status.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address of domain name.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of domain name.",
			},
		},
	}
}

func resourceTencentCloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cdnService := CdnService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cdn.NewAddCdnDomainRequest()
	domain := d.Get("domain").(string)
	request.Domain = &domain
	request.ServiceType = helper.String(d.Get("service_type").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}
	fullUrlCache := d.Get("full_url_cache").(bool)
	request.CacheKey = &cdn.CacheKey{}
	if fullUrlCache {
		request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
	} else {
		request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
	}

	// origin
	origins := d.Get("origin").([]interface{})
	if len(origins) < 1 {
		return fmt.Errorf("origin is required")
	}
	origin := origins[0].(map[string]interface{})
	request.Origin = &cdn.Origin{}
	request.Origin.OriginType = helper.String(origin["origin_type"].(string))
	originList := origin["origin_list"].([]interface{})
	request.Origin.Origins = make([]*string, 0, len(originList))
	for _, item := range originList {
		request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
	}
	if v := origin["server_name"]; v.(string) != "" {
		request.Origin.ServerName = helper.String(v.(string))
	}
	if v := origin["cos_private_access"]; v.(string) != "" {
		request.Origin.CosPrivateAccess = helper.String(v.(string))
	}
	if v := origin["origin_pull_protocol"]; v.(string) != "" {
		request.Origin.OriginPullProtocol = helper.String(v.(string))
	}
	if v := origin["backup_origin_type"]; v.(string) != "" {
		request.Origin.BackupOriginType = helper.String(v.(string))
	}
	if v := origin["backup_server_name"]; v.(string) != "" {
		request.Origin.BackupServerName = helper.String(v.(string))
	}
	if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
		backupOriginList := v.([]interface{})
		request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
		for _, item := range backupOriginList {
			request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
		}
	}

	// https config
	if v, ok := d.GetOk("https_config"); ok {
		httpsConfigs := v.([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := meta.(*TencentCloudClient).apiV3Conn.UseCdnClient().AddCdnDomain(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROE {
					return resource.NonRetryableError(err)
				}
			}
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(domain)

	time.Sleep(1 * time.Second)
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if err != nil {
			return retryError(err, InternalError)
		}
		if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	region := client.Region
	cdnService := CdnService{client: client}
	tagService := TagService{client: client}

	domain := d.Id()
	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("service_type", domainConfig.ServiceType)
	_ = d.Set("project_id", domainConfig.ProjectId)
	_ = d.Set("area", domainConfig.Area)
	_ = d.Set("status", domainConfig.Status)
	_ = d.Set("create_time", domainConfig.CreateTime)
	_ = d.Set("cname", domainConfig.Cname)
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

	httpsConfigs := make([]map[string]interface{}, 0, 1)
	httpsConfig := make(map[string]interface{}, 7)
	httpsConfig["https_switch"] = domainConfig.Https.Switch
	httpsConfig["http2_switch"] = domainConfig.Https.Http2
	httpsConfig["ocsp_stapling_switch"] = domainConfig.Https.OcspStapling
	httpsConfig["spdy_switch"] = domainConfig.Https.Spdy
	httpsConfig["verify_client"] = domainConfig.Https.VerifyClient

	oldHttpsConfigs := make([]interface{}, 0)
	if _, ok := d.GetOk("https_config"); ok {
		oldHttpsConfigs = d.Get("https_config").([]interface{})
	}
	oldHttpsConfig := make(map[string]interface{})
	if len(oldHttpsConfigs) > 0 {
		oldHttpsConfig = oldHttpsConfigs[0].(map[string]interface{})
	}
	oldServerConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["server_certificate_config"]; ok {
		oldServerConfigs = oldHttpsConfig["server_certificate_config"].([]interface{})
	}
	oldServerConfig := make(map[string]interface{})
	if len(oldServerConfigs) > 0 {
		oldServerConfig = oldServerConfigs[0].(map[string]interface{})
	}
	oldClientConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["client_certificate_config"]; ok {
		oldClientConfigs = oldHttpsConfig["client_certificate_config"].([]interface{})
	}
	oldClientConfig := make(map[string]interface{})
	if len(oldClientConfigs) > 0 {
		oldClientConfig = oldClientConfigs[0].(map[string]interface{})
	}

	if domainConfig.Https.CertInfo != nil && domainConfig.Https.CertInfo.CertName != nil {
		serverCertConfigs := make([]map[string]interface{}, 0, 1)
		serverCertConfig := make(map[string]interface{}, 5)
		serverCertConfig["certificate_id"] = domainConfig.Https.CertInfo.CertId
		serverCertConfig["certificate_name"] = domainConfig.Https.CertInfo.CertName
		serverCertConfig["certificate_content"] = oldServerConfig["certificate_content"]
		serverCertConfig["private_key"] = oldServerConfig["private_key"]
		serverCertConfig["message"] = domainConfig.Https.CertInfo.Message
		serverCertConfig["deploy_time"] = domainConfig.Https.CertInfo.DeployTime
		serverCertConfig["expire_time"] = domainConfig.Https.CertInfo.ExpireTime
		serverCertConfigs = append(serverCertConfigs, serverCertConfig)
		httpsConfig["server_certificate_config"] = serverCertConfigs
	}
	if domainConfig.Https.ClientCertInfo != nil && domainConfig.Https.ClientCertInfo.CertName != nil {
		clientCertConfigs := make([]map[string]interface{}, 0, 1)
		clientCertConfig := make(map[string]interface{}, 2)
		clientCertConfig["certificate_content"] = oldClientConfig["certificate_content"]
		clientCertConfig["certificate_name"] = domainConfig.Https.ClientCertInfo.CertName
		clientCertConfig["deploy_time"] = domainConfig.Https.ClientCertInfo.DeployTime
		clientCertConfig["expire_time"] = domainConfig.Https.ClientCertInfo.ExpireTime
		clientCertConfigs = append(clientCertConfigs, clientCertConfig)
		httpsConfig["client_certificate_config"] = clientCertConfigs
	}
	httpsConfigs = append(httpsConfigs, httpsConfig)
	_ = d.Set("https_config", httpsConfigs)

	tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
	if err != nil {
		return nil
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	cdnService := CdnService{client: client}

	d.Partial(true)
	updateAttrs := make([]string, 0)

	domain := d.Id()
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if d.HasChange("service_type") {
		request.ServiceType = helper.String(d.Get("service_type").(string))
		updateAttrs = append(updateAttrs, "service_type")
	}
	if d.HasChange("project_id") {
		request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
		updateAttrs = append(updateAttrs, "project_id")
	}
	if d.HasChange("area") {
		request.Area = helper.String(d.Get("area").(string))
		updateAttrs = append(updateAttrs, "area")
	}
	if d.HasChange("full_url_cache") {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
		updateAttrs = append(updateAttrs, "full_url_cache")
	}
	if d.HasChange("origin") {
		updateAttrs = append(updateAttrs, "origin")
		origins := d.Get("origin").([]interface{})
		if len(origins) < 1 {
			return fmt.Errorf("origin is required")
		}
		origin := origins[0].(map[string]interface{})
		request.Origin = &cdn.Origin{}
		request.Origin.OriginType = helper.String(origin["origin_type"].(string))
		originList := origin["origin_list"].([]interface{})
		request.Origin.Origins = make([]*string, 0, len(originList))
		for _, item := range originList {
			request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
		}
		if v := origin["server_name"]; v.(string) != "" {
			request.Origin.ServerName = helper.String(v.(string))
		}
		if v := origin["cos_private_access"]; v.(string) != "" {
			request.Origin.CosPrivateAccess = helper.String(v.(string))
		}
		if v := origin["origin_pull_protocol"]; v.(string) != "" {
			request.Origin.OriginPullProtocol = helper.String(v.(string))
		}
		if v := origin["backup_origin_type"]; v.(string) != "" {
			request.Origin.BackupOriginType = helper.String(v.(string))
		}
		if v := origin["backup_server_name"]; v.(string) != "" {
			request.Origin.BackupServerName = helper.String(v.(string))
		}
		if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
			backupOriginList := v.([]interface{})
			request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
			for _, item := range backupOriginList {
				request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
			}
		}
	}
	if d.HasChange("https_config") {
		updateAttrs = append(updateAttrs, "https_config")
		httpsConfigs := d.Get("https_config").([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
		}
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err := meta.(*TencentCloudClient).apiV3Conn.UseCdnClient().UpdateDomainConfig(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		for _, attr := range updateAttrs {
			d.SetPartial(attr)
		}

		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return retryError(err, InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	cdnService := CdnService{client: client}

	domain := d.Id()
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		deleteTags := make([]string, 0, len(tags))
		for key := range tags {
			deleteTags = append(deleteTags, key)
		}
		err := tagService.ModifyTags(ctx, resourceName, nil, deleteTags)
		if err != nil {
			return err
		}
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		return nil
	}

	if *domainConfig.Status == CDN_DOMAIN_STATUS_ONLINE {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet = cdnService.StopDomain(ctx, domain)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return retryError(err, InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet = cdnService.DeleteDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
