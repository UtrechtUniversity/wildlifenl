# WildlifeNL App tests – Landholder
Today you are landholder, for example an agrarian, using the WildlifeNL WildRapport app. We are running this app against a TEST environment with simulated animals, so there will only be human-wildlife interactions that are artificially created by us. All locations are within the test living lab area Botanische Tuinen in Utrecht. If you notice an issue, bug or inconvenience when performing the test scenarios below, please write them down so we can later provide them to the development team. Note that aside from the 'Logging in' feature, all features are marked with codes in brackets ([ ]) that match the user story codes in the Functional Design. You can look-up the expected behaviour of the app and the purpose of the functionality this way. Additionally, it is possible that the screens that are mentioned in single quotes in these test scenarios have a somewhat different Dutch name in the app. Please login with your email address either from Utrecht University or Fontys University of Applied Sciences.
<br/>

## Feature: Logging in
<br/>

| L1 | Unsuccessfully logging in due to leaving field empty |
|:-:|:-|
|`GIVEN`| you are on the 'login screen' |
|`WHEN`| you leave the email box empty |
|`AND`| attempt to log in |
|`THEN`| a message is shown explaining that providing an email address is mandatory |
|`AND`| the app does __NOT__ navigate beyond the 'login screen'. |
---
<br/>

| L2 | Unsuccessfully logging in due to wrong validation code |
|:-:|:-|
|`GIVEN`| you are on the 'login screen' |
|`AND`| you have provided a valid email address |
|`WHEN`| you provide a validation code that is different from the one received by email |
|`THEN`| a message is shown explaining that the code is wrong and you can try again |
|`AND`| the app does __NOT__ navigate beyond the 'login screen'. |
---
<br/>

| L3 | Successfully logging in |
|:-:|:-|
|`GIVEN`| you are on the 'login screen' |
|`AND`| you have provided a valid email address |
|`WHEN`| you provide the validation code as received by email |
|`THEN`| the app navigates beyond the 'login screen'. |
---

## Feature: Read and accept Terms & Conditions [R1]

| R1.1 | Unsuccessfully accept terms & conditions |
|:-:|:-|
|`GIVEN`| you are on the 'terms & conditions screen' |
|`WHEN`| you do __NOT__ tick the checkbox to accept the terms & conditions |
|`AND`| you try to continue |
|`THEN`| the app must __NOT__ navigate beyond the 'terms & conditions screen'. |
---
<br/>

| R1.2 | Successfully accept terms & conditions |
|:-:|:-|
|`GIVEN`| you are on the 'terms & conditions screen' |
|`WHEN`| you tick the checkbox to accept the terms & conditions |
|`AND`| you try to continue |
|`THEN`| the app navigates beyond the 'terms & conditions screen'. |
---
<br/>

## Feature: Location Sharing [R2]

| R2.1 | Successfully activate location sharing |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`AND`| you did __NOT__ activate location sharing previously |
|`WHEN`| you navigate to the 'settings screen' |
|`AND`| activate location sharing |
|`THEN`| the app enables location sharing |
|`AND`| shows a permanent indicator that location sharing is active. |
---
<br/>

| R2.2 | Successfully deactivate location sharing |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`AND`| you activated location sharing previously |
|`WHEN`| you navigate to the 'settings screen' |
|`AND`| deactivate location sharing |
|`THEN`| the app disables location sharing |
|`AND`| does __NOT__ shows a permanent indicator that location sharing is active. |
---
<br/>

## Feature: View my tracked movements on a map [R3]

| R3.1 | Successfully view tracked movements |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`AND`| you activated location sharing previously |
|`AND`| you have moved around for at least 5 minutes since then |
|`WHEN`| you navigate to the 'view my movements screen' |
|`THEN`| the app shows your previously tracked movements on a map. |
---
<br/>

## Feature: Report a sighting [R4]

| R4.1 | Successfully report a sighting |
|:-:|:-|
|`GIVEN`| you are on the 'report sighting screen' |
|`WHEN`| you selected animal species 'Hooglander' |
|`AND`| you specified two adult female animals and one infant male animal |
|`AND`| you specified as date+time yesterday 12:15
|`AND`| you picked a location about 15 meters from where you are now |
|`AND`| you wrote "Landholder R4.1" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your sighting |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor waarnemingen' |
|`AND`| if you provide the following answers 1: “Vervelend”, 2: “Het was erg koud”, “Er lag sneeuw”, “Er was onweer” |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen' |
|`AND`| a message with title "Een vervelende waarneming" and type URGENT is displayed explaining that it is showing because you selected "Vervelend" for question one. |
---
<br/>

## Feature: Report property damage [R5]

