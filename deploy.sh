#!/bin/bash

# build the application
go build -o app ./cmd

# copy to server (assuming it is configured for private/public key authentication in ~/.ssh/config as "do" hostname)
scp ./app do:/opt/backend_app/app_new

ssh do << EOF
    systemctl stop goapp

    chown root:goapp /opt/backend_app/app_new
    chmod 750 /opt/backend_app/app_new

    mv /opt/backend_app/app /opt/backend_app/app_old
    mv /opt/backend_app/app_new /opt/backend_app/app

    systemctl start goapp
EOF

# curl lekkoatletawka.pl/api/healthz and check if the response is 200
# curl lekkoatletawka.pl/api/healthz and check if the response is 200
response_code=$(curl -s -o /dev/null -w "%{http_code}" https://lekkoatletawka.pl/api/readyz)

if [ "$response_code" -eq 200 ]; then
    echo "Deployment successful"
else
    echo "Deployment failed with response code $response_code"
fi
