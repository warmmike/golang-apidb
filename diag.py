from diagrams import Diagram
from diagrams.aws.compute import Lambda
from diagrams.aws.database import Aurora
from diagrams.aws.storage import S3
from diagrams.aws.integration import Eventbridge
from diagrams.aws.compute import EKS
from diagrams.aws.network import APIGateway
from diagrams.aws.general import Users

with Diagram("Golang Project API DB", show=True):
    gateway = APIGateway("API Gateway")

    EKS_cluster = EKS("Golang API\non EKS Cluster")
    aurora_db = Aurora("Amazon Aurora Postgres DB")
    s3_bucket = S3("Amazon S3 Bucket")
    event_bridge = Eventbridge("EventBridge")
    lambda_function = Lambda("Lambda")

    usr = Users("API Users")
    provider = Users("Data Providers")

    provider >> s3_bucket
    s3_bucket >> event_bridge
    event_bridge >> lambda_function
    lambda_function >> aurora_db

    usr >> gateway >> EKS_cluster >> aurora_db
