package db

//	User 用户数据
type User struct {
	UserName        string
	Password        string
	FollowCount     int64
	FollowerCount   int64
	Avatar          string
	BackgroundImage string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}
