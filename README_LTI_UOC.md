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
