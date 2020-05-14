package service

import "cloudClass/models"

type UserService struct {
	userModel models.UserModel
}
func (this *UserService)ServiceInit()  {
	this.userModel = models.UserModel{}
}
