module metrics

go 1.21

replace github.com/jmpsec/osctrl/types => ../types

replace github.com/jmpsec/osctrl/queries => ../queries

replace github.com/jmpsec/osctrl/nodes => ../nodes

replace github.com/jmpsec/osctrl/utils => ../utils

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.8.4
	gorm.io/driver/postgres v1.4.6
	gorm.io/gorm v1.25.5
)

require (
	github.com/google/uuid v1.5.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jmpsec/osctrl/nodes v0.3.5 // indirect
	github.com/jmpsec/osctrl/utils v0.0.0-20231221224851-87ae6d485538 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
	golang.org/x/sync v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmpsec/osctrl/queries v0.3.5 // indirect
	github.com/jmpsec/osctrl/types v0.0.0-20231221224851-87ae6d485538
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
