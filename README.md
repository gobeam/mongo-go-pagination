# Golang Oauth2 with jwt OAuth 2.0
[![Build][Build-Status-Image]][Build-Status-Url] [![Go Report Card](https://goreportcard.com/badge/github.com/roshanr83/goOauth2?branch=master)](https://goreportcard.com/report/github.com/roshanr83/goOauth2) [![GoDoc][godoc-image]][godoc-url]

This project is modified version of [go-oauth2/oauth2](https://github.com/go-oauth2/oauth2). Since that didn't meet my requirement so I modified the code so I can implement oauth2 alongside with JWT.
<br>
This package uses <b>EncryptOAEP</b> which encrypts the given data with <b>RSA-OAEP</b> to encrypt token data. Two separate file <b>private.pem</b> and <b>public.pem</b> file will be created on your root folder which includes respective private and public RSA keys which is used for encryption.
<br>
This package only handles Resource Owner Password Credentials type.
<br>
Official docs: [Here](https://godoc.org/github.com/roshanr83/goOauth2)

## Install

``` bash
$ go get -u -v github.com/roshanr83/go-oauth2
```

## Usage

``` go
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/roshanr83/go-oauth2"
	"gopkg.in/go-oauth2/mysql.v3"
	"time"
)

func main() {
	//register store
	store := oauth.NewDefaultStore(
		oauth.NewConfig("root:root@tcp(127.0.0.1:8889)/goauth?charset=utf8&parseTime=True&loc=Local"),
	)
	defer store.Close()



	/* to create client
	 where 1 is user ID Which will return Oauth Clients
	 struct which include client id and secret whic is
	 later used to validate client credentials */
	store.CreateClient(userId int64)



	/* create access token alongside refresh token
	Since it will not include user authentication
	because it can be  different for everyone you will
	have to authenticate user and pass user id to Token struct.
	 Here you will authenticate user and get userID
	 you will have to provide all the field given below.
	 ClientID must be  valid uuid. AccessExpiresIn is required
	 to mark expiration time. In response you will get TokenResponse
	 including accesstoken and refeshtoken. */
	accessToken := &oauth.Token{
		ClientID:        uuid.MustParse("17d5a915-c403-487e-b41f-92fd1074bd30"),
		ClientSecret:    "UnCMSiJqxFg1O7cqL0MM",
		UserID:          userID,
		Scope:           "*",
		AccessCreateAt:  time.Now(),
		AccessExpiresIn: time.Second * 15,
		RefreshCreateAt: time.Now(),
	}
	resp, err := store.Create(accessToken TokenInfo)



	/*To check valid accessToken, you should
	pass accessToken and it will check if it is valid accesstoken
	including if it is valid and non revoked. If it is valid
	in response it will return AccessTokens data correspond to that token */
	resp, err := store.GetByAccess(accessToken string)



	/* To check valid refreshToken, you should pass
	refreshToken and it will check if it is valid
	refreshToken including if it is valid and non revoked
	and if it;s related accessToken is already revoked or
	not. If it is valid in response it will return AccessTokens
	data correspond to that token*/
	/* Note that refresh token after using one time
	will be revoked and cannot be used again */
	resp, err := store.GetByRefresh(refreshToken string)



	/*You can manually revoke access token by passing
	userId which you can get from valid token info */
	store.RevokeByAccessTokens(userId int64)



	/*You can manually revoke refresh token by passing
	accessTokenId which you can get from valid token info */
	store.RevokeRefreshToken(accessTokenId string)



	/* you can also clear all token related to
	user by passing TokenInfo from valid token */
	store.ClearByAccessToken(userId int64)
	
}


```

## Running the tests
Database config is used as "root:root@tcp(127.0.0.1:3306)/goauth?charset=utf8&parseTime=True&loc=Local" in const.go file, You may have to change that configuration according to your system config for successful test.

``` bash
$ go test
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.


## Acknowledgments
<ol>
<li> https://github.com/go-oauth2/oauth2 </li>
</ol>



## MIT License

```
Copyright (c) 2019
```

[Build-Status-Url]: https://travis-ci.org/roshanr83/goOauth2
[Build-Status-Image]: https://travis-ci.org/roshanr83/goOauth2.svg?branch=master
[godoc-url]: https://godoc.org/github.com/roshanr83/goOauth2
[godoc-image]: https://godoc.org/github.com/roshanr83/goOauth2?status.svg
