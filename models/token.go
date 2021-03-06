// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"errors"
	"time"

	"github.com/go-gitea/gitea/modules/base"
	"github.com/go-gitea/gitea/modules/uuid"
)

var (
	ErrAccessTokenNotExist = errors.New("Access token does not exist")
)

// AccessToken represents a personal access token.
type AccessToken struct {
	Id                int64
	Uid               int64
	Name              string
	Sha1              string    `xorm:"UNIQUE VARCHAR(40)"`
	Created           time.Time `xorm:"CREATED"`
	Updated           time.Time
	HasRecentActivity bool `xorm:"-"`
	HasUsed           bool `xorm:"-"`
}

// NewAccessToken creates new access token.
func NewAccessToken(t *AccessToken) error {
	t.Sha1 = base.EncodeSha1(uuid.NewV4().String())
	_, err := x.Insert(t)
	return err
}

// GetAccessTokenBySha returns access token by given sha1.
func GetAccessTokenBySha(sha string) (*AccessToken, error) {
	t := &AccessToken{Sha1: sha}
	has, err := x.Get(t)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrAccessTokenNotExist
	}
	return t, nil
}

// ListAccessTokens returns a list of access tokens belongs to given user.
func ListAccessTokens(uid int64) ([]*AccessToken, error) {
	tokens := make([]*AccessToken, 0, 5)
	err := x.Where("uid=?", uid).Desc("id").Find(&tokens)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		t.HasUsed = t.Updated.After(t.Created)
		t.HasRecentActivity = t.Updated.Add(7 * 24 * time.Hour).After(time.Now())
	}
	return tokens, nil
}

// DeleteAccessTokenById deletes access token by given ID.
func DeleteAccessTokenById(id int64) error {
	_, err := x.Id(id).Delete(new(AccessToken))
	return err
}
