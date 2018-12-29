package http

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/filebrowser/filebrowser/types"
	"github.com/gorilla/mux"
)

func getUserID(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	i, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(i), err
}

type modifyUserRequest struct {
	modifyRequest
	Data *types.User `json:"data"`
}

func getUser(w http.ResponseWriter, r *http.Request) (*modifyUserRequest, bool) {
	if r.Body == nil {
		httpErr(w, r, http.StatusBadRequest, nil)
		return nil, false
	}

	req := &modifyUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpErr(w, r, http.StatusBadRequest, err)
		return nil, false
	}

	if req.What != "user" {
		httpErr(w, r, http.StatusBadRequest, nil)
		return nil, false
	}

	return req, true
}

func (e *Env) usersGetHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := e.getUser(w, r)
	if !ok {
		return
	}

	if !user.Perm.Admin {
		httpErr(w, r, http.StatusForbidden, nil)
		return
	}

	users, err := e.Store.Users.Gets()
	if err != nil {
		httpErr(w, r, http.StatusInternalServerError, err)
		return
	}

	for _, u := range users {
		u.Password = ""
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	renderJSON(w, r, users)
}

func (e *Env) userSelfOrAdmin(w http.ResponseWriter, r *http.Request) (*types.User, uint, bool) {
	user, ok := e.getUser(w, r)
	if !ok {
		return nil, 0, false
	}

	id, err := getUserID(r)
	if err != nil {
		httpErr(w, r, http.StatusInternalServerError, err)
		return nil, 0, false
	}

	if user.ID != id && !user.Perm.Admin {
		httpErr(w, r, http.StatusForbidden, nil)
		return nil, 0, false
	}

	return user, id, true
}

func (e *Env) userGetHandler(w http.ResponseWriter, r *http.Request) {
	_, id, ok := e.userSelfOrAdmin(w, r)
	if !ok {
		return
	}

	u, err := e.Store.Users.Get(id)
	if err == types.ErrNotExist {
		httpErr(w, r, http.StatusNotFound, nil)
		return
	}

	if err != nil {
		httpErr(w, r, http.StatusInternalServerError, err)
		return
	}

	u.Password = ""
	renderJSON(w, r, u)
}

func (e *Env) userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	_, id, ok := e.userSelfOrAdmin(w, r)
	if !ok {
		return
	}

	err := e.Store.Users.Delete(id)
	if err == types.ErrNotExist {
		httpErr(w, r, http.StatusNotFound, nil)
		return
	}

	if err != nil {
		httpErr(w, r, http.StatusInternalServerError, err)
	}
}

func (e *Env) userPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: fill me
}

func (e *Env) userPutHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser, modifiedID, ok := e.userSelfOrAdmin(w, r)
	if !ok {
		return
	}

	req, ok := getUser(w, r)
	if !ok {
		return
	}

	if req.Data.ID != modifiedID {
		httpErr(w, r, http.StatusBadRequest, nil)
		return
	}

	var err error

	if len(req.Which) == 1 && req.Which[0] == "all" {
		if !sessionUser.Perm.Admin {
			httpErr(w, r, http.StatusForbidden, nil)
			return
		}

		if req.Data.Password != "" {
			req.Data.Password, err = types.HashPwd(req.Data.Password)
		} else {
			var suser *types.User
			suser, err = e.Store.Users.Get(modifiedID)
			req.Data.Password = suser.Password
		}

		if err != nil {
			httpErr(w, r, http.StatusInternalServerError, err)
			return
		}

		req.Which = []string{}
	}

	for k, v := range req.Which {
		if v == "password" {
			if !sessionUser.Perm.Admin && sessionUser.LockPassword {
				httpErr(w, r, http.StatusForbidden, nil)
				return
			}

			req.Data.Password, err = types.HashPwd(req.Data.Password)
			if err != nil {
				httpErr(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		if !sessionUser.Perm.Admin && (v == "scope" || v == "perm" || v == "username") {
			httpErr(w, r, http.StatusForbidden, nil)
			return
		}

		req.Which[k] = strings.Title(v)
	}

	err = e.Store.Users.Update(req.Data, req.Which...)
	if err != nil {
		httpErr(w, r, http.StatusInternalServerError, err)
	}
}
