package sqlbuilder

import (
	"book/internal/models"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
)

func Create(req *models.BookHotelRequest, roomPrice float64) (string, []interface{}, error) {
	// if req.CheckInDate.Before(time.Now().Truncate(24 * time.Hour)) {
	// 	return "", nil, errors.New("check-in date cannot be in the past")
	// }

	totalCost := TotalCostCalculate(req.CheckInDate, req.CheckOutDate, roomPrice)

	var status string
	today := time.Now().Truncate(24 * time.Hour)
	if req.CheckInDate.Equal(today) {
		status = "Guest Living"
	} else {
		status = "Booked"
	}
	query, args, err := squirrel.Insert("booked").
		Columns("user_id", "hotel_id", "room_id", "room_type", "enterydate", "leavingdate", "totalcost", "status").
		Values(req.UserID, req.HotelID, req.RoomID, req.RoomType, req.CheckInDate, req.CheckOutDate, totalCost, status).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	return query, args, nil
}

func Get(req *models.GetUsersBookRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("booked").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func Update(req *models.BookHotelUpdateRequest, roomPrice float64) (string, []interface{}, error) {
	setMap := make(map[string]interface{})
	if req.CheckInDate != (time.Time{}) {
		setMap["enterydate"] = req.CheckInDate
	}
	if req.CheckOutDate != (time.Time{}) {
		setMap["leavingdate"] = req.CheckOutDate
	}
	if req.RoomID != 0 {
		setMap["room_id"] = req.RoomID
	}
	if req.RoomType != "" {
		setMap["room_type"] = req.RoomType
	}
	if req.CheckInDate != (time.Time{}) && req.CheckOutDate != (time.Time{}) {
		totalCost := TotalCostCalculate(req.CheckInDate, req.CheckOutDate, roomPrice)
		setMap["totalcost"] = totalCost
	}
	setMap["status"]="updated"

	query, args, err := squirrel.Update("booked").
		SetMap(setMap).
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func Cancel(req *models.CancelRoomRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Delete("booked").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}
func TotalCostCalculate(in, out time.Time, price float64) float64 {
	days := out.Sub(in).Hours() / 24
	return days * price
}

func CreateWaitingList(req *models.CreateWaitingList) (string, []interface{}, error) {
	query, args, err := squirrel.Insert("waitinglist").
		Columns("user_id", "hotel_id", "room_type", "user_email", "enterydate", "leavingdate", "status").
		Values(req.UserID, req.HotelID, req.RoomType, req.UserEmail, req.CheckInDate, req.CheckOutDate, "waiting").
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetWaitingList(req *models.GetWaitinglistRequest) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("waitinglist").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func UpdateWaitingList(req *models.UpdateWaitingListRequest) (string, []interface{}, error) {
	setMap := make(map[string]interface{})
	if req.CheckInDate != (time.Time{}) {
		setMap["enterydate"] = req.CheckInDate
	}
	if req.CheckOutDate != (time.Time{}) {
		setMap["leavingdate"] = req.CheckOutDate
	}
	if req.HotelID != 0 {
		setMap["room_id"] = req.HotelID
	}
	if req.RoomType != "" {
		setMap["room_type"] = req.RoomType
	}
	if req.UserID != 0 {
		setMap["user_id"] = req.UserID
	}
	setMap["status"]="updated"
	query, args, err := squirrel.Update("waitinglist").
		SetMap(setMap).
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func DeleteWaitingList(req *models.DeleteWaitingList) (string, []interface{}, error) {
	query, args, err := squirrel.Delete("waitinglist").
		Where(squirrel.Eq{"id": req.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetRoomInfo(req *models.GetRoomInfo) (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("booked").
		Where(squirrel.Eq{"hotel_id": req.HotelID, "room_id": req.RoomID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}

func GetAllUsers() (string, []interface{}, error) {
	query, args, err := squirrel.Select("*").
		From("waitinglist").
		ToSql()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	return query, args, nil
}
