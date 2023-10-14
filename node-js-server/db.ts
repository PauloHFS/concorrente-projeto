import { Client } from 'pg';

/**
 * Seed the database with the table and some initial data
 * Use this function only once
 * @returns {Promise<import('pg').QueryResult>}
 * @example
 * await seed().catch(console.error);
 */
const seed = () => {
  const query = `
    CREATE TABLE IF NOT EXISTS resources (
      id SERIAL PRIMARY KEY,
      name VARCHAR(255) NOT NULL,
      description TEXT NOT NULL,
      values INT[] NOT NULL
    );
  `;

  return client.query(query);
};

const client = new Client({
  host: process.env.HOST ?? 'localhost',
  port: 5432,
  database: 'postgres',
  user: 'postgres',
  password: 'postgres',
});

export { client, seed };
