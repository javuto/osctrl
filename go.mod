module osctrl

go 1.17

replace github.com/jmpsec/osctrl/admin/auth => ./admin/auth

replace github.com/jmpsec/osctrl/admin/handlers => ./admin/handlers

replace github.com/jmpsec/osctrl/admin/sessions => ./admin/sessions

replace github.com/jmpsec/osctrl/api/handlers => ./api/handlers

replace github.com/jmpsec/osctrl/backend => ./backend

replace github.com/jmpsec/osctrl/carves => ./carves

replace github.com/jmpsec/osctrl/environments => ./environments

replace github.com/jmpsec/osctrl/logging => ./logging

replace github.com/jmpsec/osctrl/metrics => ./metrics

replace github.com/jmpsec/osctrl/nodes => ./nodes

replace github.com/jmpsec/osctrl/queries => ./queries

replace github.com/jmpsec/osctrl/settings => ./settings

replace github.com/jmpsec/osctrl/tags => ./tags

replace github.com/jmpsec/osctrl/tls/handlers => ./tls/handlers

replace github.com/jmpsec/osctrl/types => ./types

replace github.com/jmpsec/osctrl/users => ./users

replace github.com/jmpsec/osctrl/utils => ./utils

replace github.com/jmpsec/osctrl/version => ./version

require (
	github.com/crewjam/saml v0.4.6
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/gorilla/mux v1.8.0
	github.com/jmpsec/osctrl/admin/handlers v0.2.7
	github.com/jmpsec/osctrl/admin/sessions v0.2.7
	github.com/jmpsec/osctrl/backend v0.2.7
	github.com/jmpsec/osctrl/carves v0.2.7
	github.com/jmpsec/osctrl/environments v0.2.7
	github.com/jmpsec/osctrl/logging v0.2.7
	github.com/jmpsec/osctrl/metrics v0.2.7
	github.com/jmpsec/osctrl/nodes v0.2.7
	github.com/jmpsec/osctrl/queries v0.2.7
	github.com/jmpsec/osctrl/settings v0.2.7
	github.com/jmpsec/osctrl/tags v0.2.7
	github.com/jmpsec/osctrl/tls/handlers v0.2.7
	github.com/jmpsec/osctrl/types v0.2.7
	github.com/jmpsec/osctrl/users v0.2.7
	github.com/jmpsec/osctrl/utils v0.2.7
	github.com/jmpsec/osctrl/version v0.2.7
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/viper v1.10.1
	github.com/urfave/cli/v2 v2.3.0
	gorm.io/gorm v1.22.5
)

require (
	github.com/aws/aws-sdk-go v1.42.44 // indirect
	github.com/beevik/etree v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/crewjam/httperr v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.9.0 // indirect
	github.com/jackc/pgx/v4 v4.14.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/russellhaering/goxmldsig v1.1.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.66.3 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/postgres v1.2.3 // indirect
)
