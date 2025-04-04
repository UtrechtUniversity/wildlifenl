# FUNCTIONAL DESIGN

*WildlifeNL: Applications and Data Platform*

## Introduction
This functional design provides the basis for the development of the WildlifeNL system. It incorporates the content of earlier documents about user stories and functionalities as gathered from interviews with the potential end-users, as well as the results of the brainstorm session organized by the researcher team. Ultimately unifying the wishes of the end-users with the requirements of the researcher team. Moreover, this functional design provides boundaries to the scopes of the different applications. 
To address feasibility an earlier version of this document was reviewed by the WildlifeNL consortium partners that participate in the implementation (Gamr-holding and Fontys University of Applied Sciences) and review comments have been processed in this version. 

## Roles
In the document ‘User Story’s.docx’, regarding the WildlifeNL user stories, 7 different roles were defined (Agrariër, Omwonende, Recreant, Faunabeheerder, Kuddebeheerder, Terreinbeheerder and Onderzoeker). Because developing 7 different applications is outside of the scope of the WildlifeNL project and some of the roles appear to have large overlaps, the roles Faunabeheerder, Kuddebeheerder, Terreinbeheerder have been combined into the more generic role Natural Resource Manager, so that they can be assigned to the same application. However, due to time limitations, initially the functionalities that are originate from the Kuddebeheerder role will not be implemented, effectively putting the role Herd Manager on pending. The role Agrariër was split into two parts so that one part can be merged with Omwonende into Inhabitant and use the corresponding Smartphone App, and the other part, Land User, is added to the group of other natural resource managers that use a Web App.
As a result, the manageable roles structure is as follows:

- 👲 Recreationist
- 🧑‍🌾 Inhabitant
- 🧑‍💼 Natural Resource Manager
    - Nature Area Manager
    - Wildlife Manager
    - Land User
    - Herd Manager (pending)
- 🧑‍🔬 Researcher
- 🧑🏻‍💻 Administrator

## Applications & User Stories
To address the objectives of the WildlifeNL project, four end-user applications have been defined and functionalities, collected from User Stories, have been added to each application. Note that functionalities with a priority of 4 or lower have been added for archiving purposes but will not receive any resources for implementation yet. Two applications are aimed at mobile use using a smartphone or tablet, whereas the other two are aimed at large screen use in a web browser on a laptop/desktop, but also functions, albeit less comfortable, in a mobile web-browser on a smartphone or tablet. All applications share data via the WildlifeNL API (see chapter Architecture). An Administrator role is defined, but no application will be created yet for this role. Current administrators will use the default user interface of the Wildlife API to activate functionalities assigned to their role.

- 👲 Recreationist -> 📱 WildGids (smartphone)
- 🧑‍🌾 Inhabitant ->📱 Wilde Buren (smartphone)
- 🧑‍💼 Natural Resource Manager -> 💻 WildRadar (web browser)
- 🧑‍🔬 Researcher -> 💻 ResearchConnect (web browser)
- 🧑🏻‍💻 Administrator -> 💻 Administration (web browser)
 
	
### WildGids
*Being able to enjoy human-wildlife interactions and reduce possible negative encounters with wildlife to recreate confidently and well-informed in a living lab*

👲 Recreationist | 📱 Smartphone App
   
I love being in and learning about nature, but feeling safe is important to me. Learning about the wild animals around me increases my positive nature experience. Knowing how to behave around wild animals or how to interact, or not interact, with them makes me feel more secure. I enjoy the freedom to choose activities like running, cycling, horseback riding, or walking. Sometimes I like being alone in nature, other times with my family and dog, or in a larger group. It is helpful to know where wild animals are so I can plan my route accordingly. Occasionally, I find it exciting to take photos of wildlife, but I want to make sure that it is appropriate or at the right time. I find it useful to register interactions I have with wild animals as this supports researchers, and if I spot something broken, I wish to report it for maintenance. Getting a notification with information on how to behave when I am approaching a potentially dangerous animal gives me peace of mind and reduces the probability of me performing undesirable behaviour (e.g., behaviour that leads to conflict with the wild animal or that negatively affects the welfare of the wild animal). 

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|View the location of wild animals (collars, camera trap images, etc.)|… can go and see them in real life and perhaps make a photograph. … can avoid them if I do not want to have an interaction with them.|1|
|Receive a message (push) if an encounter with a wild animal is imminent (distance based).|… know how to act accordingly, or can prevent a dangerous interaction.|1|
|Report a human-wildlife interaction and fill-out the corresponding questionnaire.|… help researchers and managers regarding presence information of wild animals (sighting). … report a wild animal-vehicle collision (animal-vehicle collision). … help researchers by providing information on how I experience different human-wildlife interactions and get advice based on the answers that I gave.|2|
|View information on wild animal species and how to interact with them.|… learn more about their behaviour, history, reasons for being in this environment, etc. … know how to behave and prevent dangerous situation.|2|
|Update my profile (username, date of birth, postcode, etc.)|… help researchers by providing context information.|2|
|Play an educational game about wild animals.|… , regardless of being a young person, remain interested in wildlife management.|3|
|Report a maintenance notice.| … help in providing the information to the nature area managers and the problem gets resolved quicker.|4|
|Watch webcam/nestcam streams.|	… can enjoy nature even from home. … can observe animal sanctuaries even though, as a human, I cannot enter them.|5|



