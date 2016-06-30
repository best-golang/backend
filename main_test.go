package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	app              *App

	userOffering		string
	userRequesting		string
	userRegionAdmin 	string
	userSuperAdmin		string

	regionID         	string
	offerID				string
	requestID			string
	matchingID			string
	notificationID		string
)

// init() will always be called before TestMain or Tests
func init() {
	fmt.Println("initiating...")

	// initialize and configure server
	app = InitApp()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}


type PromoteAdminPayload struct {
	Mail string
}
/*
func TestGetRegions(t *testing.T) {

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/regions", nil)
	if err != nil {
		t.Error("Error while trying to get regions", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	app.Router.ServeHTTP(resp, req)

	regions := parseResponseToArray(resp)
	regionID = regions[0]["ID"].(string)
}*/




// ----------------------------------------------------------------- AUTH

	// [X] Login a: N - check if returns JWT
	// [] Authorize : L
	// [] Renew Token : L - check if returns new JWT
	// [] Logout : L -



func LoginTest(t *testing.T, Email string, Password string, AssertCode int) string{
	loginParams := LoginPayload{
		Email,
		Password,
	}
	resp := app.Request("POST", "/auth", loginParams)


	// expected login to work 
	if AssertCode == 200 && resp.Code != 200 {
		t.Error("User login failed", resp.Body.String())
		return ""
	}

	// expecting login to fail
	if AssertCode == 400 {
		if resp.Code != 400 {
			t.Error(fmt.Printf("User login unexpected response %d", resp.Code))
		}
		return ""
	}

	// check if access token exists
	dat := parseResponse(resp)
	if dat["AccessToken"] == nil {
		t.Error("User Access Token is empty")
		return ""
	}

	// return token
	return dat["AccessToken"].(string)
}



// ----------------------------------------------------------------- USERS

	// [X] CreateUser : U
	// [] UpdateUser : S
	// [] ListUsers : A
	// [] GetUser : A
	// [] PromoteToSystemAdmin : S


func CreateUserTest(t *testing.T, Email string, Password string, Name string, AssertCode int) {
	// Create User
	createParams := CreateUserPayload{
		Name:          Name,
		PreferredName: Name + " Pref",
		Mail:          Email,
		Password:      Password,
		PhoneNumbers:  make([]string, 1),
	}
	resp := app.Request("POST", "/users", createParams)

	if AssertCode == 200 && resp.Code != 200{
		t.Error("User creation failed")
	}

	if AssertCode == 400 && resp.Code != 400 {
		t.Error("User should already exist")
	}
}


// ----------------------------------------------------------------- GROUPS

	// [] GetGroups : S
	// [] ListSystemAdmins : S


// ----------------------------------------------------------------- OFFERS

	// [X] CreateOffer - L
	// [X] GetOffer - C
	// [X] UpdateOffer - C


func CreateOfferTest(t *testing.T, jwt string, Name string, Location GeoLocation, Validity string, AssertCode int) string{
	plCreateOffer := CreateOfferPayload {
		Name,
		Location,
		[]string{},
		Validity,
	}

	// check if offer was created
	resp := app.RequestWithJWT("POST", "/offers", plCreateOffer, jwt)
	
	if AssertCode == 201 && resp.Code != 201 {
		t.Error("Could not CreateOffer")
		return ""
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error("CreateOffer should return BadRequest, but didnt")
		}
		return ""
	}

	data := parseResponse(resp)
	return data["ID"].(string)
}

func GetOfferTest(t *testing.T, jwt string, Offer string, AssertCode int) map[string]interface{}{
	resp := app.RequestWithJWT("GET", "/offers/" + Offer, nil, jwt)

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("Could not get offer")
		return map[string]interface{}{}
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("GetOffer should return BadRequest, but didnt"))
		}
		return map[string]interface{}{}
	}
	if AssertCode == 401 {
		if resp.Code != 401{
			t.Error(fmt.Printf("GetOffer should return Unauthorized, but didnt"))
		}
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	return data
}


func UpdateOfferTest(t *testing.T, jwt string, Offer string, Name string, Location GeoLocation, Validity string, Tags []string, Matched bool, AssertCode int) map[string]interface{}{
	updateOfferParams := UpdateOfferPayload {
		Name,
		Location,
		Tags,
		Validity,
		Matched,
	}

	resp := app.RequestWithJWT("PUT", "/offers/" + Offer, updateOfferParams, jwt)

	
	if AssertCode == 200 && resp.Code != 200 {
		t.Error("Could not get offer")
		return map[string]interface{}{}
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("GetOffer should return BadRequest, but didnt"))
		}
		return map[string]interface{}{}
	}
	if AssertCode == 401 {
		if resp.Code != 401{
			t.Error(fmt.Printf("GetOffer should return Unauthorized, but didnt"))
		}
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	return data
}


