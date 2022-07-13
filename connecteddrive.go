package connected_drive

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	iosUserAgent          = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Mobile/15E148 Safari/604.1"
	androidUserAgent      = "android(SP1A.210812.016.C1);bmw;2.5.2(14945)"
	authUrl               = "https://customer.bmwgroup.com/gcdm/oauth/authenticate"
	authTokenUrl          = "https://customer.bmwgroup.com/gcdm/oauth/token"
	clientId              = "31c357a0-7a1d-4590-aa99-33b97244d048"
	clientPassword        = "c0e3393d-70a2-4f6f-9d3c-8530af64d552"
	vehiclesRequestUrl    = "https://cocoapi.bmwgroup.com/eadrax-vcs/v1/vehicles?apptimezone=%d&appDateTime=%d&tireGuardMode=ENABLED"
	contentTypeUrlEncoded = "application/x-www-form-urlencoded; charset=UTF-8"
	contentTypeJson       = "application/json; charset=UTF-8"
)

type auth struct {
	Token        string `json:"access_token"`
	Expires      int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
}

type Vehicle struct {
	Vin            string `json:"vin"`
	Model          string `json:"model"`
	Year           int    `json:"year"`
	Brand          string `json:"brand"`
	HeadUnit       string `json:"headUnit"`
	IsLscSupported bool   `json:"isLscSupported"`
	DriveTrain     string `json:"driveTrain"`
	PuStep         string `json:"puStep"`
	IStep          string `json:"iStep"`
	TelematicsUnit string `json:"telematicsUnit"`
	HmiVersion     string `json:"hmiVersion"`
	BodyType       string `json:"bodyType"`
	A4AType        string `json:"a4aType"`
	Capabilities   struct {
		IsRemoteServicesBookingRequired    bool `json:"isRemoteServicesBookingRequired"`
		IsRemoteServicesActivationRequired bool `json:"isRemoteServicesActivationRequired"`
		Lock                               struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"lock"`
		Unlock struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"unlock"`
		Lights struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"lights"`
		Horn struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"horn"`
		VehicleFinder struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"vehicleFinder"`
		SendPoi struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
		} `json:"sendPoi"`
		LastStateCall struct {
			IsNonLscFeatureEnabled bool   `json:"isNonLscFeatureEnabled"`
			LscState               string `json:"lscState"`
		} `json:"lastStateCall"`
		ClimateNow struct {
			IsEnabled                   bool   `json:"isEnabled"`
			IsPinAuthenticationRequired bool   `json:"isPinAuthenticationRequired"`
			ExecutionMessage            string `json:"executionMessage"`
			ExecutionPopup              struct {
				ExecutionMessage    string `json:"executionMessage"`
				PopupType           string `json:"popupType"`
				Title               string `json:"title"`
				PrimaryButtonText   string `json:"primaryButtonText"`
				SecondaryButtonText string `json:"secondaryButtonText"`
				IconId              int    `json:"iconId"`
			} `json:"executionPopup"`
			ExecutionStopPopup struct {
				ExecutionMessage string `json:"executionMessage"`
				Title            string `json:"title"`
			} `json:"executionStopPopup"`
		} `json:"climateNow"`
		IsRemoteHistorySupported  bool `json:"isRemoteHistorySupported"`
		CanRemoteHistoryBeDeleted bool `json:"canRemoteHistoryBeDeleted"`
		ClimateTimer              struct {
			IsEnabled                   bool `json:"isEnabled"`
			IsPinAuthenticationRequired bool `json:"isPinAuthenticationRequired"`
			Tile                        struct {
				IconId      int    `json:"iconId"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"tile"`
			Page struct {
				PrimaryButtonText   string `json:"primaryButtonText"`
				SecondaryButtonText string `json:"secondaryButtonText"`
				Title               string `json:"title"`
				Subtitle            string `json:"subtitle"`
				Description         string `json:"description"`
			} `json:"page"`
			IsToggleEnabled bool `json:"isToggleEnabled"`
		} `json:"climateTimer"`
		IsChargingHistorySupported       bool          `json:"isChargingHistorySupported"`
		IsScanAndChargeSupported         bool          `json:"isScanAndChargeSupported"`
		IsDCSContractManagementSupported bool          `json:"isDCSContractManagementSupported"`
		IsBmwChargingSupported           bool          `json:"isBmwChargingSupported"`
		IsMiniChargingSupported          bool          `json:"isMiniChargingSupported"`
		IsChargeNowForBusinessSupported  bool          `json:"isChargeNowForBusinessSupported"`
		IsDataPrivacyEnabled             bool          `json:"isDataPrivacyEnabled"`
		IsChargingPlanSupported          bool          `json:"isChargingPlanSupported"`
		IsChargingPowerLimitEnable       bool          `json:"isChargingPowerLimitEnable"`
		IsChargingTargetSocEnable        bool          `json:"isChargingTargetSocEnable"`
		IsChargingLoudnessEnable         bool          `json:"isChargingLoudnessEnable"`
		IsChargingSettingsEnabled        bool          `json:"isChargingSettingsEnabled"`
		IsChargingHospitalityEnabled     bool          `json:"isChargingHospitalityEnabled"`
		IsEvGoChargingSupported          bool          `json:"isEvGoChargingSupported"`
		IsFindChargingEnabled            bool          `json:"isFindChargingEnabled"`
		IsCustomerEsimSupported          bool          `json:"isCustomerEsimSupported"`
		IsCarSharingSupported            bool          `json:"isCarSharingSupported"`
		IsEasyChargeSupported            bool          `json:"isEasyChargeSupported"`
		IsSustainabilitySupported        bool          `json:"isSustainabilitySupported"`
		SpecialThemeSupport              []interface{} `json:"specialThemeSupport"`
		IsRemoteParkingSupported         bool          `json:"isRemoteParkingSupported"`
	} `json:"capabilities"`
	ConnectedDriveServices []interface{} `json:"connectedDriveServices"`
	Properties             struct {
		LastUpdatedAt    time.Time `json:"lastUpdatedAt"`
		InMotion         bool      `json:"inMotion"`
		AreDoorsLocked   bool      `json:"areDoorsLocked"`
		OriginCountryISO string    `json:"originCountryISO"`
		AreDoorsClosed   bool      `json:"areDoorsClosed"`
		AreDoorsOpen     bool      `json:"areDoorsOpen"`
		AreWindowsClosed bool      `json:"areWindowsClosed"`
		DoorsAndWindows  struct {
			Doors struct {
				DriverFront    string `json:"driverFront"`
				DriverRear     string `json:"driverRear"`
				PassengerFront string `json:"passengerFront"`
				PassengerRear  string `json:"passengerRear"`
			} `json:"doors"`
			Windows struct {
				DriverFront    string `json:"driverFront"`
				DriverRear     string `json:"driverRear"`
				PassengerFront string `json:"passengerFront"`
				PassengerRear  string `json:"passengerRear"`
			} `json:"windows"`
			Trunk string `json:"trunk"`
			Hood  string `json:"hood"`
		} `json:"doorsAndWindows"`
		IsServiceRequired bool `json:"isServiceRequired"`
		FuelLevel         struct {
			Value int    `json:"value"`
			Units string `json:"units"`
		} `json:"fuelLevel"`
		CombustionRange struct {
			Distance struct {
				Value int    `json:"value"`
				Units string `json:"units"`
			} `json:"distance"`
		} `json:"combustionRange"`
		CheckControlMessages []interface{} `json:"checkControlMessages"`
		ServiceRequired      []struct {
			Type     string    `json:"type"`
			Status   string    `json:"status"`
			DateTime time.Time `json:"dateTime"`
			Distance struct {
				Value int    `json:"value"`
				Units string `json:"units"`
			} `json:"distance,omitempty"`
		} `json:"serviceRequired"`
		VehicleLocation struct {
			Coordinates struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"coordinates"`
			Address struct {
				Formatted string `json:"formatted"`
			} `json:"address"`
			Heading int `json:"heading"`
		} `json:"vehicleLocation"`
		ClimateControl struct {
		} `json:"climateControl"`
	} `json:"properties"`
	IsMappingPending     bool `json:"isMappingPending"`
	IsMappingUnconfirmed bool `json:"isMappingUnconfirmed"`
	DriverGuideInfo      struct {
		Title            string `json:"title"`
		AndroidAppScheme string `json:"androidAppScheme"`
		IosAppScheme     string `json:"iosAppScheme"`
		AndroidStoreUrl  string `json:"androidStoreUrl"`
		IosStoreUrl      string `json:"iosStoreUrl"`
	} `json:"driverGuideInfo"`
	ThemeSpecs struct {
		VehicleStatusBackgroundColor struct {
			Red   int `json:"red"`
			Green int `json:"green"`
			Blue  int `json:"blue"`
		} `json:"vehicleStatusBackgroundColor"`
	} `json:"themeSpecs"`
	Status struct {
		LastUpdatedAt  time.Time `json:"lastUpdatedAt"`
		CurrentMileage struct {
			Mileage          int    `json:"mileage"`
			Units            string `json:"units"`
			FormattedMileage string `json:"formattedMileage"`
		} `json:"currentMileage"`
		Issues struct {
		} `json:"issues"`
		DoorsGeneralState                string `json:"doorsGeneralState"`
		CheckControlMessagesGeneralState string `json:"checkControlMessagesGeneralState"`
		DoorsAndWindows                  []struct {
			IconId       int    `json:"iconId"`
			Title        string `json:"title"`
			State        string `json:"state"`
			Criticalness string `json:"criticalness"`
		} `json:"doorsAndWindows"`
		CheckControlMessages []struct {
			Criticalness string `json:"criticalness"`
			IconId       int    `json:"iconId"`
			Title        string `json:"title"`
			State        string `json:"state"`
		} `json:"checkControlMessages"`
		RequiredServices []struct {
			Id              string `json:"id"`
			Title           string `json:"title"`
			IconId          int    `json:"iconId"`
			LongDescription string `json:"longDescription"`
			Subtitle        string `json:"subtitle"`
			Criticalness    string `json:"criticalness"`
		} `json:"requiredServices"`
		RecallMessages    []interface{} `json:"recallMessages"`
		RecallExternalUrl interface{}   `json:"recallExternalUrl"`
		FuelIndicators    []struct {
			SecondaryBarValue int         `json:"secondaryBarValue"`
			InfoIconId        int         `json:"infoIconId"`
			InfoLabel         string      `json:"infoLabel"`
			RangeIconId       int         `json:"rangeIconId"`
			RangeUnits        string      `json:"rangeUnits"`
			RangeValue        string      `json:"rangeValue"`
			LevelIconId       int         `json:"levelIconId"`
			IsCircleIcon      bool        `json:"isCircleIcon"`
			IconOpacity       string      `json:"iconOpacity"`
			ChargingType      interface{} `json:"chargingType"`
			MainBarValue      int         `json:"mainBarValue"`
			ShowsBar          bool        `json:"showsBar"`
			LevelUnits        string      `json:"levelUnits"`
			LevelValue        string      `json:"levelValue"`
			IsInaccurate      bool        `json:"isInaccurate"`
		} `json:"fuelIndicators"`
		TimestampMessage string `json:"timestampMessage"`
	} `json:"status"`
	ExFactoryPUStep string `json:"exFactoryPUStep"`
	ExFactoryILevel string `json:"exFactoryILevel"`
	Euiccid         string `json:"euiccid"`
}
type Vehicles []*Vehicle

