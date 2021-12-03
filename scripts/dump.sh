#!/usr/bin/env bash

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

rm -rf $SCRIPT_DIR/output/metadata
rm -rf $SCRIPT_DIR/output/definitions
mkdir -p $SCRIPT_DIR/output/metadata
mkdir -p $SCRIPT_DIR/output/definitions

services=$(curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}"  https://eve-api.unanet.io/services | jq '. | sort_by(.id)')
for row in $(echo "${services}" | jq -r '.[] | @base64'); do
    _jq() {
     echo ${row} | base64 --decode | jq -r ${1}
    }

    service_id=$(_jq '.id')
    echo "Processing $service_id"
    curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}" https://eve-api.unanet.io/services/${service_id}/metadata | jq . > $SCRIPT_DIR/output/metadata/${service_id}.json
    curl -s -H "Authorization: Bearer ${EVE_ADMIN_TOKEN}" https://eve-api.unanet.io/services/${service_id}/definitions | jq '. | sort_by(.class)' > $SCRIPT_DIR/output/definitions/${service_id}.json
done


