package service

import "cloudClass/models"

type AdminService struct {
	adminModel models.AdminModel
}
func (this *AdminService)ServiceInit()  {
	this.adminModel = models.AdminModel{}
}
