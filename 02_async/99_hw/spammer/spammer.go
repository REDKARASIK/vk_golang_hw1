package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	in := make(chan interface{})
	prevOut := in
	var wg sync.WaitGroup
	wg.Add(len(cmds))
	for _, c := range cmds {
		nextOut := make(chan interface{})
		go func(cmd cmd, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			c(in, out)
		}(c, prevOut, nextOut)
		prevOut = nextOut
	}

	donePipeline := make(chan struct{})
	go func(lastOut chan interface{}) {
		for range lastOut {
			continue
		}
		close(donePipeline)
	}(prevOut)
	close(in)
	wg.Wait()
	<-donePipeline
}

func SelectUsers(in, out chan interface{}) {
	// 	in - string
	// 	out - User
	var wg sync.WaitGroup
	var maps sync.Map

	for emailIn := range in {
		var email string
		switch emailIn.(type) {
		case string:
			email = emailIn.(string)
		default:
			panic("Invalid email type")
		}
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			user := GetUser(email)
			if _, ok := maps.LoadOrStore(user.Email, struct{}{}); !ok {
				out <- user
			}
		}(email)
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID

	batch := make([]User, 0, 2)
	var wg sync.WaitGroup

	for userIn := range in {
		var user User

		switch userIn.(type) {
		case User:
			user = userIn.(User)
		default:
			panic("Invalid user type")
		}
		if len(batch) == 2 {
			wg.Add(1)
			batchFunc := make([]User, 2)
			copy(batchFunc, batch)
			go func(b []User) {
				defer wg.Done()
				res, err := GetMessages(b...)
				if err != nil {
					panic(err)
				}
				for _, msg := range res {
					out <- msg
				}
			}(batchFunc)
			batch = batch[:0]
		}
		batch = append(batch, user)
	}
	if len(batch) > 0 {
		res, err := GetMessages(batch...)
		if err != nil {
			panic(err)
		}
		for _, msg := range res {
			out <- msg
		}
	}
	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData

	var wg sync.WaitGroup
	parallels := make(chan struct{}, 5)

	for msgIn := range in {
		var msgID MsgID
		switch msgIn.(type) {
		case MsgID:
			msgID = msgIn.(MsgID)
		default:
			panic("Invalid msg type")
		}
		wg.Add(1)
		go func(msg MsgID) {
			defer wg.Done()
			parallels <- struct{}{}
			res, err := HasSpam(msg)
			if err != nil {
				panic(err)
			}
			<-parallels
			out <- MsgData{
				ID:      msgID,
				HasSpam: res}
		}(msgID)
	}
	wg.Wait()
	close(parallels)
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string

	var msgs []MsgData

	for msgIn := range in {
		var msgData MsgData
		switch msgIn.(type) {
		case MsgData:
			msgData = msgIn.(MsgData)
		default:
			panic("Invalid msg type")
		}
		msgs = append(msgs, msgData)
	}

	sort.Slice(msgs, func(i, j int) bool {
		if msgs[i].HasSpam != msgs[j].HasSpam {
			return msgs[i].HasSpam // true пойдёт раньше false
		}
		return msgs[i].ID < msgs[j].ID
	})
	for _, msg := range msgs {
		out <- fmt.Sprintf("%t %d", msg.HasSpam, msg.ID)
	}
}
