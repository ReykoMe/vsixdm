package loader

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func getExtIDFromVSCodeStore(parsedLink *url.URL) (string, error) {
	searchParamName := "itemName"
	queryParams := parsedLink.Query()
	// Check is link from vscode marketplace: https://marketplace.visualstudio.com/items?itemName=publisher.extension
	codeMarketExtensionID := queryParams.Get(searchParamName)
	if len(codeMarketExtensionID) > 0 {
		return codeMarketExtensionID, nil
	}
	return "", fmt.Errorf("missed 'itemName' param %q, %q", searchParamName, parsedLink.RawPath)

}

func getExtIDFromOpenVSIXStore(parsedLink *url.URL) (string, error) {
	// Check is link is related to openvsx
	host := "open-vsx.org"
	if !strings.Contains(parsedLink.Host, host) {
		return "", fmt.Errorf("not an 'open-vsx.org' lin %q", parsedLink)
	}

	//split path, extract parts, remove empty elements
	parts := strings.Split(strings.Trim(parsedLink.Path, "/"), "/")
	if len(parts) != 3 {
		return "", fmt.Errorf("please check open-vsx.org link: %q", parsedLink.RawPath)
	}
	slice := parts[len(parts)-2:]
	extId := strings.Join(slice, ".")
	return extId, nil

}

func GetExtensionIdByLink(link string) (string, error) {
	parsedLink, err := url.Parse(link)

	if err != nil {
		log.Fatalf("Can't parse link %q", link)
	}

	getters := []func(*url.URL) (string, error){
		getExtIDFromVSCodeStore,
		getExtIDFromOpenVSIXStore,
	}

	for _, getExtensionId := range getters {
		if extenstionId, err := getExtensionId(parsedLink); err == nil && len(extenstionId) > 0 {
			return extenstionId, nil
		}
	}

	return "", fmt.Errorf("can't parse a link %q", link)
}
