#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/photoz"

echo "==== Releasing photoz ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm photoz
echo "  Done!"

echo "  Deleting existing code..."
ssh root@photoz.reidc.io "rm -rf /root/go/src/photoz"
echo "  Code deleted successfully!"

echo "  Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@photoz.reidc.io:/root/go/src/photoz/
echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@photoz.reidc.io "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

echo "  Building the code on remote server..."
ssh root@photoz.reidc.io 'export GOPATH=/root/go; cd /root/app; /usr/local/go/bin/go build -o ./server $GOPATH/src/photoz/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@photoz.reidc.io "cd /root/app; cp -R /root/go/src/photoz/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@photoz.reidc.io "cd /root/app; cp -R /root/go/src/photoz/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@photoz.reidc.io "cd /root/app; cp /root/go/src/photoz/Caddyfile ."
echo "  Views moved successfully!"

echo "  Restarting the server..."
ssh root@photoz.reidc.io "sudo service photoz restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@photoz.reidc.io "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing photoz ===="