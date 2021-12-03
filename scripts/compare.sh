#!/usr/bin/env bash
SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

compareService() {
    service_id=$1
    echo "Processing $service_id"
    metadata=$(curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}" https://eve-api.unanet.io/services/${service_id}/metadata | jq .)
    definitions=$(curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}" https://eve-api.unanet.io/services/${service_id}/definitions | jq '. | sort_by(.class)')

    diff <(echo "$metadata") $SCRIPT_DIR/output/metadata/$service_id.json
    retVal=$?
    if [ $retVal -ne 0 ]; then
        echo "THE METADATA FOR SERVICE: $service_id IS NOT THE SAME!"
        exit 1
    fi

    diff <(echo "$definitions") $SCRIPT_DIR/output/definitions/$service_id.json
    retVal=$?
    if [ $retVal -ne 0 ]; then
        echo "THE DEFINITIONS FOR SERVICE: $service_id ARE NOT THE SAME!" 
        exit 1
    fi
 }

if [ -z "$1" ]; then 
    services=$(curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}"  https://eve-api.unanet.io/services | jq '. | sort_by(.id)')
    for row in $(echo "${services}" | jq -r '.[] | @base64'); do
        _jq() {
        echo ${row} | base64 --decode | jq -r ${1}
        }  

        service_id=$(_jq '.id')
        compareService $service_id
    done
else
    compareService $1
fi
