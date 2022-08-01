#!/usr/bin/env bash

openssl genpkey -outform PEM -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private.key

openssl req -new -nodes -key private.key -config csrconfig.txt -nameopt utf8 -utf8 -out server.csr

openssl req -x509 -nodes -in server.csr -days 365 -key private.key \
-config csrconfig.txt -extensions req_ext -nameopt utf8 -utf8 -out server.crt