### Wilde Buren
*As an inhabitant of the living lab, I live with wild animals and want to do so in a way that safeguards the well-being of myself and my property but also that of the wild animals.*

🧑🏻‍🌾 Inhabitant	📱 Smartphone App
  
I want to be better aware of the wild animals around me so that I can enjoy the positive interactions but also mitigate the negative ones. Being better informed about the animals around me and about ways to promote coexistence with them increases my tolerance towards these wild animals. I want to ensure that my pastures, fields, animals, garden, house, and property are safe from damage by wild animals. Knowing how to protect my land is crucial. Understanding which animals are near my property or the road I often travel on and how to respond to their presence and interact with them can help me avoid risks. It is essential for me to receive honest and proactive information, like if there are crop-raiding animals nearby and how to respond. I need the ability to report interactions with wild animals (e.g., wild animal observations, vehicle collisions, damage to property) and clarity on potential compensation and the necessary steps. Learning more about my wild animal neighbours, who they are, what they do and why they do this, helps me understand them better and increase my tolerance. Sharing positive experiences with my wild animal neighbours with others in my neighbourhood increases our shared sense of place and belonging.

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|View the location of wild animals (collars, camera trap images, etc.)|… can apply preventive measures. … get the confirmation that wild animals do not belong in my street and garden. … can enjoy nature if I want to. … can use this information to support my property damage claim.|1|
|Report a human-wildlife interaction and fill-out the corresponding questionnaire.|	… help researchers and managers regarding presence information of wild animals (sighting). … report damage that was inflicted by wild animals (property damage). … report a wild animal-vehicle collision (animal-vehicle collision). … help researchers by providing information on how I experience different human-wildlife interactions and get advice based on the answers that I gave.|1|
|View information on wild animal species and how to interact with them.|	… learn more about the importance of their presence and perhaps change my opinion. … get advice on how to interact with them.|2|
|Update my profile (username, date of birth, postcode, etc.)|… help researchers by providing context information.|2|
|Designate (point with radius) an area of interest, i.e. my fields or garden, and receive a message (push) when animals are within this area.|… can take measures to prevent damage to my crops, herds, garden, etc.|4|
|View the location where wild animals cross the roads.|	… can avoid these places or be more attentive and improve traffic safety.|4|
|Report a maintenance notice.|	… help in providing the information to the nature area managers and the problem gets resolved quicker.|4|
|Receive a message if an animal disease is reported in my neighbourhood.|… take preventive measures to protect my herd or pet.|5|

 
### WildRadar
*Adaptive management of wild animals and human-wildlife interactions for more positive and less negative human-wildlife interactions*

🧑‍💼 Natural Resource Manager	💻 Web App
   
Note that this application presents a slightly different user interface and set of functionalities dependent on the role(s) of the logged in user; Nature Area Manager, Wildlife Manager, Land User and later also Herd Manager (pending). Some functionalities are associated with several roles for different reasons, and some are exclusive to a specific role. It is possible for a user to have multiple roles and the application should adapt accordingly. 

