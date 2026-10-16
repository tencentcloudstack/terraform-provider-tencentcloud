package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceDomainCreate,
		Read:   resourceTencentCloudTeoInferenceDomainRead,
		Update: resourceTencentCloudTeoInferenceDomainUpdate,
		Delete: resourceTencentCloudTeoInferenceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service ID.",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service domain.",
			},

			"auth_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Inference task request authentication switch. Valid values: Off: disable authentication; On: enable authentication. Default value: On.",
			},

			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Inference service domain list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name.",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain status. Valid values: Online: effective; Process: deploying; Offline: disabled; Init: not effective, waiting to activate the site; Forbidden: domain is forbidden.",
						},

						"auth_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Inference task request authentication switch. Valid values: Off: disable authentication; On: enable authentication. Default value: On.",
						},

						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME address.",
						},

						"ownership_verification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "When the domain needs to perform ownership verification to continue providing services, this object will carry the information required by the corresponding verification method.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Information required when using DNS resolution verification for CNAME or no-domain access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subdomain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Host record.",
												},

												"record_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record type.",
												},

												"record_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Record value.",
												},
											},
										},
									},

									"file_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Information required when using file verification for CNAME or no-domain access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "URL Path part that needs to be created.",
												},

												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Content of the verification file.",
												},
											},
										},
									},

									"ns_verification": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Information required to switch DNS servers for NS access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name_servers": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "DNS server addresses assigned to the user for NS access.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},

						"certificate": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain certificate information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain.",
									},

									"mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate configuration mode. Valid values: disable: no certificate; eofreecert: EdgeOne free certificate; sslcert: SSL certificate.",
									},

									"host_cert_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Server certificate configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cert_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Server certificate ID.",
												},

												"alias": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate alias name.",
												},

												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate type. Valid values: default: default certificate; upload: user uploaded; managed: Tencent Cloud managed.",
												},

												"expire_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate expiration time.",
												},

												"deploy_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate deployment time.",
												},

												"sign_algo": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Signature algorithm.",
												},

												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Certificate status. Valid values: deployed: deployed; processing: deploying; applying: applying; failed: application failed; issued: binding failed.",
												},
											},
										},
									},

									"apply_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application type. Valid values: apply: managed by EdgeOne; none: not managed by EdgeOne. If not specified, the default value is none.",
									},
								},
							},
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},

						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoInferenceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		request   = teo.NewCreateInferenceDomainRequest()
		zoneId    string
		serviceId string
		domain    string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("auth_switch"); ok {
		request.AuthSwitch = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateInferenceDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo inference domain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, serviceId, domain}, tccommon.FILED_SP))

	return resourceTencentCloudTeoInferenceDomainRead(d, meta)
}

func resourceTencentCloudTeoInferenceDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]
	domain := idSplit[2]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("service_id", serviceId)
	_ = d.Set("domain", domain)

	respData, err := service.DescribeTeoInferenceDomainById(ctx, zoneId, serviceId, domain)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_inference_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.AuthSwitch != nil {
		_ = d.Set("auth_switch", respData.AuthSwitch)
	}

	inferenceDomains := []*teo.InferenceDomain{respData}
	domainsList := make([]map[string]interface{}, 0, len(inferenceDomains))
	for _, inferenceDomain := range inferenceDomains {
		inferenceDomainMap := map[string]interface{}{}
		if inferenceDomain.Domain != nil {
			inferenceDomainMap["domain"] = inferenceDomain.Domain
		}

		if inferenceDomain.Status != nil {
			inferenceDomainMap["status"] = inferenceDomain.Status
		}

		if inferenceDomain.AuthSwitch != nil {
			inferenceDomainMap["auth_switch"] = inferenceDomain.AuthSwitch
		}

		if inferenceDomain.Cname != nil {
			inferenceDomainMap["cname"] = inferenceDomain.Cname
		}

		if inferenceDomain.OwnershipVerification != nil {
			ownershipVerificationMap := map[string]interface{}{}
			if inferenceDomain.OwnershipVerification.DnsVerification != nil {
				dnsVerificationMap := map[string]interface{}{}
				if inferenceDomain.OwnershipVerification.DnsVerification.Subdomain != nil {
					dnsVerificationMap["subdomain"] = inferenceDomain.OwnershipVerification.DnsVerification.Subdomain
				}

				if inferenceDomain.OwnershipVerification.DnsVerification.RecordType != nil {
					dnsVerificationMap["record_type"] = inferenceDomain.OwnershipVerification.DnsVerification.RecordType
				}

				if inferenceDomain.OwnershipVerification.DnsVerification.RecordValue != nil {
					dnsVerificationMap["record_value"] = inferenceDomain.OwnershipVerification.DnsVerification.RecordValue
				}

				ownershipVerificationMap["dns_verification"] = []interface{}{dnsVerificationMap}
			}

			if inferenceDomain.OwnershipVerification.FileVerification != nil {
				fileVerificationMap := map[string]interface{}{}
				if inferenceDomain.OwnershipVerification.FileVerification.Path != nil {
					fileVerificationMap["path"] = inferenceDomain.OwnershipVerification.FileVerification.Path
				}

				if inferenceDomain.OwnershipVerification.FileVerification.Content != nil {
					fileVerificationMap["content"] = inferenceDomain.OwnershipVerification.FileVerification.Content
				}

				ownershipVerificationMap["file_verification"] = []interface{}{fileVerificationMap}
			}

			if inferenceDomain.OwnershipVerification.NsVerification != nil {
				nsVerificationMap := map[string]interface{}{}
				if inferenceDomain.OwnershipVerification.NsVerification.NameServers != nil {
					nameServersList := make([]interface{}, 0, len(inferenceDomain.OwnershipVerification.NsVerification.NameServers))
					for _, nameServers := range inferenceDomain.OwnershipVerification.NsVerification.NameServers {
						nameServersList = append(nameServersList, nameServers)
					}
					nsVerificationMap["name_servers"] = nameServersList
				}

				ownershipVerificationMap["ns_verification"] = []interface{}{nsVerificationMap}
			}

			inferenceDomainMap["ownership_verification"] = []interface{}{ownershipVerificationMap}
		}

		if inferenceDomain.Certificate != nil {
			certificateMap := map[string]interface{}{}
			if inferenceDomain.Certificate.Host != nil {
				certificateMap["host"] = inferenceDomain.Certificate.Host
			}

			if inferenceDomain.Certificate.Mode != nil {
				certificateMap["mode"] = inferenceDomain.Certificate.Mode
			}

			if inferenceDomain.Certificate.HostCertInfo != nil {
				hostCertInfoList := make([]map[string]interface{}, 0, len(inferenceDomain.Certificate.HostCertInfo))
				for _, hostCertInfo := range inferenceDomain.Certificate.HostCertInfo {
					hostCertInfoMap := map[string]interface{}{}
					if hostCertInfo.CertId != nil {
						hostCertInfoMap["cert_id"] = hostCertInfo.CertId
					}

					if hostCertInfo.Alias != nil {
						hostCertInfoMap["alias"] = hostCertInfo.Alias
					}

					if hostCertInfo.Type != nil {
						hostCertInfoMap["type"] = hostCertInfo.Type
					}

					if hostCertInfo.ExpireTime != nil {
						hostCertInfoMap["expire_time"] = hostCertInfo.ExpireTime
					}

					if hostCertInfo.DeployTime != nil {
						hostCertInfoMap["deploy_time"] = hostCertInfo.DeployTime
					}

					if hostCertInfo.SignAlgo != nil {
						hostCertInfoMap["sign_algo"] = hostCertInfo.SignAlgo
					}

					if hostCertInfo.Status != nil {
						hostCertInfoMap["status"] = hostCertInfo.Status
					}

					hostCertInfoList = append(hostCertInfoList, hostCertInfoMap)
				}
				certificateMap["host_cert_info"] = hostCertInfoList
			}

			if inferenceDomain.Certificate.ApplyType != nil {
				certificateMap["apply_type"] = inferenceDomain.Certificate.ApplyType
			}

			inferenceDomainMap["certificate"] = []interface{}{certificateMap}
		}

		if inferenceDomain.CreateTime != nil {
			inferenceDomainMap["create_time"] = inferenceDomain.CreateTime
		}

		if inferenceDomain.UpdateTime != nil {
			inferenceDomainMap["update_time"] = inferenceDomain.UpdateTime
		}

		domainsList = append(domainsList, inferenceDomainMap)
	}

	_ = d.Set("domains", domainsList)

	return nil
}

func resourceTencentCloudTeoInferenceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	_ = tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	immutableArgs := []string{"zone_id", "service_id", "domain", "auth_switch"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudTeoInferenceDomainRead(d, meta)
}

func resourceTencentCloudTeoInferenceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]
	domain := idSplit[2]

	request := teo.NewOperateInferenceDomainRequest()
	request.ZoneId = helper.String(zoneId)
	request.ServiceId = helper.String(serviceId)
	request.Domain = helper.String(domain)
	request.Operation = helper.String("Delete")

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().OperateInferenceDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo inference domain failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
