package models

const (
	HAITON  = 1
	ARTICLE = 2
	PLC     = 3
)

// 文章
type Article struct {
	BaseModel
	Title       string          `json:"title" orm:"size(128)" doc:"标题"`
	Description string          `json:"description" orm:"size(256)" doc:"文章简介"`
	Class       *Classification `json:"tag_num" orm:"rel(fk);on_delete(set_null);null" doc:"所属分类,外键"`
	Comments    []*Comment      `json:"-" orm:"reverse(many)" doc:"其下的评论"`
	Author      *Account        `json:"author" orm:"rel(fk);on_delete(set_null);null"`
	Release     bool            `json:"release" orm:"default(true)" doc:"是否发布"`

	Html string `json:"html" orm:"type(text)" doc:"文章的HTML代码"`
	Code string `json:"code" orm:"type(text)" doc:"源码 markdown/txt"`

	Visit int64 `json:"visit" orm:"default(0)" doc:"访问次数"`

	ClsId int64 `json:"cls_id" orm:"-"`
}

// 文章分类
type Classification struct {
	BaseModel
	Name    string     `json:"name" orm:"size(16);unique" doc:"分类名"`
	Cover   string     `json:"cover" orm:"size(256);null" doc:"封面图"`
	ArtNums int32      `json:"art_nums" orm:"default(0)" doc:"其下文章数量"`
	BigCls  int8       `json:"big_cls" doc:"所属大类，1,2,3"`
	Arts    []*Article `json:"arts" orm:"reverse(many)" doc:"关联的文章"`
}

// 评论
type Comment struct {
	BaseModel
	Nick      string   `json:"nick" orm:"size(32)" doc:"评论人昵称"`
	Content   string   `json:"content" orm:"size(512)" doc:"评论内容"`
	ArtObject *Article `json:"-" orm:"rel(fk);on_delete(set_null);null"`
	ArtId     int64    `json:"art_id" orm:"-"`
}
