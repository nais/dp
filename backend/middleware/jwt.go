package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/nais/dp/backend/auth"
	log "github.com/sirupsen/logrus"
)

func mockJWTValidatorMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func JWTValidatorMiddleware(discoveryURL, clientID string, mock bool, azureGroups auth.AzureGroups, teamUUIDs map[string]string) func(http.Handler) http.Handler {
	if mock {
		return mockJWTValidatorMiddleware()
	}
	certificates, err := auth.FetchCertificates(discoveryURL)
	if err != nil {
		log.Fatalf("Fetching signing certificates from IDP: %v", err)
	}
	validator := JWTValidator(certificates, clientID)

	return TokenValidatorMiddleware(validator, azureGroups, teamUUIDs)
}

func TokenValidatorMiddleware(jwtValidator jwt.Keyfunc, azureGroups auth.AzureGroups, teamUUIDs map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var claims jwt.MapClaims

			token := jwtauth.TokenFromCookie(r)

			_, err := jwt.ParseWithClaims(token, &claims, jwtValidator)
			if err != nil {
				log.Errorf("parsing token: %v", err)
				w.WriteHeader(http.StatusForbidden)
				_, err = fmt.Fprintf(w, "Unauthorized access: %s", err.Error())
				if err != nil {
					log.Errorf("Writing http response: %v", err)
				}
				return
			}

			email := strings.ToLower(claims["preferred_username"].(string))
			r = r.WithContext(context.WithValue(r.Context(), "preferred_username", email))

			groups, err := azureGroups.GetGroupsForUser(r.Context(), token, email)
			if err != nil {
				log.Errorf("getting groups for user: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err = fmt.Fprintf(w, "Unauthorized access: %s", err.Error())
				if err != nil {
					log.Errorf("Writing http response: %v", err)
				}
				return
			}

			teams := make([]string, 0)

			for _, uuid := range groups {
				if _, found := teamUUIDs[uuid]; found {
					teams = append(teams, teamUUIDs[uuid])
				}
			}
			r = r.WithContext(context.WithValue(r.Context(), "teams", teams))

			next.ServeHTTP(w, r)
		})
	}
}

func JWTValidator(certificates map[string]auth.CertificateList, audience string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		var certificateList auth.CertificateList
		var kid string
		var ok bool

		if claims, ok := token.Claims.(*jwt.MapClaims); !ok {
			return nil, fmt.Errorf("unable to retrieve claims from token")
		} else {
			if valid := claims.VerifyAudience(audience, true); !valid {
				return nil, fmt.Errorf("the token is not valid for this application")
			}
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if kid, ok = token.Header["kid"].(string); !ok {
			return nil, fmt.Errorf("field 'kid' is of invalid type %T, should be string", token.Header["kid"])
		}

		if certificateList, ok = certificates[kid]; !ok {
			return nil, fmt.Errorf("kid '%s' not found in certificate list", kid)
		}

		for _, certificate := range certificateList {
			return certificate.PublicKey, nil
		}

		return nil, fmt.Errorf("no certificate candidates for kid '%s'", kid)
	}
}
