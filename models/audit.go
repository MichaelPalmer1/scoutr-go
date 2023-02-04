package models

// AuditUser : User object used in audit logs
type AuditUser struct {
	ID        string `json:"id" firestore:"id"`
	Username  string `json:"username" firestore:"username"`
	Name      string `json:"name" firestore:"name"`
	Email     string `json:"email" firestore:"email"`
	SourceIP  string `json:"source_ip,omitempty" firestore:"source_ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty" firestore:"user_agent,omitempty"`
}

// AuditLog : Audit log object
type AuditLog struct {
	Time        string              `json:"time" firestore:"time"`
	User        AuditUser           `json:"user" firestore:"user"`
	Action      string              `json:"action" firestore:"action"`
	Method      string              `json:"method" firestore:"method"`
	Path        string              `json:"path" firestore:"path"`
	ExpireTime  int64               `json:"expire_time,omitempty" firestore:"expire_time,omitempty"`
	QueryParams map[string][]string `json:"query_params,omitempty" firestore:"query_params,omitempty"`
	Resource    map[string]string   `json:"resource,omitempty" firestore:"resource,omitempty"`
	Body        interface{}         `json:"body,omitempty" firestore:"body,omitempty"`
}
