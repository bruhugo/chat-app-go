variable "db_host" {
    type        = string                         
    description = "Database's host"  
    default     = "t2.micro"                   
}

variable "db_port" {
    type        = string                         
    description = "Database's port"  
    default     = "3306"                   
}

variable "port" {
    type = string
    description = "Application's port"
    default = "8080"
}

variable "db_name" {
    type = string
    description = "Database's name"
    default = "go_chat_db"
}