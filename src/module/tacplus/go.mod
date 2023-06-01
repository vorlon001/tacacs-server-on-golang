module tacplus

go 1.20

require module/config v0.0.0-00010101000000-000000000000

require gopkg.in/yaml.v2 v2.4.0 // indirect

replace module/config => ../../../src/module/config
