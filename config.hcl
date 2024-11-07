account_id = "${aws_account_id}"
region     = "${aws_region}"
principal  = "${aws_principal}"
email      = split(":", aws_principal)[1]
message    = upper(from_env("message"))
