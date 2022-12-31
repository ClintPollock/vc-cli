package user

type UserData struct {
	UserID       string `json:"user_id"`
	UserLegacyID int    `json:"user_legacy_id"`
	UserName     string `json:"user_name"`
	OktaUserID   string `json:"okta_user_id"`
	Organization struct {
		OrgID   string `json:"org_id"`
		OrgName string `json:"org_name"`
		Links   struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	} `json:"organization"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	EmailAddress      string `json:"email_address"`
	IsEmailVerified   bool   `json:"is_email_verified"`
	LastLogin         string `json:"last_login"`
	LastHost          string `json:"last_host"`
	PinRequired       bool   `json:"pin_required"`
	SamlUser          bool   `json:"saml_user"`
	IPRestricted      bool   `json:"ip_restricted"`
	AgreeTerms        bool   `json:"agree_terms"`
	LoginQuestion     string `json:"login_question"`
	ShowWelcome       bool   `json:"show_welcome"`
	Active            bool   `json:"active"`
	LoginEnabled      bool   `json:"login_enabled"`
	LoginFailureCount int    `json:"login_failure_count"`
	AllowEdit         bool   `json:"allow_edit"`
	APICredentials    struct {
		APIID        string `json:"api_id"`
		ExpirationTs string `json:"expiration_ts"`
		Links        struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	} `json:"api_credentials"`
	Teams []interface{} `json:"teams"`
	Roles []struct {
		RoleID          string `json:"role_id"`
		RoleName        string `json:"role_name"`
		RoleDescription string `json:"role_description"`
		Links           struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	} `json:"roles"`
	Permissions []struct {
		PermissionID   string `json:"permission_id"`
		PermissionName string `json:"permission_name"`
		Links          struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	} `json:"permissions"`
	ProxyOrganizations []interface{} `json:"proxy_organizations"`
	Links              struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}
