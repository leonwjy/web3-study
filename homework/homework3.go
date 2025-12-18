package homework

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// task 1
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 使用Gorm定义 User 、 Post 和 Comment 模型，
// 其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
// Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。

type User struct {
	gorm.Model
	Username     string  `gorm:"not null;comment:'username'"`
	Nickname     string  `gorm:"not null;comment:'nickname'"`
	Phone        string  `gorm:"not null;unique;comment:'phone'"`
	Email        *string `gorm:"comment:'email'"`
	Salt         string  `gorm:"not null;comment:'salt'"`
	Password     string  `gorm:"not null;comment:'password'"`
	Posts        []Post  `gorm:"foreignKey:UserID;references:ID"`
	PostCount    int64   `gorm:"not null;default:0;comment:'post_count'"`
	CommentCount int64   `gorm:"not null;default:0;comment:'comment_count'"`
}

type Post struct {
	gorm.Model
	UserID        uint      `gorm:"not null;comment:user_id;index:idx_user_id"`
	User          User      `gorm:"foreignKey:UserID;references:ID"`
	Title         string    `gorm:"not null;comment:title"`
	Content       string    `gorm:"comment:content"`
	CommentCount  int       `gorm:"not null;default:0;comment:'comment_count'"`
	ForwardCount  int       `gorm:"not null;default:0;comment:'forward_count'"`
	LikeCount     int       `gorm:"not null;default:0;comment:'like_count'"`
	Comments      []Comment `gorm:"foreignKey:PostID;references:ID"`
	CommentStatus int       `gorm:"not null;default:'0';comment:'comment_status: 0-无评论 1-有评论'"`
}

type Comment struct {
	gorm.Model
	PostID       uint   `gorm:"not null;comment:post_id;index:idx_post_id"`
	Post         Post   `gorm:"foreignKey:PostID;references:ID"`
	UserID       uint   `gorm:"not null;comment:user_id;index:idx_user_id"`
	User         User   `gorm:"foreignKey:UserID;references:ID"`
	Content      string `gorm:"comment:content"`
	CommentCount int    `gorm:"not null;default:0;comment:'comment_count'"`
	ForwardCount int    `gorm:"not null;default:0;comment:'forward_count'"`
	LikeCount    int    `gorm:"not null;default:0;comment:'like_count'"`
}

// Task 2
// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func SearchPostAndCommented(db *gorm.DB, userId uint) []Post {
	var posts []Post
	db.Model(&Post{}).Where("user_id = ?", userId).Preload("Comments").Find(&posts)

	for _, post := range posts {
		fmt.Printf("Post ID: %d, Title: %s, CreatedAt: %s\n",
			post.ID, post.Title, post.CreatedAt.Format("yyyy-mm-dd hh:mm:ss"))
		for _, comment := range post.Comments {
			fmt.Printf("  Comment ID: %d, Content: %s, CreatedAt: %s\n",
				comment.ID, comment.Content, comment.CreatedAt.Format("yyyy-mm-dd hh:mm:ss"))
		}
		fmt.Println()
	}

	return posts
}

// 编写Go代码，使用Gorm查询评论数量最多的文章信息
func SearchMostCommentedPost(db *gorm.DB) []Post {
	var posts []Post
	db.Model(&Post{}).Order("comment_count DESC").First(&posts)
	return posts
}

// Task 3
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	count := int64(0)
	// 查询该用户的文章数量
	tx.Model(&Post{}).Where("user_id = ?", post.UserID).Count(&count)
	// 更新用户的文章数量统计字段
	tx.Model(&User{}).Where("id = ?", post.UserID).Update("post_count", count)
	return
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (comment *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	count := int64(0)
	// 查询该文章的评论数量
	tx.Model(&Comment{}).Where("post_id = ?", comment.PostID).Count(&count)
	// 更新文章的评论状态
	tx.Model(&Post{}).Where("id = ?", comment.PostID).Update("comment_status", count == 0)
	return
}

func Run1() {
	db, err := gorm.Open(mysql.Open("admin:123456@tcp(localhost:3306)/web3_study?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// user1 := User{
	// 	Username: "user1",
	// 	Nickname: "user1",
	// 	Phone:    "15000000001",
	// 	Salt:     "abc",
	// 	Password: "123456",
	// }
	// user2 := User{
	// 	Username: "user2",
	// 	Nickname: "user2",
	// 	Phone:    "15000000002",
	// 	Salt:     "abc",
	// 	Password: "123456",
	// }
	// result1 := db.Create(&[]User{user1, user2})
	// fmt.Println(result1.RowsAffected)

	// post1 := Post{
	// 	UserID:  1,
	// 	Title:   "post1",
	// 	Content: "post1",
	// }

	// post2 := Post{
	// 	UserID:  1,
	// 	Title:   "post2",
	// 	Content: "post2",
	// }

	// post3 := Post{
	// 	UserID:  2,
	// 	Title:   "post3",
	// 	Content: "post3",
	// }

	// result2 := db.Create([]Post{post1, post2, post3})
	// fmt.Println(result2.RowsAffected)

	// comment1 := Comment{
	// 	PostID:  1,
	// 	UserID:  1,
	// 	Content: "Comment1",
	// }

	// comment2 := Comment{
	// 	PostID:  1,
	// 	UserID:  2,
	// 	Content: "Comment2",
	// }

	// comment3 := Comment{
	// 	PostID:  2,
	// 	UserID:  1,
	// 	Content: "Comment3",
	// }

	// comment4 := Comment{
	// 	PostID:  3,
	// 	UserID:  2,
	// 	Content: "Comment4",
	// }

	// result3 := db.Create([]Comment{comment1, comment2, comment3, comment4})

	// fmt.Println(result3.RowsAffected)

	// result := SearchPostAndCommented(db, 1)
	// fmt.Println(result)

	// result2 := SearchMostCommentedPost(db)
	// fmt.Println(result2)
}
