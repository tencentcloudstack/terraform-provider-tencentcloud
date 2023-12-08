import yaml
import re

YAML_FILE = "scripts/internal-version-code.yaml"
EXTENSION_BILLING_FILE = "tencentcloud/extension_billing.go"
SERVICE_TENCENTCLOUD_BILLING_FILE = "tencentcloud/service_tencentcloud_billing.go"
GO_MOD_FILE = "go.mod"
MARK_BEGIN = "//internal version: replace {} begin.*?//internal version: replace {} end"

def replace_code(dictionary, code):
    for key, value in dictionary.items():
        mark_str = MARK_BEGIN.format(key, key)
        if key in code:
            code = re.sub(mark_str, value, code, flags=re.DOTALL)
    return code

def replace(dictionary):
    for file_name, content in dictionary.items():
        if file_name in [EXTENSION_BILLING_FILE, SERVICE_TENCENTCLOUD_BILLING_FILE]:
            with open(file_name, "w") as file:
                file.write(content["all"])
            continue

        if file_name == GO_MOD_FILE:
            with open(file_name, "a") as file:
                file.write(content)
            continue

        with open(file_name, "r") as file:
            code = file.read()

        replaced_code = replace_code(content, code)

        with open(file_name, "w") as file:
            file.write(replaced_code)

    print("Success replace")

def run():
    with open(YAML_FILE, "r") as f:
        yaml_data = yaml.safe_load(f)

    dictionary = yaml_data

    replace(dictionary)

run()