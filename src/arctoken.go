/*
	arctoken - get an oath token for google service

Copyright (c) 2014, Waitman Gobble <ns@waitman.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package main

import (
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
        "code.google.com/p/goauth2/oauth"
        "code.google.com/p/google-api-go-client/storage/v1"
)

const (
        scope      = storage.DevstorageFull_controlScope
        authURL    = "https://accounts.google.com/o/oauth2/auth"
        tokenURL   = "https://accounts.google.com/o/oauth2/token"
        entityName = "allUsers"
        redirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	
        cacheFile = flag.String("cache", "/etc/archiver.photo.cache.json", "Token cache file")
        code      = flag.String("code", "", "Authorization Code")

)

func main() {

	readconfig()

        var oconfig = &oauth.Config{
                ClientId:     config.ClientId,
                ClientSecret: config.ClientSecret,
                Scope:        scope,
                AuthURL:      authURL,
                TokenURL:     tokenURL,
                TokenCache:   oauth.CacheFile(*cacheFile),
                RedirectURL:  redirectURL,
        }

        flag.Parse()

        transport := &oauth.Transport{
                Config:    oconfig,
                Transport: http.DefaultTransport,
        }

        token, err := oconfig.TokenCache.Token()
        if err != nil {
                if *code == "" {
                        url := oconfig.AuthCodeURL("")
                        fmt.Println("Visit URL to get a code, then re-run with 'arctoken --code=CODE'")
                        fmt.Println(url)
                        os.Exit(1)
                }

                // Exchange auth code for access token
                token, err = transport.Exchange(*code)
                if err != nil {
                        log.Fatal("Exchange: ", err)
                }
                fmt.Printf("Token is cached in %v\n", oconfig.TokenCache)
        } else {
		fmt.Printf("Token is cached in %v\n", oconfig.TokenCache)
	}
	transport.Token = token

}
