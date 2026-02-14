package external

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"localhost/CIS/modules/util"
	"net/http"
	"os"
	"reflect"
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

type DataStuff struct {
	Data []MetaData `json:"data"`
}

type WebsiteInfo struct {
	Domain    string `json:"name"`
	IndexPath string `json:"indexPath"`
	State     string `json:"state"`
	Info      `json:"meta"`
}

var websites string = util.GetOSPaths().Websites

// getMetaFromGitea
//
// Sends a GET request to Gitea's /api/v1/repos route and returns a []byte
// of the body's content.
func getMetaFromGitea() ([]byte, error) {
	fmt.Printf("%v", os.Getenv("IS_DEV"))
	client := http.DefaultClient
	res, err := client.Get(util.LoadEnv("WEBSITES_GITEA_ADDR"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to API. Make sure you are connected to the network: %w", err)
	}

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	return bodyContent, nil

}

// processedMetaBytesFromGitea
//
// Takes the []byte from getMetaFromGitea and parses it into []WebsiteInfo.
func processMetaBytesFromGitea(metaBytes []byte) ([]WebsiteInfo, error) {
	var processedMetaData DataStuff
	var processedWebsiteInfo []WebsiteInfo

	err := json.Unmarshal(metaBytes, &processedMetaData)
	if err != nil {
	}

	for _, metaData := range processedMetaData.Data {
		siteData := makeWebsiteInfo(metaData)

		processedWebsiteInfo = append(processedWebsiteInfo, siteData)
	}
	return processedWebsiteInfo, nil
}

// createLocalMeta
//
// If there is no info.json file located in the websites folder,
// this fn will create it and save the website metadata to file.
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
	defer file.Close()

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

// loadMetaFromFile
//
// Checks to see if info.json exists. If not, it calls
// createLocalMeta. Otherwise, it opens the file and unmarshals
// the json into []WebsiteInfo for further use.
func loadMetaFromFile() ([]WebsiteInfo, error) {
	var localMeta []WebsiteInfo

	_, err := os.Stat(websites + "/info.json")
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

	file, err := os.ReadFile(websites + "/info.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &localMeta)
	if err != nil {
		return nil, fmt.Errorf("fromFile: %w", err)
	}
	return localMeta, nil
}

func BuildWebsiteList() ([]WebsiteInfo, error) {
	localWebsitesDir := util.MakeRootDir(websites)
	localWebsitesContents := localWebsitesDir.ListBuilder()

	err := UpdateSiteMetaData()
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}

	localData, err := loadMetaFromFile()
	if err != nil {
		return nil, err
	}

	for i, website := range localData {
		if website.Topics == nil {
			localData[i].Topics = []string{}
		}
		if slices.Contains(localWebsitesContents, website.Domain) {
			if localData[i].State == "update_needed" {
				localData[i].State = "update_needed"
			} else {
				localData[i].State = "installed"
			}
		}
	}
	return localData, nil
}

// processOrphanSites
//
// creates a list of domains available from Gitea and compares it to the contents
// of the websites folder. If it finds directories in the websites folder that does
// not match the list available on Gitea, it creates orphan WebsiteInfo and appends
// it to the current website list sent to the front end.
func processOrphanSites(siteList []WebsiteInfo) ([]WebsiteInfo, error) {
	var classSites []string
	localWebsitesDir := util.MakeRootDir(websites)

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
			orphanRoot := util.MakeRootDir(util.GetOSPaths().Websites + "/" + dirent)
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

// SendAllSites
//
// Collects all the site data and returns a []WebsiteInfo of
// All available sites.
func SendAllSites() ([]WebsiteInfo, error) {
	external, err := BuildWebsiteList()
	if err != nil {
		return nil, err
	}

	allSites, err := processOrphanSites(external)
	if err != nil {
		return nil, err
	}

	allSites = sortWebInfos(allSites)

	return allSites, nil
}

// UpdateSiteMetaData
// Checks the local list against a fresh copy of the metadata
// from Gitea. Updates info.json and sets the state of any entry
// to "need_update" as necessary.
func UpdateSiteMetaData() error {
	local, err := loadMetaFromFile()
	if err != nil {
		return err
	}

	external, err := getMetaFromGitea()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(websites+"/info.json", os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	processedExt, err := processMetaBytesFromGitea(external)
	if err != nil {
		return err
	}

	localDomainList, externalDomainList := listOfDomains(local, processedExt)

	toDelete := []string{}

	println(len(toDelete))
	for _, name := range localDomainList {
		if !slices.Contains(externalDomainList, name) {
			toDelete = append(toDelete, name)
		}
	}

	if len(toDelete) > 0 {
		for _, name := range toDelete {
			local = filter(local, func(elm WebsiteInfo) bool {
				return elm.Domain != name
			})
		}
	}

	toAdd := []string{}

	for _, name := range externalDomainList {
		if !slices.Contains(localDomainList, name) {
			toAdd = append(toAdd, name)
		}
	}

	if len(toAdd) > 0 {
		for _, name := range toAdd {
			for i, site := range processedExt {
				if site.Domain == name && name != "" {
					local = append(local, processedExt[i])
				}
			}
		}
	}

	for i, localSite := range local {
		for _, externalSite := range processedExt {
			if !reflect.DeepEqual(localSite, externalSite) {
				if localSite.Domain == externalSite.Domain {
					if localSite.Size != externalSite.Size || localSite.Updated_At != externalSite.Updated_At {
						local[i].State = "update_needed"
					}
					local[i].Description = externalSite.Description
					local[i].Info = Info{
						localSite.Size,
						externalSite.Topics,
						externalSite.Description,
						externalSite.Created_At,
						localSite.Updated_At,
					}
					break
				} else {
					continue
				}

			} else {
				break
			}

		}

	}

	toWrite, err := json.Marshal(local)
	if err != nil {
		return err
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}
	return nil
}

func listOfDomains(localSlice []WebsiteInfo, externalSlice []WebsiteInfo) ([]string, []string) {
	var local []string
	var external []string

	for _, site := range localSlice {
		local = append(local, site.Domain)
	}

	for _, site := range externalSlice {
		external = append(external, site.Domain)
	}

	return local, external
}

func getSingleSiteMeta(domain string) (WebsiteInfo, error) {
	singleSiteInfo := MetaData{}
	client := http.DefaultClient
	res, err := client.Get(util.LoadEnv("WEBSITES_REPO_ADDR") + domain)
	if err != nil {
		return WebsiteInfo{}, fmt.Errorf("unable to connect to API. Make sure you are connected to the network: %w", err)
	}

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return WebsiteInfo{}, fmt.Errorf("could not read body: %w", err)
	}

	err = json.Unmarshal(bodyContent, &singleSiteInfo)
	if err != nil {
		return WebsiteInfo{}, err
	}

	return makeWebsiteInfo(singleSiteInfo), nil
}

func UpdateSingleInfo(domainUpdated string) error {
	fullLocalList, err := loadMetaFromFile()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(websites+"/info.json", os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	singleInfoToUpdate, err := getSingleSiteMeta(domainUpdated)
	if err != nil {
		return err
	}

	for i, localToUpdate := range fullLocalList {
		if localToUpdate.Domain == domainUpdated {
			fullLocalList[i] = singleInfoToUpdate
		}
	}

	toWrite, err := json.Marshal(fullLocalList)
	if err != nil {
		return err
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

func makeWebsiteInfo(meta MetaData) WebsiteInfo {

	filePath := util.GetOSPaths().Websites + "/" + meta.Name
	index := ""
	_, err := os.Stat(filePath)
	if err == nil {
		websiteRoot := util.MakeRootDir(filePath)
		websiteRoot.FindIndex(func(path string, ent fs.DirEntry) {
			index = path + "/" + ent.Name()
		})
	}
	siteData := WebsiteInfo{
		meta.Name,
		index,
		"not_installed",
		meta.Info,
	}
	return siteData
}

// filter
//
// An example of a generic fn. It takes a slice of any type and returns a
// filtered slice of the same type. It takes a callback that returns a bool
// type. Every item in the passed slice is fed into the callback fn and is
// added to the returned slice if the callback returns true.
func filter[Type any](s []Type, fn func(elm Type) bool) []Type {
	var results []Type
	for _, item := range s {
		if fn(item) {
			results = append(results, item)
		}
	}
	return results
}

// sortWebInfos
//
// a fn that takes a []WebInfo and returns a []WebInfo that has been
// sorted by the Domain field. This helps arrange the links on the
// offline-links page
func sortWebInfos(s []WebsiteInfo) []WebsiteInfo {
	var domain []string
	var sorted []WebsiteInfo

	for _, site := range s {
		domain = append(domain, site.Domain)
	}

	slices.Sort(domain)

	for _, name := range domain {
		for _, site := range s {
			if name == site.Domain {
				sorted = append(sorted, site)
				break
			}
		}
	}
	return sorted
}
