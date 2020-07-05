# tacacs-server
example tacacs server on golang

Seample tacacs server on GoLang
 - build and testing  on 1.14.4

use
 - https://github.com/sevlyar/go-daemon
 - https://github.com/jtblin/go-ldap-client
 - https://github.com/nwaples/tacplus

```
Config cisco 2950 (example for testimg)
Example config  Cisco 2950
....
!
aaa new-model
aaa authentication login default group tacacs+ enable local
aaa authorization exec default group tacacs+ if-authenticated local
aaa authorization commands 1 admin group tacacs+
aaa authorization commands 3 admin group tacacs+
aaa authorization commands 5 admin group tacacs+
aaa authorization commands 10 admin group tacacs+
aaa authorization commands 15 admin group tacacs+
aaa accounting send stop-record authentication failure
aaa accounting exec default start-stop group tacacs+
aaa accounting exec admin start-stop group tacacs+
aaa accounting commands 1 admin start-stop group tacacs+
aaa accounting commands 3 admin start-stop group tacacs+
aaa accounting commands 5 admin start-stop group tacacs+
aaa accounting commands 10 admin start-stop group tacacs+
aaa accounting commands 15 admin start-stop group tacacs+
aaa accounting connection admin start-stop group tacacs+
aaa accounting system default start-stop group tacacs+
............
tacacs-server host 10.1.0.1
tacacs-server timeout 3
tacacs-server key 7 <REMOVE>
banner motd ^C WARNING!!! ^C
!
line vty 0 4
 exec-timeout 3000 0
 authorization commands 1 admin
 authorization commands 3 admin
 authorization commands 5 admin
 authorization commands 10 admin
 authorization commands 15 admin
 accounting connection admin
 accounting commands 1 admin
 accounting commands 3 admin
 accounting commands 5 admin
 accounting commands 10 admin
 accounting commands 15 admin
 accounting exec admin
 transport input telnet
line vty 5 15
 exec-timeout 3000 0
 authorization commands 1 admin
 authorization commands 3 admin
 authorization commands 5 admin
 authorization commands 10 admin
 authorization commands 15 admin
 accounting connection admin
 accounting commands 1 admin
 accounting commands 3 admin
 accounting commands 5 admin
 accounting commands 10 admin
 accounting commands 15 admin
 accounting exec admin
 transport input telnet
!
.....
```
