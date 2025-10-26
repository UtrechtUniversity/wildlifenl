# TECHNICAL DESIGN

*WildlifeNL: Applications and Data Platform*

## Introduction
This technical design provides an overview of the implementation details for the WildlifeNL API backend and data storage solutions. It defines the entities and logic involved.

## Roles and access scopes
Access to certain API endpoints is limited to those accounts that have the required credential scopes. Based on the role of the user, scopes are assigned by the Adminitrators. It is possible for the same account to have multiple scopes assigned.

|Role|Scope|
|----|-----|
|👲 Recreationist|[none]|
|🧑🏻‍💼 Inhabitant|[none]|
|🧑‍🌾 Land User|`land-user`|
|🧝🏻 Nature Area Manager|`nature-area-manager`|
|🧑🏻‍🔧 Wildlife Manager|`wildlife-manager`|
|🧑‍🔬 Researcher|`researcher`|
|🧑🏻‍💻 Administrator|`administrator`|

Additionally, the scope `data-system` exists for accounts that are meant to be used by automated systems that provide data ingest and the scope `herd-manager` exists but is currently not used.

## Entities

The `wildlifenl` API discloses access to the following entities using end-points that have the same name.

### Animal
A specific instance of an animal. For example: that one horse. Usually this implies that the animal in question carries an animal-borne sensor for unique identifaction.

### Animal-Info
Additional information on an animal that is involved in an `SightingReport` or `CollisionReport` as reported by the end-user.

### Answer
A possible answer for a question (choice).

### Belonging
Additional information for a `DamageReport` specifying the type of belonging being damaged. 

### Borne-Sensor-Deployment
A record associating `Borne-Sensor-Reading` records with an `Animal`, having a start and end timestamp to match records of type `Borne-Sensor-Reading` with the correct `Borne-Sensor-Deployment`.

### Borne-Sensor-Reading
A timestamped reading from an animal borne sensor including its location (latitude/longitude), type and values. Note that `Borne-Sensor-Reading` is weakly linked (using a corresponding sensorID) to `Borne-Sensor-Deployment` because it resides in another database.

### Collision-Report
Additional information about an interaction with type `animal-vehicle-collision`. The corresponding `Interaction` can have just one additional report.

### Conveyance
The record of a message being sent to a user. A `Conveyance` is the result of an `Encounter` or a given `Response`.

### Damage-Report
Additional information about an interaction with type `property-damage`. The corresponding `Interaction` can have just one additional report.

### Detection
A timestamped record of the location of a member of a specified animal species.

### Experiment
A collection of `Questionnaire`s and `Message`s with a start date and end date, that together make up a research experiment.

### Interaction
The report from a user about having had a human-wildlife interaction.

### Interaction-Type
The type of the interaction, for example a sighting, damage report or an animal-vehicle collision.

### Question 
A question that was asked the user upon registering an interaction.

### Questionnaire 
A predefined group of questions that a user should answer upon registering an interaction of a certain `Interaction-Type`. As soon as any of its questions has an associated `Response` record, the questionnaire can no longer be changed. A `Questionnaire` can be deactivated when it is no longer needed to deploy it, and a new questionnaire can be created for the new situation.

### LivingLab
A place in the world designated as a nature area.

### Message
A piece of information to send to the user upon certain conditions. As soon as a message has an associated `Conveyance` record, it can no longer be changed. A `Message` can be deactivated when it is no longer needed to send it, and a new `Message` can be created for the new situation.

### Response
The answer the user gave to the question. If the `Question` has associated possible `Answer` records, the response must refer to one of the `Answer` records of the question.

### Role 
An authorization record for a set of functionalities. A user having a certain role is authorized to call the functionalities assigned to that role.

### Sighting-Report
Additional information about an interaction with type `sighting`. The corresponding `Interaction` can have just one additional report.

### Species 
The definition for a group of animals (for example: horses), having a field that defines the encounter distance in meters and encounter time in minutes.

