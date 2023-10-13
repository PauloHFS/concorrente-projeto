let queue: any[] = [];

const insert = (resource: any) => {
  queue.push(resource);
};

const insertMany = (resources: any[]) => {
  queue.push(...resources);
};

const flush = () => {
  const resources = [...queue];
  queue = [];
  return resources;
};

export const Queue = { insert, insertMany, flush };
