package types

type ScanType interface {
	Record | AuditLog
}