### Tracking-Reading
A timestamped record of the location a user. Note that `Tracking-Reading` is weakly linked (using a corresponding userID) to `User` because it resides in another database.

### User
A human user of any role. The `Role` records associated with the user authorize him/her to more functionality. A user can also have no associated `Role` records, which referes to the most prohibited user types (Recreationist, Inhabitant). User also optionally includes the most recent tracked location.

### Visitor-Sensor-Deployment
A sensor that counts human presence in a certain area specified by a location (latitude/longitude).

### Visitor-Sensor-Reading
A timestamped reading from a visitor sensor. Note that *`SensorID`* is weakly linked to the *`SensorID`* of `Visitor-Sensor-Deployment`.

### Zone
A circular area in the world with a central postion (latitude/longitude) and radius in meters[^1], that is of interest to a user, for example including agriculture fields or a gardens.

## Entity-Relation Diagram

![entity-relation-diagram](assets/entity-relation-diagram.svg)

The blue entities are stored in a relational database, whereas the green entities are stored in a timeseries database. As a result the relationships between blue and green entities are "by convention" of having the same IDs rather than them being enforced by any rules. For `Sensor-Reading` and `Borne-Sensor-Reading` this brings an advantage in the cases where data is ingested from external automated systems as it is therefore possible to accept and store these data, regardless of there being a way to use and retreive them in a meaningful context (even if there are no `Sensor-Installation` or `Borne-Sensor-Deployment` with the referenced ID), subsequently allowing that the installation or deployment that adds meaning be added later, rather than rejecting these data for the reason of there being no meaningful context, and perhaps make external automated systems fail. Note that this advantage is not needed for `Tracking-Reading` as it is highly unlikely for a new reading to refer to a userID that is unknown.

Blue entities with a red border are created internally as a result of logic (see below) and are therefore *de facto* read-only.

## Logic

Some end-points have interal logic that does more than just serve the end-point result.

### Add Interaction -> Get Questionnaire
Upon a new `Interaction` being posted, it should be checked whether a `Questionnaire` exists that is accociated with the `InteractionType` of the newly posted `Interaction`, that is associated with a non-ended `Experiment`, and if that `Experiment` is bound to a specific `LivingLab` the newly posted `Interaction` should be within that `LivingLab`. If so, this `Questionnaire` should be in the response body of the post request. If multiple `Questionnaire`s qualify, a random single one is chosen as to not overload the end-user. 

### Add Interaction -> Create Alarms
Upon a new `Interaction` being posted, including lat/long, it should be checked whether its type is 'Sighting'. If so, it should be checked whether a `Zone` exits that is associated with the `Species` of this interaction, and that has a spatiotemperal overlap[^3] with this `Interaction`. If so, a new `Alarm` must be created being associated with the `Zone` in question and linked to this `Interaction`.

### Add Detection -> Create Alarms
Upon a new `Detection` being posted, including lat/long, it should be checked whether a `Zone` exits that is associated with the `Species` of this detection and that has a spatiotemperal overlap[^3] with this `Detection`. If so, a new `Alarm` must be created being associated with the `Zone` in question and linked to this `Detection`.

### Add Borne-Sensor-Reading -> Create Alarms
Upon a new `Borne-Sensor-Reading` being posted, including lat/long, it should be checked whether a `Borne-Sensor-Deployment` exists for this reading, if so the location of the associated `Animal` should be updated. Then, it should be checked whether a `Zone` exits that is associated with the `Species` of this animal, and that has a spatiotemperal overlap[^3] with this `Animal`. If so, a new `Alarm` must be created being associated with the `Zone` in question and linked to the `Animal` that the `Borne-Sensor-Deployment` refers to.

### Add Response -> Create Conveyance
Upon a new `Response` being posted, it should be checked whether this response refers to an `Answer`. If so, it should be checked whether a `Message` exists, in non-ended `Experiment`, that is associated with the same `Answer`. If the `Experiment` has an association with a `LivingLab`, the `Interaction` that is associated with the `Response` must have a spatiotemperal overlap[^3] with this `LivingLab`. If so, a new `Conveyance` must be created referring to that `Message` and associated with the `Response` and NOT with an `Animal` and NOT with an `Alarms`. The `Conveyance` and its `Message` should be in the response body of the post request. If multiple `Message`s qualify, a random single one is chosen as to not overload the end-user.

