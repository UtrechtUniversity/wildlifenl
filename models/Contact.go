package models

import "time"

type ContactRecord struct {
	ID string `json:"ID" format:"uuid" readOnly:"true" doc:"The ID of this contact."`
	//UserID                  string     `json:"userID" format:"uuid" writeOnly:"true" doc:"The ID of the user."`
	//BorneSensorDeploymentID string     `json:"animalID" format:"uuid" writeOnly:"true" doc:"The ID of the borne-sensor deployment."`
	Start         time.Time  `json:"start" format:"date-time" doc:"The moment this contact started."`
	End           *time.Time `json:"end" format:"date-time" required:"false" readOnly:"true" doc:"The moment this contact finished."`
	StartLocation Point      `json:"startLocation" readOnly:"true" doc:"The location of the user at the moment this contact started. The north pole (latitude 90, longitude 0) for contacts recorded before locations were stored."`
	EndLocation   Point      `json:"endLocation" readOnly:"true" doc:"The location of the user at the moment this contact finished. The north pole (latitude 90, longitude 0) while the contact has not finished, and for contacts recorded before locations were stored."`
}

type Contact struct {
	ContactRecord
	User                  User                  `json:"user" doc:"The user involved in the contact tracing event."`
	BorneSensorDeployment BorneSensorDeployment `json:"borneSensorDeployment" doc:"The borne-sensor deployment and animal involved in the contact tracing event."`
	Conveyances           []Conveyance          `json:"conveyances" doc:"The conveyances that were created for this contact tracing event."`
}
