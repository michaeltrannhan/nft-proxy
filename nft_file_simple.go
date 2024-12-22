package nft_proxy

import "strings"

type NFTMetadataSimple struct {
	Name            string              `json:"name"`
	Symbol          string              `json:"symbol"`
	Decimals        uint8               `json:"-"`
	Image           string              `json:"image"`
	AnimationURL    string              `json:"animationUrl"`
	ExternalURL     string              `json:"externalUrl"`
	Properties      NFTPropertiesSimple `json:"properties"`
	Files           []NFTFiles          `json:"files"`
	UpdateAuthority string              `json:"updateAuthority"`
}

func (m *NFTMetadataSimple) AnimationFile() *NFTFiles {
	for _, f := range m.Files {
		if f.URL == m.Image || strings.Contains(f.Type, "image") && !strings.Contains(f.Type, "gif") {
			continue
		}
		return &f
	}
	return nil
}

func (m *NFTMetadataSimple) ImageFile() *NFTFiles {
	for _, f := range m.Files {
		if f.URL == m.Image {
			return &f
		}
	}

	return nil
}

type NFTFiles struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type NFTPropertiesSimple struct {
	Category string             `json:"category"`
	Creators []NFTCreatorSimple `json:"creators"`
}

type NFTCollectionSimple struct {
	Name   string `json:"name"`
	Family string `json:"family"`
}

type NFTCreatorSimple struct {
	Address string `json:"address"`
}

type NFTAttributeSimple struct {
	TraitType string      `json:"traitType"`
	Value     interface{} `json:"value"`
}

type NFTFileSimple struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}
