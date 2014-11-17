package main 

import (
	"strings"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type Config struct {
        Connstr string
	ProjectID string
        ClientId string
        ClientSecret string
}

var config Config

func readconfig() {

        /* get configuration file */
        file, err := ioutil.ReadFile("/etc/archiver.photo.json")
        if err != nil {
                fmt.Println("Error opening config: ", err)
                os.Exit(1)
        }

        if err = json.Unmarshal(file, &config); err != nil {
                fmt.Println("Error parsing config: ", err)
                os.Exit(1);
        }

}

func escape(s string) string {
	/* escape input for MySQL query */
	re := strings.NewReplacer("'", "''")
	return "'" + re.Replace(s) + "'"
}
