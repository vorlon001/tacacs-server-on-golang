module main

go 1.20

require ldap v0.0.0-00010101000000-000000000000

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-ldap/ldap/v3/gssapi v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.7.0 // indirect
)

replace (
	github.com/go-ldap/ldap/v3/gssapi => ../../../src/ldap/v3/gssapi
	ldap => ../../../src/ldap/v3
)
