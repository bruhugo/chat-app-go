resource "aws_ecs_cluster" "go_chat_cluster" {
  name = "go_chat_cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecr_repository" "chat_app_ecr_repo" {
  name = "go_chat_repo"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_iam_role" "ecs_task_role" {
  name = "gochat-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy" "read_db_secret" {
  name = "gochat-read-db-secret"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = data.aws_secretsmanager_secret.credentials.arn
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "task_read_secret_attach" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.read_db_secret.arn
}


resource "aws_iam_role" "ecs_execution_role" {
  name = "gochat-ecs-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_exec_attach" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "ecs_secrets_policy" {
  name = "gochat-ecs-secrets-policy"
  role = aws_iam_role.ecs_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = "arn:aws:secretsmanager:us-east-1:${data.aws_caller_identity.current.account_id}:secret:go_chat_secrets*"
    }]
  })
}

# Needed to dynamically resolve your account ID in the ARN above
data "aws_caller_identity" "current" {}

resource "aws_cloudwatch_log_group" "gochat" {
  name              = "/ecs/gochat"
  retention_in_days = 7
}

resource "aws_ecs_task_definition" "go_chat_td" {
  family                   = "go_chat_family"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  memory                   = 1024
  cpu                      = 512
  container_definitions = jsonencode([
    {
      name      = "go_chat_container"
      image     = "${aws_ecr_repository.chat_app_ecr_repo.repository_url}:latest"
      memory    = 1024
      essential = true
      portMappings = [
        {
          containerPort = 8080,
          protocol      = "tcp"
        }
      ],
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.gochat.name
          awslogs-region        = "us-east-1"
          awslogs-stream-prefix = "ecs"
        }
      },
      environment = [
        { name = "db_host", value = aws_db_instance.default.address },
        { name = "db_port", value = var.db_port },
        { name = "db_database", value = var.db_name },
        { name = "port", value = var.port },
        { name = "frontend_host", value = var.frontend_host },
        { name = "SECRETS", value = "AWS" }
      ]
      secrets = [
        {
          name      = "DB_SECRETS",
          valueFrom = data.aws_secretsmanager_secret.credentials.arn
        }
      ]
    },
  ])
}

resource "aws_ecs_service" "go_chat_service" {
  name            = "go_chat"
  cluster         = aws_ecs_cluster.go_chat_cluster.id
  task_definition = aws_ecs_task_definition.go_chat_td.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.ecs_sg.id]
    subnets          = [aws_subnet.public_1.id, aws_subnet.public_2.id]
    assign_public_ip = true
  }
  load_balancer {
    target_group_arn = aws_lb_target_group.app_tg.arn
    container_name   = "go_chat_container"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.http_listener]
}

