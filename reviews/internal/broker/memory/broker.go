package memory

import "strconv"

type Broker struct {
	msgs map[string][]byte
}

func New() *Broker {
	return &Broker{
		msgs: make(map[string][]byte),
	}
}

func (b *Broker) NewReviewNotification(reviewID uint64, msg []byte) error {
	b.msgs[strconv.FormatUint(reviewID, 10)] = msg
	return nil
}
