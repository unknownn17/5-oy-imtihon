package adjsut

import (
	interfaceservices "book/internal/interface/services"
	"book/internal/models"
	"book/internal/protos/booking"
	"book/internal/protos/hotel"
	notificationss "book/internal/protos/notification"
	"book/internal/protos/user"
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Adjust struct {
	User  user.UserClient
	Hotel hotel.HotelClient
	S     *interfaceservices.Database
	N     notificationss.NotificationClient
}

var uuser int

var hotelid int

func (u *Adjust) Create(ctx context.Context, req *booking.BookHotelRequest) (*booking.GeneralResponse, error) {
	email, err1 := u.CheckUser(ctx, req)
	if err1 != nil {
		return nil, err1
	}
	_, err := u.CheckHotel(ctx, req)
	switch err {
	case models.ErrHotelNotFound:
		_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: err.Error()})
		if err != nil {
			log.Println(err)
		}
		return nil, err
	case models.ErrRoomNotFound:
		_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: err.Error()})
		if err != nil {
			log.Println(err)
		}
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
		uuser = int(newReq.UserID)
		_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: res.Message})
		if err != nil {
			log.Println(err)
		}
		return &booking.GeneralResponse{Message: res.Message}, nil
	}
	res1, err := u.Hotel.Get(ctx, &hotel.GetroomRequest{HotelId: req.HotelID, Id: req.RoomId})
	if err != nil {
		log.Println(err)
	}
	// hotelid = int(res1.HotelId)
	var newreq = models.BookHotelRequest{
		UserID:       req.UserID,
		RoomID:       req.RoomId,
		RoomType:     req.RoomType,
		HotelID:      req.HotelID,
		CheckInDate:  req.CheckInDate.AsTime(),
		CheckOutDate: req.CheckOutDate.AsTime(),
	}
	res, err := u.S.Create(ctx, &newreq, float64(res1.PricePerNight))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Hotel.UpdateRoom(ctx, &hotel.UpdateRoomRequest{Available: false, HotelId: req.HotelID, Id: req.RoomId})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.Email(ctx, &notificationss.EmailSend{Email: email, Message: fmt.Sprintf("We congratulate you with successfully booking room and your book id is %v",res.Message)})
	if err != nil {
		log.Println(err)
	}
	_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: fmt.Sprintf("We congratulate you with successfully booking room and your book id is %v", res.Message)})
	if err != nil {
		log.Println(err)
	}
	uuser = int(newreq.UserID)
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
	res1, err := u.Hotel.Get(ctx, &hotel.GetroomRequest{HotelId: int32(hotelid), Id: req.RoomId})
	if err != nil {
		log.Println(err)
	}
	res, err := u.S.Update(ctx, &models.BookHotelUpdateRequest{ID: req.Id,
		RoomID:       req.RoomId,
		RoomType:     req.RoomType,
		CheckInDate:  req.CheckInDate,
		CheckOutDate: req.CheckOutDate}, float64(res1.PricePerNight))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: int32(uuser), Message: "Your room info successfully updated"})
	if err != nil {
		log.Println(err)
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) Cancel(ctx context.Context, req *booking.CancelROomRequest) (*booking.GeneralResponse, error) {
	info, err := u.S.Get(ctx, &models.GetUsersBookRequest{ID: req.Id})
	if err != nil {
		log.Println("info",err)
		return nil, err
	}
	fmt.Println("info",info.HotelID,info.RoomID)
	_, err = u.Hotel.UpdateRoom(ctx, &hotel.UpdateRoomRequest{Available: true, HotelId: info.HotelID, Id: info.RoomID})
	if err != nil {
		log.Println("error is there",err)
		return nil, err
	}
	_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: int32(uuser), Message: "You cancelled room succesfully"})
	if err != nil {
		log.Println("in notification error",err)
	}
	res, err := u.S.Cancel(ctx, &models.CancelRoomRequest{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
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
	_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: int32(uuser), Message: "You added to the waiting users list"})
	if err != nil {
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
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) DeleteW(ctx context.Context, req *booking.DeleteWaitingList) (*booking.GeneralResponse, error) {
	res, err := u.S.DeleteW(ctx, &models.DeleteWaitingList{ID: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &booking.GeneralResponse{Message: res.Message}, nil
}

func (u *Adjust) CheckUser(ctx context.Context, req *booking.BookHotelRequest) (string, error) {
	res, err := u.User.GetUser(ctx, &user.GetUserRequest{Id: req.UserID})
	if err != nil {
		log.Println(err)
		_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: "there is no such user with this id"})
		if err != nil {
			log.Println(err)
		}
		return "", errors.New("there is no such user with this id")
	}
	if res.Age < int32(18) {
		_, err = u.N.Notification(ctx, &notificationss.ProduceMessage{UserId: req.UserID, Message: "you must be older than 18 years old to get a room"})
		if err != nil {
			log.Println(err)
		}
		return "", errors.New("you must be old enough to get a room")
	}
	return res.Email, nil
}

func (u *Adjust) CheckHotel(ctx context.Context, req *booking.BookHotelRequest) (float64, error) {
	res, err := u.Hotel.GetRooms(ctx, &hotel.GetroomRequest{HotelId: req.HotelID})
	if err != nil {
		log.Println(err)
		return 0, models.ErrHotelNotFound
	}

	checkInDate := req.CheckInDate.AsTime()
	var availableRooms []*hotel.UpdateRoomRequest
	for _, v := range res.Rooms {
		if v.RoomType == req.RoomType {
			roomInfo, err := u.S.GetRoomInfo(ctx, &models.GetRoomInfo{HotelID: req.HotelID})
			if err != nil {
				log.Println(err)
				return 0, err
			}
			if roomInfo.CheckOutDate.Before(checkInDate) || roomInfo.CheckOutDate.Equal(checkInDate) {
				availableRooms = append(availableRooms, v)
			}
		}
	}

	if len(availableRooms) > 0 {
		return float64(availableRooms[0].PricePerNight), nil
	}

	// if len(res.Rooms) == 0 {
	// 	return 0, models.ErrRoomNotFound
	// }

	return 0, models.ErrRoomNotAvailable
}
