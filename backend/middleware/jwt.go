package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/nais/dp/backend/auth"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
)

func JWTValidatorMiddleware(oAuth2Config auth.Google) func(http.Handler) http.Handler {
	jwtValidator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		log.Fatalf("Creating JWT validator: %v", err)
	}
	return TokenValidatorMiddleware(jwtValidator, oAuth2Config.ClientID)
}

func MockJWTValidatorMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func CreateJWTValidator(google auth.Google) (jwt.Keyfunc, error) {
	if len(google.ClientID) == 0 || len(google.DiscoveryURL) == 0 {
		return nil, fmt.Errorf("missing required google configuration")
	}

	certificates, err := auth.FetchCertificates(google)
	if err != nil {
		return nil, fmt.Errorf("retrieving google certificates for token validation: %v", err)
	}

	return JWTValidator(certificates, google.ClientID), nil
}

func TokenValidatorMiddleware(validator *idtoken.Validator, clientId string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := jwtauth.TokenFromCookie(r)

			payload, err := validator.Validate(context.Background(), token, clientId)
			if err != nil {
				log.Errorf("validating token: %v", err)
				w.WriteHeader(http.StatusForbidden)
				_, err = fmt.Fprintf(w, "Unauthorized access: %s", err.Error())
				if err != nil {
					log.Errorf("Writing http response: %v", err)
				}
				return
			}

			// var groups []string
			// groupInterface := payload.Claims["groups"].(interface{})
			// groups = make([]string, len(groupInterface))
			// for i, v := range groupInterface {
			// 	groups[i] = v.(string)
			// }
			// r = r.WithContext(context.WithValue(r.Context(), "groups", groups))

			username := payload.Claims["email"].(string)
			r = r.WithContext(context.WithValue(r.Context(), "email", username))

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
