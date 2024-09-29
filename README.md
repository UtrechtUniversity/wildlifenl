# wildlifenl

Project `wildlifenl` is the backend API for the WildlifeNL system.

## Entities

### Animal
A specific instance of an animal. For example: that one horse. Usually this implies that the animal in question carries an animal-borne sensor for unique identifaction.

### Answer
A possible answer for a question (choice).

### Borne-Sensor-Deployment
A record associating `Borne-Sensor-Reading` records with an `Animal`, having a start and end timestamp to match records of type `Borne-Sensor-Reading` with the correct `Borne-Sensor-Deployment`.

### Borne-Sensor-Reading
A timestamped reading from an animal borne sensor including its location (latitude/longitude), type and values. Note that `Borne-Sensor-Reading` is weakly linked (using a corresponding sensorID) to `Borne-Sensor-Deployment` because it resides in another database.

### Conveyance
The record of a message being sent to a user. A `Conveyance` is the result of an `Encounter` or a given `Response`.

### Detection
A timestamped record of the location of a member of a specified animal species.

### Experiment
A collection of `Questionnaire`s and `Message`s with a start date and end date, that together make up a research experiment.

### Encounter
An instance of the position of a `User` and the position of an `Animal` at a certain moment in time, where these positions are closer to each other than the *`EncounterDistance`* value of the `Species` of the animal in question, and that resulted in a message `Conveyance`.

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

### Sensor-Installation (pending)
A record associating `Sensor-Reading` records with a `LivingLab`, including the most recent data. Only readings that have a timestamp between *`StartTimestamp`* and *`EndTimestamp`* are considered valid readings for this installation, and *`EndTimestamp`* can be empty. Note that *`SensorID`* is weakly linked to the *`SensorID`* of `Sensor-Reading`.

### Sensor-Reading (pending)
A timestamped reading from a sensor including its location (latitude/longitude), type and value. Note that *`SensorID`* is weakly linked to the *`SensorID`* of `Sensor-Installation`.

### Species 
The definition for a group of animals (for example: horses), having a field that defines the encounter distance in meters and encounter time in minutes.

### Tracking-Reading
A timestamped record of the location a user. Note that `Tracking-Reading` is weakly linked (using a corresponding userID) to `User` because it resides in another database.

### User
A human user of any role. The `Role` records associated with the user authorize him/her to more functionality. A user can also have no associated `Role` records, which referes to the most prohibited user types (Recreationist, Inhabitant). User also optionally includes the most recent tracked location.

### Zone
A circular area in the world with a central postion (latitude/longitude) and radius in meters[^1], that is of interest to a user, for example including agriculture fields or a gardens.

## Entity-Relation Diagram

![Entity-Relation Diagram](EntityRelationDiagram.svg)

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
Upon a new `Response` being posted, it should be checked whether this response refers to an `Answer`. If so, it should be checked whether a `Message` exists that is associated with the same `Answer`, and is associated with a non-ended `Experiment`. If the `Experiment` has an association with a `LivingLab`, the `Interaction` that is associated with the `Response` must have a spatiotemperal overlap[^3] with this `LivingLab`. If so, a new `Conveyance` must be created referring to that `Message` and associated with the `Response` and NOT with an `Encounter`. The `Conveyance` and its `Message` should be in the response body of the post request. If multiple `Message`s qualify, a random single one is chosen as to not overload the end-user.

### Add Tracking-Reading -> Create Encounter and create Conveyance 
Upon a new `Tracking-Reading` being posted, including lat/long, it should be checked whether there is an `Animal` that has spatiotemperal overlap[^3] within the margins as specified by its `Species` with this `Tracking-Reading`. Per animal for which this is true an `Encounter` must be created having the user location field set to the lat/long of the `Tracking-Reading` and the animal location field set to that of the animal in question. Then, it should be checked whether a `Message` exists that is associated with a non-ended `Experiment` and with the same `Species` as the `Animal` for which the `Encounter` is created. If the `Experiment` has an association with a `LivingLab`, the `Tracking-Reading` must have a spatiotemperal overlap[^3] with this `LivingLab`. Then, a new `Conveyance` must be created and associated with the previously mentioned `Message` and `Encounter` and NOT with a `Response`. The `Conveyance` and its `Message` should be in the response body of the post request. If multiple `Message`s qualify, a random single one is chosen as to not overload the end-user.

---

[^1]: To calculate distances in meters between two points as given in latitude and longitude the conversion 1 meter = 0.00001 degree (either latitude or longitude) is used. This produces a small difference with reality as 1 degree latitude in reality is about 110 km and 1 degree longitude in reality is about 111 km on the equator and reaches zero at the poles, but this simplified conversion greatly increases calculation speeds as the problem can then be expressed in euclidean distance.

[^2]: A record is marked *Deactivated* using a nullable DateTime field. When NULL it means the record is active, and when filled with a timestamp it means the record was deactivated at that moment. A deactivated record does, for all intents and purposes, no longer exist, has no effect in logic rules and does no longer show up in lists, with the noteworthy exception those lists that explicitly state that they contain ALL records. Also, it is still possible to retrieve a deactivated record by ID.

[^3]: A spatiotemporal overlap means that both the distance as well as the timestamps of the location measurements of two different entities are below a predefined threshold. In other words: They were close to eachother at roughly the same moment in time.
