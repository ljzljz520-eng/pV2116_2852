package domain

func CanTransition(from, to Status) bool {
	switch from {
	case StatusRegistered:
		return to == StatusReviewing
	case StatusReviewing:
		return to == StatusConfirmed || to == StatusRegistered
	case StatusConfirmed:
		return to == StatusPublished || to == StatusReviewing
	case StatusPublished:
		return to == StatusArchived
	case StatusArchived:
		return false
	default:
		return false
	}
}

func Transition(batch *Batch, target Status) error {
	if batch == nil || !CanTransition(batch.Status, target) {
		return ErrTransition
	}
	batch.Status = target
	batch.Version++
	return nil
}

func ConfirmAll(batch *Batch) bool {
	if batch == nil || len(batch.Records) == 0 {
		return false
	}
	for index := range batch.Records {
		if batch.Records[index].Result != "pass" {
			return false
		}
		batch.Records[index].Confirmed = true
	}
	return true
}

func CloneBatch(batch Batch) Batch {
	clone := batch
	clone.Records = append([]StickerRecord(nil), batch.Records...)
	for i := range clone.Records {
		clone.Records[i].Divisors = append([]int(nil), batch.Records[i].Divisors...)
	}
	return clone
}
