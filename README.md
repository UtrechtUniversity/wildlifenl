# wildlifenl

Project `wildlifenl` is the backend API for the WildlifeNL system.

## Entities

### Animal
A specific instance of an animal. For example: that one horse. Usually this implies that the animal in question carries an animal-borne sensor for unique identifaction.

[ ]

### Answer
A possible answer for a question (choice).

[ *Text* ]

### Borne-Sensor-Deployment
A record associating `Borne-Sensor-Reading` records with an `Animal`, including the most recent data. Only readings that have a timestamp between *`StartTimestamp`* and *`EndTimestamp`* are considered valid readings for this deployment, and *`EndTimestamp`* can be empty. Note that *`BorneSensorID`* is weakly linked to the *`SensorID`* of `Borne-Sensor-Reading`.

[ *BorneSensorID, Type, Latitude, Longitude, StartTimestamp, EndTimestamp, Timestamp, Latitude, Longitude, AccelerometerX, AccelerometerY, AccelerometerZ, HeartRate, Temperature* ]

### Borne-Sensor-Reading
A reading for an animal borne sensor.

A timestamped reading from an animal borne sensor including its location (latitude/longitude), type and values. Note that *`SensorID`* is weakly linked to the *`BorneSensorID`* of `Borne-Sensor-Deployment`.

[ *SensorID, Timestamp, Latitude, Longitude, AccelerometerX, AccelerometerY, AccelerometerZ, HeartRate, Temperature* ]

### Conveyance
The record of a message being sent to a user. A `Conveyance` is the result of an `Encounter` or a given `Response`.

### Encounter
An instance of the position of a `User` and the position of an `Animal` at a certain moment in time, where these positions are closer to each other than the *`EncounterDistance`* value of the `Species` of the animal in question, and that resulted in a message `Conveyance`.

[ *Timestamp, UserLatitude, UserLongitude, AnimalLatitude, AnimalLongitude* ] 

### Interaction
The report from a user about having had a human-wildlife interaction. This optionally includes the species of the animal.

[ *Timestamp, Latitude, Longitude, Description* ]

### Interaction-Type
The type of the interaction, for example a sighting, damage report or an animal-vehicle collision.

[ *Name* ]

### Question 
A question that was asked the user upon registering an interaction.

[ *Text* ]

### Questionnaire 
A predefined group of questions that a user should answer upon registering an interaction of a certain `Interaction-Type`. As soon as any of its questions has an associated `Response` record, the questionnaire can no longer be changed. A `Questionnaire` can be deactivated when it is no longer needed to deploy it, and a new questionnaire can be created for the new situation.

[ *Deactivated, Name* ]

### LivingLab
A place in the world designated as a nature area.

[ *Name* ]

### Message
A piece of information to send to the user upon certain conditions. As soon as a message has an associated `Conveyance` record, it can no longer be changed. A `Message` can be deactivated when it is no longer needed to send it, and a new `Message` can be created for the new situation.

[ *Deactivated, Text*]

### Message-Type
The type of a message to indicate severity, for example a piece of information or an urgent warning.

[ *Name* ]

### Response
The answer the user gave to the question. If the `Question` has associated possible `Answer` records, the response must refer to one of the `Answer` records of the question.

[ *Text* ]

### Role 
An authorization record for a set of functionalities. A user having a certain role is authorized to call the functionalities assigned to that role.

[ *Name* ]

### Sensor-Installation
A record associating `Sensor-Reading` records with a `LivingLab`, including the most recent data. Only readings that have a timestamp between *`StartTimestamp`* and *`EndTimestamp`* are considered valid readings for this installation, and *`EndTimestamp`* can be empty. Note that *`SensorID`* is weakly linked to the *`SensorID`* of `Sensor-Reading`.

[ *ID, Type, Latitude, Longitude, StartTimestamp, EndTimestamp, TimeStamp, Value* ]

