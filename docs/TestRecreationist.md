# WildlifeNL App tests – Recreationist
Today you are recreationist using the WildlifeNL WildGids app. We are running this app against a TEST environment with simulated animals, so there will only be human-wildlife interactions that are artificially created by us. All locations are within the test living lab area Botanische Tuinen in Utrecht. If you notice an issue, bug or inconvenience when performing the test scenarios below, please write them down so we can later provide them to the development team. Note that aside from the 'Logging in' feature, all features are marked with codes in brackets ([ ]) that match the user story codes in the Functional Design. You can look-up the expected behaviour of the app and the purpose of the functionality this way. Additionally, it is possible that the screens that are mentioned in single quotes in these test scenarios have a somewhat different Dutch name in the app. Please login with your email address either from Utrecht University or Fontys University of Applied Sciences.
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

## Feature: Read and accept Terms & Conditions and location sharing [G1]

| G1.1 | Unsuccessfully accept terms & conditions |
|:-:|:-|
|`GIVEN`| you are on the 'terms & conditions screen' |
|`WHEN`| you do __NOT__ tick the checkbox to accept the terms & conditions |
|`AND`| you try to continue |
|`THEN`| the app must __NOT__ navigate beyond the 'terms & conditions screen'. |
---
<br/>

| G1.2 | Unsuccessfully accept location sharing |
|:-:|:-|
|`GIVEN`| you are on the 'terms & conditions screen' |
|`WHEN`| you tick the checkbox to accept the terms & conditions |
|`AND`| you try to continue |
|`THEN`| the app asks you to activate location sharing |
|`AND`| if you deny |
|`THEN`| the app does __NOT__ navigate beyond the 'terms & conditions screen' |
|`AND`| a message is shown explaining that location sharing is mandatory |
|`AND`| the checkbox to accept the terms & conditions is unchecked. |
---
<br/>

| G1.3 | Successfully accept terms & conditions |
|:-:|:-|
|`GIVEN`| you are on the 'terms & conditions screen' |
|`WHEN`| you tick the checkbox to accept the terms & conditions |
|`AND`| you try to continue |
|`THEN`| the app asks you to activate location sharing |
|`AND`| if you accept |
|`THEN`| the app navigates beyond the 'terms & conditions screen' |
|`AND`| shows a permanent indicator that location sharing is active. |
---
<br/>

## Feature: View my tracked movements on a map [G2]

| G3.1 | Successfully view tracked movements |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`AND`| location sharing was successfully activated previously |
|`AND`| you have moved around for at least 5 minutes since then |
|`WHEN`| you navigate to the 'view my movements screen' |
|`THEN`| the app shows your previously tracked movements on a map. |
---
<br/>

## Feature: Report a sighting [G5]

| G5.1 | Successfully report a sighting |
|:-:|:-|
|`GIVEN`| you are on the 'report sighting screen' |
|`WHEN`| you selected animal species 'Wisent' |
|`AND`| you specified one adult female animal and one adolescent male animal |
|`AND`| you specified as date+time yesterday 10:15
|`AND`| you picked a location about 15 meters from where you are now |
|`AND`| you wrote "Recreationist G5.1" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your sighting |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor waarnemingen' |
|`AND`| if you provide the following answers 1: “Super Leuk!”, 2: “Het was zonnig”, “Het waaide”, “Het was erg heet” |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen'. |
---
<br/>

| G5.2 | Successfully report a sighting with message |
|:-:|:-|
|`GIVEN`| you are on the 'report sighting screen' |
|`WHEN`| you selected animal species 'Wezel' |
|`AND`| you specified one adult animal of unknown sex and two infant animals of unknown sex |
|`AND`| you specified as date+time now
|`AND`| you specified the current location |
|`AND`| you wrote "Landholder R5.2" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your sighting |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor waarnemingen' |
|`AND`| if you provide the following answers 1: “Vervelend”, 2: “Het was erg koud”, “Er lag sneeuw”, “Er was onweer” |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen' |
|`AND`| a message with title "Een vervelende waarneming" and type URGENT is displayed explaining that it is showing because you selected "Vervelend" for question one. |
---
<br/>

## Feature: Receive a push message [G4]

| G4.1 | Successfully receive push messages |
|:-:|:-|
|`GIVEN`| you move around |
|`WHEN`| you approach an animal of species Bever, Wisent, Das, Eekhoorn or Wild zwijn |
|`THEN`| a notification of type WARNING will show explaining that you have approached the animal in question and provide advice on how to proceed. |
|`WHEN`| you move away and approach another animal of species Bever, Wisent, Das, Eekhoorn or Wild zwijn |
|`THEN`| again a notification of type WARNING will show explaining that you have approached the animal in question and provide advice on how to proceed. |
---
<br/>

