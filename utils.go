package wildlifenl

import (
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lib/pq"
)

// BM: In an ideal world these text constants would not be hard coded but would be read from somewhere else. Unfortunately, the world is not always perfect.
const (
	appName        = "WildlifeNL"
	appDescription = "This is the WildlifeNL API. Before you can start making calls to the provided end-points you should acquire a bearer token. To do so, make a POST request at /auth/ providing the required fields including a valid email address. A validation code will be send to that email address. Then, make a PUT request at /auth/ providing the same email address and the validation code. The response will include a field named \"token\" containing your bearer token. Use this bearer token in the header of any future calls you make."
	emailSubject   = "Aanmelden bij WildlifeNL"
	emailBody      = "Beste {displayName}<br/>De applicatie {appName} wil graag aanmelden bij WildlifeNL met jouw emailadres. Om dit toe te staan, voer onderstaande code in bij deze applicatie.<br/>Code: {code}<br/><br/>Met vriendelijke groet<br/>WildlifeNL<br/><br/>"
)

func generateDescription(text string, scopes []string) string {
	if len(scopes) == 0 {
		return text
	}
	result := make([]string, 0)
	for _, value := range scopes {
		result = append(result, "`"+value+"`")
	}
	return text + "<br/><br/>**Scopes**<br/>" + strings.Join(result, ", ")
}

func handleError(err error) error {
	switch typedError := err.(type) {
	case *pq.Error:
		message := typedError.Message
		detail := typedError.Detail
		switch typedError.Code {
		case "23503": // violates foreign key constraint
			text := detail[4:strings.LastIndex(detail, ")")+1] + " does not exist."
			return huma.Error400BadRequest(text)
		case "23514": // violates check constraint
			text := strings.ReplaceAll(message, "new row for relation", "cannot add or update")
			text = strings.ReplaceAll(text, "violates", "because it violates")
			text = strings.ReplaceAll(text, "\"", "'")
			return huma.Error400BadRequest(text)
		}
	}
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Println("ERROR", frame.File+" "+strconv.Itoa(frame.Line)+" "+frame.Function+": ", err)
	return huma.Error500InternalServerError("")
}

func generateNotFoundByIDError(objectType string, id string) error {
	return huma.Error404NotFound("No '" + objectType + "' with ID '" + id + "' was found")
}

func generateNotFoundForThisUserError(objectType string, id string) error {
	return huma.Error404NotFound("No '" + objectType + "' with ID '" + id + "' was found for the current user")
}
