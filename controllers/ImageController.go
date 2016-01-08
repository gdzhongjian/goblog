package controllers

import (
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

type ImageController struct {
	beego.Controller
}

func (this *ImageController) Get() {

}

/*
上传照片页面
*/
func (this *ImageController) Image() {
	this.TplNames = "upload/uploadimage.html"
}

/*
上传头像处理
*/
func (this *ImageController) Postimage() {
	_, f, err := this.GetFile("uploadimage")
	if err != nil {
		this.Data["shangchuan"] = "请上传文件！"
		this.Data["errmsg"] = true
		this.TplNames = "upload/image.html"
		return
	}
	if f == nil {
		this.Data["shangchuan"] = "请上传文件！"
		this.Data["errmsg"] = true
		this.TplNames = "upload/image.html"
		return
	}
	filename := f.Filename
	houzhui := strings.LastIndex(filename, ".")
	//后缀名
	realhouzhui := filename[houzhui:]
	// this.Data["test"] = realhouzhui
	// this.TplNames = "article/testarticle.html"
	// return
	curtime := time.Now().Unix()
	formatcurtime := this.Timetoformat(curtime)
	//获取当前用户名
	userid := this.GetSession("sessionuid")
	if userid == nil {
		//session不存在
		return
	}
	var userinfo *models.User
	userinfo, err = models.Userinfowithuid(userid)
	if err != nil {
		return
	}
	username := userinfo.Username
	//头像保存路径
	path := "static/uploadimage/" + username + "/" + formatcurtime + "/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return
	}
	err = this.SaveToFile("uploadimage", path+"original"+realhouzhui)
	if err != nil {
		return
	}
	//获取原图片的宽度和高度
	name1 := "original" + realhouzhui
	width, height := this.Getimageheight(path, name1)

	//根据获取的宽度和高度压缩比例,宽度默认压缩为200
	height1 := 200 * height / width
	height2 := uint(height1)
	//图片保存成功后对图片生成缩略图，传递图片路径,图片名称，和缩略图大小和缩略图类型
	name := "original" + realhouzhui
	err = this.Thumbnail1(path, name, height2, curtime)
	if err != nil {
		return
	}
	// this.Redirect("/admin", 301)
	//移除原图

	//获取上传后图片路径
	curtime32 := int(curtime)
	curtimestring := strconv.Itoa(curtime32)
	imagepath := path + curtimestring + realhouzhui

	//获取上传后图片高度
	name2 := curtimestring + realhouzhui
	_, height3 := this.Getimageheight(path, name2)

	//插入数据到图片表
	userid1 := userid.(int)
	err = models.Insertpicture(userid1, height3, imagepath)
	if err != nil {
		return
	}
	this.Data["success"] = true
	this.TplNames = "upload/uploadimage.html"
}

/*
后台图片处理
*/
func (this *ImageController) Pictureedit() {
	/*
		相册显示
	*/
	var status int
	status1 := this.GetString("status")
	if status1 == "" {
		status = 0
	} else {
		status, _ = strconv.Atoi(status1)
	}
	if status == 0 {
		this.Data["yifabu"] = true
	} else {
		this.Data["yifabu"] = false
	}
	userid := this.GetSession("sessionuid")
	userid1 := userid.(int)
	pictureinfo, err := models.Pictureread(userid1, status)
	//获取已发布和回收站数量
	count1, err := models.Imagereadcount(userid, 0)
	count2, err := models.Imagereadcount(userid, 1)
	if err != nil {
		return
	}
	this.Data["count1"] = count1
	this.Data["count2"] = count2
	this.Data["pictureinfo"] = pictureinfo
	this.TplNames = "message/image.html"
}

/*
删除图片
*/
func (this *ImageController) Delete() {
	pictureid, _ := this.GetInt("pictureid")
	//查找图片路径
	picture, err := models.Pictureselect(pictureid)
	if err != nil {
		return
	}
	picturelujing := picture.Picture
	err = models.Picturedelete(pictureid)
	if err != nil {
		return
	}
	//删除本地图片
	err = os.Remove(picturelujing)
	if err != nil {
		return
	}
	this.Redirect("/image/pictureedit", 301)

}

//批处理照片
func (this *ImageController) Batch() {
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
		err := models.Imagebatch(idarr, 0)
		if err != nil {
			return
		}

	case "totrash":
		//移到回收站
		if len(idarr) == 0 {
			break
		}

		err := models.Imagebatch(idarr, 1)
		if err != nil {
			return
		}
	case "delete":
		//批量删除
		if len(idarr) == 0 {
			break
		}
		picture, err := models.Imagebatchselect(idarr)
		err = models.Imagebatchdelete(idarr)
		for _, v := range picture {
			picturelujing := v.Picture
			os.Remove(picturelujing)
		}
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

/*
时间戳转换
*/
func (this *ImageController) Timetoformat(createtime int64) string {
	out := time.Unix(createtime, 0).Format("20060102")
	return out
}

/*
生成图片缩略图函数
*/
func (this *ImageController) Thumbnail1(path, name string, size uint, curtime int64) error {
	//获取图片后缀和图片完整地址
	houzhui := name[strings.LastIndex(name, "."):]
	realpath := path + name
	// 打开图片
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
	m := resize.Resize(200, size, img, resize.Lanczos3)
	curtime32 := int(curtime)
	curtimestring := strconv.Itoa(curtime32)
	out, err := os.Create(path + curtimestring + houzhui)
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

/*
获取图片高度
*/

func (this *ImageController) Getimageheight(path, name string) (width, height int) {
	//获取图片完整地址
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
	i := img.Bounds()
	width = i.Dx()
	height = i.Dy()
	return width, height
}
