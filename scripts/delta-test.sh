#!/bin/bash

pr_id=${PR_ID}

new_resources=`cat .changelog/${pr_id}.txt| grep -Poz "(?<=release-note:new-resource\n)\w+" | awk '{print "resource/"$0}'`
new_data_sources=`cat .changelog/${pr_id}.txt| grep -Poz "(?<=release-note:new-data-source\n)\w+" | awk '{print "datasource/"$0}'`
source_names=`cat .changelog/${pr_id}.txt| grep -E "^(resource|datasource)\/(\w+)" | awk -F ":" '{print $1}'`
source_names="$source_names $new_resources $new_data_sources"
source_names=`echo $source_names | xargs -n1 | sort | uniq`
test_files=""
for source_name in $source_names; do
    name=${source_name#*/}
    type=${source_name%/*}
    if [ $type == "datasource" ]; then
        type=dataSource
    fi 
    # echo $source_name $type $name
    function_name=$(cat tencentcloud/provider.go | grep "\"${name}\"" | grep "${type}")
    function_name=${function_name#*:}
    function_name=$(echo $(echo ${function_name%,*}))

    test_file=$(grep -r "func $function_name \*schema\.Resource" tencentcloud)
    test_file=${test_file#*/}
    test_file=${test_file%:*}
    test_files="$test_files $test_file"
done
echo "test files:" $test_files

for test_file in $test_files; do
    test_case_type=${test_file%_tc*}
    test_case_name=${test_file#*tc_}
    test_case_name=${test_case_name%.*}

    test_case_type=`echo $test_case_type | sed -r 's/(^|_)(\w)/\U\2/g'`
    test_case_name=`echo $test_case_name | sed -r 's/(^|_)(\w)/\U\2/g'`
   
    go_test_cmd="go test -v -run TestAccTencentCloud${test_case_name}${test_case_type} -timeout=0 ./tencentcloud/"
    echo $go_test_cmd
    $go_test_cmd
    if [ $? -ne 0 ]; then
        printf "[GO TEST FILED] ${go_test_cmd}"
        exit 1
    fi
done
