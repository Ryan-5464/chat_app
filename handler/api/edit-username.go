package api

import (
	"encoding/json"
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func EditUserName(a i.AuthService, u i.UserService) http.Handler {
	h := editUserName{
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type editUserName struct {
	userS i.UserService
}

func (h editUserName) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	var req eurequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	userId := session.UserId()

	res, err := h.handleRequest(req, userId)
	if err != nil {
		util.Log.Errorf("failed to handle contact edit chat name request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbugf("->>>> RESPONSE SENT:: %v", res)
}

func (h editUserName) handleRequest(req eurequest, userId typ.UserId) (euresponse, error) {
	util.Log.FunctionInfo()

	err := h.userS.EditUserName(req.Name, userId)
	if err != nil {
		return euresponse{}, err
	}

	return euresponse{Name: req.Name}, nil
}

type eurequest struct {
	Name string `json:"Name"`
}

type euresponse struct {
	Name string
}
