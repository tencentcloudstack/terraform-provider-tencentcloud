package artifact

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var version = flag.String("version", "", "specify provider version")
var source = flag.String("source", "tencentcloudstack/tencentcloud", "specify provider source")
var testDir = "../examples/tencentcloud-user-info"

type ProviderMeta struct {
	Version string
	Source  string
}

func TestTerraformBasicExample(t *testing.T) {
	t.Parallel()

	log.Printf("source: %s, version: %s", *source, *version)

	providerText, err := template.ParseFiles("provider.tmpl")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fmt.Sprintf("%s/provider.tf", testDir))
	if err != nil {
		log.Println("create file error: ", err)
		return
	}

	err = providerText.Execute(f, ProviderMeta{
		Source:  *source,
		Version: *version,
	})

	if err != nil {
		panic(err)
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: testDir,

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	})

	// website::tag::4::Clean up resources with "terraform destroy". Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer func() {
		terraform.Destroy(t, terraformOptions)
		if err := os.Remove("../examples/tencentcloud-vpc/provider.tf"); err != nil {
			panic(err)
		}
	}()

	// website::tag::2::Run "terraform init" and "terraform apply".
	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	actualTextExample := terraform.Output(t, terraformOptions, "output")

	// website::tag::3::Check the output against expected values.
	// Verify we're getting back the outputs we expect
	assert.NotEqual(t, "", actualTextExample)
}
