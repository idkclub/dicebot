#!/bin/bash
echo "Copy the following items from https://api.slack.com/applications"
read -p "Client ID: " id
read -p "Client Secret: " secret

cat > secrets.go <<EOL
package hackyslack

const (
	clientId     = "${id}"
	clientSecret = "${secret}"
)
EOL
echo "secrets.go generated."
