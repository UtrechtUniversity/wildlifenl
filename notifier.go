package wildlifenl

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"google.golang.org/api/option"
)

type Notifier struct {
	firebaseClient *firebase.App
}

func newNotifier(credentials []byte) (*Notifier, error) {
	app, err := firebase.NewApp(context.Background(), nil, option.WithAuthCredentialsJSON(option.ServiceAccount, credentials))
	if err != nil {
		return nil, err
	}
	return &Notifier{firebaseClient: app}, nil
}

func (n *Notifier) send(token string, data map[string]string) error {
	client, err := n.firebaseClient.Messaging(context.Background())
	if err != nil {
		return err
	}
	message := &messaging.Message{
		Token:   token,
		Data:    data,
		Android: &messaging.AndroidConfig{Priority: "high"},
		APNS:    &messaging.APNSConfig{Headers: map[string]string{"apns-priority": "10"}},
		Webpush: &messaging.WebpushConfig{Headers: map[string]string{"Urgency": "high"}},
	}
	_, err = client.Send(context.Background(), message)
	return err
}

func (n *Notifier) SendAlarms(alarmIDs []string) error {
	alarmStore := stores.NewAlarmStore(relationalDB)
	profileStore := stores.NewProfileStore(relationalDB)
	for _, alarmID := range alarmIDs {
		alarm, err := alarmStore.Get(alarmID)
		if err != nil {
			return err
		}
		profile, err := profileStore.Get(alarm.Zone.User.ID)
		if err != nil {
			return err
		}
		if profile.FirebaseCloudMessagingToken == nil {
			continue
		}
		if err := n.send(*profile.FirebaseCloudMessagingToken, map[string]string{"alarmID": alarmID}); err != nil {
			return err
		}
	}
	return nil
}
