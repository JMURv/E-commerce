package model

import cm "github.com/JMURv/e-commerce/api/pb/common"

func FavoriteToProto(f *Favorite) *cm.Favorite {
	return &cm.Favorite{
		Id:     f.ID,
		UserId: f.UserID,
		ItemId: f.ItemID,
	}
}

func FavoriteFromProto(f *cm.Favorite) *Favorite {
	return &Favorite{
		ID:     f.Id,
		UserID: f.UserId,
		ItemID: f.ItemId,
	}
}

func FavoritesToProto(favs []*Favorite) []*cm.Favorite {
	r := make([]*cm.Favorite, 0, len(favs))
	for _, f := range favs {
		r = append(r, FavoriteToProto(f))
	}
	return r
}
