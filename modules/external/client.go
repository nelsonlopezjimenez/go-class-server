package external

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	CIS "localhost/CIS/modules"
	"net/http"
	"os"
	"time"
)

type MetaData struct {
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	Topics      []string  `json:"topics"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

type ExternalData struct {
	Name string
	Meta Info
}

type Info struct {
	Size        int
	Topics      []string
	Description string
	Created_At  time.Time
	Updated_At  time.Time
}

var SiteData struct {
	DataStuff []MetaData `json:"data"`
}

var websites string = CIS.GetOSPaths().Websites

func GetSiteMetaData() ([]ExternalData, error) {
	_, err := os.Stat(websites + "/info.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("websites metadata file does not exist")
		}
		return nil, fmt.Errorf("problem opening website metadata: %w", err)
	}
	file, err := os.OpenFile(websites+"/info.json", os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fsys := os.DirFS(websites)
	file1, err := fs.ReadFile(fsys, "info.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(file1, &SiteData)
	if err != nil {
		fmt.Println(err)
	}

	MetaSlice := []ExternalData{}
	for _, metaDataObj := range SiteData.DataStuff {
		MetaSlice = append(MetaSlice, ExternalData{metaDataObj.Name, Info{metaDataObj.Size, metaDataObj.Topics, metaDataObj.Description, metaDataObj.Created_At, metaDataObj.Updated_At}})
	}
	return MetaSlice, nil
}

func UpdateSiteMetaData() error {
	file, err := os.OpenFile(websites+"/info.json", os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	client := http.DefaultClient
	res, err := client.Get("http://192.168.1.47:3000/api/v1/repos/search?uid=6")
	if err != nil {
		return fmt.Errorf("unable to connect to API. Make sure you are connected to the network: %w", err)
	}

	reader := bufio.NewReader(res.Body)
	reader.WriteTo(file)
	return nil
}
