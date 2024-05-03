package handlers

import (
	"chronokeep/results/types"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const CERTIFICATE_DOWNLOAD_SLEEP_MILLISECONDS = 150

// APIError holds information on an error from the API
type APIError struct {
	Message string `json:"message,omitempty"`
}

func getAPIError(c echo.Context, code int, message string, err error) error {
	log.WithFields(log.Fields{
		"message": message,
		"error":   err,
		"code":    code,
	}).Error("API Error.")
	return c.JSON(code, APIError{Message: message})
}

func retrieveKey(r *http.Request) (*string, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("unknown authorization header")
	}
	return &strArr[1], nil
}

func verifyToken(r *http.Request) (*types.Account, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("unknown authorization header")
	}
	token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("claims not set or token is not valid")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token claims")
	}
	account, err := database.GetAccount(email)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.Token != strArr[1] || account.Token == "" {
		return nil, errors.New("token no longer valid")
	}
	return account, nil
}

func createTokens(email string) (*string, *string, error) {
	// Create token
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(expirationWindow).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(config.SecretKey))
	if err != nil {
		return nil, nil, err
	}
	// Create refresh token
	claims = jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(refreshWindow).Unix()
	r := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh, err := r.SignedString([]byte(config.RefreshKey))
	if err != nil {
		return nil, nil, err
	}
	return &token, &refresh, nil
}

func CreateCertificate(name string, event string, timeString string, date string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(
		ctx,
		chromedp.Tasks{
			chromedp.Navigate("about:blank"),
			chromedp.ActionFunc(
				func(ctx context.Context) error {
					frameTree, err := page.GetFrameTree().Do(ctx)
					if err != nil {
						return err
					}
					return page.SetDocumentContent(frameTree.Frame.ID, GetCertificateHTML(name, event, timeString, date)).Do(ctx)
				},
			),
			chromedp.Sleep(CERTIFICATE_DOWNLOAD_SLEEP_MILLISECONDS * time.Millisecond),
			chromedp.FullScreenshot(&buf, 90),
		}); err != nil {
		return nil, err
	}
	return buf, nil
}

func GetCertificateHTML(name string, event string, time string, date string) string {
	return fmt.Sprintf(
		"<html>"+
			"<head></head>"+
			"<body style='width:800;height:565;padding:0px;background-image:url(\"%s\");background-size:cover;'>"+
			"<div style='margin:0px;width:800px;height:565px;position:relative;'>"+
			"<div style='width:100%%;margin:0;position:absolute;top:50%%;-ms-transform:translateY(-50%%);transform:translateY(-50%%);'>"+
			"<div style='font-size:60px;text-align:center;font-weight:bold;'>%s</div>"+
			"<div style='font-size:30px;text-align:center;margin-left:100px;width:600px;'>finished the %s with a time of</div>"+
			"<div style='font-size:60px;text-align:center;font-weight:bold;'>%s</div>"+
			"<div style='font-size:30px;text-align:center;'>on this day of %s</div>"+
			"</div>"+
			"</div>"+
			"</body>"+
			"</html>",
		config.CertificateURL,
		name,
		event,
		time,
		date,
	)
}
