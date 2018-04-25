package provider

import "fmt"

type errMediaNotSupported error

func mediaNotSupported(api string) errMediaNotSupported {
	return errMediaNotSupported(fmt.Errorf("media not supported %v", api))
}

// IsErrMediaNotSupported return true if the error dictates taht the media
// type was not supported by the scraper
func IsErrMediaNotSupported(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(errMediaNotSupported); ok {
		return true
	}
	return false
}
