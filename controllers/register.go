package controllers

// RegisterController handles the welcome screen that allows user to pick a technology and username.
type RegisterController struct {
	baseController // Embed to use methods that are implemented in baseController.
}

// Get implemented Get() method for RegisterController.
func (this *RegisterController) Get() {
	this.Data["IsRegister"] = true
	this.TplName = "register.html"
}

// Join method handles POST requests for RegisterController.
func (this *RegisterController) Reg() {

	this.TplName = "main.html"
}
