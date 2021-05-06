package auth

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type CertificateList []*x509.Certificate

type Google struct {
	DiscoveryURL string
	ClientID     string
	ClientSecret string
}

func FetchCertificates(google Google) (map[string]CertificateList, error) {
	log.Infof("Discover Google signing certificates from %s", google.DiscoveryURL)
	googleKeyDiscovery, err := DiscoverURL(google.DiscoveryURL)
	if err != nil {
		return nil, err
	}

	log.Infof("Decoding certificates for %d keys", len(googleKeyDiscovery.Keys))
	googleCertificates, err := googleKeyDiscovery.Map()
	if err != nil {
		return nil, err
	}
	return googleCertificates, nil
}

type EncodedCertificate string

type WellKnownResponse struct {
	JWKSURI string `json:"jwks_uri"`
}

type KeyDiscovery struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Kid string               `json:"kid"`
	X5c []EncodedCertificate `json:"x5c"`
}

func DiscoverURL(url string) (*KeyDiscovery, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var wellKnownResponse WellKnownResponse
	err = json.NewDecoder(response.Body).Decode(&wellKnownResponse)

	jwksResponse, err := http.Get(wellKnownResponse.JWKSURI)
	if err != nil {
		return nil, err
	}

	return Discover(jwksResponse.Body)
}

func Discover(reader io.Reader) (*KeyDiscovery, error) {
	document, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	keyDiscovery := &KeyDiscovery{}
	err = json.Unmarshal(document, keyDiscovery)

	return keyDiscovery, err
}

// Transform a KeyDiscovery object into a dictionary with "kid" as key
// and lists of decoded X509 certificates as values.
//
// Returns an error if any certificate does not decode.
func (k *KeyDiscovery) Map() (result map[string]CertificateList, err error) {
	result = make(map[string]CertificateList)

	for _, key := range k.Keys {
		certList := make(CertificateList, 0)
		for _, encodedCertificate := range key.X5c {
			certificate, err := encodedCertificate.Decode()
			if err != nil {
				return nil, err
			}
			certList = append(certList, certificate)
		}
		result[key.Kid] = certList
	}

	return
}

// Decode a base64 encoded certificate into a X509 structure.
func (c EncodedCertificate) Decode() (*x509.Certificate, error) {
	stream := strings.NewReader(string(c))
	decoder := base64.NewDecoder(base64.StdEncoding, stream)
	key, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(key)
}