// ----------------------------------------------------------------- REQUESTS

	// [X] CreateRequest - L
	// [X] GetRequest - C
	// [X] UpdateRequest - C


func CreateRequestTest(t *testing.T, jwt string, Name string, Location GeoLocation, Validity string, Tags []string, AssertCode int) string {
	plCreateRequest := CreateRequestPayload {
		Name,
		Location,
		Tags,
		Validity,
	}

	resp := app.RequestWithJWT("POST", "/requests", plCreateRequest, jwt)

	if AssertCode == 201 && resp.Code != 201 {
		t.Error("Could not create request")
		return ""
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("CreateRequest should return BadRequest, but did return %d", resp.Code))
		}
		return ""
	}

	data := parseResponse(resp)
	return data["ID"].(string)
}

func GetRequestTest(t *testing.T, jwt string, Request string, AssertCode int) map[string]interface{}{
	resp := app.RequestWithJWT("GET", "/requests/" + Request, nil, jwt)

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("Could not get request")
		return map[string]interface{}{}
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("GetRequest should return BadRequest, but didnt"))
		}
		return map[string]interface{}{}
	}
	if AssertCode == 401 {
		if resp.Code != 401{
			t.Error(fmt.Printf("GetRequest should return Unauthorized, but didnt"))
		}
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	return data
}

func UpdateRequestTest(t *testing.T, jwt string, Request string, Name string, Location GeoLocation, Validity string, Tags []string, Matched bool, AssertCode int) map[string]interface{}{
	updateRequestParams := UpdateRequestPayload {
		Name,
		Location,
		Tags,
		Validity,
		Matched,
	}

	resp := app.RequestWithJWT("PUT", "/requests/" + Request, updateRequestParams, jwt)

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("Could not update request" + resp.Body.String())
		return map[string]interface{}{}
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("UpdateRequest should return BadRequest, but didnt"))
		}
		return map[string]interface{}{}
	}
	if AssertCode == 401 {
		if resp.Code != 401{
			t.Error(fmt.Printf("UpdateRequest should return Unauthorized, but didnt"))
		}
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	return data
}

// ----------------------------------------------------------------- MATCHINGS

	// [X] CreateMatching - A
	// [] GetMatching - C
	// [] UpdateMatching - C

func CreateMatchingTest(t *testing.T, jwt string, Region string, Offer string, Request string, AssertCode int) string {
	plCreateMatching := CreateMatchingPayload{
		Region,
		Request,
		Offer,
	}

	resp := app.RequestWithJWT("POST", "/matchings", plCreateMatching, jwt)
	
	if AssertCode == 201 && resp.Code != 201 {
		t.Error("Could not create matching")
		return ""
	}
	if AssertCode == 400 {
		if resp.Code != 400{
			t.Error(fmt.Printf("CreateMatching should return BadRequest, but did return %d", resp.Code))
		}
		return ""
	}
	if AssertCode == 401 {
		if resp.Code != 401 {
			t.Error("CreateMatching should return UnAuthorized but didnt")
		}
		return "" 
	}

	data := parseResponse(resp)
	return data["ID"].(string)
}

// ----------------------------------------------------------------- REGIONS

	// [X] CreateRegion - L
	// [] ListRegions - U
	// [X] GetRegion - U
	// [] UpdateRegion - A
	// [] ListOffersForRegion - A
	// [] ListRequestsForRegion - A
	// [] ListMatchingsForRegion - A
	// [X] PromoteUserToAdminForRegion - A
	// [] ListAdminsForRegion - A

func CreateRegionTest(t *testing.T, jwt string, Name string, Desc string, Locations []Location, AssertCode int) string {
	// create region
	plCreateRegion := CreateRegionPayload {
		Name,
		Desc,
		Boundaries{
			Locations,
		},
	}

	resp := app.RequestWithJWT("POST", "/regions", plCreateRegion, jwt)
	if AssertCode == 201 && resp.Code != 201{
		t.Error("could not create region")
		return ""
	}

	if AssertCode == 400 { 
		if resp.Code != 400 {
			t.Error("CreateRegion should return BadRequest but didnt")
		}
		return ""
	}

	// check if ID exists
	dat := parseResponse(resp)
	if dat["ID"] == nil {
		t.Error("CreateRegion did not return ID parameter")
		return ""
	}

	return dat["ID"].(string)
}

