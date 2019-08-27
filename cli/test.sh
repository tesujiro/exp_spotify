#!/bin/bash

TESTCASE_DIR=./testcase
TESTCASE_EXT="case"
COMMAND="./cli"
PRINT_INFO=true
PRINT_DEBUG=false
TEST_ERROR=false
CASE_OK=0
CASE_FAIL=0

get_case(){
    head -1 "$@"
}

get_want(){
    tail -n +2 "$@"
}

check(){
    :
    # check env vars
    # check reverse proxy running
}

debug(){
    if "$PRINT_DEBUG"
    then
       	if [ -p /dev/stdin ]
	then
	    cat -
	else
	    echo "$@"
       	fi
    fi
}

info(){
    if "$PRINT_INFO"
    then
       	if [ -p /dev/stdin ]
	then
	    cat -
	else
	    echo "$@"
       	fi
    fi
}

main() {
    local number=1
    local EVAL_RC_WANT=0
    for file in $( ls $TESTCASE_DIR/*.$TESTCASE_EXT ) 
    do
	RESULT=$(mktemp)
	eval $(get_case $file) > $RESULT
	EVAL_RC=$?

	diff -s <(cat $RESULT) <(get_want $file) > tempfile.$$
	DIFF_RC=$?
	rm -f $RESULT

	if [ $EVAL_RC == $EVAL_RC_WANT -a $DIFF_RC == 0 ]
	then
	    CASE_OK=$((++CASE_OK))
	    debug ======================================
	    debug test $number : $file
	    get_case $file | debug
	    debug ======================================
	    cat tempfile.$$ | debug
	else
	    CASE_FAIL=$((++CASE_FAIL))
	    TEST_ERROR=true
	    info ======================================
	    info test $number : $file
	    get_case $file | info
	    info ======================================
	    cat tempfile.$$ | info
	fi

	rm tempfile.$$
	number=$((++number))
    done
}

while getopts ":dq" OPT
do
    case $OPT in
       	d) PRINT_DEBUG=true;;
       	q) PRINT_INFO=false;;
       	\?) echo "[ERROR] Undefined options.";exit 1;;
    esac
done

check
main

echo "CASE_OK  : $CASE_OK" | info
echo "CASE_FAIL: $CASE_FAIL" | info

if "$TEST_ERROR"
then
    exit 1
fi
exit 0
