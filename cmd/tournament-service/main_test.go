package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var conn *pgxpool.Pool
var r *gin.Engine
var ctx context.Context = context.Background()

func TestPreFn(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	databaseUrl := "postgresql://postgres:postgres@localhost:5433/tournament" // lets assume you are using docker database
	var err error
	conn, err = connectDB(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Error on postgres database: %v\n", err)
	}

	r = gin.New()
	r.Use(gin.Recovery()) // recovery middleware
	setupRouter(r, false)
}

func TestCleanup(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeams(t *testing.T) {
	result := make(map[string]interface{})

	// new team
	bodyData := map[string]interface{}{
		"name": "AC Milan",
	}

	jsonValue, _ := json.Marshal(bodyData)
	req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code) // ok

	// new team with same name
	{
		jsonValue, _ := json.Marshal(bodyData)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code) // name is unique
	}

	// new another team
	{
		bodyData := map[string]interface{}{
			"name": "Juventus",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code) // ok
	}

	// team generation on non-empty teams table
	{
		req, _ := http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, allowed only into empty table

		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "generation_only_allowed_into_empty_table", result["err"].(string))
	}

	// team generation on empty teams table
	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok, allowed into empty table
	}

	// new another team on teams table having 16 records(max)
	{
		bodyData := map[string]interface{}{
			"name": "Nottingham Forest",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, max 16

		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "max_16_teams_allowed", result["err"].(string))
	}
}

func TestDivisions(t *testing.T) {
	result := make(map[string]interface{})

	req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, no teams
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, "must_have_16_teams_to_prepare_divisions", result["err"].(string))

	{
		bodyData := map[string]interface{}{
			"name": "Liverpool",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ := http.NewRequest("POST", "/api/teams", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, not enough teams, need 16
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "must_have_16_teams_to_prepare_divisions", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, not enough teams and not prepared
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "must_have_16_teams_to_start_divisions", result["err"].(string))
	}

	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok, we have 16 teams, divided teams into two divisions and generated matches for each division
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, already prepared divisions
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_already_prepared", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok, we have prepared division, and then started(generated matches scores and we have now top 4 on both divisions)
		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, already started division, we have already all matches played and have top 4
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_already_started", result["err"].(string))
	}

	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, not have prepared division, we only have 16 teams, no divisions, no matches
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_not_prepared", result["err"].(string))
	}
}

func TestPlayoffs(t *testing.T) {
	// playoff quarter final prepare - gets top 4 from both divisions, randomly picks matches as best team plays with worst.
	// playoff quarter final start - generates matches scores
	result := make(map[string]interface{})
	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// try playoff prepare on empty db
		bodyData := map[string]interface{}{
			"stage": "quarter",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have division fully ended
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_not_started", result["err"].(string))

		// try playoff prepare on 16 teams, no division
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have division fully ended, having 16 teams is not playoff condition, it is division condition
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_not_started", result["err"].(string))

		// try playoff prepare on 16 teams, division prepared
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have division fully ended, division is prepared(have teams and matches but none played)
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "division_is_not_started", result["err"].(string))

		// try playoff start on playoff non prepared
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, to generates matches scores you have to have teams and matches for quarter final stage
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_quarter_is_not_prepared", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok
	}

	// playoff semi prepare - picks winners of quarter stage and randomly prepares matches.
	// playoff semi start - generates matches scores
	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// now we have division complete but non quarter

		// try semi without quarter
		bodyData := map[string]interface{}{
			"stage": "semi",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have quarter stage complete
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_quarter_is_not_started", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have semi prepared stage
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_semi_is_not_prepared", result["err"].(string))

		// complete quarter
		bodyData = map[string]interface{}{
			"stage": "quarter",
		}
		jsonValue, _ = json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// try again
		bodyData = map[string]interface{}{
			"stage": "semi",
		}
		jsonValue, _ = json.Marshal(bodyData)
		// try start semi without prepare
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have semi prepared stage, we have complete quarter, semi start looks at semi prepared, semi prepared looks at quarter complete
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_semi_is_not_prepared", result["err"].(string))

		// try again
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok
	}

	// playoff final prepare - picks winners of semi stage and prepares final match.
	// playoff final start - generates match scores
	{
		req, _ := http.NewRequest("POST", "/api/cleanup", bytes.NewBuffer([]byte{}))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/teams/generate", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/prepare", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/divisions/start", bytes.NewBuffer([]byte{}))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// now we have division complete but non quarter

		// try final without quarter
		bodyData := map[string]interface{}{
			"stage": "final",
		}
		jsonValue, _ := json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have quarter stage complete
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_quarter_is_not_started", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have final prepared stage
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_final_is_not_prepared", result["err"].(string))

		// complete quarter
		bodyData = map[string]interface{}{
			"stage": "quarter",
		}
		jsonValue, _ = json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// try again
		bodyData = map[string]interface{}{
			"stage": "final",
		}
		jsonValue, _ = json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, semi is not started
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_semi_is_not_started", result["err"].(string))
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, final is not prepared
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_final_is_not_prepared", result["err"].(string))

		// complete semi
		bodyData = map[string]interface{}{
			"stage": "semi",
		}
		jsonValue, _ = json.Marshal(bodyData)
		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// try again
		bodyData = map[string]interface{}{
			"stage": "final",
		}
		jsonValue, _ = json.Marshal(bodyData)
		// try start final without prepare
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code) // not allowed, must have final prepared stage, we have complete quarter, complete semi, non prepared final
		json.Unmarshal([]byte(w.Body.String()), &result)
		assert.Equal(t, "playoff_final_is_not_prepared", result["err"].(string))

		req, _ = http.NewRequest("POST", "/api/playoff/prepare", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok
		req, _ = http.NewRequest("POST", "/api/playoff/start", bytes.NewBuffer(jsonValue))
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) // ok, we have a winner here
	}
}

func TestPostFn(t *testing.T) {
	defer conn.Close()
}