func GetRegionTest(t *testing.T, jwt string, Region string, AssertCode int) map[string]interface{} {
	resp := app.RequestWithJWT("GET", "/regions/" + Region, nil, jwt)

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("GetRegion failed")
		return map[string]interface{}{}
	}

	if AssertCode == 400 {
		if resp.Code != 400  {
			t.Error("GetRegion should return bad request, but didnt")
		}
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	if(data["ID"] != Region) {
		t.Error("Wrong region was returned")
	}

	return data
}

func PromoteUserToAdminForRegionTest(t *testing.T, jwt string, Email string, Region string, AssertCode int) bool{
	promoteParams := PromoteAdminPayload{Email}
	resp := app.RequestWithJWT("POST", "/regions/" + Region + "/admins", promoteParams, jwt)
	

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("Promoting User to Admin did not work, but should: ", resp.Body.String())
		return false
	}
	if AssertCode == 400 {
		if resp.Code != 400 {
			t.Error("PromoteUserToAdmin should return BadRequest but didnt")
		}
		return false 
	}
	if AssertCode == 401 {
		if resp.Code != 401 {
			t.Error("PromoteUserToAdmin should return UnAuthorized but didnt")
		}
		return false 
	}
	if AssertCode == 404 {
		if resp.Code != 404 {
			t.Error("PromoteUserToAdmin should return NotFound but didnt")
		}
		return false 
	}

	return true
}


// ------------------------------------------------------------------------------- ME

	// [X] GetMe - L
	// [] UpdateMe - L
	// [] ListUserOffers - L
	// [] ListUserRequests - L
	// [] ListUserMatchings - L

func GetMeTest(t *testing.T, jwt string, AssertCode int) map[string]interface{}{
	resp := app.RequestWithJWT("GET", "/me", nil, jwt)

	if AssertCode == 200 && resp.Code != 200 {
		t.Error("GetMe fail ", resp.Body.String())
		return map[string]interface{}{}
	}

	data := parseResponse(resp)
	return data
}

// ------------------------------------------------------------------------------- Notifications

	// [] ListNotifications - L
	// [] UpdateNotification - C



// ----------------------------------------------------- SCENARIO ALPHA

func TestSetupAlpha(t *testing.T) {
	fmt.Println("\n--------------------- SetupTestAlpha ---------------------\n")

	// HACKY : CreateUserTest AssertCode set to zero because no prior knowledge exists about database
	// response 200 and 400 could both be valid, but our testframework does not support multi-asserts
	// VALID: CreateUser + Login
	emailRegionAdmin := "regionadmin@test.org"
	CreateUserTest(t, "offering@test.org", "ICanOfferAllThemHelp666!", "OfferBoy", 0)
	userOffering = LoginTest(t, "offering@test.org", "ICanOfferAllThemHelp666!", 200)
	CreateUserTest(t, "requesting@test.org", "INeedAllThemHelp666!", "RequestDude", 0)
	userRequesting = LoginTest(t, "requesting@test.org", "INeedAllThemHelp666!", 200)
	CreateUserTest(t, emailRegionAdmin, "LetMeAdminAllYourHelp666!", "AdminMate", 0)
	userRegionAdmin = LoginTest(t, emailRegionAdmin, "LetMeAdminAllYourHelp666!", 200)

	// VALID GetMe
	regionAdminResp := GetMeTest(t, userRegionAdmin, 200)
	if regionAdminResp["Mail"] != emailRegionAdmin {
		t.Error("CreateUser followed by GetUser: comparing email for region admin failed")
	}

	// INVALID: Login superadmin
	LoginTest(t, "admin@example.org", "nonononooo", 400)
	// VALID: Login superadmin
	userSuperAdmin = LoginTest(t, "admin@example.org", "CaTUstrophyAdmin123$", 200)

	// INVALID: CreateRegion
	CreateRegionTest(t, userRegionAdmin, "", "", []Location{}, 400)
	// VALID: CreateRegion
	regionName := "Milkshake Region"
	regionID = CreateRegionTest(t, userRegionAdmin, regionName, "Ma Region brings all the boys in the yard", 
		[]Location{ Location{ 10.0, 0.0, }, Location{ 11.0, 0.0, }, Location{ 10.0, 1.0, }, Location{ 10.0, 0.0, }, }, 201)

	// INVALID: GetRegion 
	GetRegionTest(t, userOffering, regionID + "a", 400)
	// VALID: GetRegion 
	region := GetRegionTest(t, userOffering, regionID, 200)
	// compare region names of Create & Get
	if region["Name"] != regionName {
		t.Error("CreateRegion followed by GetRegion returns not the same value for field Region.Name")
	}

	// INVALID: PromoteUserToAdminForRegion
	PromoteUserToAdminForRegionTest(t, userSuperAdmin, "", regionID, 400)
	PromoteUserToAdminForRegionTest(t, userRegionAdmin, emailRegionAdmin, regionID, 401)
	PromoteUserToAdminForRegionTest(t, userSuperAdmin, "nobody@donotexist.com", regionID, 404)
	// VALID: PromoteUserToAdminForRegion
	PromoteUserToAdminForRegionTest(t, userSuperAdmin, emailRegionAdmin, regionID, 200)
}

