# Accessing an AWS account using AssumeRoleWithWebIdentity from GCP
## 1. Create IAM Role on AWS

create iam role with following trust policy
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "accounts.google.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "accounts.google.com:sub": "service account unique ID"
                }
            }
        }
    ]
}
```

## 2. Get JWT Token from metadata server

```
curl http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token -H "Metadata-Flavor: Google"
```

## 3. Access to s3 using AssumeRoleWithWebIdentity
