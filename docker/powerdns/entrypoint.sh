#!/usr/bin/env bash

CONFIG_FILE=/etc/powerdns/pdns.conf

sed -i 's/{{PDNS_REMOTE_HTTP_HOST}}/'"${PDNS_REMOTE_HTTP_HOST}"'/g' ${CONFIG_FILE}

if [ "${PDNS_CARBON_SERVER}" ]; then
  echo "carbon-server=${PDNS_CARBON_SERVER}" >> "${CONFIG_FILE}"
fi

if [ "${PDNS_CARBON_OURNAME}" ]; then
  echo "carbon-ourname=${PDNS_CARBON_OURNAME}" >> "${CONFIG_FILE}"
fi

if [ "${PDNS_ZONE_CACHE}" = false ]; then
  echo "zone-cache-refresh-interval=0" >> ${CONFIG_FILE}
fi

exec "$@"