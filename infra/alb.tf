# ALB
resource "aws_lb" "app_alb" {
  name               = "go-chat-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public_1.id, aws_subnet.public_2.id]  # Public subnets

  tags = {
    Name = "go_chat_alb"
  }
}

# Target Group
resource "aws_lb_target_group" "app_tg" {
  name        = "go-chat-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.go_chat_vpc.id
  target_type = "ip" 

  health_check {
    path                = "/health"  
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  stickiness {
    type            = "lb_cookie"
    cookie_duration = 86400  # 1 day
  }
}

# Listener
resource "aws_lb_listener" "http_listener" {
  load_balancer_arn = aws_lb.app_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app_tg.arn
  }
}