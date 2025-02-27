package validator

import (
	"errors"
	"os"
	"strings"

	"github.com/cloudposse/atmos/pkg/downloader"
	"github.com/cloudposse/atmos/pkg/schema"
)

// Fetcher is an interface for fetching data from various sources.
type Fetcher interface {
	Fetch() ([]byte, error)
}

var ErrUnsupportedSource = errors.New("unsupported source type")

// getDataFetcher returns the appropriate Fetcher based on the input source.
func (d *dataFetcher) getDataFetcher(atmosConfig *schema.AtmosConfiguration, source string) (Fetcher, error) {
	switch {
	case strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://"):
		return &URLFetcher{URL: source, fileDownloader: downloader.NewGoGetterDownloader(atmosConfig)}, nil
	case strings.HasPrefix(source, "atmos://"):
		return &AtmosFetcher{Key: source}, nil
	default:
		if _, err := os.Stat(source); err == nil {
			return &FileFetcher{FilePath: source}, nil
		}
		return nil, ErrUnsupportedSource
	}
}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
type DataFetcher interface {
	GetData(atmosConfig *schema.AtmosConfiguration, source string) ([]byte, error)
}

type dataFetcher struct {
	fileDownloader downloader.FileDownloader
}

// NewDataFetcher creates a new dataFetcher instance.
func NewDataFetcher() *dataFetcher {
	return &dataFetcher{
		fileDownloader: downloader.NewGoGetterDownloader(&schema.AtmosConfiguration{}),
	}
}

// GetData returns the data based on source.
func (d *dataFetcher) GetData(atmosConfig *schema.AtmosConfiguration, source string) ([]byte, error) {
	fetcher, err := d.getDataFetcher(atmosConfig, source)
	if err != nil {
		return nil, err
	}
	return fetcher.Fetch()
}
