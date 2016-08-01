// This file contains code for supporting ACME challenges.
// Let's Encrypt uses ACME challenges for domain verification.

package mgae

import (
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const (
	ACME_CHALLENGE_ENTITY_ID = "acme-secret-singleton"
	ACME_CHALLENGE_URI_PATH  = "/.well-known/acme-challenge/"
)

const (
	entityName             = "AcmeSecret"
	challengeUriPathLength = len(ACME_CHALLENGE_URI_PATH)
)

type AcmeSecret struct {
	Challenge string // This is part of the URL.
	Response  string // This is the expected response.
	Timestamp time.Time
	UpdatedBy string
}

func SaveAcmeSecret(w http.ResponseWriter, r *http.Request) *Error {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if !u.Admin {
		return NewError(nil, "Only admin users can save ACME secrets.", http.StatusForbidden)
	}

	secret := AcmeSecret{
		Challenge: r.FormValue("chal"),
		Response:  r.FormValue("resp"),
		Timestamp: time.Now().UTC(),
		UpdatedBy: u.Email,
	}

	key := datastore.NewKey(ctx, entityName, ACME_CHALLENGE_ENTITY_ID, 0, nil)
	_, err := datastore.Put(ctx, key, &secret)
	if err != nil {
		return NewError(err, "datastore.Put failed.", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))

	return nil
}

func ServeAcmeSecret(w http.ResponseWriter, r *http.Request) *Error {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, entityName, ACME_CHALLENGE_ENTITY_ID, 0, nil)
	var secret AcmeSecret
	err := datastore.Get(ctx, key, &secret)
	if err != nil {
		return NewError(err, "datastore.Get failed.", http.StatusInternalServerError)
	}

	requestedChallenge := r.URL.Path[challengeUriPathLength:]
	if requestedChallenge != secret.Challenge {
		return NewError(nil, "Not found", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(secret.Response))

	return nil
}
