provider "aws" {
  region = "us-east-1"
}

resource "aws_security_group" "appsg"{
  name = "appsg"

  ingress{
    from_port=3000
    to_port=3000
    protocol="tcp"
    cidr_blocks=["0.0.0.0/0"]
  }

  ingress{
    from_port=22
    to_port=22
    protocol="tcp"
    cidr_blocks=["0.0.0.0/0"]
  }

  ingress{
    from_port=-1
    to_port=-1
    protocol="icmp"
    cidr_blocks=["0.0.0.0/0"]
  }
}

resource "aws_security_group" "dbsg"{
  name = "dbsg"

  ingress {
    from_port = 3306
    to_port = 3306
    protocol = "tcp"
    security_groups  = ["${aws_security_group.appsg.name}"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

}


resource "aws_db_instance" "default" {
  allocated_storage = "5"
  storage_type = "gp2"
  engine = "mysql"
  engine_version = "5.6.17"
  instance_class = "db.t1.micro"
  name = "mydb"
  username = "mlabouardy"
  password = "12345678"
  multi_az = false
  security_group_names = ["${aws_security_group.dbsg.name}"]
}

resource "aws_key_pair" "default"{
  key_name = "appgokp"
  public_key = "${file("/home/core/.ssh/id_rsa.pub")}"

}

resource "aws_instance" "default"{
  ami="ami-4fffc834"
  instance_type="t2.micro"
  key_name="${aws_key_pair.default.id}"
  security_groups = ["${aws_security_group.appsg.name}"]
  user_data = "${file("bootstrap.sh")}"


  tags {
    Name="app"
  }
}
