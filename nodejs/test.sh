yarn build
node gen_dataset.js

workers=(1 2 4 6 8 10 12)

for numWorkers in "${workers[@]}"
do
  echo "$numWorkers - Starting Test with $numWorkers workers"
  output_path="tests/node-$numWorkers-workers.csv"
  mkdir -p tests

  yarn start $numWorkers

  docker exec -it concorrente-projeto-db bash -c "psql -U postgres -c \"COPY resources TO '/tmp/test.csv' WITH CSV HEADER;\""
  docker cp concorrente-projeto-db:/tmp/test.csv $output_path
done
