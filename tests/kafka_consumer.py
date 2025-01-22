from kafka import KafkaConsumer

consumer = KafkaConsumer('random_data')

for message in consumer:
    print(message)