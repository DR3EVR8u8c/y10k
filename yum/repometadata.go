package yum

import (
	"encoding/xml"
	"fmt"
	"io"
)

// RepoMetadata represents the metadata XML file for a RPM/Yum repository. It
// contains pointers to database files which describe the packages available in
// the repository.
type RepoMetadata struct {
	XMLName  xml.Name `xml:"repomd"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSRPM string   `xml:"xmlns:rpm,attr"`

	Revision  int            `xml:"revision"`
	Databases []RepoDatabase `xml:"data"`
}

// RepoDatabase represents an entry in a repository metadata file for an
// individual database file such as primary_db or filelists_db.
type RepoDatabase struct {
	Type            string               `xml:"type,attr"`
	Location        RepoDatabaseLocation `xml:"location"`
	Timestamp       int64                `xml:"timestamp"`
	Size            int                  `xml:"size"`
	Checksum        RepoDatabaseChecksum `xml:"checksum"`
	OpenSize        int                  `xml:"open-size"`
	OpenChecksum    RepoDatabaseChecksum `xml:"open-checksum"`
	DatabaseVersion int                  `xml:"database_version"`
}

// RepoDatabaseLocation represents the URI, relative to a package repository,
// of a repository database.
type RepoDatabaseLocation struct {
	Href string `xml:"href,attr"`
}

func (c *RepoDatabase) String() string {
	return c.Type
}

// ReadRepoMetadata loads a repomd.xml file from the given io.Reader and returns
// a pointer to the resulting RepoMetadata struct.
func ReadRepoMetadata(r io.Reader) (*RepoMetadata, error) {
	md := RepoMetadata{
		Databases: make([]RepoDatabase, 0),
	}

	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&md)

	if err != nil {
		return nil, fmt.Errorf("Error decoding repository metadata: %v", err)
	}

	return &md, nil
}

// Write encodes a RepoMetadata struct in the repomd.xml format to the given
// io.Writer stream.
func (c *RepoMetadata) Write(w io.Writer) error {
	c.XMLNS = "http://linux.duke.edu/metadata/repo"
	c.XMLNSRPM = "http://linux.duke.edu/metadata/rpm"

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("Error encoding repository metadata: %v", err)
	}

	return nil
}
