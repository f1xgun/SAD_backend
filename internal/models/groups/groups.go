package groupsModels

import (
	"database/sql"
	"sad/internal/models/users"
)

type Group struct {
	Id     string `json:"id"`
	Number string `json:"number"`
}

type GroupRepoModel struct {
	Id     sql.NullString
	Number sql.NullString
}

type UserGroup struct {
	UserId  string `json:"user_id"`
	GroupId string `json:"group_id"`
}

type GroupDetails struct {
	Group
	Users []usersModels.UserInfo `json:"users"`
}

type GroupDetailsRepo struct {
	GroupRepoModel
	Users []usersModels.UserInfoRepoModel
}
