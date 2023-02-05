package types

import "time"

// AuditUser : User object used in audit logs
type AuditUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SourceIP  string `json:"source_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

// AuditLog : Audit log object
type AuditLog struct {
	Time        string                 `json:"time"`
	User        AuditUser              `json:"user"`
	Action      string                 `json:"action"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	ExpireTime  int64                  `json:"expire_time,omitempty"`
	QueryParams map[string][]string    `json:"query_params,omitempty"`
	Resource    map[string]interface{} `json:"resource,omitempty"`
	Body        interface{}            `json:"body,omitempty"`
}

type AuditEvent struct {
	EventData AuditEventData `json:"eventData"`
}

type AuditEventData struct {
	Version             string                 `json:"version"`
	UserIdentity        AuditEventUserIdentity `json:"userIdentity"`
	UserAgent           string                 `json:"userAgent"`
	EventSource         string                 `json:"eventSource"`
	EventName           string                 `json:"eventName"`
	EventTime           time.Time              `json:"eventTime"`
	UID                 string                 `json:"UID"`
	RequestParameters   map[string]interface{} `json:"requestParameters"`
	ResponseElements    map[string]interface{} `json:"responseElements"`
	ErrorCode           string                 `json:"errorCode"`
	ErrorMessage        string                 `json:"errorMessage"`
	SourceIPAddress     string                 `json:"sourceIPAddress"`
	RecipientAccountId  string                 `json:"recipientAccountId"`
	AdditionalEventData map[string]interface{} `json:"additionalEventData"`
}

type AuditEventUserIdentity struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	Details     *User  `json:"details"`
}
