package main

import (
	"fmt"
	"os"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	accountsid string
	authtoken  string
	fromphone  string
	tophone    string
	client     *twilio.RestClient
)

func init() {
	initializers.Getenv()

	accountsid = os.Getenv("ACCOUNT_SID")
	authtoken = os.Getenv("AUTH_TOKEN")
	fromphone = os.Getenv("FROM_PHONE")
	tophone = os.Getenv("TO_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountsid,
		Password: authtoken,
	})
}

func Sendmessege(msg string) {
	params := openapi.CreateMessageParams{}
	params.SetTo(tophone)
	params.SetFrom(fromphone)
	params.SetBody(msg)
	response, err := client.Api.CreateMessage(&params)
	if err != nil {
		fmt.Println("error creating messege", err.Error())
		return
	}
	fmt.Printf("messege SID:%s\n", *response.Sid)
}
func main(){
	msg:=fmt.Sprintf(os.Getenv("MSG"),"safwan") 
	Sendmessege(msg)
}
