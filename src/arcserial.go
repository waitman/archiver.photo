/*

arcserial - list drives with serial numbers

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
//      "bytes"
        "fmt"
        "log"
	"os"
        "os/exec"
        "strings"
)

func main() {
        var news string
        var disks []string
        var serial string

	out,err := exec.Command("uname","-s").Output()
	if (strings.TrimSpace(string(out))!="FreeBSD") {
		fmt.Println("This command only works on FreeBSD")
		os.Exit(1)
	}
        out,err = exec.Command("sysctl", "kern.disks").Output()
        if err != nil {
                log.Fatal(err)
        }

        news = strings.TrimSpace(strings.Replace(string(out),"kern.disks: ","",1))
        disks = strings.Split(news," ")

        for _, disk := range disks {
                out,err := exec.Command("camcontrol","identify","/dev/"+disk).Output()
                if err != nil {
                        //drive doesnt exist, ie sdhc card
                } else {
                        // get serial number
                        var ns []string
                        ns = strings.Split(string(out),"serial number")
                        ns = strings.Split(ns[1],"\n")
                        serial = strings.TrimSpace(ns[0])
                        fmt.Printf("%s %s\n",disk,serial)
                }
        }
}

