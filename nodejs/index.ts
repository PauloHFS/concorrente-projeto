import cluster from 'cluster';
import { once } from 'events';
import { Database } from 'sqlite3';

type Resource = {
  nome: string;
  descricao: string;
  valores: string[];
};

const masterCluster = async () => {
  let running = true;
  const numWorks = process.argv[2];

  if (!numWorks || isNaN(parseInt(numWorks))) {
    console.log('Error - please specify number of workers');
    process.exit(1);
  }

  cluster.on('exit', (worker, code, signal) => {
    if (running) return;
    console.log(`Processo filho ${worker.process.pid} morreu. Reiniciando...`);
    cluster.fork();
  });

  const dataset: Resource[] = require(process.cwd() + '/dataset.json');

  console.log('Dataset loaded with %s rows', dataset.length);

  const db = new Database('db.sqlite3');

  db.serialize(() => {
    db.run('DROP TABLE IF EXISTS resources');
    db.run(
      `CREATE TABLE IF NOT EXISTS resources (
      id INTEGER PRIMARY KEY, 
      nome TEXT NOT NULL, 
      descricao TEXT NOT NULL, 
      valores TEXT NOT NULL, 
      createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    )`
    );
  });
  db.close();

  console.log('Starting cluster with %s workers', numWorks);

  for (let i = 0; i < parseInt(numWorks); i++) {
    const worker = cluster.fork();
    await once(worker, 'message');
  }

  const queue: Resource[] = [];

  function distributeTask() {
    for (const id in cluster.workers) {
      const worker = cluster.workers[id];
      const task = queue.shift();
      if (task) {
        worker?.send({
          type: 'insert',
          resource: task,
        });
      }
    }
  }

  for (const resource of dataset) {
    queue.push(resource);
    distributeTask();
  }

  running = true;
  cluster.disconnect();
};

const childCluster = () => {
  const db = new Database('db.sqlite3');

  process.on('message', (message: { type: 'insert'; resource: Resource }) => {
    if (message.type == 'insert') {
      const resource: Resource = message.resource;
      db.run(
        `
        INSERT INTO resources (nome, descricao, valores)
        VALUES (?, ?, ?)
        `,
        [resource.nome, resource.descricao, resource.valores.join(', ')],
        function (err) {
          if (err) {
            console.log(err);
            return;
          }
          console.log('Inserted row with id %s', this.lastID);
        }
      );
    }
  });

  process.send?.({ type: 'ready' });
};

if (cluster.isPrimary) {
  masterCluster();
} else {
  childCluster();
}
