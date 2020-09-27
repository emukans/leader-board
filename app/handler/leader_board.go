package handler

import (
	"encoding/json"
	"leader-board/app/model"
	"net/http"
	"time"
)

type leaderBoardResponse struct {
	Results  []model.PlayerScore `json:"results"`
	NextPage int                 `json:"next_page"`
	AroundMe []model.PlayerScore `json:"around_me"`
}

// I have to introduce a context struct, because golang don't have generics or inheritance
type leaderBoardContext struct {
	pageNumber      int
	limit           int
	offset          int
	totalScoreCount int
	periodFrom      time.Time
	sentName        string
	payload         []byte
}

func LeaderBoard(writer http.ResponseWriter, request *http.Request) {
	context := leaderBoardContext{}
	response := leaderBoardResponse{}

	ErrWriter{writer: writer}.Then(func(self ErrWriter) ErrWriter {
		pageNumber, err := parsePageNumber(request)
		self.err = err
		context.pageNumber = pageNumber

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		periodFrom, err := parsePeriod(request)
		context.periodFrom = periodFrom
		self.err = err

		return self
	}).MaybeHandleError(func(self ErrWriter) ErrWriter {
		http.Error(writer, self.err.Error(), http.StatusBadRequest)

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		context.limit, self.err = model.FindPageLimit()

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		context.offset = (context.pageNumber - 1) * context.limit

		response.Results, self.err = model.FindAllScores(context.limit, context.offset, context.periodFrom)

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		for rank, _ := range response.Results {
			response.Results[rank].Rank = context.offset + rank + 1
		}

		context.sentName = parseName(request)
		if context.sentName != "" {
			isNameInTheList := isNameInScoreList(response.Results, context.sentName)
			if !isNameInTheList {
				response.AroundMe, self.err = model.FindScoresAroundName(context.sentName, context.limit)
			}
		}

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		context.totalScoreCount, self.err = model.FindScoreCount()

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		if context.totalScoreCount > (context.offset + context.limit) {
			response.NextPage = context.pageNumber + 1
		}
		return self
	}).Then(func(self ErrWriter) ErrWriter {
		self.writer.Header().Set("Content-Type", "application/json")
		context.payload, self.err = json.Marshal(response)

		return self
	}).Then(func(self ErrWriter) ErrWriter {
		writer.WriteHeader(http.StatusOK)
		writer.Write(context.payload)

		return self
	}).MaybeHandleInternalError()
}
