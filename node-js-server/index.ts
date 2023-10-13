import express from 'express';
import morgan from 'morgan';
import z from 'zod';
import { bgJob } from './bgJob';
import { client, seed } from './db';
import { Queue } from './queue';

const app = express();

app.use(express.json());
app.use(morgan('dev'));

const resourceSchema = z.object({
  name: z.string().min(5),
  description: z.string().min(10),
  values: z.array(z.number()),
});

app.post('/resources', (req, res) => {
  try {
    const resource = resourceSchema.parse(req.body);
    Queue.insert(resource);
    return res.status(202).json({
      message: 'Resource added to queue',
    });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return res.status(400).json({
        message: 'Invalid resource',
        errors: error.issues,
      });
    }
    return res.status(500).json({
      message: 'Internal server error',
    });
  }
});

client
  .connect()
  .then(() => {
    console.log('Connected to database!');
    return seed();
  })
  .then(() => {
    console.log('Database seeded!');
    app.listen(8080, () => {
      console.log('Server is listening on port 8080');
      setInterval(bgJob, 2000);
    });
  })
  .catch(error => {
    console.error(error);
    console.error('Could not start server!!!');
  });
