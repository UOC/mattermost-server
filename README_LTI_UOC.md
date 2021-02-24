# LTI UOC

new file uoc.go
Modified file model/lti.go
Modified file web/lti.go

## How to sync teams

uoc.go `GetTeam` must return the slug of the Team to join. **It must exist**.

*TODO: Should we create them dinamically?*

### Create teams

Teams are created in function `a.CreateTeam(team)` of the file `app/team.go`.  Team struct:

```go
  type Team struct {
    Id                 string  `json:"id"`
    CreateAt           int64   `json:"create_at"`
    UpdateAt           int64   `json:"update_at"`
    DeleteAt           int64   `json:"delete_at"`
    DisplayName        string  `json:"display_name"`
    Name               string  `json:"name"`
    Description        string  `json:"description"`
    Email              string  `json:"email"`
    Type               string  `json:"type"`
    CompanyName        string  `json:"company_name"`
    AllowedDomains     string  `json:"allowed_domains"`
    InviteId           string  `json:"invite_id"`
    AllowOpenInvite    bool    `json:"allow_open_invite"`
    LastTeamIconUpdate int64   `json:"last_team_icon_update,omitempty"`
    SchemeId           *string `json:"scheme_id"`
  }
```

Required fields:

* DisplayName
* Name
* Description
* Type
* AllowOpenInvite

## How to sync channels

`SyncLTIChannels` function changes name of channnel by LTI's data.

## How to sync user

function `OnboardLTIUser` in file `app/lti.go` adds users to channels via `createChannelsIfRequired` and `joinChannelsIfRequired`. They call uoc's `GetPublicChannelsToJoin` and `GetPrivateChannelsToJoin` and create the channel if not exist yet.

## Deploy

More information https://developers.mattermost.com/contribute/server/developer-setup/

Install go:
```
brew install go
```

Update your shell’s initialization script (e.g. .bashrc or .zshrc) and add the following:

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
ulimit -n 8096
```

Re-source your shell’s initialization script to update GOPATH and PATH in your current shell:

```
source $HOME/.bashrc
OR 
source $HOME/.zshrc
```

Start the server
```
make run-server
```

If you get errors install go packages
```
go get github.com/hako/durafmt
go get github.com/dgryski/dgoogauth
go get github.com/go-ldap/ldap
go get github.com/hashicorp/memberlist
go get github.com/prometheus/client_golang/prometheus
go get gopkg.in/olivere/elastic.v5
go get github.com/mattermost/mattermost-server/cmd/mattermost/commands
go get github.com/mattermost/mattermost-server/model/gitlab
go get github.com/mattermost/mattermost-server/imports
go get github.com/mattermost/rsc/qr
go get github.com/tylerb/graceful
go get gopkg.in/hash/maphash
go get github.com/avct/uasurfer
go get github.com/dyatlov/go-opengraph/opengraph
go get github.com/minio/minio-go
go get github.com/minio/minio-go/pkg/credentials
go get github.com/gorilla/schema
go get github.com/NYTimes/gziphandler
go get github.com/mattermost/viper
go get github.com/fsnotify/fsnotify
go get github.com/stretchr/testify/assert
go get github.com/mattermost/gorp
go get github.com/lib/pq
go get github.com/go-redis/redis
go get gopkg.in/gomail.v2
go get github.com/nicksnyder/go-i18n/i18n
go get github.com/disintegration/imaging
go get github.com/golang/freetype
go get github.com/golang/freetype/truetype
go get github.com/gorilla/handlers
go get github.com/disintegration/imaging
go get github.com/golang/freetype
go get github.com/golang/freetype/truetype
go get github.com/gorilla/handlers
go get github.com/rs/cors
go get github.com/rwcarlsen/goexif/exif
go get github.com/segmentio/analytics-go
go get github.com/throttled/throttled
go get github.com/segmentio/analytics-go
go get github.com/icrowley/fake
go get github.com/spf13/cobra
go get go.uber.org/zap
go get go.uber.org/zap/zapcore
go get gopkg.in/natefinch/lumberjack.v2
go get github.com/pborman/uuid
go get github.com/hashicorp/go-hclog
go get github.com/hashicorp/go-plugin
go get github.com/stretchr/objx
go get 
```
