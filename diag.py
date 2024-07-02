from diagrams import Diagram
from diagrams.aws.compute import Lambda
from diagrams.aws.database import Aurora
from diagrams.aws.storage import S3
from diagrams.aws.integration import Eventbridge,SQS
from diagrams.aws.compute import EKS
from diagrams.aws.network import APIGateway
from diagrams.aws.network import ELB
from diagrams.aws.general import Users

with Diagram("Golang Project API DB", show=True):
    gateway = APIGateway("API Gateway")
    lb = ELB("Application Load Balancer")

    EKS_cluster = EKS("Golang API\non EKS Cluster")
    aurora_db = Aurora("Amazon Aurora Postgres DB")
    s3_bucket = S3("Amazon S3 Bucket")
    event_bridge = Eventbridge("EventBridge")
    sqs = SQS("SQS")
    lambda_function = Lambda("Lambda")

    usr = Users("API Users")
    provider = Users("Data Providers")

    provider >> s3_bucket
    s3_bucket >> event_bridge
    event_bridge >> sqs
    sqs >> lambda_function
    lambda_function >> aurora_db

    usr >> gateway >> lb >> EKS_cluster >> aurora_db
