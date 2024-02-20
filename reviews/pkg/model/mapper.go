package model

import "github.com/JMURv/e-commerce/api/pb/common"

func ReviewFromProto(r *common.Review) *Review {
	return &Review{
		ID:             r.ReviewId,
		UserID:         r.UserId,
		ItemID:         r.ItemId,
		ReviewedUserID: r.ReviewedUserId,
		Advantages:     r.Advantages,
		Disadvantages:  r.Disadvantages,
		ReviewText:     r.ReviewText,
		Rating:         r.Rating,
	}
}

func ReviewToProto(r *Review) *common.Review {
	return &common.Review{
		ReviewId:       r.ID,
		UserId:         r.UserID,
		ItemId:         r.ItemID,
		ReviewedUserId: r.ReviewedUserID,
		Advantages:     r.Advantages,
		Disadvantages:  r.Disadvantages,
		ReviewText:     r.ReviewText,
		Rating:         r.Rating,
	}
}

func ReviewsToProto(reviews []*Review) []*common.Review {
	var protoReviews []*common.Review

	for i := range reviews {
		protoReviews = append(protoReviews, ReviewToProto(reviews[i]))
	}

	return protoReviews
}
