/*
Use this data source to query detailed information of ssl describe_certificate

Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate" "describe_certificate" {
                                                                      }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSslDescribeCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCertificateRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},

			"owner_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Account UIN.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"project_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Project ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"from": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate source: Trustasia,uploadNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"certificate_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate type: CA = CA certificate, SVR = server certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"package_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Types of Certificate Package: 1 = Geotrust DV SSL CA -G3, 2 = Trustasia TLS RSA CA, 3 = SecureSite Enhanced Enterprise Edition (EV Pro), 4 = SecureSite enhanced (EV), 5 = SecureSite Enterprise Professional Edition (OVPro), 6 = SecureSite Enterprise (OV), 7 = SecureSite Enterprise (OV) compatriots, 8 = Geotrust enhanced type (EV), 9 = Geotrust Enterprise (OV), 10 = Geotrust Enterprise (OV) pass,11 = Trustasia Domain Multi -domain SSL certificate, 12 = Trustasia domain model (DV) passing, 13 = Trustasia Enterprise Passing Character (OV) SSL certificate (D3), 14 = Trustasia Enterprise (OV) SSL certificate (D3), 15= Trustasia Enterprise Multi -domain name (OV) SSL certificate (D3), 16 = Trustasia enhanced (EV) SSL certificate (D3), 17 = Trustasia enhanced multi -domain name (EV) SSL certificate (D3), 18 = GlobalSign enterprise type enterprise type(OV) SSL certificate, 19 = GlobalSign Enterprise Type -type STL Certificate, 20 = GlobalSign enhanced (EV) SSL certificate, 21 = Trustasia Enterprise Tongzhi Multi -domain name (OV) SSL certificate (D3), 22 = GlobalSignignMulti -domain name (OV) SSL certificate, 23 = GlobalSign Enterprise Type -type multi -domain name (OV) SSL certificate, 24 = GlobalSign enhanced multi -domain name (EV) SSL certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"product_zh_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate issuer name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"domain": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"alias": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Remark name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "= Submitted information, to be uploaded to confirmation letter, 9 = Certificate is revoked, 10 = revoked, 11 = Re -issuance, 12 = Upload and revoke the confirmation letter.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"status_msg": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"verify_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Verification type: DNS_AUTO = Automatic DNS verification, DNS = manual DNS verification, file = file verification, email = email verification.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"vulnerability_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Vulnerability scanning status.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"cert_begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate takes effect time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"cert_end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The certificate is invalid time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"validity_period": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Validity period: unit (month).Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"insert_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"order_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Order ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"certificate_extra": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate extension information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate can be configured in the number of domain names.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"origin_certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original certificate ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"replaced_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Re -issue the original ID of the certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"replaced_for": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Re -issue a new ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"renew_order": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "New order certificate ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"s_m_cert": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it a national secret certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"dv_auth_detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DV certification information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dv_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_key_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification sub -domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auths": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DV certification information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dv_auth_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"dv_auth_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"dv_auth_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"dv_auth_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"dv_auth_sub_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification sub -domain name,Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"dv_auth_verify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"vulnerability_report": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Vulnerability scanning evaluation report.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"package_type_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Certificate type name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"status_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Status description.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"subject_alt_name": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The certificate contains multiple domain names (containing the main domain name).Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_vip": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is a VIP customer.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_wildcard": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is a pan -domain certificate certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_dv": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it is the DV version.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_vulnerability": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the vulnerability scanning function is enabled.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"renew_able": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether you can issue a certificate.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"submitted_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Submitted information information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"csr_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CSR type, (online = online CSR, PARSE = paste CSR).Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"csr_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CSR content.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"certificate_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"domain_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "DNS information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"key_password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private key password.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise or unit name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_division": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Department.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Nation.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "City.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"organization_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Province.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"postal_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Postal code.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"phone_area_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Local region code.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"phone_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Landline number.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"admin_first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"admin_last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The surname of the administrator.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"admin_phone_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator phone number.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"admin_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator mailbox address.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"admin_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator position.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"contact_first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"contact_last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact surname.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"contact_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact phone number.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"contact_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact mailbox address,Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"contact_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact position.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"verify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Verification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"deployable": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether it can be deployed.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"c_a_encrypt_algorithms": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "All encryption methods of CA certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"c_a_common_names": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "All general names of the CA certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"c_a_end_times": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CA certificate all maturity timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"dv_revoke_auth_detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DV certificate revoking verification valueNote: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dv_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification key.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification value.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication value domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV authentication value path.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification sub -domain name,Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"dv_auth_verify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DV certification type.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_describe_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ownerUin = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ownerUin))
	if certificateId != nil {
		_ = d.Set("certificate_id", certificateId)
	}

	if ownerUin != nil {
		_ = d.Set("owner_uin", ownerUin)
	}

	if projectId != nil {
		_ = d.Set("project_id", projectId)
	}

	if from != nil {
		_ = d.Set("from", from)
	}

	if certificateType != nil {
		_ = d.Set("certificate_type", certificateType)
	}

	if packageType != nil {
		_ = d.Set("package_type", packageType)
	}

	if productZhName != nil {
		_ = d.Set("product_zh_name", productZhName)
	}

	if domain != nil {
		_ = d.Set("domain", domain)
	}

	if alias != nil {
		_ = d.Set("alias", alias)
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	if statusMsg != nil {
		_ = d.Set("status_msg", statusMsg)
	}

	if verifyType != nil {
		_ = d.Set("verify_type", verifyType)
	}

	if vulnerabilityStatus != nil {
		_ = d.Set("vulnerability_status", vulnerabilityStatus)
	}

	if certBeginTime != nil {
		_ = d.Set("cert_begin_time", certBeginTime)
	}

	if certEndTime != nil {
		_ = d.Set("cert_end_time", certEndTime)
	}

	if validityPeriod != nil {
		_ = d.Set("validity_period", validityPeriod)
	}

	if insertTime != nil {
		_ = d.Set("insert_time", insertTime)
	}

	if orderId != nil {
		_ = d.Set("order_id", orderId)
	}

	if certificateExtra != nil {
		certificateExtraMap := map[string]interface{}{}

		if certificateExtra.DomainNumber != nil {
			certificateExtraMap["domain_number"] = certificateExtra.DomainNumber
		}

		if certificateExtra.OriginCertificateId != nil {
			certificateExtraMap["origin_certificate_id"] = certificateExtra.OriginCertificateId
		}

		if certificateExtra.ReplacedBy != nil {
			certificateExtraMap["replaced_by"] = certificateExtra.ReplacedBy
		}

		if certificateExtra.ReplacedFor != nil {
			certificateExtraMap["replaced_for"] = certificateExtra.ReplacedFor
		}

		if certificateExtra.RenewOrder != nil {
			certificateExtraMap["renew_order"] = certificateExtra.RenewOrder
		}

		if certificateExtra.SMCert != nil {
			certificateExtraMap["s_m_cert"] = certificateExtra.SMCert
		}

		ids = append(ids, *certificateExtra.CertificateId)
		_ = d.Set("certificate_extra", certificateExtraMap)
	}

	if dvAuthDetail != nil {
		dvAuthDetailMap := map[string]interface{}{}

		if dvAuthDetail.DvAuthKey != nil {
			dvAuthDetailMap["dv_auth_key"] = dvAuthDetail.DvAuthKey
		}

		if dvAuthDetail.DvAuthValue != nil {
			dvAuthDetailMap["dv_auth_value"] = dvAuthDetail.DvAuthValue
		}

		if dvAuthDetail.DvAuthDomain != nil {
			dvAuthDetailMap["dv_auth_domain"] = dvAuthDetail.DvAuthDomain
		}

		if dvAuthDetail.DvAuthPath != nil {
			dvAuthDetailMap["dv_auth_path"] = dvAuthDetail.DvAuthPath
		}

		if dvAuthDetail.DvAuthKeySubDomain != nil {
			dvAuthDetailMap["dv_auth_key_sub_domain"] = dvAuthDetail.DvAuthKeySubDomain
		}

		if dvAuthDetail.DvAuths != nil {
			dvAuthsList := []interface{}{}
			for _, dvAuths := range dvAuthDetail.DvAuths {
				dvAuthsMap := map[string]interface{}{}

				if dvAuths.DvAuthKey != nil {
					dvAuthsMap["dv_auth_key"] = dvAuths.DvAuthKey
				}

				if dvAuths.DvAuthValue != nil {
					dvAuthsMap["dv_auth_value"] = dvAuths.DvAuthValue
				}

				if dvAuths.DvAuthDomain != nil {
					dvAuthsMap["dv_auth_domain"] = dvAuths.DvAuthDomain
				}

				if dvAuths.DvAuthPath != nil {
					dvAuthsMap["dv_auth_path"] = dvAuths.DvAuthPath
				}

				if dvAuths.DvAuthSubDomain != nil {
					dvAuthsMap["dv_auth_sub_domain"] = dvAuths.DvAuthSubDomain
				}

				if dvAuths.DvAuthVerifyType != nil {
					dvAuthsMap["dv_auth_verify_type"] = dvAuths.DvAuthVerifyType
				}

				dvAuthsList = append(dvAuthsList, dvAuthsMap)
			}

			dvAuthDetailMap["dv_auths"] = []interface{}{dvAuthsList}
		}

		ids = append(ids, *dvAuthDetail.CertificateId)
		_ = d.Set("dv_auth_detail", dvAuthDetailMap)
	}

	if vulnerabilityReport != nil {
		_ = d.Set("vulnerability_report", vulnerabilityReport)
	}

	if packageTypeName != nil {
		_ = d.Set("package_type_name", packageTypeName)
	}

	if statusName != nil {
		_ = d.Set("status_name", statusName)
	}

	if subjectAltName != nil {
		_ = d.Set("subject_alt_name", subjectAltName)
	}

	if isVip != nil {
		_ = d.Set("is_vip", isVip)
	}

	if isWildcard != nil {
		_ = d.Set("is_wildcard", isWildcard)
	}

	if isDv != nil {
		_ = d.Set("is_dv", isDv)
	}

	if isVulnerability != nil {
		_ = d.Set("is_vulnerability", isVulnerability)
	}

	if renewAble != nil {
		_ = d.Set("renew_able", renewAble)
	}

	if submittedData != nil {
		submittedDataMap := map[string]interface{}{}

		if submittedData.CsrType != nil {
			submittedDataMap["csr_type"] = submittedData.CsrType
		}

		if submittedData.CsrContent != nil {
			submittedDataMap["csr_content"] = submittedData.CsrContent
		}

		if submittedData.CertificateDomain != nil {
			submittedDataMap["certificate_domain"] = submittedData.CertificateDomain
		}

		if submittedData.DomainList != nil {
			submittedDataMap["domain_list"] = submittedData.DomainList
		}

		if submittedData.KeyPassword != nil {
			submittedDataMap["key_password"] = submittedData.KeyPassword
		}

		if submittedData.OrganizationName != nil {
			submittedDataMap["organization_name"] = submittedData.OrganizationName
		}

		if submittedData.OrganizationDivision != nil {
			submittedDataMap["organization_division"] = submittedData.OrganizationDivision
		}

		if submittedData.OrganizationAddress != nil {
			submittedDataMap["organization_address"] = submittedData.OrganizationAddress
		}

		if submittedData.OrganizationCountry != nil {
			submittedDataMap["organization_country"] = submittedData.OrganizationCountry
		}

		if submittedData.OrganizationCity != nil {
			submittedDataMap["organization_city"] = submittedData.OrganizationCity
		}

		if submittedData.OrganizationRegion != nil {
			submittedDataMap["organization_region"] = submittedData.OrganizationRegion
		}

		if submittedData.PostalCode != nil {
			submittedDataMap["postal_code"] = submittedData.PostalCode
		}

		if submittedData.PhoneAreaCode != nil {
			submittedDataMap["phone_area_code"] = submittedData.PhoneAreaCode
		}

		if submittedData.PhoneNumber != nil {
			submittedDataMap["phone_number"] = submittedData.PhoneNumber
		}

		if submittedData.AdminFirstName != nil {
			submittedDataMap["admin_first_name"] = submittedData.AdminFirstName
		}

		if submittedData.AdminLastName != nil {
			submittedDataMap["admin_last_name"] = submittedData.AdminLastName
		}

		if submittedData.AdminPhoneNum != nil {
			submittedDataMap["admin_phone_num"] = submittedData.AdminPhoneNum
		}

		if submittedData.AdminEmail != nil {
			submittedDataMap["admin_email"] = submittedData.AdminEmail
		}

		if submittedData.AdminPosition != nil {
			submittedDataMap["admin_position"] = submittedData.AdminPosition
		}

		if submittedData.ContactFirstName != nil {
			submittedDataMap["contact_first_name"] = submittedData.ContactFirstName
		}

		if submittedData.ContactLastName != nil {
			submittedDataMap["contact_last_name"] = submittedData.ContactLastName
		}

		if submittedData.ContactNumber != nil {
			submittedDataMap["contact_number"] = submittedData.ContactNumber
		}

		if submittedData.ContactEmail != nil {
			submittedDataMap["contact_email"] = submittedData.ContactEmail
		}

		if submittedData.ContactPosition != nil {
			submittedDataMap["contact_position"] = submittedData.ContactPosition
		}

		if submittedData.VerifyType != nil {
			submittedDataMap["verify_type"] = submittedData.VerifyType
		}

		ids = append(ids, *submittedData.CertificateId)
		_ = d.Set("submitted_data", submittedDataMap)
	}

	if deployable != nil {
		_ = d.Set("deployable", deployable)
	}

	if cAEncryptAlgorithms != nil {
		_ = d.Set("c_a_encrypt_algorithms", cAEncryptAlgorithms)
	}

	if cACommonNames != nil {
		_ = d.Set("c_a_common_names", cACommonNames)
	}

	if cAEndTimes != nil {
		_ = d.Set("c_a_end_times", cAEndTimes)
	}

	if dvRevokeAuthDetail != nil {
		for _, dvAuths := range dvRevokeAuthDetail {
			dvAuthsMap := map[string]interface{}{}

			if dvAuths.DvAuthKey != nil {
				dvAuthsMap["dv_auth_key"] = dvAuths.DvAuthKey
			}

			if dvAuths.DvAuthValue != nil {
				dvAuthsMap["dv_auth_value"] = dvAuths.DvAuthValue
			}

			if dvAuths.DvAuthDomain != nil {
				dvAuthsMap["dv_auth_domain"] = dvAuths.DvAuthDomain
			}

			if dvAuths.DvAuthPath != nil {
				dvAuthsMap["dv_auth_path"] = dvAuths.DvAuthPath
			}

			if dvAuths.DvAuthSubDomain != nil {
				dvAuthsMap["dv_auth_sub_domain"] = dvAuths.DvAuthSubDomain
			}

			if dvAuths.DvAuthVerifyType != nil {
				dvAuthsMap["dv_auth_verify_type"] = dvAuths.DvAuthVerifyType
			}

			ids = append(ids, *dvAuths.CertificateId)
			tmpList = append(tmpList, dvAuthsMap)
		}

		_ = d.Set("dv_revoke_auth_detail", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