func TestMatchingAlpha(t *testing.T) {
	fmt.Println("\n--------------------- MatchingTestAlpha ---------------------\n")

	// INVALID CreateOffer
	CreateOfferTest(t, userOffering, "", GeoLocation{10.2, .0}, "2017-11-01T22:08:41+00:00", 400)
	// VALID CreateOffer
	offerID = CreateOfferTest(t, userOffering, "Milk x10", GeoLocation{10.2, .0}, "2017-11-01T22:08:41+00:00", 201)

	// INVALID GetOffer
	GetOfferTest(t, userOffering, offerID + "a", 400)
	GetOfferTest(t, userRequesting, offerID, 401)
	// VALID GetOffer
	offer := GetOfferTest(t, userOffering, offerID, 200)
	GetOfferTest(t, userRegionAdmin, offerID, 200)

	// INVALID CreateRequest
	CreateRequestTest(t, userRequesting, "", GeoLocation{10.3, 0.2}, "2016-11-01T22:08:41+00:00", []string{}, 400)
	// VALID CreateRequest
	requestName := "Me thirsty"
	requestID = CreateRequestTest(t, userRequesting, requestName, GeoLocation{10.3, 0.2}, "2016-11-01T22:08:41+00:00", []string{"Food"}, 201)

	// INVALID GetRequest
	GetRequestTest(t, userRequesting, requestID + "a", 400)
	GetRequestTest(t, userOffering, requestID, 401)
	// VALID GetRequest
	request := GetRequestTest(t, userRequesting, requestID, 200)
	GetRequestTest(t, userRegionAdmin, requestID, 200)
	GetRequestTest(t, userSuperAdmin, requestID, 200)
	if request["Name"] != requestName {
		t.Error("CreateRequest followed by GetRequest dont seem to return same values")
	}

	// INVALID CreateMatching
	CreateMatchingTest(t, userOffering, regionID, offerID, requestID, 401) // TODO : somehow this seems to work just fine
	CreateMatchingTest(t, userRegionAdmin, "", offerID, requestID, 400)
	// VALID CreateMatching
	matchingID = CreateMatchingTest(t, userRegionAdmin, regionID, offerID, requestID, 201)


	// VALID UpdateOfferTest
	offer = UpdateOfferTest(t, userOffering, offerID,
		offer["Name"].(string) + " Updated",  
		GeoLocation{0,0}, 
		offer["ValidityPeriod"].(string), 
		[]string{"Tool"}, false, 
		200,
	)
	// GetOffer and check if updates were propagated
	offer =  GetOfferTest(t, userOffering, offerID, 200)
	// check if tags were updated
	tags := offer["Tags"].([]interface{})
	if len(tags) == 0 {
		t.Error("UpdateOffer failed to update")
	}
	if tags[0].(map[string]interface{})["Name"].(string) != "Tool" {
		t.Error("UpdateOffer failed to insert tag Tool")
	}

	// check if location was updated
	location := offer["Location"].(map[string]interface{})
	if location["lat"].(float64) != 0 || location["lng"].(float64) != 0{
		//t.Error("UpdateOfffer didnt update Location")
	}


	// VALID UpdateRequestTest
	request = UpdateRequestTest(t, userRequesting, requestID,
		request["Name"].(string) + " Updated",  
		GeoLocation{0,0}, 
		request["ValidityPeriod"].(string), 
		[]string{"Water"}, false, 
		200,
	)

	// GetRequest and check if updates were propagated
	request =  GetRequestTest(t, userRequesting, requestID, 200)
	// check if tags were updated
	tags = request["Tags"].([]interface{})
	if len(tags) == 0 {
		t.Error("UpdateRequest failed to update")
	}
	if tags[0].(map[string]interface{})["Name"].(string) != "Water" {
		t.Error("UpdateRequest failed to insert tag Water")
	}

	// check if location was updated
	location = request["Location"].(map[string]interface{})
	if location["lat"].(float64) != 0 || location["lng"].(float64) != 0{
		//t.Error("UpdateRequest didnt update Location")
	}

}
