module main

go 1.20

require (
	ldap v0.0.0-00010101000000-000000000000 // indirect
	module/config v0.0.0-00010101000000-000000000000
	module/go-daemon v0.0.0-00010101000000-000000000000
	module/ldap-client v0.0.0-00010101000000-000000000000
	module/tacplus v0.0.0-00010101000000-000000000000
)

replace (
	github.com/go-ldap/ldap/v3/gssapi => ./src/ldap/v3/gssapi
	ldap => ./src/ldap/v3
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-ldap/ldap/v3 v3.4.4 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	ldap/v3 => ./src/ldap/v3
	module/config => ./src/module/config
	module/go-daemon => ./src/module/go-daemon
	module/ldap-client => ./src/module/ldap-client
	module/tacplus => ./src/module/tacplus
)