## Feature: Report an animal-vehicle collision [G6]

| G6.1 | Successfully report an animal-vehicle collision |
|:-:|:-|
|`GIVEN`| you are on the 'report an animal-vehicle collision screen' |
|`WHEN`| you selected animal species 'Damhert' |
|`AND`| you specified one adult male animal |
|`AND`| you specified 1000 euros as estimated damage value |
|`AND`| you specified HIGH intensity and HIGH urgency |
|`AND`| you specified the current date+time |
|`AND`| you picked a location about 25 meters from where you are now |
|`AND`| you wrote "Recreationist G6.1" as comment |
|`AND`| you submit the form |
|`THEN`| the app successfully submits your animal-vehicle collision |
|`AND`| presents the questionnaire 'Bas Vragenlijst voor dieraanrijdingen' |
|`AND`| if you provide the following answers 1: “Een vrachtwagen”, 2: “6”. |
|`THEN`| a message shows explaining that 6 is not a valid value for question two |
|`BUT`| if you provide the following answers 1: “Een vrachtwagen”, 2: “4”. |
|`THEN`| the app confirms successful submission of your answers and allows you to navigate back to the 'main screen'. |
---
<br/>

## Feature: View my interactions [G7]

| G7.1 | Successfully list previously added interactions |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'view my interactions screen' |
|`THEN`| the interactions you just successfully reported for G5 and G6 are displayed with all provided information visible. |
---
<br/>

## Feature: View information on wild animal species [G8]

| G8.1 | Successfully view information on wild animals |
|:-:|:-|
|`GIVEN`| you are on the 'animal information screen' |
|`WHEN`| you select the 'Haas' animal |
|`THEN`| information about the Haas animal's Latin Name, Common Name, Description, Advice, Behaviour, Category and Role in Nature is displayed, and none of them contains placeholder data |
|`AND`| if you close the Haas screen and you select the 'Egel' animal |
|`THEN`| information about the Egel animal's Latin Name, Common Name, Description, Advice, Behaviour, Category and Role in Nature is displayed, and none of them contains placeholder data. |
---
<br/>

## Feature: View the location of animals etc. [G3]

| G3.1 | Successfully view vicinity |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'view vicinity screen' (map) |
|`THEN`| a view of your vicinity including animals, detections and interactions is shown |
|`AND`| you can move the view around |
|`BUT`| this does __NOT__ allow you to obtain information about areas outside of your direct vicinity. |
---
<br/>

## Feature: View animal detections and interactions reported by others [G11]

| G11.1 | Successfully view a detection |
|:-:|:-|
|`GIVEN`| you are on the 'view vicinity screen' |
|`WHEN`| you select a Detection |
|`THEN`| information about that detection is shown that includes at least the animals involved (species, lifestage, sex), the type of detection device (for example camera or acoustic) and the user that is the owner of this detection device |
|`AND`| optionally an URL to observe the detection using the web browser. |
---
<br/>

| G11.2 | Successfully view an interaction |
|:-:|:-|
|`GIVEN`| you are on the 'view vicinity screen' |
|`WHEN`| you select an Interaction |
|`THEN`| information about that interaction is shown that includes at least the animals involved (species, lifestage, sex), the type of interaction  (sighting, property damage, animal-vehicle collision) and the user that reported this interaction. |
---
<br/>

## Feature: Update my profile [G9]

| G9.1 | Successfully update profile |
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

## Feature: View questionnaires assigned to me and answers given [G12]

| G12.1 | Successfully view questionnaires |
|:-:|:-|
|`GIVEN`| you are on the 'main menu screen' |
|`WHEN`| you navigate to the 'questionnaires screen' |
|`THEN`| a list of questionnaires that were previously assigned to you is shown.  |
---
<br/>

| G12.2 | Successfully view answers given |
|:-:|:-|
|`GIVEN`| you are on the 'questionnaires screen' |
|`WHEN`| you select a questionnaire |
|`THEN`| all answers you gave on the questions in this questionnaire are shown. |
---
<br/>

## Feature: Close my profile [G10]

| G10.1 | Successfully close profile |
|:-:|:-|
|`GIVEN`| you are on the 'profile screen' |
|`WHEN`| you select 'close my profile' |
|`AND`| confirm the closure upon a warning message being shown |
|`THEN`| the app logs you out |
|`AND`| when you login with the same email address all previously entered information is no longer available. |
---
<br/>