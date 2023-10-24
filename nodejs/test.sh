yarn build

workers=(1 2 4 6 8 10 12)

for num in "${workers[@]}"
do
  echo "$num - Starting Test with $num workers"
  numWorkers=$num
  database_path="db.sqlite3"
  output_path="tests/node-$numWorkers-workers.csv"
  mkdir -p tests

  yarn start $numWorkers

  sqlite3 -header -csv $database_path "SELECT * FROM resources;" > $output_path
done
