package models

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleSuperAdmin UserRole = "SUPERADMIN"
	UserRoleAdmin      UserRole = "ADMIN"
	UserRoleCashier    UserRole = "CASHIER"
	UserRoleMember     UserRole = "MEMBER"
	UserRolePT         UserRole = "PT"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
)

// MembershipStatus represents the status of a membership
type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "ACTIVE"
	MembershipStatusExpired   MembershipStatus = "EXPIRED"
	MembershipStatusCancelled MembershipStatus = "CANCELLED"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusExpired   OrderStatus = "EXPIRED"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

// AttendanceDirection represents the direction of attendance (in/out)
type AttendanceDirection string

const (
	AttendanceDirectionIn  AttendanceDirection = "IN"
	AttendanceDirectionOut AttendanceDirection = "OUT"
)

// AttendanceResult represents the result of an attendance scan
type AttendanceResult string

const (
	AttendanceResultAccepted AttendanceResult = "ACCEPTED"
	AttendanceResultRejected AttendanceResult = "REJECTED"
)

// OTVPassStatus represents the status of a one-time visit pass
type OTVPassStatus string

const (
	OTVPassStatusUnused  OTVPassStatus = "UNUSED"
	OTVPassStatusUsed    OTVPassStatus = "USED"
	OTVPassStatusExpired OTVPassStatus = "EXPIRED"
)
