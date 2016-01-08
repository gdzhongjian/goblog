package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/nfnt/resize"
	"goblog/models"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) Post() {

}

//保存
func (this *ArticleController) Save() {
	var (
		id      int    = 0
		title   string = strings.TrimSpace(this.GetString("title"))
		content string = this.GetString("content")
		tags    string = strings.TrimSpace(this.GetString("tags"))
		color   string = strings.TrimSpace(this.GetString("color"))
		picture string = strings.TrimSpace(this.GetString("articleimage"))
		brief   string = strings.TrimSpace(this.GetString("brief"))
		// timestr int    = time.Now().Unix()
		status1 int = 0
		istop1  int = 0
		// article models.Article
	)
	if title == "" {
		//标题不能为空
		return
	}
	if brief == "" {
		return
	}

	id1 := this.GetString("id")
	id2, _ := strconv.ParseInt(id1, 10, 64)
	id = int(id2)
	status1, _ = this.GetInt("status")
	status := int8(status1)
	istop1, _ = this.GetInt("istop")
	istop := int8(istop1)
	if istop1 == 1 {
		istop = 1
	} else {
		istop = 0
	}

	if status != 1 && status != 2 {
		status = 0
	}

	addtags := make([]string, 0)
	//标签过滤
	if tags != "" {
		tagarr := strings.Split(tags, ",")
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tag != "" {
				exists := false
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}
				if !exists {
					addtags = append(addtags, tag)
				}
			}
		}
	}

	if id < 1 {
		//id 小于1为添加文章操作
		uid := this.GetSession("sessionuid")
		userinfo, err := models.Userinfowithuid(uid)
		if err != nil {
			return
		}

		//插入操作
		author := userinfo.Username
		var articleid int
		articleid, err = models.Articleadd(author, title, content, tags, color, picture, brief, status, istop)
		if err != nil {
			return
		}

		/*
			处理上传文章封面图片,保存在static/articleimage/用户名/文章ID/封面图片
		*/
		_, f, _ := this.GetFile("articleimage")
		if f != nil {
			filename := f.Filename
			houzhui := strings.LastIndex(filename, ".")
			//后缀名
			realhouzhui := filename[houzhui:]
			userid := this.GetSession("sessionuid")
			userinfo, err = models.Userinfowithuid(userid)

			if err != nil {
				return
			}
			username := userinfo.Username
			articleid1 := strconv.Itoa(articleid)
			//保存路径
			path := "static/articleimage/" + username + "/" + articleid1 + "/"
			err = os.MkdirAll(path, os.ModePerm)
			err = this.SaveToFile("articleimage", path+"original"+realhouzhui)
			if err != nil {
				return
			}
			// this.Redirect("/admin/articlelist", 301)
			//图片保存成功后对图片生成缩略图，传递图片路径,图片名称和缩略图大小
			name := "original" + realhouzhui
			var size1 uint = 200
			var size2 uint = 123
			err = this.Thumbnailarticle(path, name, size1, size2)
			if err != nil {
				return
			}
			// this.Redirect("/admin/articlelist", 301)

			//获取保存后路径
			articleimage := path + "articleimage" + realhouzhui
			//更新数据库信息
			err = models.Articleimageupdate(articleid, articleimage)
			if err != nil {
				return
			}
		}

		this.Redirect("/admin/articlelist", 301)
	} else {
		//id大于等于1时为编辑文章操作
		// uid := this.GetSession("sessionuid")
		// userinfo, err := models.Userinfowithuid(uid)
		// if err != nil {
		// 	return
		// }

		//表单传递的时间格式转换为时间戳

		timevalue := this.GetString("posttime")
		timestampvalue := int(Datetotimestamp(timevalue))

		//文章编辑操作
		err := models.Articleedit(title, color, tags, content, brief, status, istop, timestampvalue, id)
		if err != nil {
			return
		}

		/*
			处理上传文章封面图片,保存在static/articleimage/用户名/文章ID/封面图片
		*/
		//如果编辑页面重新上传了新的封面，就进行封面更新，否则不处理

		_, f, _ := this.GetFile("articleimage")
		if f != nil {
			// this.Redirect("/admin", 301)
			filename := f.Filename
			houzhui := strings.LastIndex(filename, ".")
			//后缀名
			realhouzhui := filename[houzhui:]
			userid := this.GetSession("sessionuid")
			userinfo, err := models.Userinfowithuid(userid)

			if err != nil {
				return
			}
			username := userinfo.Username
			//保存路径
			articleid := strconv.Itoa(id)
			path := "static/articleimage/" + username + "/" + articleid + "/"
			err = os.MkdirAll(path, os.ModePerm)
			err = this.SaveToFile("articleimage", path+"original"+realhouzhui)
			if err != nil {
				return
			}
			// this.Redirect("/admin/articlelist", 301)
			//图片保存成功后对图片生成缩略图，传递图片路径,图片名称和缩略图大小
			name := "original" + realhouzhui
			var size1 uint = 200
			var size2 uint = 123
			err = this.Thumbnailarticle(path, name, size1, size2)
			if err != nil {
				return
			}
			// this.Redirect("/admin/articlelist", 301)

			//获取保存后路径
			articleimage := path + "articleimage" + realhouzhui
			//更新数据库信息
			err = models.Articleimageupdate(id, articleimage)
			if err != nil {
				return
			}
		}

		this.Redirect("/admin/articlelist", 301)
	}
	// uid := this.GetSession("uid")
	// userinfo, _ := models.Userinfowithuid(uid)
	// err := models.Articleadd(title, color, content)
	// if err != nil {
	// 	return
	// }
	this.Redirect("/admin/userinfo", 301)
}

