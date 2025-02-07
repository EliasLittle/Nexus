./nexus-client publish value /testing/a 123
./nexus-client publish value /testing/b 456
./nexus-client publish value /testing/c 789
./nexus-client publish value /testing/d 101
./nexus-client publish value /testing/e 102
./nexus-client publish value /testing/f 103
./nexus-client publish value /testing/g 104
./nexus-client publish value /testing/h 105
./nexus-client publish value /testing/i 106
./nexus-client publish value /testing/j 107
./nexus-client publish value /testing/alpha 108
./nexus-client publish value /testing/alpha_beta 109
./nexus-client publish value /testing/alpha_gamma 110
./nexus-client publish value /testing/alphabet 'abc'

./nexus-client publish file /testing/dataset/a ./tests/example_a.csv
./nexus-client publish file /testing/dataset/b ./tests/example_b.csv
./nexus-client publish file /testing/dataset/c ./tests/example_c.csv
./nexus-client publish file /testing/dataset/d ./tests/example_d.csv

./nexus-client publish directory /testing/dataset/dir_a ./tests/dir_a
./nexus-client publish directory /testing/dataset/dir_b ./tests/dir_b

./nexus-client publish DBTable /testing/dataset/db_a postgres box-1 5432 tradingdb temperature_data

./nexus-client publish event /testing/events/random_data 'random_data'