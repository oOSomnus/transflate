# Build with Node.js
FROM node:18 as build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json to install dependencies
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the source code to the working directory
COPY . .

# Build the React project
RUN npm run build

# Use Nginx to serve static files
FROM nginx:alpine

RUN apk add --no-cache bash

# Copy the generated static files to the default Nginx static directory
COPY --from=build /app/build /usr/share/nginx/html

# certs
COPY certs/translaterequest.com.key /etc/nginx/certs/translaterequest.com.key
COPY certs/translaterequest.com.pem /etc/nginx/certs/translaterequest.com.pem

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose Nginx's default port 80
EXPOSE 80 443

# Start the Nginx service
CMD ["nginx", "-g", "daemon off;"]