//上传文件
func (this *ArticleController) Upload() {
	_, header, err := this.GetFile("upfile")
	ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, "."):])
	out := make(map[string]string)
	out["url"] = ""
	out["fileType"] = ext
	out["original"] = header.Filename
	out["state"] = "SUCCESS"
	if err != nil {
		out["state"] = err.Error()
	} else {
		savepath := "./static/upload/" + time.Now().Format("20060102")
		if err := os.MkdirAll(savepath, os.ModePerm); err != nil {
			out["state"] = err.Error()
		} else {
			filename := fmt.Sprintf("%s/%d%s", savepath, time.Now().UnixNano(), ext)
			if err := this.SaveToFile("upfile", filename); err != nil {
				out["state"] = err.Error()
			} else {
				out["url"] = filename[1:]
			}
		}
	}
	this.Data["json"] = out
	this.ServeJson()
}

//删除文章
func (this *ArticleController) Delete() {
	//获取文章id
	articleid1 := this.Input().Get("articleid")
	status := this.Input().Get("status")
	articleid2, err := strconv.ParseInt(articleid1, 10, 32)
	articleid := int(articleid2)
	// articleid := int(articleid1)
	err = models.Articledelete(articleid)
	if err != nil {
		return
	}
	this.Redirect("/admin/articlelist?status="+status, 301)
}

//批处理文章
func (this *ArticleController) Batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}

	switch op {
	case "topub":
		//移到已发布
		if len(idarr) == 0 {
			break
		}
		err := models.Articlebatch(idarr, 0)
		if err != nil {
			return
		}
	case "todrafts":
		//移到草稿箱
		if len(idarr) == 0 {
			break
		}
		err := models.Articlebatch(idarr, 1)
		if err != nil {
			return
		}
	case "totrash":
		//移到回收站
		if len(idarr) == 0 {
			break
		}
		err := models.Articlebatch(idarr, 2)
		if err != nil {
			return
		}
	case "delete":
		//批量删除
		if len(idarr) == 0 {
			break
		}
		err := models.Articlebatchdelete(idarr)
		if err != nil {
			return
		}
	default:
		if len(idarr) == 0 {
			break
		}
	}
	this.Redirect(this.Ctx.Request.Referer(), 302)
}

// //显示文章列表
// func (this *ArticleController) List() {
// 	//获取用户uid
// 	uid := this.GetString("sessionuid")
// 	articlelist, err := models.Articleread(uid)
// 	if err != nil {
// 		return
// 	}

// 	this.TplNames = "article/list.html"
// 	this.Data["article"] = articlelist
// }

/*
日期字符串转换成时间戳
*/
func Datetotimestamp(datetime string) (timestampvalue int64) {
	toBeCharge := datetime
	timeLayOut := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local") //获取时区
	theTime, _ := time.ParseInLocation(timeLayOut, toBeCharge, loc)
	timestampvalue = theTime.Unix()
	return timestampvalue
}

/*
文章封面图片处理
*/

func (this *ArticleController) Thumbnailarticle(path, name string, size1, size2 uint) error {
	//获取图片后缀和图片完整地址
	houzhui := name[strings.LastIndex(name, "."):]
	realpath := path + name
	// open "test.jpg"
	file, err := os.Open(realpath)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(size1, size2, img, resize.Lanczos3)

	out, err := os.Create(path + "articleimage" + houzhui)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, m, nil)
	if err != nil {
		return err
	}
	return nil
}
