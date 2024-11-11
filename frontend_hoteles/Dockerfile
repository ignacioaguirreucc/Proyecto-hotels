# Usa una imagen de Node que incluya npm, adecuada para React
FROM node:18-alpine

# Define el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el package.json y el package-lock.json si existe
COPY package*.json ./

# Instala las dependencias
RUN npm install

# Copia el resto de la aplicación
COPY . .

# Expone el puerto en el que React servirá la aplicación
EXPOSE 3000

# Comando para iniciar la aplicación
CMD ["npm", "start"]
