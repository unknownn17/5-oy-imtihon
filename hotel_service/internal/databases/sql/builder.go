package sqlbuilder

import (
	"hotel/internal/models"
	"log"

	"github.com/Masterminds/squirrel"
)

func CreateHotel(req *models.CreateHotelRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Insert("hotels").
		Columns("name", "location", "rating", "address").
		Values(req.Name, req.Location, req.Rating, req.Address).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetHotel(req *models.GetHotelRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("hotels").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetsHotel(req *models.GetsRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("hotels").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func UpdateHotel(req *models.UpdateHotelRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Update("hotels").
		SetMap(map[string]interface{}{
			"name":     req.Name,
			"location": req.Location,
			"rating":   req.Rating,
			"address":  req.Address,
		}).
		Where(squirrel.Eq{"id": req.ID}).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func DeleteHotel(req *models.GetHotelRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Delete("hotels").
		Where(squirrel.Eq{"id": req.ID}).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func CreateRoom(req *models.CreateRoomRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Insert("rooms").
		Columns("hotel_id", "room_type", "price_per_night").
		Values(req.HotelID, req.RoomType, req.PricePerNight).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetRoom(req *models.GetRoomRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("rooms").
		Where(squirrel.Eq{"id": req.ID, "hotel_id": req.HotelID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetsRoom(req *models.GetsRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("rooms").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}
func UpdateRoom(req *models.UpdateRoomRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Update("rooms").
		SetMap(map[string]interface{}{
			"available":     req.Available,
			"room_type": req.RoomType,
			"price_per_night":   req.PricePerNight,
		}).
		Where(squirrel.Eq{"id": req.ID,"hotel_id":req.HotelID}).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func DeleteRoom(req *models.GetRoomRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Delete("rooms").
		Where(squirrel.Eq{"id": req.ID,"hotel_id":req.HotelID}).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}
