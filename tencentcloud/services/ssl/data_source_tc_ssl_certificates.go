package ssl

import (
	"context"
	"log"
	"strconv"
	"strings"

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
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain of the SSL certificate to be queried. Matches the primary domain or any domain in the subject alternative names exactly.",
			},
			"sort_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SslCertificatesSortBy),
				Description:  "Field used to sort the returned certificates. Currently only `end_time` is supported. If omitted, the API default order (creation time descending) is used.",
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
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account UIN.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"validity_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Validity period: unit (month).Note: This field may return NULL, indicating that the valid value cannot be obtained.",
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
		domain   *string
		sortBy   *string
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

	if raw, ok := d.GetOk("domain"); ok {
		domain = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("sort_by"); ok {
		sortBy = helper.String(raw.(string))
	}

	sslService := SSLService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
	certificateList, err := GetCertificateList(ctx, sslService, id, name, domain, certType, sortBy)
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
			log.Printf("the certificate is invalid, certificate %v are nil", nilNames)
			continue
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
			"owner_uin":       *certificate.OwnerUin,
			"validity_period": *certificate.ValidityPeriod,
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
			if describeResponse.Response.CertificatePublicKey != nil {
				m["cert"] = *describeResponse.Response.CertificatePublicKey
			}
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

func GetCertificateList(ctx context.Context, sslService SSLService, id, name, domain, certType, sortBy *string) (certificateList []*ssl.Certificates, errRet error) {
	logId := tccommon.GetLogId(ctx)

	// The API only accepts a single fuzzy SearchKey, so every exact filter is
	// resolved with its own request and the results are intersected below.
	searchKeys := make([]*string, 0, 3)
	for _, key := range []*string{id, name, domain} {
		if key != nil {
			searchKeys = append(searchKeys, key)
		}
	}
	if len(searchKeys) == 0 {
		searchKeys = append(searchKeys, nil)
	}

	lists := make([][]*ssl.Certificates, 0, len(searchKeys))
	for _, searchKey := range searchKeys {
		describeRequest := buildDescribeCertificatesRequest(certType, searchKey, sortBy)

		var (
			outErr, inErr error
			list          []*ssl.Certificates
		)
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			list, inErr = sslService.DescribeCertificates(ctx, describeRequest)
			if inErr != nil {
				if sdkErr := helper.UnwarpSDKError(inErr); sdkErr != nil && tccommon.IsContains("LimitExceeded", sdkErr.Code) {
					return resource.RetryableError(inErr)
				}
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read certificates failed, reason: %v", logId, outErr)
			errRet = outErr
			return
		}

		if domain != nil {
			filtered := make([]*ssl.Certificates, 0, len(list))
			for _, certificate := range list {
				if certificateMatchesDomain(certificate, *domain) {
					filtered = append(filtered, certificate)
				}
			}
			list = filtered
		}

		lists = append(lists, list)
	}

	certificateList = GetCommonCertificates(lists...)
	return
}

// buildDescribeCertificatesRequest builds one page request for
// DescribeCertificates. sort_by = "end_time" maps to the API's ExpirationSort;
// DESC returns the certificate expiring last (the newest one after a renewal)
// first.
func buildDescribeCertificatesRequest(certType, searchKey, sortBy *string) *ssl.DescribeCertificatesRequest {
	request := ssl.NewDescribeCertificatesRequest()
	request.CertificateType = certType
	request.SearchKey = searchKey
	if sortBy != nil {
		request.ExpirationSort = helper.String(SSL_CERTIFICATE_EXPIRATION_SORT_DESC)
	}
	return request
}

// certificateMatchesDomain reports whether the certificate covers the given
// domain as its primary domain or as one of its subject alternative names.
func certificateMatchesDomain(certificate *ssl.Certificates, domain string) bool {
	if certificate == nil {
		return false
	}
	if certificate.Domain != nil && strings.EqualFold(*certificate.Domain, domain) {
		return true
	}
	for _, subjectName := range certificate.SubjectAltName {
		if subjectName != nil && strings.EqualFold(*subjectName, domain) {
			return true
		}
	}
	return false
}

// GetCommonCertificates intersects the given certificate lists by certificate
// ID, preserving the order of the first list.
func GetCommonCertificates(lists ...[]*ssl.Certificates) (result []*ssl.Certificates) {
	if len(lists) == 0 {
		return
	}
	if len(lists) == 1 {
		return lists[0]
	}

	others := make([]map[string]struct{}, 0, len(lists)-1)
	for _, list := range lists[1:] {
		ids := make(map[string]struct{}, len(list))
		for _, certificate := range list {
			if certificate != nil && certificate.CertificateId != nil {
				ids[*certificate.CertificateId] = struct{}{}
			}
		}
		others = append(others, ids)
	}

	seen := make(map[string]struct{}, len(lists[0]))
	for _, certificate := range lists[0] {
		if certificate == nil || certificate.CertificateId == nil {
			continue
		}
		if _, ok := seen[*certificate.CertificateId]; ok {
			continue
		}

		matched := true
		for _, ids := range others {
			if _, ok := ids[*certificate.CertificateId]; !ok {
				matched = false
				break
			}
		}
		if !matched {
			continue
		}

		seen[*certificate.CertificateId] = struct{}{}
		result = append(result, certificate)
	}
	return
}
