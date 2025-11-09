package wildlifenl

import (
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-mail/mail"
	"github.com/lib/pq"
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

// handleError uses some "fuzzy hocus pocus" to convert a Go error of different sources and natures to an huma REST error (of HTTP status 4xx) if it can infer the most likely type and a reasonable end-user message, otherwise it returns an huma REST error (of HTTP status 500) and logs the Go error.
func handleError(err error) error {
	switch typedError := err.(type) {
	case *mail.SendError:
		return huma.Error504GatewayTimeout("could not send email because the SMTP server is unavailable, please try again")
	case *stores.ErrRecordInattainable:
		return huma.Error404NotFound("record is inattainable: " + err.Error())
	case *stores.ErrRecordImmutable:
		return huma.Error409Conflict("record is immutable: " + err.Error())
	case *pq.Error:
		message := typedError.Message
		detail := typedError.Detail
		switch typedError.Code {
		case "23503": // violates foreign key constraint
			text := detail
			if strings.Contains(detail, "is still referenced from table") {
				text = detail[4:]
				text = strings.ReplaceAll(text, "referenced from table", "used by")
				text = strings.ReplaceAll(text, "\"", "'")
			} else {
				text = detail[4:strings.LastIndex(detail, ")")+1] + " does not exist."
			}
			return huma.Error400BadRequest(text)
		case "23505": // violates primary key constraint
			text := detail
			return huma.Error400BadRequest(text)
		case "23514": // violates check constraint
			text := strings.ReplaceAll(message, "new row for relation", "cannot add or update")
			text = strings.ReplaceAll(text, "violates", "because it violates")
			text = strings.ReplaceAll(text, "\"", "'")
			return huma.Error400BadRequest(text)
		case "P0001": // raised exception
			return huma.Error400BadRequest(message)
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
