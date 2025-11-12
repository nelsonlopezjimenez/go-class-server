package external

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"localhost/CIS/modules/util"
	"net/http"
	"os"
	"slices"
	"time"
)

type Info struct {
	Size        int       `json:"size"`
	Topics      []string  `json:"topics"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}
type MetaData struct {
	Name string `json:"name"`
	Info
}

type ExternalData struct {
	Name string
	Meta Info
}

type DataStuff struct {
	Data []WebsiteInfo `json:"data"`
}

type WebsiteInfo struct {
	Domain    string `json:"name"`
	IndexPath string `json:"indexPath"`
	State     string `json:"state"`
	Info      `json:"meta"`
}

var websites string = util.GetOSPaths().Websites

func getMetaFromGitea() ([]byte, error) {
	client := http.DefaultClient
	res, err := client.Get("http://localhost:3000/api/v1/repos/search?uid=5&limit=200")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to API. Make sure you are connected to the network: %w", err)
	}

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	return bodyContent, nil

}

func processMetaBytesFromGitea(metaBytes []byte) ([]WebsiteInfo, error) {
	var processedMetaData DataStuff
	var processedWebsiteInfo []WebsiteInfo

	err := json.Unmarshal(metaBytes, &processedMetaData)
	if err != nil {
	}

	for _, metaData := range processedMetaData.Data {
		filePath := util.GetOSPaths().Websites + "/" + metaData.Domain
		index := ""
		_, err := os.Stat(filePath)
		if err == nil {
			websiteRoot := util.MakeRootDir(filePath)
			websiteRoot.FindIndex(func(path string, ent fs.DirEntry) {
				index = path + "/" + ent.Name()
			})
		}
		siteData := WebsiteInfo{
			metaData.Domain,
			index,
			"not_installed",
			metaData.Info,
		}

		processedWebsiteInfo = append(processedWebsiteInfo, siteData)
	}
	return processedWebsiteInfo, nil
}

func createLocalMeta() error {
	metaBytes, err := getMetaFromGitea()
	if err != nil {
		return err
	}

	processedWebInfo, err := processMetaBytesFromGitea(metaBytes)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(websites+"/info.json", os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	toWrite, err := json.Marshal(processedWebInfo)
	if err != nil {
		return err
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}
	return nil
}

func loadMetaFromFile() ([]WebsiteInfo, error) {
	var localMeta []WebsiteInfo
	file, err := os.ReadFile(websites + "/info.json")
	if err != nil {
		if os.IsNotExist(err) {
			err := createLocalMeta()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	err = json.Unmarshal(file, &localMeta)
	if err != nil {
		return nil, fmt.Errorf("fromFile: %w", err)
	}

	return localMeta, nil
}

func BuildWebsiteList() ([]WebsiteInfo, error) {
	// var websiteInfoList []WebsiteInfo
	localWebsitesDir := util.RootDir{Root: websites}
	localWebsitesContents := localWebsitesDir.ListBuilder()

	localData, err := loadMetaFromFile()
	if err != nil {
		return nil, err
	}

	_, err = getMetaFromGitea()
	if err != nil {
		goto offline
	}

	fmt.Printf("localData: %v\n", localData)

offline:
	for i, website := range localData {
		if website.Topics == nil {
			localData[i].Topics = []string{}
		}
		if slices.Contains(localWebsitesContents, website.Domain) {
			localData[i].State = "installed"
		}
	}
	return localData, nil
}

func processOrphanSites(siteList []WebsiteInfo) ([]WebsiteInfo, error) {
	var classSites []string
	localWebsitesDir := util.RootDir{Root: websites}

	localWebsitesContents := localWebsitesDir.ListBuilder()
	for _, external := range siteList {
		classSites = append(classSites, external.Domain)
	}

	for _, dirent := range localWebsitesContents {
		if !slices.Contains(classSites, dirent) {
			orphanMeta := Info{}
			fileInfo, err := os.Stat(util.GetOSPaths().Websites + "/" + dirent)
			if err == nil {
				orphanMeta = Info{
					int(fileInfo.Size()),
					[]string{},
					"No data Available",
					time.Time{},
					fileInfo.ModTime(),
				}
			}
			index := ""
			orphanRoot := util.RootDir{Root: util.GetOSPaths().Websites + "/" + dirent}
			orphanRoot.FindIndex(func(path string, ent fs.DirEntry) {
				index = path + "/" + ent.Name()

			})
			orphanedSiteInfo := WebsiteInfo{
				dirent,
				index,
				"orphan",
				orphanMeta,
			}
			siteList = append(siteList, orphanedSiteInfo)
		}
	}
	return siteList, nil

}

func SendAllSites() ([]WebsiteInfo, error) {
	external, err := BuildWebsiteList()
	if err != nil {
		return nil, err
	}

	allSites, err := processOrphanSites(external)
	if err != nil {
		return nil, err
	}

	return allSites, nil
}

func UpdateSiteMetaData() error {
	local, err := loadMetaFromFile()
	if err != nil {
		return err
	}

	external, err := getMetaFromGitea()
	if err != nil {
		return err
	}

	processedExt, err := processMetaBytesFromGitea(external)
	if err != nil {
		return err
	}

	if len(local) == len(processedExt) {
		for i, localSite := range local {
			for _, externalSite := range processedExt {
				if localSite.Domain == externalSite.Domain {
					local[i].Description = externalSite.Description
					local[i].Topics = externalSite.Topics
					if localSite.Size != externalSite.Size || localSite.Updated_At != externalSite.Updated_At {
						local[i].State = "update_needed"
					}
				}
			}
		}
	}
	return nil
}
