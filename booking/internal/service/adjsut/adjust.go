package adjsut

import (
	"book/internal/brokers/notification"
	interfaceservices "book/internal/interface/services"
	"book/internal/models"
	"book/internal/protos/booking"
	"book/internal/protos/hotel"
	"book/internal/protos/user"
	"context"
	"errors"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Adjust struct {
	User  user.UserClient
	Hotel hotel.HotelClient
	S     *interfaceservices.Database
}

func (u *Adjust) Create(ctx context.Context, req *booking.BookHotelRequest) (*booking.GeneralResponse, error) {
	email, err1 := u.CheckUser(ctx, req)
	if err1 != nil {
		return nil, err1
	}
	price, err := u.CheckHotel(ctx, req)
	switch err {
	case models.ErrHotelNotFound:
		return nil, err
	case models.ErrRoomNotFound:
		return nil, err
	case models.ErrRoomNotAvailable:
		var newReq = models.CreateWaitingList{
			UserID:       req.UserID,
			UserEmail:    email,
			RoomType:     req.RoomType,
			HotelID:      req.HotelID,
			CheckInDate:  req.CheckInDate.AsTime(),
			CheckOutDate: req.CheckOutDate.AsTime(),
		}
		res, err := u.S.CreateW(ctx, &newReq)
		if err != nil {
			log.Println("waiting list error", err)
			return nil, err
		}
		if err := notification.Producer([]byte(res.Message)); err != nil {
			log.Println(err)
		}
		return &booking.GeneralResponse{Message: res.Message}, nil
	}
	var newreq = models.BookHotelRequest{
		UserID:       req.UserID,
		RoomID:       req.RoomId,
		RoomType:     req.RoomType,
		HotelID:      req.HotelID,
		CheckInDate:  req.CheckInDate.AsTime(),
		CheckOutDate: req.CheckOutDate.AsTime(),
	}
	res, err := u.S.Create(ctx, &newreq, price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Hotel.UpdateRoom(ctx, &hotel.UpdateRoomRequest{Available: false, HotelId: req.HotelID, Id: req.RoomId})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) Get(ctx context.Context, req *booking.GetUsersBookRequest) (*booking.GetUsersBookResponse, error) {
	res, err := u.S.Get(ctx, &models.GetUsersBookRequest{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &booking.GetUsersBookResponse{Id: res.ID,
		UserID:       res.UserID,
		HotelID:      res.HotelID,
		RoomId:       res.RoomID,
		RoomType:     res.RoomType,
		CheckInDate:  timestamppb.New(res.CheckInDate),
		CheckOutDate: timestamppb.New(res.CheckOutDate),
		TotalAmount:  res.TotalAmount,
		Status:       res.Status}, nil
}

func (u *Adjust) Update(ctx context.Context, req *booking.BookHotelUpdateRequest) (*booking.GeneralResponse, error) {
	res, err := u.S.Update(ctx, &models.BookHotelUpdateRequest{ID: req.Id,
		RoomID:       req.RoomId,
		RoomType:     req.RoomType,
		CheckInDate:  req.CheckInDate.AsTime(),
		CheckOutDate: req.CheckOutDate.AsTime()})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) Cancel(ctx context.Context, req *booking.CancelROomRequest) (*booking.GeneralResponse, error) {
	res, err := u.S.Cancel(ctx, &models.CancelRoomRequest{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}
func (u *Adjust) CreateW(ctx context.Context, req *booking.CreateWaitingList) (*booking.GeneralResponse, error) {
	var newReq = models.CreateWaitingList{
		UserID:       req.UserId,
		UserEmail:    req.UserEmail,
		RoomType:     req.RoomType,
		HotelID:      req.HotelId,
		CheckInDate:  req.CheckInDate.AsTime(),
		CheckOutDate: req.CheckOutDate.AsTime(),
	}
	res, err := u.S.CreateW(ctx, &newReq)
	if err != nil {
		log.Println("waiting list error", err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}
func (u *Adjust) GetW(ctx context.Context, req *booking.GetWaitinglistRequest) (*booking.GetWaitinglistResponse, error) {
	res, err := u.S.GetW(ctx, &models.GetWaitinglistRequest{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &booking.GetWaitinglistResponse{
		UserId:       res.UserID,
		UserEmail:    res.UserEmail,
		RoomType:     res.RoomType,
		HotelId:      res.HotelID,
		CheckInDate:  timestamppb.New(res.CheckInDate),
		CheckOutDate: timestamppb.New(res.CheckOutDate),
		Status:       res.Status,
		Id:           res.ID,
	}, nil
}

func (u *Adjust) UpdateW(ctx context.Context, req *booking.UpdateWaitingListRequest) (*booking.GeneralResponse, error) {
	res, err := u.S.UpdateW(ctx, &models.UpdateWaitingListRequest{
		UserID:       req.UserId,
		RoomType:     req.RoomType,
		HotelID:      req.HotelId,
		CheckInDate:  req.CheckInDate.AsTime(),
		CheckOutDate: req.CheckOutDate.AsTime(),
		ID:           req.Id,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) DeleteW(ctx context.Context, req *booking.DeleteWaitingList) (*booking.GeneralResponse, error) {
	res, err := u.S.DeleteW(ctx, &models.DeleteWaitingList{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := notification.Producer([]byte(res.Message)); err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) CheckUser(ctx context.Context, req *booking.BookHotelRequest) (string, error) {
	res, err := u.User.GetUser(ctx, &user.GetUserRequest{Id: req.UserID})
	if err != nil {
		log.Println(err)
		return "", errors.New("there is no such user with this id")
	}
	if res.Age < int32(18) {
		return "", errors.New("you must be old enough to get a room")
	}
	return res.Email, nil
}

func (u *Adjust) CheckHotel(ctx context.Context, req *booking.BookHotelRequest) (float64, error) {
	res, err := u.Hotel.GetHotel(ctx, &hotel.GetHotelRequest{Id: req.HotelID})
	if err != nil {
		log.Println(err)
		return 0, models.ErrHotelNotFound
	}
	for _, v := range res.Rooms {
		if v.RoomType == req.RoomType {
			if v.Available {
				return float64(v.PricePerNight), nil
			} else {
				roomInfo, err := u.S.GetRoomInfo(ctx, &models.GetRoomInfo{HotelID: req.HotelID, RoomID: req.RoomId})
				if err != nil {
					log.Println(err)
					return 0, err
				}
				checkInDate := req.CheckInDate.AsTime()
				if roomInfo.CheckOutDate.Before(checkInDate) || roomInfo.CheckOutDate.Equal(checkInDate) {
					return float64(v.PricePerNight), nil
				} else {
					return 0, models.ErrRoomNotAvailable
				}
			}
		}
	}

	return 0, models.ErrRoomNotFound
}
