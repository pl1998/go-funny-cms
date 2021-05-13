package controllers

import (
	"github.com/jordan-wright/email"
	"github.com/spf13/cast"
	"gocms/app/http/admin/validates"
	"gocms/app/validates/validate"
	"gocms/pkg/config"
	"gocms/pkg/help"
	"gocms/pkg/mail"
	"gocms/pkg/response"
	"gocms/wrap"
	"net/textproto"
)

type MailController struct{}

func (m *MailController) List(c *wrap.ContextWrapper) {
	page := c.DefaultQuery("page", 1)
	pageSize := c.DefaultQuery("pageSize", 10)
	keyword := c.DefaultQuery("keyword", "")
	total := 0

	list := []mail.MailerModel{}
	query := config.Db.Model(&mail.MailerModel{})

	if len(keyword) > 0 {
		query = query.Where("name like ?", "%"+keyword+"%")
	}

	query.Limit(pageSize).Offset(page).Scan(&list)
	query.Count(&total)

	response.SuccessResponse(list).WriteTo(c)
	return
}

// Store 处理邮件发送
func (m *MailController) Store(c *wrap.ContextWrapper) {
	var params validates.EmailValidate
	_ = c.ShouldBind(&params)

	if !validate.WithResponseMsg(params, c) {
		return
	}

	for _, to := range params.Emails {
		express := mail.NewMailerExpress()
		express.Mailer.Mail = &email.Email{
			To:      []string{to},
			Subject: params.Subject,
			HTML:    []byte(params.Content),
			Headers: textproto.MIMEHeader{},
		}
		express.Attachments = params.Attachments
		SendAt := help.ParseTime(params.SendAt)
		express.Options.SendAt = SendAt

		task := mail.NewTaskExpress()
		err := task.Dispatch(express)
		if err != nil {
			response.ErrorResponse(501, err.Error()).WriteTo(c)
			return
		}
	}

	response.SuccessResponse().WriteTo(c)
	return
}

func (m *MailController) Mailer(c *wrap.ContextWrapper) {
	mailers := config.GetMailerLabels()

	response.SuccessResponse(mailers).WriteTo(c)
	return
}

// 重新发送
func (m *MailController) Resend(c *wrap.ContextWrapper) {
	id := c.DefaultQuery("id", 0)

	mailModel := mail.MailerModel{}
	config.Db.Model(mailModel).Find(&mailModel, id)

	task := mail.NewTaskExpress()
	express := task.ForParse(&mailModel)

	express.Send(express, cast.ToInt(id))

	response.SuccessResponse().WriteTo(c)
	return
}

// 重新发送
func (m *MailController) Delete(c *wrap.ContextWrapper) {
	id := c.DefaultQuery("id", 0)

	mailModel := mail.MailerModel{}
	config.Db.Model(mailModel).Where("id = ?", id).Update(map[string]int{
		"status": mail.TASK_DELETE,
	})

	response.SuccessResponse().WriteTo(c)
	return
}
