FROM node:lts-buster
WORKDIR /usr/src/app
COPY . .
RUN npm install
EXPOSE 8080
RUN npm run build
CMD ["npm", "start"]

