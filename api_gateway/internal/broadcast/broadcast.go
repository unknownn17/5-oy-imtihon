package broadcast17

import (
	email1 "api/internal/email"
	jwttoken "api/internal/jwt"
	"api/internal/kafka/producer"
	"api/internal/models"
	"api/internal/protos/booking"
	"api/internal/protos/hotel"
	"api/internal/protos/user"
	redismethod "api/internal/redis/method"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Adjust struct {
	U user.UserClient
	R *redismethod.Redis
	B booking.BookHotelClient
	H hotel.HotelClient
	// N   notificationss.NotificationClient
	Ctx context.Context
}

func (u *Adjust) Register(req *models.RegisterUserRequest) error {
	_, err := u.R.RegisterGet(req.Email)
	if err != nil {
		if err := u.R.Register(req); err != nil {
			log.Println(err)
			return err
		}
		code := email1.Sent(req.Email)
		if err := u.R.VerifyCodeRequest(&models.VerifyRequest{Email: req.Email, Code: code}); err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
	return errors.New("this email is already exsted")
}

func (u *Adjust) Verify(req *models.VerifyRequest) error {
	res, err := u.R.VerifyCodeResponse(req)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(res)
	if res != req.Code {
		return errors.New("password or email isn't match")

	}
	user, err := u.R.RegisterGet(req.Email)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := u.CreateUser(user); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *Adjust) CreateUser(req *models.RegisterUserRequest) error {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := producer.Producer("create", byted); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *Adjust) Login(req *models.LogInRequest) (map[string]string, error) {
	res, err := u.U.LogIn(u.Ctx, &user.LogInRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if res.Status {
		token, err := jwttoken.CreateToken(req)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res, err := u.U.LastInserted(u.Ctx, &user.LastInsertedUser{})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		id := res.Id
		return map[string]string{fmt.Sprintf("your account is created with this id %v", id): token}, err
	}
	return nil, errors.New("password or email isn't match or missing")
}

func (u *Adjust) GetUser(req *models.GetUserRequest) (*models.GetUserResponse, error) {
	user1, err := u.R.GetUser(req)
	if err != nil {
		res, err := u.U.GetUser(u.Ctx, &user.GetUserRequest{Id: req.ID})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if err := u.R.SetUser(&models.GetUserResponse{ID: res.Id, Username: res.Username, Age: res.Age, Email: res.Email, LogOut: res.Logout}); err != nil {
			log.Println(err)
		}
		return &models.GetUserResponse{ID: res.Id, Username: res.Username, Age: res.Age, Email: res.Email, LogOut: res.Logout}, nil
	}
	fmt.Printf("id %v\nusername %v\nage %v\nemail %v\nlogout %v\n", user1.ID, user1.Username, user1.Age, user1.Email, user1.LogOut)
	return &models.GetUserResponse{ID: user1.ID, Username: user1.Username, Age: user1.Age, Email: user1.Email, LogOut: user1.LogOut}, nil

}

func (u *Adjust) UpdateUser(req *models.UpdateUserRequest) error {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := producer.Producer("update", byted); err != nil {
		log.Println(err)
		return err
	}
	res, err := u.U.GetUser(u.Ctx, &user.GetUserRequest{Id: req.ID})
	if err != nil {
		log.Println(err)
	}
	if err := u.R.SetUser(&models.GetUserResponse{ID: res.Id, Username: res.Username, Age: res.Age, Email: res.Email, LogOut: res.Logout}); err != nil {
		log.Println(err)
	}
	return nil
}

func (u *Adjust) DeleteUser(req *models.GetUserRequest) error {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := producer.Producer("delete", byted); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *Adjust) Logout(req *models.GetUserRequest) error {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := producer.Producer("logout", byted); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *Adjust) CreateHotel(req *models.CreateHotelRequest) error {
	_, err := u.H.CreateHotel(u.Ctx, &hotel.CreateHotelRequest{Name: req.Name, Location: req.Location, Rating: req.Rating, Address: req.Address})
	if err != nil {
		return err
	}
	return nil
}

func (u *Adjust) GetHotel(req *models.GetHotelRequest) (*models.GetHotelResponse, error) {
	res, err := u.H.GetHotel(u.Ctx, &hotel.GetHotelRequest{Id: req.ID})
	if err != nil {
		return nil, err
	}
	var rooms []*models.UpdateRoomRequest

	for _, v := range res.Rooms {
		var all = models.UpdateRoomRequest{
			Available:     v.Available,
			RoomType:      v.RoomType,
			PricePerNight: v.PricePerNight,
			ID:            v.Id,
			HotelID:       v.HotelId,
		}
		rooms = append(rooms, &all)
	}
	return &models.GetHotelResponse{ID: req.ID, Name: res.Name, Location: res.Location, Rating: res.Rating, Address: res.Address, Rooms: rooms}, nil
}

func (u *Adjust) GetHotels(req *models.GetsRequest) ([]*models.UpdateHotelRequest, error) {
	res, err := u.H.Gets(u.Ctx, &hotel.GetsRequest{})
	if err != nil {
		return nil, err
	}
	var rooms []*models.UpdateHotelRequest

	for _, v := range res.Hotels {
		var all = models.UpdateHotelRequest{
			ID:       v.Id,
			Name:     v.Name,
			Location: v.Location,
			Rating:   v.Rating,
			Address:  v.Address,
		}
		rooms = append(rooms, &all)
	}
	return rooms, nil
}

func (u *Adjust) UpdateHotel(req *models.UpdateHotelRequest) error {
	_, err := u.H.Update(u.Ctx, &hotel.UpdateHotelRequest{Id: req.ID, Name: req.Name, Location: req.Location, Rating: req.Rating, Address: req.Address})
	if err != nil {
		return err
	}
	return nil
}

func (u *Adjust) DeleteHotel(req *models.GetHotelRequest) error {
	_, err := u.H.Delte(u.Ctx, &hotel.GetHotelRequest{Id: req.ID})
	if err != nil {
		return err
	}
	return nil
}

func (u *Adjust) CreateRoom(req *models.CreateRoomRequest) error {
	_, err := u.H.CreateRoom(u.Ctx, &hotel.CreateRoomRequest{HotelId: req.HotelID, RoomType: req.RoomType, PricePerNight: req.PricePerNight})
	if err != nil {
		return err
	}
	return nil
}

func (u *Adjust) GetRoom(req *models.GetRoomRequest) (*models.UpdateRoomRequest, error) {
	res, err := u.H.Get(u.Ctx, &hotel.GetroomRequest{HotelId: req.HotelID, Id: req.ID})
	if err != nil {
		return nil, err
	}
	return &models.UpdateRoomRequest{Available: res.Available, RoomType: res.RoomType, PricePerNight: res.PricePerNight, ID: res.Id, HotelID: res.Id}, nil
}

func (u *Adjust) GetRooms(req *models.GetRoomRequest) (*models.GetRoomResponse, error) {
	res, err := u.H.GetRooms(u.Ctx, &hotel.GetroomRequest{HotelId: req.HotelID, Id: req.ID})
	if err != nil {
		return nil, err
	}
	var rooms []*models.UpdateRoomRequest

	for _, v := range res.Rooms {
		var all = models.UpdateRoomRequest{
			ID:            v.Id,
			HotelID:       v.HotelId,
			Available:     v.Available,
			RoomType:      v.RoomType,
			PricePerNight: v.PricePerNight,
		}
		rooms = append(rooms, &all)
	}
	return &models.GetRoomResponse{Rooms: rooms}, nil
}

func (u *Adjust) UpdateRoom(req *models.UpdateRoomRequest) error {
	_, err := u.H.UpdateRoom(u.Ctx, &hotel.UpdateRoomRequest{Available: req.Available, RoomType: req.RoomType, PricePerNight: req.PricePerNight, Id: req.ID, HotelId: req.HotelID})
	if err != nil {
		return err
	}
	return nil
}

func (u *Adjust) DeleteRoom(req *models.GetRoomRequest) error {
	_, err := u.H.DeleteRoom(u.Ctx, &hotel.GetroomRequest{HotelId: req.HotelID, Id: req.ID})
	if err != nil {
		return err
	}
	return nil
}

// booking

func (u *Adjust) CreateBooking(req *models.BookHotelRequest) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.Create(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}




func (u *Adjust) GetBooking(req *models.GetUsersBookRequest) (*models.GetUsersBookResponse, error) {
	res, err := u.B.Get(u.Ctx, &booking.GetUsersBookRequest{Id: req.ID})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GetUsersBookResponse{ID: res.Id, UserID: res.UserID, HotelID: res.HotelID, RoomID: res.RoomId, RoomType: res.RoomType, CheckInDate: res.CheckInDate.AsTime(), CheckOutDate: res.CheckOutDate.AsTime(), TotalAmount: res.TotalAmount, Status: res.Status}, nil
}

func (u *Adjust) UpdateBooking(req *models.BookHotelUpdateRequest) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.Update(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) DeleteBooking(req *models.CancelRoomRequest) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.Delete(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) CreateWaitinglist(req *models.CreateWaitingList) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.CreateWaiting(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) GetWaiting(req *models.GetWaitinglistRequest) (*models.GetWaitinglistResponse, error) {
	res, err := u.B.GetWaitinglist(u.Ctx, &booking.GetWaitinglistRequest{Id: req.ID})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GetWaitinglistResponse{UserID: res.UserId, UserEmail: res.UserEmail, RoomType: res.RoomType, HotelID: res.HotelId, CheckInDate: res.CheckInDate.AsTime(), CheckOutDate: res.CheckOutDate.AsTime(), Status: res.Status, ID: res.Id}, nil
}

func (u *Adjust) UpdateWaiting(req *models.UpdateWaitingListRequest) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.UpdateWaiting(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) DeleteWaiting(req *models.DeleteWaitingList) (*models.GeneralResponse, error) {
	byted, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := u.B.UpdateWaiting(u.Ctx, &booking.Bytes{All: byted})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &models.GeneralResponse{Message: res.Message}, nil
}
