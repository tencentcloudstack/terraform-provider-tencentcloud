package tencentcloud

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
)

func resourceTencentCloudSslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCertificateCreate,
		Read:   resourceTencentCloudSslCertificateRead,
		Delete: resourceTencentCloudSslCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"CA", "SVR"}),
				ForceNew:     true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ForceNew: true,
			},
			"cert": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("cert can't have \\n suffix"))
					}
					return
				},
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				ValidateFunc: func(v interface{}, k string) (wss []string, errs []error) {
					value := v.(string)
					if strings.HasPrefix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n prefix"))
						return
					}

					if strings.HasSuffix(value, "\n") {
						errs = append(errs, errors.New("key can't have \\n suffix"))
					}
					return
				},
			},

			// computed
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
	}
}

func resourceTencentCloudSslCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	name := d.Get("name").(string)
	certType := d.Get("type").(string)
	projectId := d.Get("project_id").(int)
	cert := d.Get("cert").(string)

	var key *string
	if raw, ok := d.GetOk("key"); ok {
		key = stringToPointer(raw.(string))
	}

	service := SslService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateCertificate(ctx, certType, cert, name, projectId, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudSslCertificateRead(d, m)
}

func resourceTencentCloudSslCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := SslService{client: m.(*TencentCloudClient).apiV3Conn}

	certificates, err := service.DescribeCertificates(ctx, &id, nil, nil)
	if err != nil {
		return err
	}

	var certificate *ssl.SSLCertificate
	for _, c := range certificates {
		if c.Id == nil {
			return errors.New("certificate id is nil")
		}

		if *c.Id == id {
			certificate = c
			break
		}
	}

	if certificate == nil {
		d.SetId("")
		return nil
	}

	if certificate.Alias == nil {
		return errors.New("certificate name is nil")
	}
	d.Set("name", certificate.Alias)

	if certificate.CertType == nil {
		return errors.New("certificate type is nil")
	}
	d.Set("type", certificate.CertType)

	if certificate.ProjectId == nil {
		return errors.New("certificate project id is nil")
	}
	projectId, err := strconv.Atoi(*certificate.ProjectId)
	if err != nil {
		return err
	}
	d.Set("project_id", projectId)

	if certificate.Cert == nil {
		return errors.New("certificate cert is nil")
	}
	d.Set("cert", certificate.Cert)

	if certificate.ProductZhName == nil {
		return errors.New("certificate product zh name is nil")
	}
	d.Set("product_zh_name", certificate.ProductZhName)

	if certificate.Domain == nil {
		return errors.New("certificate domain is nil")
	}
	d.Set("domain", certificate.Domain)

	if certificate.Status == nil {
		return errors.New("certificate status is nil")
	}
	d.Set("status", certificate.Status)

	if certificate.CertBeginTime == nil {
		return errors.New("certificate begin time is nil")
	}
	d.Set("begin_time", certificate.CertBeginTime)

	if certificate.CertEndTime == nil {
		return errors.New("certificate end time is nil")
	}
	d.Set("end_time", certificate.CertEndTime)

	if certificate.InsertTime == nil {
		return errors.New("certificate create time is nil")
	}
	d.Set("create_time", certificate.InsertTime)

	subjectAltNames := make([]string, 0, len(certificate.SubjectAltName))
	for _, subjectAltName := range certificate.SubjectAltName {
		subjectAltNames = append(subjectAltNames, *subjectAltName)
	}
	d.Set("subject_names", subjectAltNames)

	return nil
}

func resourceTencentCloudSslCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_certificate.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := SslService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteCertificate(ctx, id)
}
