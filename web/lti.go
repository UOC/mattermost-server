// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package web

import (
//	"fmt"
//	b64 "encoding/base64"
	"net/http"
//  "strings"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

func (w *Web) InitLti() {
	mlog.Info("Initializing web LTI")
	w.MainRouter.Handle("/login/lti", w.NewHandler(loginWithLti)).Methods("GET")
}

func loginWithLti(c *Context, w http.ResponseWriter, r *http.Request) {

	mlog.Info("Logging in using LTI")

	// If there is a user currently logged in, log them out

	// Check if the LTI user exists, if not create them, then log them in

	email := "mike@mailinator.com"
	var user *model.User
	var err *model.AppError

	if user, err = c.App.GetUserByEmail(email); err != nil {
		c.Err = err
		return
	}
	

	//w.Header().Set("Content-Type", "text/html")
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s = %s</div>", "LTI login test", email, user.Id)

	var session *model.Session
	var deviceId  string // no deviceId
	session, err = c.App.DoLogin(w, r, user, deviceId)
	if err != nil {
		c.Err = err
		return
	}

	c.Session = *session

	http.Redirect(w, r, c.GetSiteURLHeader(), http.StatusFound)

	/* From api4/user.go/login() - we need to do much of this here as well -mjl

	c.LogAuditWithUserId(id, "attempt - login_id="+loginId)
	user, err := c.App.AuthenticateUserForLogin(id, loginId, password, mfaToken, ldapOnly)

	if err != nil {
		c.LogAuditWithUserId(id, "failure - login_id="+loginId)
		c.Err = err
		return
	}

	c.LogAuditWithUserId(user.Id, "authenticated")

	var session *model.Session
	session, err = c.App.DoLogin(w, r, user, deviceId)
	if err != nil {
		c.Err = err
		return
	}

	c.LogAuditWithUserId(user.Id, "success")

	c.Session = *session

	user.Sanitize(map[string]bool{})

	w.Write([]byte(user.ToJson()))
	*/




	/*
	samlInterface := c.App.Saml

	if samlInterface == nil {
		c.Err = model.NewAppError("completeLti", "api.user.saml.not_available.app_error", nil, "", http.StatusFound)
		return
	}

	//Validate that the user is with SAML and all that
	encodedXML := r.FormValue("SAMLResponse")
	relayState := r.FormValue("RelayState")

	relayProps := make(map[string]string)
	if len(relayState) > 0 {
		stateStr := ""
		if b, err := b64.StdEncoding.DecodeString(relayState); err != nil {
			c.Err = model.NewAppError("completeLti", "api.user.authorize_oauth_user.invalid_state.app_error", nil, err.Error(), http.StatusFound)
			return
		} else {
			stateStr = string(b)
		}
		relayProps = model.MapFromJson(strings.NewReader(stateStr))
	}

	action := relayProps["action"]
	if user, err := samlInterface.DoLogin(encodedXML, relayProps); err != nil {
		if action == model.OAUTH_ACTION_MOBILE {
			err.Translate(c.T)
			w.Write([]byte(err.ToJson()))
		} else {
			c.Err = err
			c.Err.StatusCode = http.StatusFound
		}
		return
	} else {
		if err := c.App.CheckUserAllAuthenticationCriteria(user, ""); err != nil {
			c.Err = err
			c.Err.StatusCode = http.StatusFound
			return
		}

		switch action {
		case model.OAUTH_ACTION_SIGNUP:
			teamId := relayProps["team_id"]
			if len(teamId) > 0 {
				c.App.Go(func() {
					if err := c.App.AddUserToTeamByTeamId(teamId, user); err != nil {
						mlog.Error(err.Error())
					} else {
						c.App.AddDirectChannels(teamId, user)
					}
				})
			}
		case model.OAUTH_ACTION_EMAIL_TO_SSO:
			if err := c.App.RevokeAllSessions(user.Id); err != nil {
				c.Err = err
				return
			}
			c.LogAuditWithUserId(user.Id, "Revoked all sessions for user")
			c.App.Go(func() {
				if err := c.App.SendSignInChangeEmail(user.Email, strings.Title(model.USER_AUTH_SERVICE_SAML)+" SSO", user.Locale, c.App.GetSiteURL()); err != nil {
					mlog.Error(err.Error())
				}
			})
		}

		session, err := c.App.DoLogin(w, r, user, "")
		if err != nil {
			c.Err = err
			return
		}

		c.Session = *session

		if val, ok := relayProps["redirect_to"]; ok {
			http.Redirect(w, r, c.GetSiteURLHeader()+val, http.StatusFound)
			return
		}

		switch action {
		case model.OAUTH_ACTION_MOBILE:
			ReturnStatusOK(w)
		case model.OAUTH_ACTION_CLIENT:
			err = c.App.SendMessageToExtension(w, relayProps["extension_id"], c.Session.Token)

			if err != nil {
				c.Err = err
				return
			}
		case model.OAUTH_ACTION_EMAIL_TO_SSO:
			http.Redirect(w, r, c.GetSiteURLHeader()+"/login?extra=signin_change", http.StatusFound)
		default:
			http.Redirect(w, r, c.GetSiteURLHeader(), http.StatusFound)
		}
	}
	*/
}
