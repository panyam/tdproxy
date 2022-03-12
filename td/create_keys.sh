
# Key considerations for algorithm "RSA" ≥ 2048-bit
echo "Generate private key (.key)"
openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
echo "Generate private key (.key) - Part2"
openssl ecparam -genkey -name secp384r1 -out server.key

echo "Generating self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)"
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
