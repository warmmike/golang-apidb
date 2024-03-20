import boto3
import psycopg2
import json
import os
from urllib.parse import urlparse
import pprint

AWS_ACCESS_KEY_ID = None
AWS_SECRET_ACCESS_KEY = None
DB_PASSWORD = None

s3 = boto3.client(
    's3',
    aws_access_key_id=AWS_ACCESS_KEY_ID,
    aws_secret_access_key=AWS_SECRET_ACCESS_KEY,
    endpoint_url='http://localhost:4566',
)
sqs = boto3.client(
    'sqs',
    aws_access_key_id=AWS_ACCESS_KEY_ID,
    aws_secret_access_key=AWS_SECRET_ACCESS_KEY,
    endpoint_url='http://localhost:4566',
)

def get_objectmethods(object):
    object_methods = [method_name for method_name in dir(object) if callable(getattr(object, method_name))]
    pprint.pprint(object_methods)

def get_s3content(s3_path):
    o = urlparse(s3_path, allow_fragments=False)
    s3_object = s3.get_object(Bucket=o.netloc, Key=o.path.strip("/"))
    body = s3_object['Body']
    json_body = json.loads(body.read())
    return json_body

def get_sqsmessage(queue_url):
    # Receive message from SQS queue
    response = sqs.receive_message(
        QueueUrl=queue_url,
        AttributeNames=[
            'SentTimestamp'
        ],
        MaxNumberOfMessages=10,
        MessageAttributeNames=[
            'All'
        ],
        VisibilityTimeout=0,
        WaitTimeSeconds=0
    )
    
    message = None
    try:
        message = response['Messages'][0]
    except:
        pass
    if message != None:
        ret = get_s3content(message['Body'])
        put_db_record(ret)
        receipt_handle = message['ReceiptHandle']
    
        # Delete received message from queue
        sqs.delete_message(
            QueueUrl=queue_url,
            ReceiptHandle=receipt_handle
        )
        print('Received and deleted message: %s' % message)

def put_db_record(json_body):
    conn = psycopg2.connect(
        database="postgres",
        user='postgres',
        password=DB_PASSWORD,
        host='localhost',
        port='5432'
    )
    cursor = conn.cursor()
    cursor.execute('INSERT INTO movies(title,year,"cast",genres,href,extract,thumbnail,thumbnail_width,thumbnail_height) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)',(json_body['title'],json_body['year'],json_body['cast'],json_body['genres'],json_body['href'],json_body['extract'],json_body['thumbnail'],json_body['thumbnail_width'],json_body['thumbnail_height']))
    conn.commit()
    cursor.close()
    conn.close()

# Main Function
if __name__ == "__main__":
    if not ("AWS_ACCESS_KEY_ID" in os.environ) or not ("AWS_SECRET_ACCESS_KEY" in os.environ) or not ("DB_PASSWORD" in os.environ):
        print("Set AWS_ACCESSKEY_ID and AWS_SECRET_ACCESS_KEY and DB_PASSWORD")
        quit()
    else:
        AWS_ACCESS_KEY_ID = os.getenv("AWS_ACCESS_KEY_ID")
        AWS_SECRET_ACCESS_KEY = os.getenv("AWS_SECRET_ACCESS_KEY") 
        DB_PASSWORD = os.getenv("DB_PASSWORD")
    #get_objectmethods(s3)
    #get_objectmethods(sqs)
    queue_url = 'http://sqs.eu-west-2.localhost.localstack.cloud:4566/000000000000/moviedata'
    get_sqsmessage(queue_url)