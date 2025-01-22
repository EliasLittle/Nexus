from kafka import KafkaProducer
import random
import time
import json

# Initialize the Kafka producer
producer = KafkaProducer(bootstrap_servers='localhost:9092')
producer2 = KafkaProducer(bootstrap_servers='localhost:9092')

# Specify the Kafka topic to produce data to
topic = 'random_data'

# Specify the second Kafka topic to produce data to
topic2 = 'structured_data'

# Function to generate random data
def generate_random_data():
    return str(random.randint(1, 100))

# Function to generate structured random data
def generate_structured_data():
    return json.dumps({
        'id': random.randint(1, 100),
        'value': random.uniform(1.0, 100.0)
    }).encode()

# Produce random data to Kafka topic every second
try:
    while True:
        data = generate_random_data().encode()
        producer.send(topic, value=data)
        
        structured_data = generate_structured_data()
        producer2.send(topic2, value=structured_data)
        
        time.sleep(1)
except KeyboardInterrupt:
    producer.close()
    producer2.close()