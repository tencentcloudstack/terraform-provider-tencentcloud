package ssl

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslCertificatesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SSL certificate to be queried.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the SSL certificate to be queried. Available values includes: `CA` and `SVR`.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the SSL certificate to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of certificate. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SSL certificate.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SSL certificate.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the SSL certificate.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID of the SSL certificate.",
						},
						"cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Content of the SSL certificate.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key of the SSL certificate.",
						},
						"product_zh_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate authority.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary domain of the SSL certificate.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the SSL certificate.",
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Beginning time of the SSL certificate.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ending time of the SSL certificate.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the SSL certificate.",
						},
						"subject_names": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "ALL domains included in the SSL certificate. Including the primary domain name.",
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Order ID returned.",
						},
						"dv_auths": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DV certification information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dv_auth_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication key.",
									},
									"dv_auth_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication value.",
									},
									"dv_auth_verify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication type.",
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

func dataSourceTencentCloudSslCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_certificates.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		name     *string
		certType *string
		id       *string
	)

	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("type"); ok {
		certType = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("id"); ok {
		id = helper.String(raw.(string))
	}

	sslService := SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
	certificateList, err := GetCertificateList(ctx, sslService, id, name, certType)
	if err != nil {
		return err
	}

	certificates := make([]map[string]interface{}, 0, len(certificateList))
	ids := make([]string, 0, len(certificateList))
	for _, certificate := range certificateList {
		if nilNames := tccommon.CheckNil(certificate, map[string]string{
			"CertificateId":   "id",
			"Alias":           "name",
			"CertificateType": "type",
			"ProjectId":       "project id",
			"ProductZhName":   "product zh name",
			"Domain":          "domain",
			"Status":          "status",
			"CertBeginTime":   "begin time",
			"CertEndTime":     "end time",
			"InsertTime":      "create time",
		}); len(nilNames) > 0 {
			return fmt.Errorf("certificate %v are nil", nilNames)
		}

		ids = append(ids, *certificate.CertificateId)

		projectId, err := strconv.Atoi(*certificate.ProjectId)
		if err != nil {
			return err
		}

		m := map[string]interface{}{
			"id":              *certificate.CertificateId,
			"name":            *certificate.Alias,
			"type":            *certificate.CertificateType,
			"project_id":      projectId,
			"product_zh_name": *certificate.ProductZhName,
			"domain":          *certificate.Domain,
			"status":          *certificate.Status,
			"begin_time":      *certificate.CertBeginTime,
			"end_time":        *certificate.CertEndTime,
			"create_time":     *certificate.InsertTime,
		}

		if len(certificate.SubjectAltName) > 0 {
			subjectAltNames := make([]string, 0, len(certificate.SubjectAltName))
			for _, name := range certificate.SubjectAltName {
				subjectAltNames = append(subjectAltNames, *name)
			}
			m["subject_names"] = subjectAltNames
		}

		describeRequest := ssl.NewDescribeCertificateDetailRequest()
		describeRequest.CertificateId = certificate.CertificateId
		var outErr, inErr error
		var describeResponse *ssl.DescribeCertificateDetailResponse
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			describeResponse, inErr = sslService.DescribeCertificateDetail(ctx, describeRequest)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read certificate failed, reason: %v", logId, outErr)
			return outErr
		}

		if describeResponse != nil && describeResponse.Response != nil {
			m["cert"] = *describeResponse.Response.CertificatePublicKey
			if describeResponse.Response.CertificatePrivateKey != nil {
				m["key"] = *describeResponse.Response.CertificatePrivateKey
			}

			if describeResponse.Response.OrderId != nil {
				m["order_id"] = *describeResponse.Response.OrderId
			}
			if describeResponse.Response.DvAuthDetail != nil && len(describeResponse.Response.DvAuthDetail.DvAuths) != 0 {
				dvAuths := make([]map[string]string, 0)
				for _, item := range describeResponse.Response.DvAuthDetail.DvAuths {
					dvAuth := make(map[string]string)
					dvAuth["dv_auth_key"] = *item.DvAuthKey
					dvAuth["dv_auth_value"] = *item.DvAuthValue
					dvAuth["dv_auth_verify_type"] = *item.DvAuthVerifyType
					dvAuths = append(dvAuths, dvAuth)
				}

				m["dv_auths"] = dvAuths
			}
		}

		certificates = append(certificates, m)
	}

	_ = d.Set("certificates", certificates)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), certificates); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}

func GetCertificateList(ctx context.Context, sslService SSLService, id, name, certType *string) (certificateList []*ssl.Certificates, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		outErr, inErr                        error
		certificatesById, certificatesByName []*ssl.Certificates
	)

	if id == nil && name == nil {
		describeRequest := ssl.NewDescribeCertificatesRequest()
		describeRequest.CertificateType = certType
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			certificateList, inErr = sslService.DescribeCertificates(ctx, describeRequest)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read certificates failed, reason: %v", logId, outErr)
			errRet = outErr
			return
		}
		return
	}

	if id != nil {
		describeRequest := ssl.NewDescribeCertificatesRequest()
		describeRequest.CertificateType = certType
		describeRequest.SearchKey = id
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			certificatesById, inErr = sslService.DescribeCertificates(ctx, describeRequest)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read certificates failed, reason: %v", logId, outErr)
			errRet = outErr
			return
		}
	}
	if name != nil {
		describeRequest := ssl.NewDescribeCertificatesRequest()
		describeRequest.CertificateType = certType
		describeRequest.SearchKey = name
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			certificatesByName, inErr = sslService.DescribeCertificates(ctx, describeRequest)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read certificates failed, reason: %v", logId, outErr)
			errRet = outErr
			return
		}
	}

	certificateList = GetCommonCertificates(certificatesById, certificatesByName)
	return
}

func GetCommonCertificates(certificatesById, certificatesByName []*ssl.Certificates) (result []*ssl.Certificates) {
	if len(certificatesById) == 0 {
		return certificatesByName
	} else if len(certificatesByName) == 0 {
		return certificatesById
	}
	certificateMap := make(map[string]bool)
	for _, certificate := range certificatesById {
		if _, ok := certificateMap[*certificate.CertificateId]; ok {
			continue
		}
		certificateMap[*certificate.CertificateId] = true
	}

	for _, certificate := range certificatesByName {
		if _, ok := certificateMap[*certificate.CertificateId]; ok {
			result = append(result, certificate)
		}
	}
	return
}
