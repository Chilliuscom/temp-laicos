package cmd

import (
	"SocialNetwork/DB"
	"SocialNetwork/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func displayFeed(w http.ResponseWriter, r *http.Request) {
	log.Println("FEED")
	username := getUserBySession(r)

	var input models.Post
	var query []models.Post
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postID := input.ID
	log.Println(postID, username)
	if result :=
		DB.ConnectDB().
			Raw("SELECT * FROM post p WHERE id < ? ( group_id = ? AND user_id IN (SELECT second_user FROM following WHERE first_user = ? ) AND (privacy = 'PUBLIC' OR privacy ='PRIVATE' OR EXISTS(SELECT 1 FROM privacy pr WHERE pr.post_id = p.id AND pr.user_id = ?))) OR p.group_id IN(SELECT mg.group_id FROM members_of_group mg WHERE mg.user_id = ?) ORDER BY p.id DESC LIMIT 10", "", postID, username, username, username).Scan(&query); result != nil {
		fmt.Println("Bad following")
	}

	data := make(map[string][]models.Post)
	data["post"] = query
	/*

			SELECT
		    p.id,
		    p.text,
		    p.header,
		    p.user_id,
		    p.privacy,
		    p.created_at,
		    p.group_id
		FROM
		    post p
		WHERE
			p.id < **LAST_POST_ID**
		    (
		        p.group_id IS NULL
		        AND p.user_id IN (
		            SELECT
		                f1.second_user
		            FROM
		                friendship f1
		            WHERE
		                f1.first_user = **USERNAME**
		        )
		        AND (
		            p.privacy = 'PUBLIC'
		            OR p.privacy = 'PRIVATE'
		            OR EXISTS (
		                SELECT 1
		                FROM privacy pr
		                WHERE pr.post_id = p.id AND pr.user_id = **USERNAME**
		            )
		        )
		    )
		    OR
		    p.group_id IN (
		        SELECT
		            mg.group_id
		        FROM
		            members_of_group mg
		        WHERE
		            mg.user_id = **USERNAME**
		    )
		ORDER BY p.id DESC
		LIMIT 10;

	*/
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(data)
}
