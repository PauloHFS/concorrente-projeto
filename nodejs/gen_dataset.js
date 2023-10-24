const PATH = './dataset.json';
const SIZE = 100000;
const valoresSize = 4;

const fs = require('fs');

const timestamp = new Date().getTime();

const exist = fs.existsSync(PATH);

if (exist) {
  console.log('Arquivo já existe');
  console.log('Done (' + (new Date().getTime() - timestamp) / 1000 + 's)');
  process.exit(0);
}

const dataset = [];

for (let i = 0; i < SIZE; i++) {
  console.log(`Gerando registro ${i + 1} de ${SIZE}`);
  const id = i + 1;
  const nome = `User ${id}`;
  const descricao = `Descrição do usuário ${id}`;
  let valores = [];
  for (let j = 0; j < valoresSize; j++) {
    valores.push(Math.floor(100 + Math.random() * 100));
  }

  dataset.push({
    id,
    nome,
    descricao,
    valores,
  });
}

console.log('Gravando arquivo');
fs.writeFileSync(PATH, JSON.stringify(dataset));
console.log('Arquivo gravado com sucesso');
console.log('Done (' + (new Date().getTime() - timestamp) / 1000 + 's)');