### Sensor-Reading
A timestamped reading from a sensor including its location (latitude/longitude), type and value. Note that *`SensorID`* is weakly linked to the *`SensorID`* of `Sensor-Installation`.

[ *SensorID, Timestamp, Value* ]

### Species 
The definition for a group of animals. For example: horses. Note that the EncounterDistance is in meters.

[ *Name, EncounterDistance* ]

### Tracking-Reading
A record of the location a user. Note that *`UserID`* is weakly linked to the *`ID`* of `User`.

[ *UserID, Timestamp, Latitude, Longitude* ]

### User
A human user of any role. The `Role` records associated with the user authorize him/her to more functionality. A user can also have no associated `Role` records, which referes to the most prohibited user types (Recreationist, Inhabitant). User also optionally included the most recent tracked location.

[ *Name, Email, Timestamp, Latitude, Longitude* ]

### Zone
A circular area in the world with a central postion (latitude/longitude) and radius in meters, that is of interest to a user, for example including agriculture fields or a gardens.

[ *Latitude, Longitude, Radius* ]

## Entity-Relation Diagram

![Entity-Relation Diagram](EntityRelationDiagram.svg)

The blue entities are stored in a relational database, whereas the green entities are stored in a timeseries database. As a result the relationships between blue and green entities are "by convention" of having the same IDs rather than them being enforced by any rules. For `Sensor-Reading` and `Borne-Sensor-Reading` this brings an advantage in the cases where data is ingested from external automated systems as it is therefore possible to accept and store these data, regardless of there being a way to use and retreive them in a meaningful context (even if there are no `Sensor-Installation` or `Borne-Sensor-Deployment` with the referenced ID), subsequently allowing that the installation or deployment that adds meaning be added later, rather than rejecting these data for the reason of there being no meaningful context, and perhaps make external automated systems fail. Note that this advantage is not needed for `Tracking-Reading` as it is highly unlikely for a new reading to refer to a userID that is unknown.

Blue entities with a red border are created internally as a result of logic rules (see below) and are therefore *de facto* read-only.

## Logic rules

Upon a new `Response` being posted, it should be checked whether this response refers to an `Answer`. If so, it should be checked whether a non-deactivated `Message` exists that is associated with the same `Answer`. If so, a new `Conveyance` must be created referring to that `Message` and associated with the `Response` and NOT with an `Encounter`. The `Conveyance` and its `Message` should be in the response body of the post request.

Upon a new `Tracking-Reading` being posted, including lat/long, it should be checked whether a non-deactivated `Message` exists that is associated to a `Species` for which there is an `Animal` that has a lat/long that is closer by the lat/long of the `Tracking-Reading` than the *`EncounterDistance`* value of that `Species`. Per animal for which this is true an `Encounter` must be created having the user location fields set to the lat/long of the `Tracking-Reading` and the animal location fields set to that of the animal in question. Then, a new `Conveyance` must be created and associated with the previously mentioned `Message` and `Encounter` and NOT with a `Response`. The `Conveyance` and its `Message` should be in the response body of the post request.

Upon a new `Borne-Sensor-Reading` being posted, including lat/long, it should be checked whether a `Zone` exits that this lat/lang is within and whether a `Borne-Sensor-Deployment` exists for this reading. If so, a new `Alarm` must be created being associated with the `Zone` in question and the `Animal` that the `Borne-Sensor-Deployment` refers to. **TODO: Find a way to push this alarm to the user device.**

## Notes

To calculate distances in meters between two points as given in latitude and longitude the conversion 1 meter = 0.00001 degree (either latitude or longitude) is used. This produces a small difference with reality as 1 degree latitude in reality is about 110 km and 1 degree longitude in reality is about 111 km on the equator and reaches zero at the poles, but this simplified conversion greatly increases calculation speeds as the problem can then be expressed in euclidean distance.

The often present field *Deactivated* is a nullable DateTime field. When NULL it means the record is active, and when filled with a timestamp it means the record is deactivated.
