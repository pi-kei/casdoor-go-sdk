// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casdoorsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Token has the same definition as https://github.com/casdoor/casdoor/blob/master/object/token.go#L45
type Token struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Application  string `xorm:"varchar(100)" json:"application"`
	Organization string `xorm:"varchar(100)" json:"organization"`
	User         string `xorm:"varchar(100)" json:"user"`

	Code          string `xorm:"varchar(100) index" json:"code"`
	AccessToken   string `xorm:"mediumtext" json:"accessToken"`
	RefreshToken  string `xorm:"mediumtext" json:"refreshToken"`
	ExpiresIn     int    `json:"expiresIn"`
	Scope         string `xorm:"varchar(100)" json:"scope"`
	TokenType     string `xorm:"varchar(100)" json:"tokenType"`
	CodeChallenge string `xorm:"varchar(100)" json:"codeChallenge"`
	CodeIsUsed    bool   `json:"codeIsUsed"`
	CodeExpireIn  int64  `json:"codeExpireIn"`
}

func (c *Client) GetTokens() ([]*Token, error) {
	queryMap := map[string]string{
		"owner": c.OrganizationName,
	}

	url := c.GetUrl("get-tokens", queryMap)

	bytes, err := c.DoGetBytes(url)
	if err != nil {
		return nil, err
	}

	var tokens []*Token
	err = json.Unmarshal(bytes, &tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (c *Client) GetPaginationTokens(p int, pageSize int, queryMap map[string]string) ([]*Token, int, error) {
	queryMap["owner"] = c.OrganizationName
	queryMap["p"] = strconv.Itoa(p)
	queryMap["pageSize"] = strconv.Itoa(pageSize)

	url := c.GetUrl("get-tokens", queryMap)

	response, err := c.DoGetResponse(url)
	if err != nil {
		return nil, 0, err
	}

	tokens, ok := response.Data.([]*Token)
	if !ok {
		return nil, 0, errors.New("response data format is incorrect")
	}

	return tokens, int(response.Data2.(float64)), nil
}

func (c *Client) GetToken(name string) (*Token, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", c.OrganizationName, name),
	}

	url := c.GetUrl("get-token", queryMap)

	bytes, err := c.DoGetBytes(url)
	if err != nil {
		return nil, err
	}

	var token *Token
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c *Client) UpdateToken(token *Token) (bool, error) {
	_, affected, err := c.modifyToken("update-token", token, nil)
	return affected, err
}

func (c *Client) UpdateTokenForColumns(token *Token, columns []string) (bool, error) {
	_, affected, err := c.modifyToken("update-token", token, columns)
	return affected, err
}

func (c *Client) AddToken(token *Token) (bool, error) {
	_, affected, err := c.modifyToken("add-token", token, nil)
	return affected, err
}

func (c *Client) DeleteToken(token *Token) (bool, error) {
	_, affected, err := c.modifyToken("delete-token", token, nil)
	return affected, err
}
