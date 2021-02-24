// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/base64"
	"net/http"
	"strings"
)

const (
	uocLaunchDataEmailKey               = "lis_person_contact_email_primary"
	uocLaunchDataUsernameKey            = "custom_username" //"lis_person_sourcedid"
	uocLaunchDataFirstNameKey           = "lis_person_name_given"
	uocLaunchDataLastNameKey            = "lis_person_name_family"
	uocLaunchDataPositionKey            = "roles"
	uocLaunchDataLTIUserIdKey           = "lis_person_sourcedid"
	uocLaunchDataChannelRedirectKey     = "custom_channel_redirect"
	uocLaunchDataFullNameKey            = "lis_person_name_full"
	uocLaunchDataContextId              = "context_id"
	uocLaunchDataContextTitle           = "context_title"
	uocLaunchDataPresentationLocale     = "launch_presentation_locale"
	uocLaunchDataTeamIsTransversalParam = "custom_transversal_team"
	uocLaunchDataTransversalTeamParam   = "custom_domain_coditercers"
	uocBase64Encoded                    = "custom_base64Encoded"

	uocRedirectChannelLookupKeyword = "lookup"
)

type UocChannel struct {
	IdProperty   string
	NameProperty string
}

type UocPersonalChannels struct {
	Type        string
	ChannelList map[string]UocChannel
}

type UocDefaultChannel struct {
	Name        string
	DisplayName string
}

type UocLMS struct {
	Name                string
	Type                string
	OAuthConsumerKey    string
	OAuthConsumerSecret string
	Teams               map[string]string

	PersonalChannels UocPersonalChannels
	DefaultChannels  map[string]UocDefaultChannel
}

func returnValueBase64Encoded(value string, launchData map[string]string) string {

	if val, ok := launchData[uocBase64Encoded]; ok {
		if val == "1" {
			sDec, _ := base64.StdEncoding.DecodeString(value)
			return string(value)
		}
	}

	return value
}

func (e *UocLMS) GetEmail(launchData map[string]string) string {
	return returnValueBase64Encoded(launchData[uocLaunchDataEmailKey], launchData)
}

func (e *UocLMS) GetName() string {
	return e.Name
}

func (e *UocLMS) GetType() string {
	return e.Type
}

func (e *UocLMS) GetOAuthConsumerKey() string {
	return e.OAuthConsumerKey
}

func (e *UocLMS) GetOAuthConsumerSecret() string {
	return e.OAuthConsumerSecret
}

func (e *UocLMS) GetUserId(launchData map[string]string) string {
	return returnValueBase64Encoded(launchData[uocLaunchDataLTIUserIdKey], launchData)
}

func (e *UocLMS) ValidateLTIRequest(url string, request *http.Request) bool {
	return baseValidateLTIRequest(e.OAuthConsumerSecret, e.OAuthConsumerKey, url, request)
}

func (e *UocLMS) BuildUser(launchData map[string]string, password string) (*User, *AppError) {
	//checking if all required fields are present

	if launchData[uocLaunchDataEmailKey] == "" {
		return nil, NewAppError("Uoc_BuildUser", "Uoc.build_user.email_missing", nil, "", http.StatusBadRequest)
	}

	if launchData[uocLaunchDataUsernameKey] == "" {
		return nil, NewAppError("Uoc_BuildUser", "Uoc.build_user.username_missing", nil, "", http.StatusBadRequest)
	}

	props := StringMap{}
	props[LTI_USER_ID_PROP_KEY] = e.GetUserId(launchData)

	if props[LTI_USER_ID_PROP_KEY] == "" {
		return nil, NewAppError("Uoc_BuildUser", "Uoc.build_user.lti_user_id_missing", nil, "", http.StatusBadRequest)
	}

	firstName := strings.Trim(returnValueBase64Encoded(launchData[uocLaunchDataFirstNameKey], launchData), " ")
	lastName := strings.Trim(returnValueBase64Encoded(launchData[uocLaunchDataLastNameKey], launchData), " ")

	if firstName == "" || lastName == "" {
		name := strings.SplitN(strings.Trim(launchData[uocLaunchDataFullNameKey], " "), " ", 2)
		if firstName == "" && len(name) > 0 {
			firstName = name[0]
		}

		if lastName == "" && len(name) > 1 {
			lastName = name[1]
		}
	}

	if firstName == "" {
		firstName = returnValueBase64Encoded(launchData[uocLaunchDataUsernameKey], launchData)
	}

	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     returnValueBase64Encoded(launchData[uocLaunchDataEmailKey], launchData),
		Username:  transformLTIUsername(returnValueBase64Encoded(launchData[uocLaunchDataUsernameKey], launchData)),
		Position:  returnValueBase64Encoded(launchData[uocLaunchDataPositionKey], launchData),
		Password:  password,
		Props:     props,
		Locale:    returnValueBase64Encoded(launchData[uocLaunchDataPresentationLocale], launchData),
	}

	return user, nil
}

func (e *UocLMS) GetTeam(launchData map[string]string) string {
	// team depends on an LTI param
	if launchData[uocLaunchDataTeamIsTransversalParam] != "" {
		return returnValueBase64Encoded(launchData[uocLaunchDataTransversalTeamParam], launchData)
	}
	return returnValueBase64Encoded(launchData[uocLaunchDataContextId], launchData)
}

func (e *UocLMS) GetPublicChannelsToJoin(launchData map[string]string) map[string]string {
	return map[string]string{}
}

func (e *UocLMS) GetPrivateChannelsToJoin(launchData map[string]string) map[string]string {
	channels := map[string]string{}

	for personalChannelName, channelConfig := range e.PersonalChannels.ChannelList {
		channelDisplayName := returnValueBase64Encoded(launchData[channelConfig.NameProperty], launchData)
		channelSlug := GetLMSChannelSlug(personalChannelName, returnValueBase64Encoded(launchData[channelConfig.IdProperty], launchData))

		if channelDisplayName != "" && channelSlug != "" {
			channels[channelSlug] = channelDisplayName
		}
	}

	return channels
}

func (e *UocLMS) GetChannel(launchData map[string]string) (string, *AppError) {
	privateChannels := e.GetPrivateChannelsToJoin(launchData)
	if len(privateChannels) == 0 {
		return "", nil
	}

	for channelSlug, _ := range privateChannels {
		return truncateLMSChannelSlug(channelSlug), nil
	}

	return "", nil
}

func (e *UocLMS) SyncUser(user *User, launchData map[string]string) *User {
	if launchData[uocLaunchDataEmailKey] != "" {
		user.Email = returnValueBase64Encoded(launchData[uocLaunchDataEmailKey], launchData)
	}

	if launchData[uocLaunchDataUsernameKey] != "" {
		user.Username = transformLTIUsername(returnValueBase64Encoded(launchData[uocLaunchDataUsernameKey], launchData))
	}

	if launchData[uocLaunchDataPositionKey] != "" {
		user.Position = returnValueBase64Encoded(launchData[uocLaunchDataPositionKey], launchData)
	}

	if user.Props == nil {
		user.Props = StringMap{}
	}

	user.Props[LTI_USER_ID_PROP_KEY] = e.GetUserId(launchData)
	return user
}

func (e *UocLMS) BuildTeam(launchData map[string]string) (*Team, *AppError) {
	team := &Team{
		DisplayName:     returnValueBase64Encoded(launchData[uocLaunchDataContextTitle], launchData),
		Name:            returnValueBase64Encoded(launchData[uocLaunchDataContextId], launchData),
		AllowOpenInvite: false,
		Type:            TEAM_INVITE,
	}

	return team, nil
}
