package enums

type AccrualOrderStatus string

const (
	AccrualOrderStatusNew        AccrualOrderStatus = "NEW"
	AccrualOrderStatusProcessing AccrualOrderStatus = "PROCESSING"
	AccrualOrderStatusInvalid    AccrualOrderStatus = "INVALID"
	AccrualOrderStatusProcessed  AccrualOrderStatus = "PROCESSED"
)

func (aos AccrualOrderStatus) String() string {
	switch aos {
	case AccrualOrderStatusNew:
		return "NEW"
	case AccrualOrderStatusProcessing:
		return "PROCESSING"
	case AccrualOrderStatusInvalid:
		return "INVALID"
	case AccrualOrderStatusProcessed:
		return "PROCESSED"
	default:
		return "UNDEFINED"
	}
}
