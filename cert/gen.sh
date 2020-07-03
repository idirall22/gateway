# rm -rf *.pem

# Generate CA private key and self signed key
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=DZ/ST=idir/L=Oran/O=idir/OU=ir/CN=FQDN/emailAddress=idirall22@gmail.com"
# Print CA cert
openssl x509 -in ca-cert.pem -noout -text
#  Generate server key and request.
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=DZ/ST=Oran/L=Oran/O=idir/OU=gateway/CN=FQDN/emailAddress=idirall22@gmail.com"
# Sign server request and generate server cert
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

openssl x509 -in server-cert.pem -noout -text