package repository

import (
	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/utils"
	"github.com/fatkhur1960/goauction/system/core"
)

// NotifRepository init repo
type NotifRepository struct {
	NotifQs models.UserNotifQuerySet
}

// NewNotifRepository create instance
func NewNotifRepository() *NotifRepository {
	return &NotifRepository{
		NotifQs: models.NewUserNotifQuerySet(app.DB),
	}
}

// CreateNotif create user notification
func (n *NotifRepository) CreateNotif(targetUser int64, title string, content string, notifType core.NotifType, targetID int64) (models.UserNotif, error) {
	notif := models.UserNotif{
		UserID:    targetUser,
		Title:     title,
		Content:   content,
		NotifType: int(notifType),
		Target:    int(targetID),
		CreatedAT: &utils.NOW,
	}

	if err := notif.Create(app.DB); err != nil {
		return models.UserNotif{}, err
	}

	return notif, nil
}

// GetUserNotif list user notification
func (n *NotifRepository) GetUserNotif(userID int64, offset int, limit int) ([]models.UserNotif, int, error) {
	notifs := []models.UserNotif{}
	dao := n.NotifQs.UserIDEq(userID)
	err := dao.Offset(offset).Limit(limit).All(&notifs)
	if err != nil {
		return notifs, 0, nil
	}
	count, err := dao.Count()
	if err != nil {
		return notifs, 0, nil
	}
	return notifs, count, nil
}

// MarkAsRead mark as read notif
func (n *NotifRepository) MarkAsRead(ids []int64, userID int64) error {
	err := n.NotifQs.IDIn(ids...).UserIDEq(userID).GetUpdater().SetRead(true).Update()
	if err != nil {
		return err
	}
	return nil
}
