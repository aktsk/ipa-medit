#!/bin/bash

CERT="ipa-medit-codesign"

function error() {
    echo error: "$@"
    exit 1
}

function cleanup {
    # Remove generated files
    rm -f "/tmp/$CERT.tmpl" "/tmp/$CERT.cer" "/tmp/$CERT.key" > /dev/null 2>&1
}

function sign {
    codesign -f -s ipa-medit-codesign --entitlements ./scripts/entitlements.plist ipa-medit
}

trap cleanup EXIT

# Check if the certificate is already present in the system keychain
security find-certificate -Z -p -c "$CERT" /Library/Keychains/System.keychain > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo Certificate has already been generated and installed

    # sign the binary
    sign
    exit 0
fi

# Create the certificate template
cat <<EOF >/tmp/$CERT.tmpl
[ req ]
default_bits       = 2048        # RSA key size
encrypt_key        = no          # Protect private key
default_md         = sha512      # MD to use
prompt             = no          # Prompt for DN
distinguished_name = codesign_dn # DN template
[ codesign_dn ]
commonName         = "$CERT"
[ codesign_reqext ]
keyUsage           = critical,digitalSignature
extendedKeyUsage   = critical,codeSigning
EOF

echo Generating and installing ipa-medit-codesign certificate

# Generate a new certificate
openssl req -new -newkey rsa:2048 -x509 -days 3650 -nodes -config "/tmp/$CERT.tmpl" -extensions codesign_reqext -batch -out "/tmp/$CERT.cer" -keyout "/tmp/$CERT.key" > /dev/null 2>&1
[ $? -eq 0 ] || error Something went wrong when generating the certificate

# Install the certificate in the system keychain
sudo security add-trusted-cert -d -r trustRoot -p codeSign -k /Library/Keychains/System.keychain "/tmp/$CERT.cer" > /dev/null 2>&1
[ $? -eq 0 ] || error Something went wrong when installing the certificate

# Install the key for the certificate in the system keychain
sudo security import "/tmp/$CERT.key" -A -k /Library/Keychains/System.keychain > /dev/null 2>&1
[ $? -eq 0 ] || error Something went wrong when installing the key

# Kill task_for_pid access control daemon
sudo pkill -f /usr/libexec/taskgated > /dev/null 2>&1

# sign the binary
sign

# Exit indicating the certificate is now generated and installed, signed
exit 0
