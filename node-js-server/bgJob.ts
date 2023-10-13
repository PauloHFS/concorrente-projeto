import { client } from './db';
import { Queue } from './queue';

export const bgJob = () => {
  const resources = Queue.flush();

  if (resources.length > 0) {
    console.log('Background process running!');

    const query = `INSERT INTO resources (name, description, values) VALUES ${resources
      .map(
        resource =>
          `('${resource.name}', '${
            resource.description
          }', '{${resource.values.join(',')}}')`
      )
      .join(',')}`;

    client
      .query(query)
      .then(() => {
        console.log('Background process finished! Inserted resources');
      })
      .catch(error => {
        Queue.insertMany(resources);
        console.error("Background process finished! Couldn't insert resources");
      });
  }
};