#### Nature Area Manager
I serve as a manager of nature areas and my responsibility is providing space and developing habitat for wild animals but also for humans to recreate. Within this role, I am responsible for human-wildlife interactions in my areas, where I aim to reduce negative interactions and impacts on wildlife and humans present within these areas. I am also a neighbour, e.g. with adjacent agricultural lands or villages, and in that context responsible for maintaining respectful relations with these neighbours. Since human-wildlife interactions often cross the boundaries of my areas, I also need to work with my neighbours to manage these interactions. I manage wildlife and their habitat for the benefit of biodiversity, incl. the wildlife, but also people inside and around my areas. This means I also need to manage people. I offer education and guided tours. Additionally, I survey the flora and fauna, providing advice for nature management. I may be specialized in public engagement, management, ecology, or supervision.

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|View the location of wild animals (collars, camera trap images, etc.)|	… can protect wildlife sanctuaries or foraging areas against overgrazing or underutilization. … compare the utilization of areas with nature results. … can influence biodiversity in a positive way for example by giving more space.… can estimate the ratio of population size to the damage in the surrounding area, and conclude on the most suitable intervention (with the least impact).|1|
|View the movement of a group of wild animals.|	… can influence vegetation development through additional management. For example, more/less grazing/mowing or placing/removing fences. … see what effect events have on large animals, and can choose more consciously which events to allow and/or organize.|1|
|Designate (point with radius) an area of interest, and receive a message when animals are within this area.| … can take preventive measures against animals moving into areas that they should not.|1|
|View reported human-wildlife interactions (from smartphone apps).| … can adjust my management strategy accordingly. … can estimate the impact of the wild animal population.|2|
|View the number of humans counted in a certain area during a defined period.|… can organize recreation in a way that has the least negative impact on the wildlife. … can inform recreationists on how to behave in the vicinity of wild animals.	|2|
|View the tracked locations of humans.|	… can see where recreationists leave the paths and can intervene in order to protect wildlife sanctuaries.	|3|
|View reported maintenance notices i.e. fences etc.|	… can resolve the problem.	|4|
|Get insight in the behaviour of a wild animal at a certain moment in time.|	… can estimate which part of the population get preyed on by wolves and can adjust stock to the target stock.	|5|
|View interactions between animals.|	… can optimize my management by using different ratios of grazers, of applying fauna management clustering for the benefit of other species … know the effects of the presence of wild boar and/or wolves.	|5|

#### Wildlife Manager
My responsibility is to prevent or reduce damage or nuisance caused by certain wild animals while at the same time safeguarding, and sometimes promoting, the well-being of animals and their populations. With other wildlife managers, I am part of  a Wildbeheereenheid (Wildlife Management Unit), which works together with a Faunabeheereenheid at the provincial level. Our tasks may include the monitoring of wild animal populations and of the damage they cause, managing the numbers and behaviour of wild animals to prevent them to cause damage, ensuring management is conducted responsibly, and focusing on improving wildlife habitat and biodiversity.

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|View the location of wild animals (collars, camera trap images, etc.)|	… can apply measures where needed (fences, ecology of fear, culling, etc.)	|1|
|View the movement of a group of wild animals.|	… can correlate this to other nature values in the area and can intervene if needed.	|1|
|Designate (point with radius)  an area of interest, and receive a message (push) when animals are within this area.|	… can take preventive measures against animals moving into areas that they should not.	|1|
|View the paths that animals take to get into my area of interest.|	… investigate whether a preventive measure is needed, or can apply ecology of fear.	|2|
|View reported human-wildlife interactions (from smartphone apps).|	… get insight in how interactions with large wild animals (i.e. deer) are experienced, ultimately helping me adjust my management strategy. … can estimate the property damage impact of the presence of wild animals.	|2|
|View the number of humans counted in a certain area during a defined period.|	… infer the spatial reaction of wild animals on the presence of hunters and/or recreationists with or without pets (dog).	|3|
|View the number of wild animals counted (per species) in a certain area during a defined period.|	… infer the spatial reaction of wild animals on the presence of hunters and/or recreationists with or without pets (dog).	|3|
|Activate an existing animal behaviour influencing device (i.e. repellent).|	… can apply ecology of fear and/or culling more effectively.	|3|
|View food availability for wild animals in certain areas.|	… know where animals are likely going to move to or from and can estimate where they are going for reproduction.	|5|
|View the population sizes of wild animals.|	… can change my management strategy accordingly in order to maintain the populations as well as prevent damage.	|5|
|Report the population sizes of wild animals.|	… can change my management strategy accordingly in order to maintain the populations as well as prevent damage.	|5|

