node gen_dataset.js

#workers=(1 2 4 6 8 10 12)
workers=(1)

for num in "${workers[@]}"
do
  echo "$num - Starting Test with $num workers"
  numWorkers=$num
  database_path="db.sqlite3"
  output_path="tests/node-$numWorkers-workers.csv"
  mkdir -p tests

  go run main.go $numWorkers

  sqlite3 -header -csv $database_path "SELECT * FROM resources;" > $output_path
done