type Client struct {
	username   string
	password   string
	authStore  io.ReadWriter
	auth       *auth
	httpClient *http.Client
}

func NewClient(username string, password string, authStore io.ReadWriter, httpClient *http.Client) *Client {
	httpClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Client{username: username, password: password, authStore: authStore, httpClient: httpClient}
}

func (c *Client) GetVehicles(ctx context.Context) (vehicles Vehicles, err error) {
	err = c.refreshAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing auth while fetching vehicles list: %w", err)
	}

	_, offset := time.Now().Local().Zone()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(vehiclesRequestUrl, offset, time.Now().Unix()), nil)
	req.Header = http.Header{
		"x-user-agent":  {androidUserAgent},
		"Authorization": {fmt.Sprintf("Bearer %s", c.auth.Token)},
		"Content-Type":  {contentTypeJson},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching vehicles list: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&vehicles)
	if err != nil || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error decoding vehicles list: %w, %v", err, resp)
	}

	return vehicles, nil
}

func (c *Client) saveAuth() {
	if c.authStore == nil {
		return
	}

	e := json.NewEncoder(c.authStore)
	err := e.Encode(c.auth)
	if err != nil {
		return
	}
}

func (c *Client) loadAuth() {
	if c.authStore == nil {
		return
	}

	d := json.NewDecoder(c.authStore)
	err := d.Decode(&c.auth)
	if err != nil {
		return
	}
}