### Add Tracking-Reading -> Create Conveyance 
Upon a new `Tracking-Reading` being posted, including lat/long, it should be checked whether there is an `Animal`, having a `Species` for which there a `Message` exists, in a non-ended `Experiment`, that has spatiotemperal overlap[^3], within the margins as specified by this `Message`, with this `Tracking-Reading`. If the `Experiment` has an association with a `LivingLab`, the `Tracking-Reading` must have a spatiotemperal overlap[^3] with this `LivingLab`. Then, a new `Conveyance` must be created referring to that `Message` and associated with `Animal` and NOT with a `Response` and NOT with `Alarm`. The `Conveyance` and its `Message` should be in the response body of the post request. If multiple `Message`s qualify, a random single one is chosen as to not overload the end-user.


### From created Alarms -> Create Conveyances
Whenever a new `Alarm` is created it should be checked whether a `Message` exists with trigger type `alarm` that is for the same `Species` as the animal that triggered the alarm. If the `Message` is associated with an `Experiment` that is bound by a `LivingLab` there should be a spetiotemperal overlap[^3] between the `Zone` that the alarm is for and this `LivingLab`. If so, a new  `Conveyance` must be created that referring to that `Message` and is associated with the `Alarm` and NOT with an `Animal`, and NOT with a `Response`. Any `Conveyance` instances that are created this way and their `Message` should be in the response body of the post request.

## Location Tracking
The smartphone apps provide functionality on tracking the location of the user. The current way to implement this is as follows.

### Wild Gids
The user accepts 'terms and conditions' prior to using the app and actively agrees to be a participant in a study. The app then directly requests permission for `ACCESS_FINE_LOCATION` and `ACCESS_BACKGROUND_LOCATION`. If this permission is not given or is later revoked, the app cannot be used. The app then asks for permission again. As long as the app is running (foreground or background), the location is transmitted. However, to guard the user's privacy, it is possible that the app closes itself if it has not been in the foreground for longer than a certain period of time (for example an hour). And it is possible to show a notification first: "This app is still actively tracking your location, do you want to continue press 'yes', otherwise it closes in 10 seconds" or something similar.

### Wild Rapport
The user accepts 'terms and conditions' prior to using the app and actively agrees to be a participant in a study. The app can now be used. There is an option for the user to enable 'Location Tracking'. For example, to walk around a field, or for fauna/nature managers to prove they were in the field. The first time this is enabled, the app asks permission for `ACCESS_FINE_LOCATION` and `ACCESS_BACKGROUND_LOCATION`. If this permission is not given, 'Location Tracking' cannot be enabled and remains disabled. If this permission is later revoked, the app will ask for permission again when the user enables 'Location Teacking'. At any time the user can disable 'Location Tracking'.

---

[^1]: To calculate distances in meters between two points as given in latitude and longitude the conversion 1 meter = 0.00001 degree (either latitude or longitude) is used. This produces a small difference with reality as 1 degree latitude in reality is about 110 km and 1 degree longitude in reality is about 111 km on the equator and reaches zero at the poles, but this simplified conversion greatly increases calculation speeds as the problem can then be expressed in euclidean distance.

[^2]: A record is marked *Deactivated* using a nullable DateTime field. When NULL it means the record is active, and when filled with a timestamp it means the record was deactivated at that moment. A deactivated record does, for all intents and purposes, no longer exist, has no effect in logic rules and does no longer show up in lists, with the noteworthy exception those lists that explicitly state that they contain ALL records. Also, it is still possible to retrieve a deactivated record by ID.

[^3]: A spatiotemporal overlap means that both the distance as well as the timestamps of the location measurements of two different entities are below a predefined threshold. In other words: They were close to eachother at roughly the same moment in time.
