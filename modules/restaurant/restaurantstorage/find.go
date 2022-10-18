package restaurantstorage

import (
	"context"
	"simple-rest-api/modules/restaurant/restaurantmodel"
)

func (s *sqlStore) FindDataByCondition(
	ctx context.Context, 
	conditions map[string]interface{}, 
	moreKeys ...string,
	) (*restaurantmodel.Restaurant, error) {
		var result restaurantmodel.Restaurant
	
		db := s.db

		for i := range moreKeys {
			db = db.Preload(moreKeys[i])
		}

		if err := db.Where(conditions).
			First(&result).Error; err != nil {
			return nil, err
		}

		return &result, nil
	}