func (c *Client) refreshAuth(ctx context.Context) error {
	if c.auth == nil {
		c.loadAuth()
		if c.auth == nil {
			c.auth = new(auth)
		}
	}

	if c.auth.Token == "" {
		err := c.getToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth token: %w", err)
		}
	} else if c.auth.Token != "" && time.Now().Unix() > c.auth.Expires {
		err := c.getRefreshToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get refresh token: %w", err)
		}
	}

	c.saveAuth()

	return nil
}

func (c *Client) getToken(ctx context.Context) error {
	// stage1
	state := randomString(22)
	codeChallenge := randomString(86)
	form := url.Values{
		"client_id":             {clientId},
		"response_type":         {"code"},
		"scope":                 {"openid profile email offline_access smacc vehicle_data perseus dlm svds cesim vsapi remote_services fupo authenticate_user"},
		"redirect_uri":          {"com.bmw.connected://oauth"},
		"state":                 {state},
		"nonce":                 {"login_nonce"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"plain"},
		"username":              {c.username},
		"password":              {c.password},
		"grant_type":            {"authorization_code"},
	}
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", authUrl, body)
	if err != nil {
		return fmt.Errorf("getToken stage1 error: can't create request: %w", err)
	}

	req.Header = http.Header{
		"Content-Type": {contentTypeUrlEncoded},
		"User-Agent":   {iosUserAgent},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return fmt.Errorf("getToken stage1 error: can't send request: %w, %v", err, resp)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var stage1Response struct {
		RedirectTo string `json:"redirect_to,omitempty"`
	}
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&stage1Response)
	if err != nil {
		return fmt.Errorf("getToken stage1 error: can't decode response: %w", err)
	}

	re := regexp.MustCompile(`(?im).*authorization=(.*)`)
	authString := ""
	if len(re.FindStringSubmatch(stage1Response.RedirectTo)) > 1 {
		authString = re.FindStringSubmatch(stage1Response.RedirectTo)[1]
	}
	if authString == "" {
		return fmt.Errorf("getToken stage1 error: can't find redirect link")
	}

	// stage 2
	form = url.Values{
		"client_id":             {clientId},
		"response_type":         {"code"},
		"scope":                 {"openid profile email offline_access smacc vehicle_data perseus dlm svds cesim vsapi remote_services fupo authenticate_user"},
		"redirect_uri":          {"com.bmw.connected://oauth"},
		"state":                 {state},
		"nonce":                 {"login_nonce"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"plain"},
		"authorization":         {authString},
	}
	req, err = http.NewRequestWithContext(ctx, "POST", authUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("getToken stage2 error: can't create request: %w", err)
	}
	req.Header = http.Header{
		"Content-Type": {contentTypeUrlEncoded},
		"User-Agent":   {iosUserAgent},
		"Cookie":       {fmt.Sprintf("GCDMSSO=%s", authString)},
	}

	resp, err = c.httpClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return fmt.Errorf("getToken stage2 error: can't send request: %w, %v", err, resp)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var code string
	re = regexp.MustCompile(`(?im).*code=(.*?)&`)
	for k, vs := range resp.Header {
		if strings.ToLower(k) == "location" {
			for _, v := range vs {
				if re.MatchString(v) && len(re.FindStringSubmatch(v)) > 1 {
					code = re.FindStringSubmatch(v)[1]

					break
				}
			}
		}

		if code != "" {
			break
		}
	}
	if code == "" {
		return fmt.Errorf("getToken stage2 error: can't find code")
	}

	// stage 3
	form = url.Values{
		"code":          {code},
		"code_verifier": {codeChallenge},
		"redirect_uri":  {"com.bmw.connected://oauth"},
		"grant_type":    {"authorization_code"},
	}
	req, err = http.NewRequestWithContext(ctx, "POST", authTokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("getToken stage3 error: can't create request: %w", err)
	}
	req.Header = http.Header{
		"Content-Type":  {contentTypeUrlEncoded},
		"User-Agent":    {iosUserAgent},
		"Cookie":        {fmt.Sprintf("GCDMSSO=%s", authString)},
		"Authorization": {getBasicAuthHeader()},
	}

	resp, err = c.httpClient.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		return fmt.Errorf("getToken stage3 error: can't send request: %w, %v", err, resp)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	d = json.NewDecoder(resp.Body)
	err = d.Decode(c.auth)
	if err != nil {
		return fmt.Errorf("getToken stage3 error: can't decode response: %w", err)
	}

	return nil
}

func (c *Client) getRefreshToken(ctx context.Context) error {
	form := url.Values{
		"redirect_uri":  {"com.bmw.connected://oauth"},
		"refresh_token": {c.auth.RefreshToken},
		"grant_type":    {"refresh_token"},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", authUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("getRefreshToken error: can't create request: %w", err)
	}

	req.Header = http.Header{
		"Content-Type":  {contentTypeUrlEncoded},
		"Authorization": {getBasicAuthHeader()},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("getRefreshToken error: can't send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	d := json.NewDecoder(resp.Body)
	err = d.Decode(c.auth)
	if err != nil {
		return fmt.Errorf("getRefreshToken error: can't decode response: %w", err)
	}

	return nil
}

func getBasicAuthHeader() string {
	return fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientPassword))),
	)
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	characters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~"
	randomString := ""

	for i := 0; i < length; i++ {
		randomString += string(characters[rand.Intn(len(characters)-1)])
	}

	return randomString
}
