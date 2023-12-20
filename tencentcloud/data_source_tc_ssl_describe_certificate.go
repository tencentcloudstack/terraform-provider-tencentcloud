package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func dataSourceTencentCloudSslDescribeCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCertificateRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "result list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "domain name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Description: "status information.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Description: "application time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
									"company_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Type of company. Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
							Description: "status description.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
										Description: "department.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"organization_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "address.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"organization_country": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "nation.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"organization_city": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "city.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
					},
				}},

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

	service := SslService{client: meta.(*TencentCloudClient).apiV3Conn}
	responese := ssl.DescribeCertificateResponseParams{}
	CertificateId := d.Get("certificate_id").(string)
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateByID(ctx, CertificateId)
		if e != nil {
			return retryError(e)
		}
		responese = *result
		return nil
	})
	if err != nil {
		return err
	}
	sslResponseMap := map[string]interface{}{}
	if responese.OwnerUin != nil {
		sslResponseMap["owner_uin"] = responese.OwnerUin
	}

	if responese.ProjectId != nil {
		sslResponseMap["project_id"] = responese.ProjectId
	}

	if responese.From != nil {
		sslResponseMap["from"] = responese.From
	}

	if responese.CertificateType != nil {
		sslResponseMap["certificate_type"] = responese.CertificateType
	}

	if responese.PackageType != nil {
		sslResponseMap["package_type"] = responese.PackageType
	}

	if responese.ProductZhName != nil {
		sslResponseMap["product_zh_name"] = responese.ProductZhName
	}

	if responese.Domain != nil {
		sslResponseMap["domain"] = responese.Domain
	}

	if responese.Alias != nil {
		sslResponseMap["alias"] = responese.Alias
	}

	if responese.Status != nil {
		sslResponseMap["status"] = responese.Status
	}

	if responese.StatusMsg != nil {
		sslResponseMap["status_msg"] = responese.StatusMsg
	}

	if responese.VerifyType != nil {
		sslResponseMap["verify_type"] = responese.VerifyType
	}

	if responese.VulnerabilityStatus != nil {
		sslResponseMap["vulnerability_status"] = responese.VulnerabilityStatus
	}

	if responese.CertBeginTime != nil {
		sslResponseMap["cert_begin_time"] = responese.CertBeginTime
	}

	if responese.CertEndTime != nil {
		sslResponseMap["cert_end_time"] = responese.CertEndTime
	}

	if responese.ValidityPeriod != nil {
		sslResponseMap["validity_period"] = responese.ValidityPeriod
	}

	if responese.InsertTime != nil {
		sslResponseMap["insert_time"] = responese.InsertTime
	}

	if responese.OrderId != nil {
		sslResponseMap["order_id"] = responese.OrderId
	}

	if responese.CertificateExtra != nil {
		certificateExtraMap := map[string]interface{}{}

		if responese.CertificateExtra.DomainNumber != nil {
			certificateExtraMap["domain_number"] = responese.CertificateExtra.DomainNumber
		}

		if responese.CertificateExtra.OriginCertificateId != nil {
			certificateExtraMap["origin_certificate_id"] = responese.CertificateExtra.OriginCertificateId
		}

		if responese.CertificateExtra.ReplacedBy != nil {
			certificateExtraMap["replaced_by"] = responese.CertificateExtra.ReplacedBy
		}

		if responese.CertificateExtra.ReplacedFor != nil {
			certificateExtraMap["replaced_for"] = responese.CertificateExtra.ReplacedFor
		}

		if responese.CertificateExtra.RenewOrder != nil {
			certificateExtraMap["renew_order"] = responese.CertificateExtra.RenewOrder
		}

		if responese.CertificateExtra.SMCert != nil {
			certificateExtraMap["s_m_cert"] = responese.CertificateExtra.SMCert
		}

		if responese.CertificateExtra.CompanyType != nil {
			certificateExtraMap["company_type"] = responese.CertificateExtra.CompanyType
		}

		sslResponseMap["certificate_extra"] = []interface{}{certificateExtraMap}
	}

	if responese.DvAuthDetail != nil {
		DvAuthDetailMap := map[string]interface{}{}

		if responese.DvAuthDetail.DvAuthKey != nil {
			DvAuthDetailMap["dv_auth_key"] = responese.DvAuthDetail.DvAuthKey
		}

		if responese.DvAuthDetail.DvAuthValue != nil {
			DvAuthDetailMap["dv_auth_value"] = responese.DvAuthDetail.DvAuthValue
		}

		if responese.DvAuthDetail.DvAuthDomain != nil {
			DvAuthDetailMap["dv_auth_domain"] = responese.DvAuthDetail.DvAuthDomain
		}

		if responese.DvAuthDetail.DvAuthPath != nil {
			DvAuthDetailMap["dv_auth_path"] = responese.DvAuthDetail.DvAuthPath
		}

		if responese.DvAuthDetail.DvAuthKeySubDomain != nil {
			DvAuthDetailMap["dv_auth_key_sub_domain"] = responese.DvAuthDetail.DvAuthKeySubDomain
		}

		if responese.DvAuthDetail.DvAuths != nil {
			dvAuthsList := []interface{}{}
			for _, dvAuths := range responese.DvAuthDetail.DvAuths {
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

			DvAuthDetailMap["dv_auths"] = []interface{}{dvAuthsList}
		}

		sslResponseMap["dv_auth_detail"] = []interface{}{DvAuthDetailMap}
	}

	if responese.VulnerabilityReport != nil {
		sslResponseMap["vulnerability_report"] = responese.VulnerabilityReport
	}

	if responese.PackageTypeName != nil {
		sslResponseMap["package_type_name"] = responese.PackageTypeName
	}

	if responese.StatusName != nil {
		sslResponseMap["status_name"] = responese.StatusName
	}

	if responese.SubjectAltName != nil {
		sslResponseMap["subject_alt_name"] = responese.SubjectAltName
	}

	if responese.IsVip != nil {
		sslResponseMap["is_vip"] = responese.IsVip
	}

	if responese.IsWildcard != nil {
		sslResponseMap["is_wildcard"] = responese.IsWildcard
	}

	if responese.IsDv != nil {
		sslResponseMap["is_dv"] = responese.IsDv
	}

	if responese.IsVulnerability != nil {
		sslResponseMap["is_vulnerability"] = responese.IsVulnerability
	}

	if responese.RenewAble != nil {
		sslResponseMap["renew_able"] = responese.RenewAble
	}

	if responese.SubmittedData != nil {
		submittedDataMap := map[string]interface{}{}

		if responese.SubmittedData.CsrType != nil {
			submittedDataMap["csr_type"] = responese.SubmittedData.CsrType
		}

		if responese.SubmittedData.CsrContent != nil {
			submittedDataMap["csr_content"] = responese.SubmittedData.CsrContent
		}

		if responese.SubmittedData.CertificateDomain != nil {
			submittedDataMap["certificate_domain"] = responese.SubmittedData.CertificateDomain
		}

		if responese.SubmittedData.DomainList != nil {
			submittedDataMap["domain_list"] = responese.SubmittedData.DomainList
		}

		if responese.SubmittedData.KeyPassword != nil {
			submittedDataMap["key_password"] = responese.SubmittedData.KeyPassword
		}

		if responese.SubmittedData.OrganizationName != nil {
			submittedDataMap["organization_name"] = responese.SubmittedData.OrganizationName
		}

		if responese.SubmittedData.OrganizationDivision != nil {
			submittedDataMap["organization_division"] = responese.SubmittedData.OrganizationDivision
		}

		if responese.SubmittedData.OrganizationAddress != nil {
			submittedDataMap["organization_address"] = responese.SubmittedData.OrganizationAddress
		}

		if responese.SubmittedData.OrganizationCountry != nil {
			submittedDataMap["organization_country"] = responese.SubmittedData.OrganizationCountry
		}

		if responese.SubmittedData.OrganizationCity != nil {
			submittedDataMap["organization_city"] = responese.SubmittedData.OrganizationCity
		}

		if responese.SubmittedData.OrganizationRegion != nil {
			submittedDataMap["organization_region"] = responese.SubmittedData.OrganizationRegion
		}

		if responese.SubmittedData.PostalCode != nil {
			submittedDataMap["postal_code"] = responese.SubmittedData.PostalCode
		}

		if responese.SubmittedData.PhoneAreaCode != nil {
			submittedDataMap["phone_area_code"] = responese.SubmittedData.PhoneAreaCode
		}

		if responese.SubmittedData.PhoneNumber != nil {
			submittedDataMap["phone_number"] = responese.SubmittedData.PhoneNumber
		}

		if responese.SubmittedData.AdminFirstName != nil {
			submittedDataMap["admin_first_name"] = responese.SubmittedData.AdminFirstName
		}

		if responese.SubmittedData.AdminLastName != nil {
			submittedDataMap["admin_last_name"] = responese.SubmittedData.AdminLastName
		}

		if responese.SubmittedData.AdminPhoneNum != nil {
			submittedDataMap["admin_phone_num"] = responese.SubmittedData.AdminPhoneNum
		}

		if responese.SubmittedData.AdminEmail != nil {
			submittedDataMap["admin_email"] = responese.SubmittedData.AdminEmail
		}

		if responese.SubmittedData.AdminPosition != nil {
			submittedDataMap["admin_position"] = responese.SubmittedData.AdminPosition
		}

		if responese.SubmittedData.ContactFirstName != nil {
			submittedDataMap["contact_first_name"] = responese.SubmittedData.ContactFirstName
		}

		if responese.SubmittedData.ContactLastName != nil {
			submittedDataMap["contact_last_name"] = responese.SubmittedData.ContactLastName
		}

		if responese.SubmittedData.ContactNumber != nil {
			submittedDataMap["contact_number"] = responese.SubmittedData.ContactNumber
		}

		if responese.SubmittedData.ContactEmail != nil {
			submittedDataMap["contact_email"] = responese.SubmittedData.ContactEmail
		}

		if responese.SubmittedData.ContactPosition != nil {
			submittedDataMap["contact_position"] = responese.SubmittedData.ContactPosition
		}

		if responese.SubmittedData.VerifyType != nil {
			submittedDataMap["verify_type"] = responese.SubmittedData.VerifyType
		}

		sslResponseMap["submitted_data"] = []interface{}{submittedDataMap}
	}

	if responese.Deployable != nil {
		sslResponseMap["deployable"] = responese.Deployable
	}

	if responese.CAEncryptAlgorithms != nil {
		sslResponseMap["c_a_encrypt_algorithms"] = responese.CAEncryptAlgorithms
	}

	if responese.CACommonNames != nil {
		sslResponseMap["c_a_common_names"] = responese.CACommonNames
	}

	if responese.CAEndTimes != nil {
		sslResponseMap["c_a_end_times"] = responese.CAEndTimes
	}

	if responese.DvRevokeAuthDetail != nil {
		tmpList := []interface{}{}
		for _, dvAuths := range responese.DvRevokeAuthDetail {
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

			tmpList = append(tmpList, dvAuthsMap)
		}

		sslResponseMap["dv_revoke_auth_detail"] = tmpList
	}
	_ = d.Set("result", []interface{}{sslResponseMap})
	d.SetId(CertificateId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), sslResponseMap); e != nil {
			return e
		}
	}
	return nil
}
