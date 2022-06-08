#!/bin/bash

#service files
update_services=`git diff --name-status origin/master| awk '{print $2}' | grep "^tencentcloud/" | xargs egrep "type .+Service struct" | awk '{print $2}'`
need_test_files=""
for update_service in $update_services; do
    tmp_files=`grep -r --with-filename $update_service ./tencentcloud | awk -F ":" '{print $1}' | grep -v "service_tencent*" | awk -F "/" '{print $3}' | sort | uniq | egrep "^resource_tc_|^data_source_tc" | awk -F "." '{print $1}' | awk '/_test$/{print "tencentcloud/"$0".go"} !/_test$/{print "tencentcloud/"$0"_test.go"}'`
    need_test_files="$need_test_files $tmp_files"
done

# resource&&data_source files
update_sources=`git diff --name-status origin/master| awk '{print $2}' | egrep "^tencentcloud/resource_tc|^tencentcloud/data_source" | egrep -v "_test.go" | awk -F "." '{print $1"_test.go"}'`
# test files
delta_test_files=`git diff --name-status origin/master | egrep "_test\.go$" | awk '{print $2}'`
# all test files
delta_test_files="$delta_test_files $need_test_files $update_sources"
delta_test_files=`echo $delta_test_files | xargs -n1 | sort | uniq`
for delta_test_file in ${delta_test_files}; do
    test_casts=`egrep "func TestAcc.+\(" ${delta_test_file} | awk -F "(" '{print $1}' | awk '{print $2}'`
    echo "[$delta_test_file]: $test_casts"
    for test_cast in ${test_casts}; do
        go_test_cmd="go test -v -run ${test_cast} -timeout=0 ./tencentcloud/"
        echo ${go_test_cmd}
        $go_test_cmd
        if [ $? -ne 0 ]; then
            printf "[GO TEST FILED] ${go_test_cmd}"
            exit 1
        fi
    done
done