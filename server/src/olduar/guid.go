package olduar

type GUID int64
var currentGuid GUID = 0
var queue chan (chan GUID) = make(chan (chan GUID))

func GenerateGUID() GUID {
	value := make(chan GUID)
	queue <- value
	return <- value
}

func init() {
	go func(){
		for {
			request := <-queue
			currentGuid++
			request <- currentGuid
			close(request)
		}
	}()
}
