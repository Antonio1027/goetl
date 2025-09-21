  #!/usr/bin/env bash
  echo "Creating MongoDB users..."
  mongosh --authenticationDatabase admin --host localhost -u root -p rootpassword --eval "db.getSiblingDB('goetldb').createUser({user: 'goetluser', pwd: 'goetlpass', roles: [{role: 'readWrite', db: 'goetldb'}]});"
  echo "MongoDB users created."