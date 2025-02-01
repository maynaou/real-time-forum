package models

type Liked_Post struct {
	Post_ID string `json:"post_id"`
	User_ID string `json:"user_id"`
}

func createLike(like Liked_Post, post Post, user RegisterRequest){

}

func getLikes() {

}
