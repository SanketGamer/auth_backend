package notifier

import(
	"fmt"
	"log"
	"os"
	"time"
)

type EventType string

const (
    EventUserCreated EventType = "USER_CREATED"
    EventUserLoggedIn EventType = "USER_LOGGED_IN"
)

type UserEvent struct{
	Type EventType
	UserID string
	Email string
	Timestamp time.Time
}

type Notifier struct{
	events chan UserEvent
	logger *log.Logger
}

//creates and open a logfile
//creates a buffred event channel
//starts a background worker(listens to events channel,processes events,logs them)
func NewNotifier(bufferSize int)*Notifier{
  file,err:=os.OpenFile("app.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
  if err!=nil{
	log.Fatal(err)
  }
  //Creates a buffered channel,Used to send events (like user signup, login, etc.)
  n:=&Notifier{
	events: make(chan UserEvent,bufferSize),
	logger: log.New(file,"",log.LstdFlags),
  }
  go n.worker()
  
  return n
}

//listen events and logs them
func (n *Notifier)worker(){
  for event:=range n.events{
	msg:=fmt.Sprintf("[%s] UserId=%s Email=%s Time=%s",
	event.Type,
	event.UserID,
	event.Email,
	event.Timestamp.Format(time.RFC3339),
   )
   n.logger.Println(msg)
   log.Println(msg)
  }
}

//send an event into channel
func (n *Notifier) Send(event UserEvent){
	select{
	case n.events<-event:
	default:
		log.Println("[NOTIFIER] dropped event, channel full")
	}
}