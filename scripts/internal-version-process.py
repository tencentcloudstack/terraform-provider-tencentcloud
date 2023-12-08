import yaml
import re
import json

YAML_FILE = "scripts/internal-version-code.yaml"
EXTENSION_BILLING_FILE = "tencentcloud/extension_billing.go"
SERVICE_TENCENTCLOUD_BILLING_FILE = "tencentcloud/service_tencentcloud_billing.go"
GO_MOD_FILE = "go.mod"
MARK_MATCH = "//internal version: replace (\w+) begin.*?//internal version: replace \w+ end"
MARK_REPLACE = "//internal version: replace %s begin.*?//internal version: replace %s end"

def replace_code(dictionary, code):

    matches = re.finditer(r"%s"%MARK_MATCH, code, flags=re.DOTALL)

    for match in matches:
        key = match.group(1)
        if key in dictionary:
            replacement_code = dictionary[key]
        else:
            replacement_code=""
        mark_str=MARK_REPLACE%(key,key)
        code = re.sub(r"%s"%mark_str, replacement_code, code, flags=re.DOTALL)
    return code

def replace(dictionary):
    for file_name, content in dictionary.items():
        if file_name in [SERVICE_TENCENTCLOUD_BILLING_FILE,EXTENSION_BILLING_FILE]:
            with open(file_name, "w") as file:
                file.write(content["all"])
            continue

        if file_name in GO_MOD_FILE:
            with open(file_name, "a") as file:
                file.write(content)
        with open(file_name, "r") as file:
            code = file.read()

        replaced_code = replace_code(content, code)

        with open(file_name, "w") as file:
            file.write(replaced_code)

    print("Success replace")

def run():

    yaml_file = YAML_FILE
    with open(yaml_file, "r") as f:
        yaml_data = yaml.safe_load(f)
    json_data = json.dumps(yaml_data)

    dictionary = json.loads(json_data)

    replace(dictionary)

run()
