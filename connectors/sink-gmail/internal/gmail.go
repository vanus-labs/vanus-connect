package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type gmailService struct {
	client *gmail.Service
}

func NewGmailService(credentialsJSON string, oauthCfg *OAuth) (*gmailService, error) {
	svc := &gmailService{}
	client, err := svc.init(credentialsJSON, oauthCfg)
	if err != nil {
		return nil, err
	}
	svc.client = client
	return svc, nil
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	tok, err := getTokenFromWeb(config)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (svc *gmailService) init(credentialsJSON string, oauthCfg *OAuth) (*gmail.Service, error) {
	var opts []option.ClientOption
	if oauthCfg != nil {
		config := oauth2.Config{
			ClientID:     oauthCfg.ClientID,
			ClientSecret: oauthCfg.ClientSecret,
			Endpoint:     google.Endpoint,
		}
		tokenSource := config.TokenSource(context.Background(), oauthCfg.GetToken())
		opts = append(opts, option.WithTokenSource(tokenSource))
	} else {
		config, err := google.ConfigFromJSON([]byte(credentialsJSON), gmail.GmailSendScope)
		if err != nil {
			return nil, err
		}
		c, err := getClient(config)
		if err != nil {
			return nil, err
		}
		opts = append(opts, option.WithHTTPClient(c))
	}
	client, err := gmail.NewService(context.Background(), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "new gmail service error")
	}
	return client, nil
}

func (svc *gmailService) Send(em *EmailMessage) error {
	var bf bytes.Buffer
	bf.WriteString("To: ")
	bf.WriteString(em.Recipients)
	bf.WriteString("\r\n")

	bf.WriteString("Subject: ")
	bf.WriteString(em.Subject)
	bf.WriteString("\r\n")

	bf.WriteString("\r\n")
	bf.WriteString(em.Body)
	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(bf.Bytes()),
	}
	messageSend := svc.client.Users.Messages.Send("me", message)
	_, err := messageSend.Do()
	return err
}