#### Land User
I use, and often own, land to produce food (crops, meat, dairy) or other natural resources (e.g., wood). In doing this, I affect the landscapes in which wild animals live. Through these effects, I may positively and/or negatively affect wild animal individuals and populations. Wild animals may also affect my work and production and use of natural resources. Because of this I am interested to learn about ways for us to share the landscapes we live in in ways that allow both of us to make an honest living. For this, I am interested to learn more about the wild animals on and around my land, where they live, with how many they are, why they do what they do, and how the actions that I take influence all of this.

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|View the location of wild animals (collars, camera trap images, etc.)|	… can apply preventive measures.|1|
|Designate (point with radius)  an area of interest, i.e. my fields or farm, and receive a message (push) when animals are within this area.|	… can take measures to prevent damage to my crops, herds, etc.	|1|
|View the movement of a group of wild animals.|	… see how my preventive measures, and those of my neighbours, affect the wild animal movements.	|1|
|View the paths that animals take to get into my area of interest.|	… investigate whether a preventive measure is needed, or can detect and report property damage.	|2|
 
### ResearchConnect
*Being able to access the data from the data platform to use it in my research*

🧑‍🔬 Researcher	💻 Web App
 
Regarding my research, I wish to use applications of WildlifeNL to run experiments in order to gather data that I can use to address my research objectives and/or hypotheses. To do so, I would like to create questionnaires that end-users of the applications can fill-out, and I would like to send push messages to the end-users, hopefully nudging them into certain behaviours. Lastly, I want to access the data platform in a way that allows me to extract complete datasets and/or aggregations that I can use in my own research. The data platform does not need to present the data in aesthetically appealing ways, I just need the raw data. 

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|Create an experiment.|	… can address my research questions and find evidence for any hypothesis that I have.	|1|
|Create and assign a questionnaire.|	… can gather data from the end-users.	|1|
|Create and assign messages.|	… can inform the end users and perhaps nudge them into a certain behaviour.	|1|
|Download data that is collected by the WildlifeNL project in a raw format.|	… can use this data in my research.|1|
|Download data that is collected by the WildlifeNL project, providing filters like start and end date, or for specific users or animals, and/or in aggregations.|	… can use this data in my research.	|2|


 
### Administration
*Assign roles to other users and manage the content of the API backend.*

🧑🏻‍💻 Administrator	💻 Web App
 
To ensure that the WildifeNL API and its backend contain the correct information and this information is provided to the correct users in the other applications, I can add, edit, and disable data elements. Also, I can assign and revoke roles for other users.
Note that no specific admin tool will be created, for now the administrators will use the default user interface of the WildlifeNL API.

|Functionality|So that I …|Priority|
|-------------|-----------|--------|
|Assign or revoke user roles for other users.|	… can manage authorisation.	|1|
|Manage (add, edit, disable?) the animal species in the system.|	… ensure the correct information on animal species is available for the other applications.	|1|
Manage (add, edit, disable?) the living labs in the system.|	… ensure the correct information on living labs is available for the other applications.	|2|

 
## Data Platform
Because the WildlifeNL system incorporates multiple applications addressing several different roles and devices, data needs to be exchanged between those applications. A central Data Platform is responsible for accepting, retrieving and storing data for the WildlifeNL project and is disclosed via the WildlifeNL API (application programming interface). The functionalities that the Data Platform offers are the direct result of the user stories and requirements of the applications that it needs to support, and thus its capacity to store data is bound by the requirements for those data by the WildlifeNL applications. In more technical terms: the data platform serves as a data back-end for the applications, rather than as a repository for data in general. However, the data platform can be extended as to integrate with existing data repositories in such a way that those data can be retrieved by the WildlifeNL applications as well.

## Requirements
WildlifeNL is a research project that aims primarily for societal impact. The purpose of the project is not to integrate the scattered wildlife management data landscape in the Netherlands, nor to explore the latest digital technologies, but rather to apply existing technologies in novel ways as possible solutions for issues in wildlife management in the Netherlands. Consequently, the goal of the data platform is primarily to implement the user stories as collected within the WildlifeNL consortium, and secondarily to support the multidisciplinary team of researchers (ecological, social, legal, governance, and communication perspective) in both data collection and "nudging" of behaviour of actors in the wildlife management domain.
Therefore, the data platform should be ‘lightweight’, in particular it should adhere to the following general principles:
- Loose coupling: The architecture should be open. It should be easily to add new existing monitoring systems (camera-traps, collars, etc.). This also implies that the project should actively avoid vendor-lock-in.
- Flexible: The data to and from the apps should be configurable such that researchers should be able to make selections (queries) in the collected pool of data, and they should be able to set up their own questionnaires that can be send to the various end-users.
- Open-source: The code that will be developed in this project should be open-source (both for the apps and the data platforms). This fits with the open access philosophy of NWO. 
- Privacy-by-design: Due to the potentially sensitive nature of some of the data that will be collected in the WildlifeNL project, the data storage solution will be provided by Utrecht University, with the ultimate goal to prevent misuse of personal data.

