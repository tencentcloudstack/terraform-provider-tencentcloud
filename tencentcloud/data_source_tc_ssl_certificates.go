package tencentcloud

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSslCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslCertificatesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// computed
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_zh_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_names": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSslCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_certificates.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var (
		name     *string
		certType *string
		id       *string
	)

	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}

	if raw, ok := d.GetOk("type"); ok {
		certType = stringToPointer(raw.(string))
	}

	if raw, ok := d.GetOk("id"); ok {
		id = stringToPointer(raw.(string))
	}

	service := SslService{client: m.(*TencentCloudClient).apiV3Conn}

	respCertificates, err := service.DescribeCertificates(ctx, id, name, certType)
	if err != nil {
		return err
	}

	certificates := make([]map[string]interface{}, 0, len(respCertificates))
	ids := make([]string, 0, len(respCertificates))

	for _, certificate := range respCertificates {
		if certificate.Id == nil {
			return errors.New("certificate id is nil")
		}

		ids = append(ids, *certificate.Id)

		if certificate.Alias == nil {
			return errors.New("certificate name is nil")
		}

		if certificate.CertType == nil {
			return errors.New("certificate type is nil")
		}

		if certificate.ProjectId == nil {
			return errors.New("certificate project id is nil")
		}
		projectId, err := strconv.Atoi(*certificate.ProjectId)
		if err != nil {
			return err
		}

		if certificate.Cert == nil {
			return errors.New("certificate cert is nil")
		}

		if certificate.ProductZhName == nil {
			return errors.New("certificate product zh name is nil")
		}

		if certificate.Domain == nil {
			return errors.New("certificate domain is nil")
		}

		if certificate.Status == nil {
			return errors.New("certificate status is nil")
		}

		if certificate.CertBeginTime == nil {
			return errors.New("certificate begin time is nil")
		}

		if certificate.CertEndTime == nil {
			return errors.New("certificate end time is nil")
		}

		if certificate.InsertTime == nil {
			return errors.New("certificate create time is nil")
		}

		m := map[string]interface{}{
			"id":              *certificate.Id,
			"name":            *certificate.Alias,
			"type":            *certificate.CertType,
			"project_id":      projectId,
			"cert":            *certificate.Cert,
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

		certificates = append(certificates, m)
	}

	d.Set("certificates", certificates)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
