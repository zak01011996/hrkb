package ctrls

import (
	"errors"
	"github.com/astaxie/beego"
	M "hrkb/models"
	"hrkb/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Upload struct {
	BaseController
	files string
	dest  string
}

//Upload file with any type
func (c *Upload) Index() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
	}

	c.files = "file"
	c.dest = beego.AppConfig.String("upload_dir") + strconv.Itoa(id) + "/"

	file, err := c.upload()

	if err != nil {
		beego.Error(err)
		return
	}

	file.Cand = id
	file.Active = true
	file.User = c.GetSession("uid").(int)

	if _, err := DM.Insert(&file); err != nil {
		beego.Error(err)

		c.Data["json"] = struct {
			Error bool
			Msg   string
		}{true, T("internal")}
		return
	}

	c.Data["json"] = struct {
		Id   int
		Name string
		Size int64
	}{file.Id, file.Name, file.Size / 1024}
}

func (c *Upload) Img() {

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
	}

	c.files = "file"
	c.dest = beego.AppConfig.String("upload_dir") + strconv.Itoa(id) + "/"

	imgTypes := strings.Split(beego.AppConfig.String("img_types"), ",")

	file, err := c.upload(imgTypes...)

	if err != nil {
		beego.Error(err)
		return
	}

	file.Cand = id
	file.Active = true
	file.User = c.GetSession("uid").(int)

	if _, err := DM.Insert(&file); err != nil {
		beego.Error(err)

		c.Data["json"] = struct {
			Error bool
			Msg   string
		}{true, T("internal")}

		return
	}

	c.Data["json"] = struct {
		Id   int
		Name string
		Size int64
	}{file.Id, file.Name, file.Size / 1024}
}

/**
@param r allowed []string, recieves slice of types or * for  any type

@example c.upload("image/jpeg", "image/jpg", "image/png", "image/gif")
@example c.upload()
*/
func (c *Upload) upload(allowed ...string) (M.File, error) {
	ufile := M.File{}
	tmpDir := beego.AppConfig.String("tmp_dir")

	_, header, err := c.GetFile(c.files)
	if err != nil {
		return ufile, err
	}

	if err := os.Mkdir(tmpDir, os.ModePerm); err != nil {
		beego.Error(err)
	}

	//todo make reader
	if err := c.SaveToFile(c.files, tmpDir+header.Filename); err != nil {
		beego.Error(err)
		return M.File{}, err
	}

	defer os.Remove(tmpDir + header.Filename)

	if err := c.SaveToFile(c.files, tmpDir+header.Filename); err != nil {
		return M.File{}, err
	}

	file, err := os.Open(tmpDir + header.Filename)
	if err != nil {
		return M.File{}, err
	}

	defer file.Close()

	buff := make([]byte, 512) // why 512 bytes according mime spec
	file.Read(buff)

	filetype := http.DetectContentType(buff)
	fileInfo, err := file.Stat()

	ufile.Mime = filetype

	if err != nil {
		return M.File{}, err
	}

	ufile.Size = fileInfo.Size()

	allow := false
	if len(allowed) == 0 {
		allow = true
	}

	for _, at := range allowed {
		if at == filetype {
			allow = true
			break
		}
	}

	if !allow {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = struct {
			Error string `json:"error"`
		}{T("file_invalid_type")}
		return M.File{}, errors.New("fileupload: File type now allowed")
	}

	ftype := strings.Split(header.Filename, ".")
	filepath := c.dest + utils.RandString(15) + "." + ftype[len(ftype)-1]

	if len(header.Filename) > 25 {
		ufile.Name = header.Filename[:20] + "..."
	} else {
		ufile.Name = header.Filename
	}

	if err := os.MkdirAll(c.dest, os.ModePerm); err != nil {
		return M.File{}, err
	}

	if err := os.Rename(tmpDir+header.Filename, filepath); err != nil {
		return M.File{}, err
	}

	ufile.Url = filepath

	return ufile, nil
}