| R5.1 | Successfully report property damage |
|:-:|:-|
|`GIVEN`| you are on the 'report property damage screen' |
|`WHEN`| you selected animal species 'Wolf' |
|`AND`| you specified "mijn kippen" as damaged belonging |
|`AND`| you specified "eenheden" as type and 5 as value |
|`AND`| you specified 500 euros as damage value and 6500 as estimated loss |
|`AND`| you specified as date+time yesterday 4:45
|`AND`| you picked your current location |
|`AND`| you wrote "Landholder R5.1" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your property damage |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor schademeldingen' |
|`AND`| if you provide the following answers 1: “De kippen zijn dood”, 2: “Nee dat weet ik niet, omdat …” and then enter the free text “het de eerste keer is”. |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen'. |
---
<br/>

## Feature: Report an animal-vehicle collision [R6]

| R6.1 | Successfully report an animal-vehicle collision |
|:-:|:-|
|`GIVEN`| you are on the 'report an animal-vehicle collision screen' |
|`WHEN`| you selected animal species 'Ree' |
|`AND`| you specified one adolescent female animal |
|`AND`| you specified 4000 euros as estimated damage value |
|`AND`| you specified MEDIUM intensity and LOW urgency |
|`AND`| you specified the current date+time |
|`AND`| you picked a location about 25 meters from where you are now |
|`AND`| you wrote "Landholder R6.1" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your animal-vehicle collision |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor dieraanrijdingen' |
|`AND`| if you provide the following answers 1: “Een personenauto”, 2: “5”. |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen'. |
---
<br/>

## Feature: View my interactions [R11]

| R11.1 | Successfully list previously added interactions |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'view my interactions screen' |
|`THEN`| the interactions you just successfully reported for R3, R4 and R5 are displayed with all provided information visible. |
---
<br/>

## Feature: Designate an area of interest and receive alarms [R13]

| R13.1 | Successfully designate an area of interest |
|:-:|:-|
|`GIVEN`| you are on the 'areas of interest screen' |
|`WHEN`| you add an area of interest |
|`AND`| you specified as animal species 'Hooglander' |
|`AND`| you specified a centroid somewhere in the middle of the living lab |
|`AND`| you specified a radius of about 800 metres |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your area of interest. |
---
<br/>

| R13.2 | Successfully receive an alarm and message |
|:-:|:-|
|`GIVEN`| you have previously successfully designated an area of interest |
|`WHEN`| an animal moves into your area, or a new detection is added within your area or a new interaction is reported within your area |
|`THEN`| the app shows a notification of type ALARM explaining what was added and in which area of interest |
|`AND`| if the animal species for this alarm is 'Hooglander' |
|`THEN`| a message of type URGENT is shown explaining what to do to avoid 'Hooglander' animals from entering your area of interest. |
---
<br/>

## Feature: View the location of animals etc. [R7]

| R7.1 | Successfully view vicinity |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'view vicinity screen' (map) |
|`THEN`| a view of your vicinity including animals, detections and interactions is shown |
|`AND`| you can move the view around |
|`BUT`| this does __NOT__ allow you to obtain information about areas outside of your direct vicinity. |
---
<br/>

## Feature: View animal detections and interactions reported by others [R8]

| R8.1 | Successfully view a detection |
|:-:|:-|
|`GIVEN`| you are on the 'view vicinity screen' |
|`WHEN`| you select a Detection |
|`THEN`| information about that detection is shown that includes at least the animals involved (species, lifestage, sex), the type of detection device (for example camera or acoustic) and the user that is the owner of this detection device |
|`AND`| optionally an URL to observe the detection using the web browser. |
---
<br/>

| R8.2 | Successfully view an interaction |
|:-:|:-|
|`GIVEN`| you are on the 'view vicinity screen' |
|`WHEN`| you select an Interaction |
|`THEN`| information about that interaction is shown that includes at least the animals involved (species, lifestage, sex), the type of interaction  (sighting, property damage, animal-vehicle collision) and the user that reported this interaction. |
---
<br/>

## Feature: Update my profile [R9]

| R9.1 | Successfully update profile |
|:-:|:-|
|`GIVEN`| you are on the 'update profile screen' |
|`WHEN`| you provide a full name |
|`AND`| you provide a date-of-birth |
|`AND`| you provide a gender |
|`AND`| you provide a postcode |
|`AND`| you provide a description |
|`AND`| you submit the form |
|`THEN`| the app successfully submits this information  |
|`AND`| shows the exact information as provided on the 'profile screen'. |
---
<br/>

## Feature: View questionnaires assigned to me and answers given [R11]

| R11.1 | Successfully view questionnaires |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'questionnaires screen' |
|`THEN`| a list of questionnaires that were previously assigned to you is shown.  |
---
<br/>

| R11.2 | Successfully view answers given |
|:-:|:-|
|`GIVEN`| you are on the 'questionnaires screen' |
|`WHEN`| you select a questionnaire |
|`THEN`| all answers you gave on the questions in this questionnaire are shown. |
---
<br/>

## Feature: Close my profile [R10]

| R10.1 | Successfully close profile |
|:-:|:-|
|`GIVEN`| you are on the 'profile screen' |
|`WHEN`| you select 'close my profile' |
|`AND`| confirm the closure upon a warning message being shown |
|`THEN`| the app logs you out |
|`AND`| when you login with the same email address all previously entered information is no longer available. |
---
<br/>