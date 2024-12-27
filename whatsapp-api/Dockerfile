FROM node:18-alpine

# Definir o diretório de trabalho
WORKDIR /app

# Instalar Chromium e suas dependências
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Criar um usuário não-root
RUN addgroup -S pptruser && adduser -S pptruser -G pptruser

# Definir permissões apropriadas para o diretório de trabalho
RUN mkdir -p /home/pptruser/Downloads /app && \
    chown -R pptruser:pptruser /home/pptruser /app

# Definir variáveis de ambiente para o Chromium
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

# Copiar o package.json e o package-lock.json e instalar as dependências do npm
COPY package*.json ./
RUN npm install

# Copiar o restante do código
COPY . .

# # Trocar para o usuário não-root
# USER pptruser

# Expor a porta 3000
EXPOSE 3000

# Comando padrão para iniciar a aplicação
CMD ["npm", "start"]
