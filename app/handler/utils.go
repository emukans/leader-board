package handler

import (
	"leader-board/app/model"
	"log"
	"net/http"
)


type ErrWriter struct {
	writer http.ResponseWriter
	err error
	isErrorHandled bool
}

func (receiver ErrWriter) Then(callback func(ErrWriter) ErrWriter) ErrWriter {
	if receiver.err != nil {
		return receiver
	}

	return callback(receiver)
}

func (receiver ErrWriter) MaybeHandleError(callback func(ErrWriter) ErrWriter) ErrWriter {
	if receiver.err != nil {
		receiver.isErrorHandled = true

		return callback(receiver)
	}

	return receiver
}

func (receiver ErrWriter) MaybeHandleInternalError() ErrWriter  {
	return receiver.MaybeHandleError(func(ErrWriter) ErrWriter {
		if !receiver.isErrorHandled {
			http.Error(receiver.writer, "Oops...something went wrong", http.StatusInternalServerError)
			log.Println(receiver.err)
		}

		return receiver
	})
}

func isNameInScoreList(scoreList []model.PlayerScore, sentName string) bool {
	isNameInTheList := false
	for _, score := range scoreList {
		if score.Name == sentName {
			isNameInTheList = true
		}
	}
	return isNameInTheList
}
