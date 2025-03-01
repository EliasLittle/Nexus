FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Add the missing Babel plugin
RUN npm install --save-dev @babel/plugin-proposal-private-property-in-object

# Add the missing grpc-web package
RUN npm install @improbable-eng/grpc-web

# Copy source code
COPY . .

# Debugging: List contents of the scripts directory to a file
RUN ls -la /app/scripts > /app/scripts_directory_listing.txt

# Generate protobuf files
#RUN apk add --no-cache protobuf protobuf-dev
#RUN chmod +x ./scripts/generate-proto.sh
#RUN ./scripts/generate-proto.sh

# Build the app
RUN npm run build

# Install serve to run the built app
RUN npm install -g serve

# Start the app
CMD ["serve", "-s", "build", "-l", "3000"] 