package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"project-cahillaw/servers/gateway/models/users"
	"project-cahillaw/servers/gateway/sessions"
	"strconv"
	"strings"
	"time"
)

func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")
	correct := strings.HasPrefix(ct, "application/json")
	if !correct {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Request body must be in JSON form"))
	}

	nu := &users.NewUser{}
	err := json.NewDecoder(r.Body).Decode(&nu)
	if err != nil {
		http.Error(w, "Incorrect formmating", http.StatusBadRequest)
		return
	}

	_, errMarshal := json.Marshal(nu)
	if errMarshal != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	us, errUser := nu.ToUser()
	if errUser != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}

	user, errInsert := ctx.UserStore.Insert(us)
	if errInsert != nil {
		http.Error(w, "Error creating user", http.StatusBadRequest)
		return
	}

	_, errMarshalUser := json.Marshal(user)
	if errMarshalUser != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ss := &SessionState{}
	ss.SessionStart = time.Now()
	ss.User = user
	_, errSess := sessions.BeginSession(ctx.Key, ctx.SessionsStore, ss, w)
	if errSess != nil {
		http.Error(w, "Error creating session", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoded, errEncode := json.Marshal(user)
	if errEncode != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
		return
	}
	w.Write(encoded)
}


func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	ss := &SessionState{}
	_, err := sessions.GetState(r, ctx.Key, ctx.SessionsStore, ss)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	userid := ss.User.ID
	pathid := path.Base(r.URL.Path)

	if r.Method == http.MethodGet {
		i, errConv := strconv.Atoi(pathid)
		if errConv != nil && pathid != "me" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		} else {
			if pathid == "me" {
				i = int(userid)
			}
			user, errGetUser := ctx.UserStore.GetByID(int64(i))
			if errGetUser != nil {
				http.Error(w, "Not a valid userid", http.StatusNotFound)
				return
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				encoded, errEncode := json.Marshal(user)
				if errEncode != nil {
					http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
					return
				}
				w.Write(encoded)
			}
		}
	} else if r.Method == http.MethodPost {
		i, errConv := strconv.Atoi(pathid)
		if errConv != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		errFollow := ctx.UserStore.Follow(userid, int64(i))
		if errFollow != nil {
			http.Error(w, errFollow.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Changed Follow Status"))
	} else if r.Method == http.MethodPatch {
		i, errConv := strconv.Atoi(pathid)
		if errConv != nil && pathid != "me" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if pathid == "me" {
			i = int(userid)
		}

		if int64(i) != userid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ct := r.Header.Get("Content-Type")
		correct := strings.HasPrefix(ct, "application/json")
		if !correct {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Request body must be in JSON form"))
		}

		up := &users.Updates{}
		errDecode := json.NewDecoder(r.Body).Decode(&up)
		if errDecode != nil {
			http.Error(w, "Incorrect formmating", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		errApply := ss.User.ApplyUpdates(up)
		_, errUpdate := ctx.UserStore.Update(ss.User.ID, up)
		if errApply == nil && errUpdate == nil {
			encoded, errEncode := json.Marshal(ss.User)
			if errEncode != nil {
				http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
				return
			}
			w.Write(encoded)
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ct := r.Header.Get("Content-Type")
		correct := strings.HasPrefix(ct, "application/json")
		if !correct {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Request body must be in JSON form"))
		}

		cr := &users.Credentials{}
		errDecode := json.NewDecoder(r.Body).Decode(&cr)
		if errDecode != nil {
			http.Error(w, "Incorrect formmating", http.StatusBadRequest)
			return
		}

		dummy := &users.User{}
		dummy.PassHash = []byte("randomstuff123")
		success := true
		us, errGetUser := ctx.UserStore.GetByEmail(cr.Email)
		if errGetUser != nil {
			//will always fail
			errDummy := dummy.Authenticate("notthesameashash")
			if errDummy != nil {
				success = false
			}
		} else {
			errAuth := us.Authenticate(cr.Password)
			if errAuth != nil {
				success = false
			}
		}

		if success == true {
			ss := &SessionState{}
			ss.SessionStart = time.Now()
			ss.User = us
			_, errSess := sessions.BeginSession(ctx.Key, ctx.SessionsStore, ss, w)
			if errSess != nil {
				http.Error(w, "Error creating session", http.StatusBadRequest)
				return
			} else {

				xfor := r.Header.Get("X-Forwarded-For")
				address := r.RemoteAddr
				if len(xfor) > 0 {
					ip := strings.Split(xfor, ", ")
					address = ip[0]
				}

				trackLogin := ctx.UserStore.TrackLogin(ss.User.ID, ss.SessionStart, address)
				if trackLogin != nil {
					fmt.Println(trackLogin.Error())
					http.Error(w, trackLogin.Error(), http.StatusBadRequest)
					return
				}

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				encoded, errEncode := json.Marshal(us)
				if errEncode != nil {
					http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
					return
				}
				w.Write(encoded)
			}
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		path := path.Base(r.URL.Path)
		if strings.ToLower(path) != "mine" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			_, errEnd := sessions.EndSession(r, ctx.Key, ctx.SessionsStore)
			if errEnd != nil {
				http.Error(w, "Failed to end session", http.StatusBadRequest)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Signed out"))
			}
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
