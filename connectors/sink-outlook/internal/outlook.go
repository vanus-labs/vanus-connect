package internal

import (
	"context"
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type gmailService struct {
	client *msgraphsdk.GraphServiceClient
}

func NewOutlookService(oauthCfg OAuth) (*gmailService, error) {
	svc := &gmailService{}
	client, err := svc.init(oauthCfg)
	if err != nil {
		return nil, err
	}
	svc.client = client
	return svc, nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token",
		oauth2.SetAuthURLParam("response_mode", "query"),
	)
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

func (svc *gmailService) init(oauthCfg OAuth) (*msgraphsdk.GraphServiceClient, error) {
	config := &oauth2.Config{
		ClientID:     oauthCfg.ClientID,
		ClientSecret: oauthCfg.ClientSecret,
		Endpoint:     microsoft.AzureADEndpoint(oauthCfg.TenantID),
		Scopes:       []string{"offline_access", "Mail.Send", "User.Read", "Mail.ReadWrite"},
	}
	var token *oauth2.Token
	if oauthCfg.RefreshToken != "" {
		token = oauthCfg.GetToken()
	} else {
		var err error
		token, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
	}
	tokenSource := config.TokenSource(context.Background(), token)
	cred := &AzureToken{
		Token: tokenSource,
	}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"Mail.Send"})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (svc *gmailService) Send(ctx context.Context, em *EmailMessage) error {
	requestBody := users.NewItemSendMailPostRequestBody()
	message := models.NewMessage()
	message.SetSubject(&em.Subject)

	body := models.NewItemBody()
	contentType := em.GetBodyType()
	body.SetContentType(&contentType)
	body.SetContent(&em.Body)
	message.SetBody(body)

	recipient := models.NewRecipient()
	emailAddress := models.NewEmailAddress()
	emailAddress.SetAddress(&em.To)
	recipient.SetEmailAddress(emailAddress)
	message.SetToRecipients([]models.Recipientable{recipient})

	requestBody.SetMessage(message)
	save := true
	requestBody.SetSaveToSentItems(&save)
	err := svc.client.Me().SendMail().Post(ctx, requestBody, nil)
	return err
}
