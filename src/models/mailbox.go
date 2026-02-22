package models

type Mailbox struct {
	Name                  string   `json:"name"`
	Address               string   `json:"address"`
	ExpiresOn             *string  `json:"expires_on"`
	IsActive              bool     `json:"is_active"`
	LastLoginAt           *string  `json:"last_login_at"`
	PasswordRecoveryEmail string   `json:"password_recovery_email"`
	StorageUsage          float64  `json:"storage_usage"`
	ActivatedAt           string   `json:"activated_at"`
	DomainName            string   `json:"domain_name"`
	Identities            []string `json:"identities"`
	LocalPart             string   `json:"local_part"`
	Expireable            bool     `json:"expireable"`
	RemoveUponExpiry      bool     `json:"remove_upon_expiry"`
	Delegations           []string `json:"delegations"`
	IsInternal            bool     `json:"is_internal"`
	MaySend               bool     `json:"may_send"`
	MayReceive            bool     `json:"may_receive"`
	MayAccessIMAP         bool     `json:"may_access_imap"`
	MayAccessPOP3         bool     `json:"may_access_pop3"`
	MayAccessManageSieve  bool     `json:"may_access_managesieve"`
	SpamAction            string   `json:"spam_action"`
	Forwardings           []string `json:"forwardings"`
	RecipientDenylist     []string `json:"recipient_denylist"`
	SenderAllowlist       []string `json:"sender_allowlist"`
	SenderDenylist        []string `json:"sender_denylist"`
	SpamAggressiveness    string   `json:"spam_aggressiveness"`
	ChangedAt             string   `json:"changed_at"`
}
