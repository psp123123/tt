package Web
//
//import (
//	"log"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//func Sendmsg() {
//	es := eventsource.New(
//		eventsource.DefaultSettings(),
//		func(req *http.Request) [][]byte {
//			return [][]byte{
//				[]byte("X-Accel-Buffering: no"),
//				[]byte("Access-Control-Allow-Origin: *"),
//			}
//		},
//	)
//	defer es.Close()
//	http.Handle("/eventsend", es)
//	go func() {
//		id := 1
//		for {
//			es.SendEventMessage("tick", "tick-event", strconv.Itoa(id))
//			id++
//			time.Sleep(2 * time.Second)
//		}
//	}()
//	log.Fatal(http.ListenAndServe(":3000", nil))
//}
//
//
//
func web() {

}