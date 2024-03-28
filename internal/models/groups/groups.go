package groupsModels

import (
	usersModels "sad/internal/models/user"
)

type Group struct {
	Id     string `json:"id"`
	Number string `json:"number"`
}

type UserGroup struct {
	UserId  string `json:"user_id"`
	GroupId string `json:"group_id"`
}

type GroupWithUsers struct {
	Group
	Users []usersModels.UserInfo `json:"users"`
}
