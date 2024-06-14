# wildlifenl

Project `wildlifenl` is the backend api for the WildlifeNL project.

## Concepts

`Animal` is a specific instance of an animal. For example: that one horse. Usually this implies that the animal in question carries a biosensor for uniqiue identification.

`Species` is a definition for a group of animals. For example: horses.

`User` is a human user of any role. The roles that the user has authorize him/her to more functionality. A user can also have no role, which referes to the most prohibited user types (Recreationist and Landowner).

`Interaction` is a report from the user about having had a human-wildlife interaction. This optionally includes the species of the animal.

`Tracking` is a session user tracking, started and stopped by user.

`Tracking-Reading` is a record of the position of the user at a certain moment in time, plus optional fields for other values. Often collected automatically and transparently.

`Role` is an authorization for a set of functionalities. A user having a certain role is authorized to call the functionalities associated with that role.

`Notice` is a report of something that is not normal and/or should be fixed. For example: something is damaged by animals or an animal was killed by a vehicle.

`Questionnaire` is a predefined group of questions that a user should answer upon registering an interaction.

`Question` is a question that was asked the user upon registering an interaction.

`Answer` is the answer given for a question by the user upon registering an interaction.

`Area` is a place in the world with a specified boundary, that is of interest to a User, for example an agriculture field or a garden.

`Animal-Reading` is a record of the position of the animal at a certain moment in time, plus optional fields for other values from sensors like accelerometer, heartbeat-meter, etc.

`Park` is a place in the world designated as a nature area, for example a living lab.

`Sensor` is a device at a location that produces readings, for example a visitor counter.

`Sensor-Reading` is a record at a certain moment in time, plus optional fields for other values.

`Questionnaire Deployment` is an designation of a questionnaire to an interaction type together with a timestamp, such that a newer questionnaire can be designated to an interaction type without losing historical information.

`Questionnaire Deployment` is an designation of a questionnaire to an interaction type together with a timestamp, such that a newer questionnaire can be designated to an interaction type without losing historical information.

```mermaid
---
title: Conceptual model
---
erDiagram 
    ANIMAL ||--o{ ANIMAL-READING : has
    ANIMAL }o--|| SPECIES : is
    INTERACTION }o--o| SPECIES : with
    INTERACTION }o--|| INTERACTION-TYPE : has
    INTERACTION-TYPE }o--o{ QUESTIONNAIRE-DEPLOYMENT : has
    QUESTIONNAIRE-DEPLOYMENT }o--|| QUESTIONNAIRE : has
    NOTICE }o--|| USER : reports
    RESPONSE }o--|| INTERACTION : has
    USER }o--o{ ROLE : has
    ANSWER }o--|| QUESTION : has
    RESPONSE }o--|| QUESTION : has
    QUESTION }o--o{ QUESTIONNAIRE : contains 
    USER ||--o{ TRACKING : has
    TRACKING ||--o{ TRACKING-READING : has
    USER ||--o{ INTERACTION : has
    USER ||--o{ AREA : has
    USER }o--o| PARK : works_at
    PARK |o--o{ SENSOR : contains
    SENSOR ||--o{ SENSOR-READING : has
```

