// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
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

func (e *UocLMS) GetEmail(launchData map[string]string) string {
	return launchData[uocLaunchDataEmailKey]
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
	return launchData[uocLaunchDataLTIUserIdKey]
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

	firstName := strings.Trim(launchData[uocLaunchDataFirstNameKey], " ")
	lastName := strings.Trim(launchData[uocLaunchDataLastNameKey], " ")

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
		firstName = launchData[uocLaunchDataUsernameKey]
	}

	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     launchData[uocLaunchDataEmailKey],
		Username:  transformLTIUsername(launchData[uocLaunchDataUsernameKey]),
		Position:  launchData[uocLaunchDataPositionKey],
		Password:  password,
		Props:     props,
		Locale:    launchData[uocLaunchDataPresentationLocale],
	}

	return user, nil
}

func (e *UocLMS) GetTeam(launchData map[string]string) string {
	// team depends on an LTI param
	if launchData[uocLaunchDataTeamIsTransversalParam] != "" {
		return launchData[uocLaunchDataTransversalTeamParam]
	}
	return launchData[uocLaunchDataContextId]
}

func (e *UocLMS) GetPublicChannelsToJoin(launchData map[string]string) map[string]string {
	return map[string]string{}
}

func (e *UocLMS) GetPrivateChannelsToJoin(launchData map[string]string) map[string]string {
	channels := map[string]string{}

	for personalChannelName, channelConfig := range e.PersonalChannels.ChannelList {
		channelDisplayName := launchData[channelConfig.NameProperty]
		channelSlug := GetLMSChannelSlug(personalChannelName, launchData[channelConfig.IdProperty])

		if channelDisplayName != "" && channelSlug != "" {
			channels[channelSlug] = channelDisplayName
		}
	}

	return channels
}

func (e *UocLMS) GetChannel(launchData map[string]string) (string, *AppError) {
	customChannelRedirect, ok := launchData[uocLaunchDataChannelRedirectKey]
	if !ok {
		return "", nil
	}

	var channelSlug string

	components := strings.Split(customChannelRedirect, ":")
	if len(components) == 1 {
		channelSlug = components[0]
	} else if components[0] == uocRedirectChannelLookupKeyword {
		UocChannel, ok := e.PersonalChannels.ChannelList[components[1]]
		if !ok {
			return "", NewAppError("GetChannel", "get_channel.redirect_lookup_channel.not_found", nil, "", http.StatusBadRequest)
		}

		channelSlug = GetLMSChannelSlug(components[1], launchData[UocChannel.IdProperty])
	}

	return truncateLMSChannelSlug(channelSlug), nil
}

func (e *UocLMS) SyncUser(user *User, launchData map[string]string) *User {
	if launchData[uocLaunchDataEmailKey] != "" {
		user.Email = launchData[uocLaunchDataEmailKey]
	}

	if launchData[uocLaunchDataUsernameKey] != "" {
		user.Username = transformLTIUsername(launchData[uocLaunchDataUsernameKey])
	}

	if launchData[uocLaunchDataPositionKey] != "" {
		user.Position = launchData[uocLaunchDataPositionKey]
	}

	if user.Props == nil {
		user.Props = StringMap{}
	}

	user.Props[LTI_USER_ID_PROP_KEY] = e.GetUserId(launchData)
	return user
}

func (e *UocLMS) BuildTeam(launchData map[string]string) (*Team, *AppError) {
	team := &Team{
		DisplayName:     launchData[uocLaunchDataContextTitle],
		Name:            launchData[uocLaunchDataContextId],
		AllowOpenInvite: false,
		Type:            TEAM_INVITE,
	}

	return team, nil
}