## Quality Attributes

### Performance
The WildlifeNL Data Platform will be hosted on the Utrecht University’s RedHat OpenShift platform (OpenShift), which is a containerized virtual environment that can automatically scale system resources based on load.

### Usability
The Data Platform discloses a REST API using JSON formatted data, serving multiple endpoints that address functional requirements as specified for the WildlifeNL applications. The API documentation is generated automatically from its code base as to ensure that it does not get outdated when new functionality is added or improvements are made to existing functionality. This API documentation is available at a public URL.

### Security
Before any application can connect with the Data Platform it needs its user to pass an authentication and authorization procedure. Authentication happens by the user offering an email address. An authentication code will be send to that address and if the user can provide that code back to the Data Platform, the email address is validated and the user is authenticated. This is form of Simple Authentication for the Web (SAW) an implementation of the ‘One Time Password (OTP) by email’ scheme ([Bonneau et al., 2012](https://doi.org/10.48456/tr-817)) which allows authentication without users needing to possess and remember a password. Authorization is based on the Roles as mentioned above in this document. By default any user can utilize the functionalities associated with the roles Recreationist and Inhabitant without any further need for authorization. An administrator can assign additional roles to registered users to allow them to utilize more functionalities. This split of default roles and additional roles allows for quick registration of new users but protects sensitive functionalities behind a layer of security using a flat Role Based Access Control (RBAC) mechanism ([Sandhu et al., 2000](https://doi.org/10.1145/344287.344301)).

### Maintainability
Components of the Data Platform run as containerized pods on the OpenShift platform and can be pulled up and torn down using the OpenShift management console. Moreover, OpenShift offers container builders that allow direct rebuilding of a specific container from a git repository when a newer version is available.

### Extensibility
The Data Platform will be developed in an object oriented modular design, facilitating industry standard ways for extending its functionalities. Given that the development process is iterative, earlier versions of the Data Platform will have limited functionality and later renditions will extend upon the earlier versions. This approach allows for future extensions based on evaluations of earlier versions even for functionalities that have
 
## Architecture
The applications work together in an architecture that defines the WildlifeNL system as specified in the image below, and cater for all the roles that are being addressed. The Data Platform discloses an API for the application to connect to and exchange data. Moreover, the WildlifeNL system allows for extension by accepting data from external systems. Multiple of such systems can be connected to the WildlifeNL system for example to acquire information on the location of animals from different sources (collar, camera trap, acoustic sensor, etc.) as well as visitor counts or additional geographic information. As an initial setup an existing Trapper ([Bubnicki et al., 2016](https://doi.org/10.1111/2041-210X.12571)) installation is added to allow Natural Resource Managers to deploy camera trap projects. Trapper offers an API for data exchange with the WildlifeNL system and thus can be connected to the WildlifeNL API. Also, an existing system for collar data from the consortium member Smart Parks will be connected to the WildlifeNL API.

![System-Architecture](SystemArchitecture.svg)
 
*existing system 


## Conceptual Model
The conceptual model describes the main concepts as used in the WildlifeNL system and how they are associated in a high-level overview.

![Conceptual-Model](ConceptualModel.svg)
 
A User has one or more Roles, and can report an Interaction with wildlife of a certain Species. As a result of the Interaction the user is possibly presented with a Questionnaire to fill out. Based on the answers given in the questionnaire it is possible that the user receives a Message. Moreover, the user can have an Encounter (being close to each other in terms of time, latitude and longitude) with an Animal of a certain Species and receive a Message about this encounter. Lastly the User can specify a Zone of interest and receive an Alarm when an Animal enters that zone, when an Interaction is reported within that zone, or when a Detection is reported within that zone. A Detection is a record of a sensor, for example a camera trap, having detected a member of a certain Species.
