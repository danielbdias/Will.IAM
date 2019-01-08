package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ghostec/Will.IAM/usecases"
	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
)

func serviceAccountsGetHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		// qs := r.URL.Query()
		// withRoles := qs["withRoles"]
		// if len(withRoles) == 1 {
		// 	return
		// }
		saID := mux.Vars(r)["id"]
		sa, err := sasUC.Get(saID)
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data := map[string]string{
			"id":         sa.ID,
			"name":       sa.Name,
			"baseRoleId": sa.BaseRoleID,
		}
		bts, err := json.Marshal(data)
		WriteBytes(w, 200, bts)
	}
}

func serviceAccountsCreateHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(body, &m)
		name, ok := m["name"].(string)
		if !ok || name == "" {
			Write(w, http.StatusUnprocessableEntity,
				`{ "error": { "name": "required" } }`)
			return
		}
		authenticationType, ok := m["authenticationType"].(string)
		if !ok || authenticationType == "" {
			Write(w, http.StatusUnprocessableEntity,
				`{ "error": { "authenticationType": "required" } }`)
			return
		}
		if authenticationType == "oauth2" {
			email, ok := m["email"].(string)
			if !ok || email == "" {
				Write(w, http.StatusUnprocessableEntity,
					`{ "error": { "email": "required" } }`)
				return
			}
			_, err = sasUC.CreateOAuth2Type(name, email)
			if err != nil {
				l.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if authenticationType == "keypair" {
			_, err = sasUC.CreateKeyPairType(name)
			if err != nil {
				l.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			Write(w, http.StatusUnprocessableEntity,
				`{ "error": { "authenticationType": "invalid" } }`)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func serviceAccountsListHandler(
	sasUC usecases.ServiceAccounts,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.GetLogger(r.Context())
		saSl, err := sasUC.List()
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bts, err := keepJSONFieldsBytes(saSl, "id", "name", "email")
		if err != nil {
			l.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		WriteBytes(w, 200, bts)
	}
}
