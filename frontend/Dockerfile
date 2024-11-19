# Usa una imagen base de Node.js
FROM node:16-alpine

# Configura el directorio de trabajo
WORKDIR /app

# Copia el código fuente y archivos de configuración
COPY package*.json ./
COPY public ./public
COPY src ./src

# Instala dependencias y construye la app
RUN npm install
RUN npm run build

# Exponer el puerto
EXPOSE 3000

# Iniciar la app
CMD ["npm", "start"]
