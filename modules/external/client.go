package external

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	CIS "localhost/CIS/modules"
	"net/http"
	"os"
	"slices"
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

type DataStuff struct {
	Data []MetaData `json:"data"`
}

type WebsiteInfo struct {
	Domain    string
	IndexPath string
	Installed bool
	IsCurrent bool
	Meta      Info
	IsOrphan  bool
}

var websites string = CIS.GetOSPaths().Websites

func GetSiteMetaData() ([]ExternalData, error) {
	_, err := os.Stat(websites + "/info.json")
	if err != nil {
		if os.IsNotExist(err) {
			err := createSiteMetaData()
			if err != nil {
				return nil, fmt.Errorf("could not create the meta file: %w", err)
			}
		}
		return nil, fmt.Errorf("problem opening website metadata: %w", err)
	}
	file, err := os.OpenFile(websites+"/info.json", os.O_RDWR, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	metaSlice, err := loadMetaFromFile()
	if err != nil {
		fmt.Println(err)
	}

	MetaSlice := []ExternalData{}
	for _, metaDataObj := range metaSlice {
		MetaSlice = append(MetaSlice, ExternalData{metaDataObj.Name, Info{metaDataObj.Size, metaDataObj.Topics, metaDataObj.Description, metaDataObj.Created_At, metaDataObj.Updated_At}})
	}
	return MetaSlice, nil
}

func loadMetaFromFile() ([]MetaData, error) {
	SiteData := DataStuff{}
	fsys := os.DirFS(websites)
	file1, err := fs.ReadFile(fsys, "info.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read from file: %w", err)
	}

	err = json.Unmarshal(file1, &SiteData)
	if err != nil {
		return nil, fmt.Errorf("could not decode file: %w", err)
	}

	return SiteData.Data, nil
}

func UpdateSiteMetaData() error {
	networkData := []MetaData{}
	fileData, err := loadMetaFromFile()
	if err != nil {
		return fmt.Errorf("could not update site list:  %w", err)
	}

	reader, err := getMetaFromGitea()
	if err != nil {
		return fmt.Errorf("could not update site data: %w", err)
	}
	networkBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("could not update site data: %w", err)
	}

	err = json.Unmarshal(networkBytes, &networkData)
	if err != nil {
		return fmt.Errorf("could not read network data to memory: %w", err)
	}

	for _, localSite := range fileData {
		for _, networkSite := range networkData {
			if localSite.Name == networkSite.Name {
				if localSite.Updated_At != networkSite.Updated_At || localSite.Size != networkSite.Size {
					fmt.Println("need update")
				}
			}
		}
	}

	return err
}

func createSiteMetaData() error {
	file, err := os.OpenFile(websites+"/info.json", os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	metaReader, err := getMetaFromGitea()
	if err != nil {
		return fmt.Errorf("could not get the data from Gitea: %w", err)
	}
	metaReader.WriteTo(file)

	return nil
}

func getMetaFromGitea() (*bufio.Reader, error) {
	client := http.DefaultClient
	res, err := client.Get("http://192.168.1.47:3000/api/v1/repos/search?uid=6&limit=200")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to API. Make sure you are connected to the network: %w", err)
	}

	reader := bufio.NewReader(res.Body)
	return reader, nil
}

func SendSiteList() ([]WebsiteInfo, error) {
	externalSiteList := []string{}
	localWebsitesDir := CIS.RootDir{Root: CIS.GetOSPaths().Websites}
	localWebsitesContents := localWebsitesDir.ListBuilder()
	websiteMetaData, err := GetSiteMetaData()
	if err != nil {
		return nil, fmt.Errorf("did not get metadata: %w", err)
	}
	websiteInfoSlice := []WebsiteInfo{}
	// for _, item := range websitesList {
	// metaInfo := external.Info{}
	for _, website := range websiteMetaData {
		externalSiteList = append(externalSiteList, website.Name)
		index := ""
		websiteRoot := CIS.MakeRootDir(CIS.GetOSPaths().Websites + "/" + website.Name)
		websiteRoot.FindIndex(func(path string, ent fs.DirEntry) {
			index = path + "/" + ent.Name()
		})
		singlePage := WebsiteInfo{}
		singlePage.Installed = slices.Contains(localWebsitesContents, website.Name)
		if singlePage.Installed {
			singlePage.IndexPath = index
			singlePage.IsCurrent = IsSiteCurrent()
		} else {
			singlePage.IsCurrent = false
		}
		singlePage.IsOrphan = false

		singlePage.Domain = website.Name
		singlePage.Meta = website.Meta

		websiteInfoSlice = append(websiteInfoSlice, singlePage)

	}

	for _, dirent := range localWebsitesContents {
		if !slices.Contains(externalSiteList, dirent) {
			orphanMeta := Info{}
			fileInfo, err := os.Stat(CIS.GetOSPaths().Websites + "/" + dirent)
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
			orphanRoot := CIS.RootDir{Root: CIS.GetOSPaths().Websites + "/" + dirent}
			orphanRoot.FindIndex(func(path string, ent fs.DirEntry) {
				index = path + "/" + ent.Name()

			})
			orphanedSiteInfo := WebsiteInfo{
				dirent,
				index,
				true,
				true,
				orphanMeta,
				true,
			}
			websiteInfoSlice = append(websiteInfoSlice, orphanedSiteInfo)

		}
	}

	return websiteInfoSlice, nil
}

func IsOrphanSite(local ExternalData) bool {
	localWebsitesDir := CIS.RootDir{Root: CIS.GetOSPaths().Websites}
	localWebsitesContents := localWebsitesDir.ListBuilder()
	networkSites, err := loadMetaFromFile()
	if err != nil {
		return true
	}
	// return !slices.Contains(localWebsitesContents, local.Name)
	for _, local := range localWebsitesContents {
		for _, network := range networkSites {
			if local == network.Name {
				return false
			}
		}
	}
	return true
}

func IsSiteCurrent() bool {
	data := DataStuff{}

	local, err := loadMetaFromFile()
	if err != nil {
		return true
	}
	reader, err := getMetaFromGitea()
	if err != nil {
		return true
	}
	resSlice, err := io.ReadAll(reader)
	if err != nil {
		return true
	}
	err = json.Unmarshal(resSlice, &data)
	if err != nil {
		return true
	}

	for _, localMeta := range local {
		for _, extMeta := range data.Data {
			if localMeta.Name == extMeta.Name {
				return localMeta.Updated_At.Equal(extMeta.Updated_At)
			}
		}
	}

	return true
}
