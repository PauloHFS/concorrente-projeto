import cluster from 'cluster';
import { once } from 'events';
import { Pool } from 'pg';

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

  const pool = new Pool({
    connectionString: 'postgresql://postgres:postgres@localhost:5432/postgres',
  });

  await pool.query('DROP TABLE IF EXISTS resources');
  await pool.query(`
    CREATE TABLE IF NOT EXISTS resources (  
      id SERIAL PRIMARY KEY,
      nome VARCHAR(255) NOT NULL,
      descricao TEXT NOT NULL,
      valores TEXT[] NOT NULL,
      created_at TIMESTAMP DEFAULT NOW()
    )
  `);

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
  const pool = new Pool({
    connectionString: 'postgresql://postgres:postgres@localhost:5432/postgres',
  });

  process.on('message', (message: { type: 'insert'; resource: Resource }) => {
    if (message.type == 'insert') {
      const resource: Resource = message.resource;
      pool
        .query(
          `
        INSERT INTO resources (nome, descricao, valores)
        VALUES ($1, $2, $3)
        `,
          [resource.nome, resource.descricao, resource.valores]
        )
        .then(() => {
          console.log('Inserted %s', resource.nome);
        })
        .catch(err => {
          console.log('Error inserting %s', resource.nome);
          console.log(err);
          process.exit(1);
        });
    }
  });

  process.send?.({ type: 'ready' });
};

if (cluster.isPrimary) {
  masterCluster();
} else {
  childCluster();
}
