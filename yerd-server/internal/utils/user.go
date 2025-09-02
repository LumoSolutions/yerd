package utils

import (
	"fmt"
	"os/user"
	"strconv"
)

type UserContext struct {
	User      *user.User
	Username  string
	GroupName string
	HomeDir   string
	UID       int
	GID       int
}

func GetUser() (*UserContext, error) {
	realUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %v", err)
	}

	uid, err := strconv.Atoi(realUser.Uid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UID: %v", err)
	}

	gid, err := strconv.Atoi(realUser.Gid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GID: %v", err)
	}

	var groupName string
	group, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		groupName = "nobody"
	} else {
		groupName = group.Name
	}

	return &UserContext{
		User:      realUser,
		Username:  realUser.Username,
		GroupName: groupName,
		HomeDir:   realUser.HomeDir,
		UID:       uid,
		GID:       gid,
	}, nil
}
