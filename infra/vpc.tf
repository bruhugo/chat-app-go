resource "aws_vpc" "go_chat_vpc"{
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
    enable_dns_support = true

    tags = {
        Name: "go_chat_vpc"
    }
}

resource "aws_subnet" "public_1"{
    vpc_id = aws_vpc.go_chat_vpc.id
    cidr_block = "10.0.1.0/24"
    availability_zone = "us-east-1a"
    map_public_ip_on_launch = true

    tags = {
        Name = "public_1"
    }
}

resource "aws_subnet" "public_2"{
    vpc_id = aws_vpc.go_chat_vpc.id
    cidr_block = "10.0.2.0/24"
    availability_zone = "us-east-1b"
    map_public_ip_on_launch = true

    tags = {
        Name = "public_2"
    }
}

resource "aws_internet_gateway" "go_chat_ig"{
    vpc_id = aws_vpc.go_chat_vpc.id

    tags = {
      Name = "go_chat_ig"
    }
}

resource "aws_route_table" "go_chat_rt"{
    vpc_id = aws_vpc.go_chat_vpc.id

    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = aws_internet_gateway.go_chat_ig.id
    }
}

resource "aws_route_table_association" "public_1_association" {
    subnet_id = aws_subnet.public_1.id
    route_table_id = aws_route_table.go_chat_rt.id
}

resource "aws_route_table_association" "public_2_association" {
    subnet_id = aws_subnet.public_2.id
    route_table_id = aws_route_table.go_chat_rt.id
}

### ECS VPC

resource "aws_security_group" "ecs_sg" {
  name        = "ecs_sg"
  vpc_id      = aws_vpc.go_chat_vpc.id

  tags = {
    Name = "ecs_sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4" {
  security_group_id = aws_security_group.ecs_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_vpc_security_group_ingress_rule" "allow_app_traffic" {
  security_group_id = aws_security_group.ecs_sg.id
  cidr_ipv4         = "0.0.0.0/0"  
  ip_protocol       = "tcp"
  from_port         = 8080
  to_port           = 8080
}


### RDS VPC

resource "aws_security_group" "rds_sg" {
  name        = "rds_sg"
  vpc_id      = aws_vpc.go_chat_vpc.id

  tags = {
    Name = "ecs_sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4_rds" {
  security_group_id = aws_security_group.rds_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_vpc_security_group_ingress_rule" "allow_database_connection" {
  security_group_id = aws_security_group.rds_sg.id
  cidr_ipv4 = aws_vpc.go_chat_vpc.cidr_block
  ip_protocol = "tcp"
  from_port = 0
  to_port = 3306
}


# ALB Security Group
resource "aws_security_group" "alb_sg" {
  name        = "alb_sg"
  vpc_id      = aws_vpc.go_chat_vpc.id

  tags = {
    Name = "alb_sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_ingress_rule" "allow_https" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
}

resource "aws_vpc_security_group_egress_rule" "alb_egress" {
  security_group_id = aws_security_group.alb_sg.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}
