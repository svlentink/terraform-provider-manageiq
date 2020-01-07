provider "manageiq" {
  hostname = "cfme.example.com"
}

variable "resource_params" {
  default = {
    examplekey = "example value"
    tag_0_backup = "yes"
    tag_0_department = "IT services"
    subnet_name = "prod_dmz_EU_01"
    vm_memory = 2048
  }
}

resource "manageiq_vm" "my-vm" {
  tags = var.resource_params
  
  provisioner "local-exec" {
    command = "echo ${manageiq_vm.my-vm[count.index].name}"
  }
}
